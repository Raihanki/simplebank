docker-postgres:
	docker run --name postgres-bank -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres

createdb:
	docker exec -it postgres-bank createdb --username=postgres --owner=postgres simplebank

dropdb:
	docker exec -it postgres-bank dropdb --username=postgres simplebank

migrateup:
	goose -dir=db/migrations postgres "postgres://postgres:secret@localhost:5432/simplebank" up

migratedown:
	goose -dir=db/migrations postgres "postgres://postgres:secret@localhost:5432/simplebank" down

sqlc:
	sqlc generate

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/raihanki/simplebank/db/sqlc Store

.PHONY: createdb dropdb docker-postgres migrateup migratedown sqlc server mock
