IMAGE ?= sgt-report-service
PORT ?= 8086

.PHONY: run test lint format tidy docs image container

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

image:
	docker build -t $(IMAGE) .

container: image
	docker run --rm --env-file .env -p $(PORT):8080 $(IMAGE)
