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

if ! command -v inotifywait &>/dev/null; then
  echo "Error: inotifywait is not installed. Please install inotify-tools."
  exit 1
fi

echo "========================================"
echo "Watching directory: $DIR"
echo "Press Ctrl+C to stop."
echo "========================================"

go run "$DIR"

while true; do
  inotifywait -rq -e modify,create,delete "$DIR" --format '%w%f' |
    while read -r file; do
      echo
      echo "========================================"
      echo "Detected change in: $file"
      echo "Running: go run $DIR"
      echo "========================================"

      go run "$DIR"
    done
done
