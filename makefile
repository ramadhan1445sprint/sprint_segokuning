include .env
export

.phony: build migrate rollback

build:
	GOOS=linux GOARCH=amd64 go build -o main main.go

migrate:
	migrate -database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=verify-full&rootcert=ap-southeast-1-bundle.pem" -path db/migrations up

rollback:
	migrate -database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=verify-full&rootcert=ap-southeast-1-bundle.pem" -path db/migrations down

migrate-dev:
	migrate -database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -path db/migrations up

rollback-dev:
	migrate -database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -path db/migrations down