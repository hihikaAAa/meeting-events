BINARY := meeting-svc
CMD_DIR := ./cmd/meeting-events

.PHONY: test build

test:
	go test ./... -v -cover

build:
	go build -o bin/$(BINARY) $(CMD_DIR)

