#!/usr/bin/env bash

set -e
echo "" > coverage.txt

for d in $(go list ./... | grep -v vendor); do
    additionalArgs=""
    if [[ "$d" =~ test$ ]]; then
        additionalArgs="-coverpkg=${d%test}"
    fi

    go test -race $additionalArgs -coverprofile=profile.out -covermode=atomic $d
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done