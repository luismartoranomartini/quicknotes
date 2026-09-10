include .env
export

server:
	@go run ./cmd/http/

dbup:
	@docker compose up -d

dbdown:
	@docker compose down

migrate-up:
	@migrate -database "postgres://${DB_USER}:${DB_PASSWORD}@127.0.0.1/${DB_USER}?sslmode=disable" -path db/migrations up

migrate-down:
	@migrate -database "postgres://${DB_USER}:${DB_PASSWORD}@127.0.0.1/${DB_USER}?sslmode=disable" -path db/migrations down

postgres:
	@docker compose exec db psql -U postgres


.PHONY: server dbup dbdown
