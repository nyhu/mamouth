#!/bin/bash -e

# The -o pipefail option is important for the trap to be executed if the "go test" command fails
set -o pipefail

: ${TEST_RESULTS:=/tmp/test-results}

mkdir -p "${TEST_RESULTS}"

trap "go-junit-report <${TEST_RESULTS}/go-test.out > ${TEST_RESULTS}/go-test-report.xml" EXIT
go test ./... -v -race | tee ${TEST_RESULTS}/go-test.out
