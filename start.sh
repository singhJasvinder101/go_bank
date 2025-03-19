#!/bin/sh
echo "Starting the application"

set -e

echo "run db migration"
/usr/local/bin/migrate -path /app/db/migrations -database "postgresql://postgres:123@postgres:5432/go_bank?sslmode=disable" -verbose up

echo "run the main application"
# run argument passed to docker run command(entrypoint) that is 
# cmd (./main) executable file
exec "$@"
