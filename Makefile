postgres:
	docker run --name postgres12 --network sample-bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable" --verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable" --verbose down

migrateup1:
	migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable" --verbose up 1

migratedown1:
	migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable" --verbose down 1

sqlc-win:
	docker run --rm -v $(shell pwd):/src -w /src kjconroy/sqlc generate

sqlc:
	sqlc generate

opendb:
	docker exec -it postgres12 psql -U root simple_bank
server:
	go run main.go
test:
	go test -v -cover ./...
mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/lamdangtung/golang-sample-bank/db/sqlc Store
build-image:
	docker build -t sample-bank:latest -f Dockerfile .

db_docs:
	dbdocs build docs/db.dbml

db_schema:
	dbml2sql docs\db.dbml -o docs\db.sql

sample-bank: 
	docker run --name sample-bank --network sample-bank-network -e GIN_MODE=release -e DB_SOURCE="postgresql://root:123456@postgres12:5432/simple_bank?sslmode=disable" -p 8080:8080 sample-bank:latest
.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc opendb server mock build-image sample-bank db_docs db_schema