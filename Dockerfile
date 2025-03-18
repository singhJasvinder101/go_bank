# Build Stage
from golang:1.24.1-alpine3.21 as build
workdir /app
copy . .
run go build -o main .

# Run Stage
from golang:1.24.1-alpine3.21
workdir /app
copy --from=build /app/main .
copy app.env /app/app.env


expose 3000
cmd ["./main"]