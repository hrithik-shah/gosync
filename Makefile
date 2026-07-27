.PHONY: swag
swag:
	swag init -g cmd/server/main.go
	swagger generate markdown -f ./docs/swagger.yaml --output=./docs/api_docs.md
