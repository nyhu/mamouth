#!/bin/sh

gofiles=$(find . -name '*.go' | grep -v -e '^./vendor')


unformatted=$(gofmt -s -l $gofiles)
[ -z "$unformatted" ] && exit 0

echo >&2 "Some files are not formatted, please run docker-compose run --rm fmt or edit the files"
echo >&2 "Unformatted files:"
echo >&2 $unformatted
exit 1
