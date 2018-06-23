#!/bin/bash -e

test -z "$(golint ./... | grep -v -e "should have comment" -e "vendor/" -e "mocks/" | tee /dev/stderr )"
