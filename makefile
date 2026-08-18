server:
	go run ./cmd/http/

dbup:
	docker compose up -d

dbdown:
	docker compose down 

migrate-up:
	@migrate -database "postgres://postgres:secret@localhost/postgres?sslmode=disable" -path db/migrations up 


.PHONY: server exp dbup dbdown 
