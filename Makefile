run:
	go run ./cmd/providerHub/main.go \
    		--config="./config/providerHub/config.yaml"
build:
	go build -o ./bin/providerHub ./cmd/providerHub/main.go
migrate:
	go run ./cmd/migrator/main.go \
			--config="./config/migrator/config.yaml" \
    		--path="./migrations/postgresql"
dockerRun:
	docker run -p 8080:8080 --name main provider-hub
dockerBuild:
	docker build ./deployments -t provider-hub
swagger:
	swag init -g ./cmd/providerHub/main.go -o ./api