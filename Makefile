-include .env
export $(shell sed 's/=.*//' .env)

start:
	go run main.go

start-servers:
	go run cmd/server/main.go --port=$(SERVER_1_PORT) &
	go run cmd/server/main.go --port=$(SERVER_2_PORT) &
	go run cmd/server/main.go --port=$(SERVER_3_PORT) &
	wait

stop-servers:
	@lsof -ti tcp:$(SERVER_1_PORT) | xargs kill -9 || true
	@lsof -ti tcp:$(SERVER_2_PORT) | xargs kill -9 || true
	@lsof -ti tcp:$(SERVER_3_PORT) | xargs kill -9 || true

lint:
	golangci-lint run