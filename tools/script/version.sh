#!/usr/bin/env bash
set -e

TAG=`git tag | tail -n1`

commit_id="`git rev-parse HEAD`"

echo "$TAG-$commit_id"
