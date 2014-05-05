GOOS=
GOARCH=

SRC_DIR=$(wildcard *.go)
SRC=$(wildcard *.go)

BINARIES=pgstatsd
BLDDIR=build

all: $(BINARIES)

fmt:
		go fmt ${SRC_DIR}

$(BLDDIR)/%:
		mkdir -p $(dir $@)
			GOOS=${GOOS} GOARCH=${GOARCH} go build ${GOFLAGS} -o $(abspath $@)

$(BINARIES): %: $(BLDDIR)/%

clean:
		rm -rf $(BLDDIR)

