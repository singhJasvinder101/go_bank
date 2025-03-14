postgres:
	docker run --name postgres -dp 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=123 postgres:17.4-alpine3.21

createDB:
	docker exec -it postgres createdb --username=postgres --owner=postgres go_bank

dropDB:
	docker exec -it postgres dropdb go_bank

migrateup:
	migrate -path db/migrations -database "postgresql://postgres:123@localhost:5432/go_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migrations -database "postgresql://postgres:123@localhost:5432/go_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migrations -database "postgresql://postgres:123@localhost:5432/go_bank?sslmode=disable" -verbose down    

migratedown1:
	migrate -path db/migrations -database "postgresql://postgres:123@localhost:5432/go_bank?sslmode=disable" -verbose down 1    

migratedirty:
	migrate -path db/migrations -database "postgresql://postgres:123@localhost:5432/go_bank?sslmode=disable" -verbose force 1

migratestatus:
	migrate -path db/migrations -database "postgresql://postgres:123@localhost:5432/go_bank?sslmode=disable" -verbose version

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createDB dropDB migrateUp sqlc