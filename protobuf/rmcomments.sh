#!/usr/bin/env bash

if [ -z "$1" ]; then
    echo "usage: rmcomments <path>"
    exit 2
fi

for FILE in "$@"; do
    START=$(($(grep -n 'Package protein' $FILE | cut -d ':' -f 1)-1))
    END=$(($(grep -n 'package protein' $FILE | cut -d ':' -f 1)-1))
    echo "removing lines $START to $END from $FILE"
    perl -i -ne "print if not $. == $START..$END" "$FILE"
done
