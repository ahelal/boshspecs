#!/bin/bash

set -eu

export PATH=$GOPATH/bin:$PATH

mkdir -p /go/src/github.com/ahelal/ || echo "OK" 

ln -fs $(pwd) /go/src/github.com/ahelal/boshspecs
pushd /go/src/github.com/ahelal/boshspecs

make clean
make deps
make lint
make test

popd
