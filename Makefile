GO=go
DOCKER=docker
MODULE=github.com/AlonMell/grovelog

.PHONY: help run build migrate docker-run docker-build lint clean

run:
	$(GO) run ./cmd/auth/main.go \
						-config="./config/auth/config.yaml"

build:
	$(GO) build -o ./bin/auth ./cmd/auth/main.go

migrate:
	$(GO) run ./cmd/migrator/main.go \
			-config="./config/migrator/config.yaml" \
						-migrations="./migrations/postgresql"

docker-run:
	$(DOCKER) run -p 8080:8080 --name main provider-hub

docker-build:
	$(DOCKER) build ./deployments -t provider-hub

lint:
	golangci-lint run

clean:
	$(GO) clean
	rm -f coverage.out

help:
	@echo "Available targets:"
	@echo "  help         - Display this help message"
	@echo "  run          - Run the auth service with the specified config"
	@echo "  build        - Build the auth service binary"
	@echo "  migrate      - Run database migrations"
	@echo "  docker-run   - Run the provider-hub container"
	@echo "  docker-build - Build the provider-hub Docker image"
	@echo "  lint         - Run golangci-lint"
	@echo "  clean        - Clean build artifacts"