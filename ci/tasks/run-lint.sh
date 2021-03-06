#!/bin/bash

set -eu

export PATH=$GOPATH/bin:$PATH

mkdir -p /go/src/github.com/ahelal/ || echo "OK" 

ln -fs $(pwd) /go/src/github.com/ahelal/boshspecs
pushd /go/src/github.com/ahelal/boshspecs

# hack to support caching. find a better away
VENDOR_CONCOURSE="$(pwd)/../vendor-concourse"
cp -r "${VENDOR_CONCOURSE}"/. vendor/
make deps
cp -r vendor/. "${VENDOR_CONCOURSE}"

make lint
popd
