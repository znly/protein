#!/usr/bin/env bash

if [ -z "$COVERALLS_TOKEN" ]; then
    exit 0 ## ignore coverage for PRs
fi

echo "mode: set" > coverage.out
cat *.out | grep -v mode: | sort -r | awk '{if($1 != last) {print $0;last=$1}}' >> coverage.out


goveralls -coverprofile=coverage.out -service travis-ci -repotoken "$COVERALLS_TOKEN"
