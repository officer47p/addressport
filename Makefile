build_api:
	@ go build -o bin/api cmd/api/main.go

run_api: build_api
	@./bin/api --listenAddr :3001

build_client:
	@ go build -o bin/client cmd/client/main.go

run_client: build_client
	@./bin/client

build_listener:
	@ go build -o bin/listener cmd/listener/main.go

run_listener: build_listener
	@./bin/listener --listenAddr :3001

remove_database:
	@docker-compose -f mongodb-docker-compose.yml rm -f

start_new_database: remove_database
	@docker-compose -f mongodb-docker-compose.yml up --build

seed:
	@go run scripts/seed.go

test:
	@go test -v ./...