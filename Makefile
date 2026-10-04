postgres:
	docker run --name postgres_db -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:latest

createdb:
	docker exec -it postgres_db createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres_db dropdb simple_bank


server:
	go run main.go

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down


migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./db/sqlc

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/thanaphat2005/simple-bank/db/sqlc Store

mysql:
	docker run --name mysql_db -p 3306:3306 -e MYSQL_ROOT_PASSWORD=secret -d mysql:latest

.PHONY: postgres createdb dropdb migratedown migrateup mysql server mock migrateup1 migratedown1