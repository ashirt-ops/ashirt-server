#! /usr/bin/env bash

# exit on error
set -e

# set cwd to ashirt root
cd "$(dirname "$0")/.."

now=$(date -u +%Y%m%d%H%M%S)

desc="$*"

if [ -z "$desc" ]; then
	read -p 'Migration name: ' desc
fi

# sanitize description (spaces -> dashes)
desc=${desc// /-}

migrationsPath="./backend/migrations"

filename=$migrationsPath/$now-$desc.sql


touch $filename
echo "-- +migrate Up" >> $filename
echo "" >> $filename
echo "-- +migrate Down" >> $filename
