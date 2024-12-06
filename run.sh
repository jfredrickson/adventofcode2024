#!/usr/bin/env bash

# Run the Go program in the given directory, watching for changes.

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <directory>"
  exit 1
fi

DIR="$1"

if [ ! -d "$DIR" ]; then
  echo "Error: $DIR is not a valid directory."
  exit 1
fi

if [[ "$(uname -s)" == "Linux" ]] ; then
  if ! command -v inotifywait &>/dev/null; then
    echo "Error: inotifywait is not installed. Please install inotify-tools."
    exit 1
  fi
  CMD="inotifywait -rq -e modify,create,delete $DIR --format '%w%f'"
fi

if [[ "$(uname -s)" == "Darwin" ]] ; then
  if ! command -v fswatch &>/dev/null; then
    echo "Error: fswatch is not installed. Please install fswatch."
    exit 1
  fi
  CMD="fswatch --one-event --no-defer -r $DIR"
fi

echo "========================================"
echo "Watching directory: $DIR"
echo "Press Ctrl+C to stop."
echo "========================================"

trap "echo 'Okay, done!'; exit 0" SIGINT

go run "$DIR"

while true; do
  $CMD | while read -r file; do
    echo
    echo "========================================"
    echo "Detected change in: $file"
    echo "Running: go run $DIR"
    echo "========================================"

    go run "$DIR"
  done
done
