#!/usr/bin/env bash

echo "mode: set" > coverage.out
cat *.out | grep -v mode: | sort -r | awk '{if($1 != last) {print $0;last=$1}}' >> coverage.out
