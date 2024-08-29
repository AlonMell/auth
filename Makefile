run:
	go run ./cmd/providerHub/main.go \
    		--config="./config/config.yaml"
migrate:
	go run ./cmd/migrator/main.go \
    		--path="./internal/" \
    		--table="migrations_history" \
    		--major=0 \
    		--minor=0