build:
	@ go build -o bin/api cmd/api/main.go

start: build
	@./bin/api

remove_database:
	@docker-compose -f mongodb-docker-compose.yml rm -f

start_new_database: remove_database
	@docker-compose -f mongodb-docker-compose.yml up --build

test:
	@go test -v ./...
