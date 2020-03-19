#!/usr/bin/env bash

set -e
echo "" > coverage.txt


go test ./... -mod=vendor -timeout=2m -parallel=4 -coverprofile=profile.out -covermode=atomic ${TESTARGS} $d
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
