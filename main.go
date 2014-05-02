/*
Grab postgres statement stats and shove them into statsd
*/
package main

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // underscore means import for side-effects only, brings no symbols into scope
)

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

func (db *Database) GetStats() []StatRow {
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
	config := ReadConfig("./conf.json")
	db := DBInit(config.ConnectionURL)
	defer db.con.Close()
	stats := db.GetStats()
	for _, s := range stats {
		log.Printf("hit_percent: %f average_time: %f calls: %d total_time: %f", s.HitPercent.Float64, s.AverageTimeMS/1000, s.CallCount, s.TotalTimeMS)
	}
}
