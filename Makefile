server:
	@go run ./cmd/http/

dbup:
	@docker compose up -d

dbdown:
	@docker compose down 

migrate-up:
	@migrate -database "postgres://postgres:secret@localhost/postgres?sslmode=disable" -path db/migrations up 

migrate-down:
	@migrate -database "postgres://postgres:secret@localhost/postgres?sslmode=disable" -path db/migrations down

postgres:
	@docker compose exec db psql -U postgres


.PHONY: server exp dbup dbdown 
