postgres:
	docker run --name postgres -dp 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=123 postgres:17.4-alpine3.21

createDB:
	docker exec -it postgres createdb --username=postgres --owner=postgres go_bank

dropDB:
	docker exec -it postgres dropdb go_bank

migrateUp:
	migrate -path db/migrations -database "postgresql://postgres:123@localhost:5432/go_bank?sslmode=disable" -verbose up

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createDB dropDB migrateUp sqlc