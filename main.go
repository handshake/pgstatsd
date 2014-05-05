/*
Grab postgres statement stats and shove them into statsd
*/
package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // underscore means import for side-effects only, brings no symbols into scope
)

// Path to config file
var configPath string

func init() {
	flag.StringVar(&configPath, "config", "/etc/pgstatsd/conf.json", "path to config.json")
}

type Database struct {
	con *sqlx.DB
}

func DBInit(connStr string) *Database {
	// move this to some db wrapper with nicer methods
	_db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return &Database{con: _db}
}

type StatRow struct {
	Query         string          `db:"query"`
	CallCount     int             `db:"calls"`
	TotalTimeMS   float32         `db:"total_time"`
	RowCount      int             `db:"rows"`
	AverageTimeMS float32         `db:"average_time"`
	HitPercent    sql.NullFloat64 `db:"hit_percent"`
}

type SizeStat struct {
	Relation string `db:"relation"`
	Bytes    int64  `db:"bytes"`
}

func (db *Database) GetBiggestRelation() int64 {
	stats := db.GetSizeStats(1)
	return stats[0].Bytes
}

func (db *Database) GetSizeStats(limit int) []SizeStat {
	// http://www.postgresql.org/docs/9.3/static/catalog-pg-class.html
	stats := []SizeStat{}
	err := db.con.Select(&stats, "SELECT nspname || '.' || relname AS relation, "+
		"pg_relation_size(C.oid) AS bytes "+
		"FROM pg_class C "+
		"LEFT JOIN pg_namespace N ON (N.oid = C.relnamespace) "+
		"WHERE nspname NOT IN ('pg_catalog', 'pg_toast', 'information_schema') "+
		"ORDER BY pg_relation_size(C.oid) DESC "+
		"LIMIT $1", limit)
	if err != nil {
		log.Fatal(err)
	}
	return stats
}

func (db *Database) GetStatementStats() []StatRow {
	// http://www.postgresql.org/docs/9.3/static/pgstatstatements.html
	stats := []StatRow{}
	err := db.con.Select(&stats, "SELECT query, calls, total_time, rows, "+
		"(total_time/calls) AS average_time, "+
		"(100.0 * shared_blks_hit / nullif(shared_blks_hit + shared_blks_read, 0)) AS hit_percent "+
		"FROM pg_stat_statements ORDER BY total_time DESC LIMIT 10")
	if err != nil {
		log.Fatal(err)
	}
	return stats
}

func main() {
	flag.Parse()

	// read our configuration
	config := ReadConfig(configPath)

	// connect to the database
	db := DBInit(config.PG.ConnectionString)
	defer db.con.Close()

	// connect to statsd
	statsd, err := statsd.New(config.ST.ConnectionString, config.ST.Prefix)
	// handle any errors
	if err != nil {
		log.Fatal(err)
	}
	// make sure to clean up
	defer statsd.Close()

	err = statsd.Gauge("biggest_relation_bytes", db.GetBiggestRelation(), 1.0)
	if err != nil {
		log.Printf("Error sending metric: %+v", err)
	}
	/*
		statementStats := db.GetStatementStats()
		for _, s := range statementStats {
			log.Printf("%#v", s)
		}

		sizeStats := db.GetSizeStats(50)
		for _, s := range sizeStats {
			log.Printf("%#v", s)
		}
	*/

}
