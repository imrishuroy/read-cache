postgres:
	docker run --name read-cache -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=IWSIWDF2024 -d postgres

createdb:
	docker exec -it read-cache createdb --username=root --owner=root read_cache_db

dropdb:
	docker exec -it read-cache dropdb read_cache_db

migrateup:
	migrate -path db/migration -database "postgres://root:IWSIWDF2024@localhost:5432/read_cache_db?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://root:IWSIWDF2024@localhost:5432/read_cache_db?sslmode=disable" -verbose down

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/imrishuroy/read-cache/db/sqlc Store


.PHONY: postgres createdb dropdb migrateup migratedown test server mock new_migration