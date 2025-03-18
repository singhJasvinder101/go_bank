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

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/singhJasvinder101/go_bank/db/sqlc Store

# explicityl writing DB_Source bcoz IP of postgres container is not static
# instead of rebuilding everytime image by changing app.env add 
# it as env argument when running contianer
container:
	docker run --name go_bank -dp 3000:3000 -e DB_SOURCE=postgresql://postgres:123@postgres:5432/go_bank?sslmode=disable go_bank:latest

disconnect_custom_network:
	docker network --help
	docker network ls
	docker network disconnect go_bank postgres
	docker network inspect go_bank
	docker network disconnect go_bank go_bank
	docker network rm go_bank

custom_network:
	docker network --help
	docker network create go_bank
	docker network connect go_bank postgres
	docker inspect go_bank
	docker inspect postgres
	docker stop go_bank
	docker rm go_bank
	docker run --name go_bank -dp 3000:3000 \
		--network go_bank \
		-e DB_SOURCE=postgresql://postgres:123@postgres:5432/go_bank?sslmode=disable \
		go_bank:latest
# Automatic dns resolution by container name in my custom network


.PHONY: postgres createDB dropDB migrateUp sqlc test server mock
