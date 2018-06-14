#!/bin/bash

set -eu

mkdir -p /go/src/github.com/ahelal/ || echo "OK" 

ln -fs $(pwd) /go/src/github.com/ahelal/boshspecs
pushd /go/src/github.com/ahelal/boshspecs

make clean
make deps
make test

popd
