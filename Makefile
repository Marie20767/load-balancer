PORTS=$(shell go run scripts/ports/main.go)

start:
	go run main.go

start-servers:
	@echo "Starting servers..."
	@for port in $(PORTS); do \
		echo "Starting server on $$port"; \
		go run cmd/server/main.go --port=$$port & \
	done
	@wait

stop-servers:
	@echo "Stopping servers..."
	@for port in $(PORTS); do \
		lsof -ti tcp:$$port | xargs kill -9 || true; \
	done

lint:
	golangci-lint run
