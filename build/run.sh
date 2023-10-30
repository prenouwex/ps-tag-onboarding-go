#!/usr/bin/env bash

source ./build/env.sh

while getopts "b" opt; do
  case $opt in
  b) SHOULD_BUILD=true ;;
  esac
done

if [ ! -f $BINARY ] || [ -n "$SHOULD_BUILD" ]; then
	echo "Building $PROJECT_NAME..."
    ./build/build.sh || exit 1
fi

echo "Attempting to stop if already running..."
pkill -f $BINARY

echo "Press CTRL+C to exit..."
$BINARY