.PHONY: run test lint format tidy docs

run:
	swag init -g main.go -o docs --parseInternal
	go run .

test:
	go test ./... -cover

format:
	golangci-lint fmt

lint:
	golangci-lint run

tidy:
	go mod tidy

docs:
	swag init -g main.go -o docs --parseInternal
