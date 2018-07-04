#!/bin/bash

set -eu

PWD=$PWD
if [ "${INTEGRATION_TEST}x" = "x" ]; then
    echo "No integration test passed"
    exit 1
fi

if ! [ -d "${PWD}/boshspecs-repo/test/integration/${INTEGRATION_TEST}" ]; then
    echo "No integration test ${INTEGRATION_TEST} found"
    exit 1
fi

# Copy binary
cp "$PWD/build-bucket/boshspecs*" /usr/local/bin/boshspecs

cd ${PWD}/boshspecs-repo/test/integration/${INTEGRATION_TEST}
inspec exec .
