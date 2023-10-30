#!/usr/bin/env bash

source ./build/env.sh

if [ ! -d $OUTPUT_DIRECTORY ]; then
    echo "Creating Directory \`$OUTPUT_DIRECTORY\`.."
    mkdir -p $OUTPUT_DIRECTORY
fi

go get -d ./...
go build -o $BINARY ./cmd/ps-tag-onboarding-go