#!/bin/sh

set -e
cd "$(dirname $0)"

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

while true; do
  echo "Building dev.go"
  if ! go build -o /tmp/dev bin/dev/dev.go; then
    echo "Failed to build. Waiting for file change before continuing"
    inotifywait -r -e modify -e create -e delete . 2> /dev/null
    continue
  fi

  echo "Starting dev.go"
  /tmp/dev &
  PID=$!
  inotifywait -r -e modify -e create -e delete .
  kill $PID
done
