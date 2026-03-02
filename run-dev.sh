#!/bin/sh

set -e
cd "$(dirname "$0")"

touch /typescript-dtos/dtos.ts

# Wait for DB
while ! nc -z db 3306 2>/dev/null; do
  sleep 1
  echo Waiting for mysql
done

# add seed assets if content store is empty (and we have default assets)
echo "Checking if content store is available and empty"
if [ -z "$(ls -A -- "/tmp/contentstore")" ]; then
  echo "Content store is empty. Trying to add seed content"
  if test -e /app/dev_seed_data/images; then
    ln -s /app/dev_seed_data/images/* /tmp/contentstore
    echo "Seeded Content added"
  fi
else
  echo "Content has been populated. Assuming that includes seed content..."
fi

# Generate frontend typescript types from dtos.go anytime dtos.go changes
(while true; do
  echo "Generating typescript types from DTOs"
  if ! go run ./internal/dtos/gentypes > /typescript-dtos/dtos.ts; then
    echo "Failed to generate frontend DTOs. Waiting for file change before continuing"
  fi
  inotifywait -r -e modify -e create -e delete ./internal/dtos 2> /dev/null
done) &
DTO_LOOP_PID=$!

# Recompile & run the server on change
(while true; do
  echo "Building dev.go"
  if ! go build -o /tmp/dev ./cmd/ashirt-server/; then
    echo "Failed to build. Waiting for file change before continuing"
    inotifywait -r -e modify -e create -e delete . 2> /dev/null
    continue
  fi

  echo "Starting dev.go"
  /tmp/dev &
  SERVER_PID=$!
  inotifywait -r -e modify -e create -e delete .
  kill $SERVER_PID
done) &
SERVER_LOOP_PID=$!

trap "kill $DTO_LOOP_PID $SERVER_LOOP_PID; exit 1" INT
wait
