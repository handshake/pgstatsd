#!/bin/sh

GO_VERSION=1.2.1
INSTALL_DIR=/usr/local
GOROOT=/usr/local/go

function install_go {
    # Download Go 1.2 source
    wget https://go.googlecode.com/files/go${GO_VERSION}.src.tar.gz -O /tmp/go${GO_VERSION}.src.tar.gz

    # Untar into $INSTALL_DIR
    pushd /tmp
    tar -zxf go${GO_VERSION}.src.tar.gz -C $INSTALL_DIR
    popd

    # Compile Go
    pushd /usr/local/go/src
    ./all.bash

    # Cleanup the tarball
    rm /tmp/go${GO_VERSION}.src.tar.gz
}

# Build for 64-bit OSX and Linux
PLATFORMS="darwin/amd64 linux/amd64"

function go_crosscompile_build {
    echo "Installing Go for $1"
    GOOS=${1%/*}
    GOARCH=${1#*/}
    cd ${GOROOT}/src ; GOOS=${GOOS} GOARCH=${GOARCH} ./make.bash --no-clean 2>&1
}

# Check for Go installation
which go > /dev/null 2>&1
if [[ ! $? -eq 0 ]]; then
    install_go
fi

for PLATFORM in $PLATFORMS; do
    TARGET_DIR=${PLATFORM//\//_}
    if [[ ! -d $GOROOT/pkg/$TARGET_DIR ]]; then
        go_crosscompile_build $PLATFORM
    fi
done

