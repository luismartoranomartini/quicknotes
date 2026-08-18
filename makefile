server:
	go run ./cmd/http/

exp:
	go run ./cmd/exp/

dbup:
	docker compose up

dbdown:
	docker compose down 

.PHONY: server exp dbup dbdown 
