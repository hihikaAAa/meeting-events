BINARY := meeting-svc
CMD_DIR := ./cmd/meeting-events

.PHONY: all build test run-prod docker-up docker-down docker-logs test-e2e clean

all: build

build:
	go build -o bin/$(BINARY) $(CMD_DIR)

test-e2e:
	go test ./tests -v -count=1

run-prod: 
	CONFIG_PATH=./config/prod.yaml go run $(CMD_DIR)

docker-up: 
	docker compose up -d --build

docker-down:
	docker compose down -v

docker-logs:
	docker compose logs -f app migrate db

clean:
	rm -rf bin
