# Build Stage
from golang:1.24.1-alpine3.21 as build
workdir /app
copy . .
run go build -o main .
run apk add curl
run apk add --no-cache curl && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz \
    -o migrate.tar.gz && \
    tar -xvzf migrate.tar.gz -C /app && \
    chmod +x /app/migrate && \
    rm migrate.tar.gz


# Run Stage
from golang:1.24.1-alpine3.21
workdir /app
copy --from=build /app/main .
copy app.env /app/app.env
copy --from=build /app/migrate /usr/local/bin/migrate
copy db/migrations /app/db/migrations

copy start.sh /app/start.sh
run chmod +x /app/start.sh

expose 3000

cmd ["/app/main"]
entrypoint ["/app/start.sh"]


# cmd when used with docker main entrypoint act as additional default
# parameter of entrypoint script & it can be overriden during runtime 
# e.g docker run <Container Name> version
# similar to entrypoint ["/app/start.sh", "./main"]