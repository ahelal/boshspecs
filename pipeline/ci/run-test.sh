#!/bin/sh

set -eu

make clean
make deps
make test
