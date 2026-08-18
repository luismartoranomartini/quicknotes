server:
	go run ./cmd/http/

dbup:
	docker compose up

dbdown:
	docker compose down 

.PHONY: server exp dbup dbdown 
