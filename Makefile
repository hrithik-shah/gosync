.PHONY: swag migrate gen
swag:
	swag init -g cmd/server/main.go
	swagger generate markdown -f ./docs/swagger.yaml --output=./docs/api_docs.md

migrate:
	go run ./cmd/server migrate

gen:
	go run ./cmd/gen