DB_USER ?= postgres
DB_PASSWORD ?= changeme
DB_HOST ?= localhost
DB_PORT ?= 5432
#VERSION=`git rev-parse --short HEAD`
VERSION?=1.0.0

# DATABASE
recreate_database:
	psql postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/postgres -f ./database/migrations/todo-database.sql

add_data_to_database:
	psql postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/postgres -f ./database/seeds/data.sql


# GO APP
integration-test:
	go test ./... -tags integration

tidy:
	go mod tidy

race_test:
	 go test -race -timeout=60s -count 1 ./...

go-build:
	 GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.Version=$(VERSION)" -o bin/todo-app cmd/main.go

go-run:
	./bin/todo-app

linter:
	 golangci-lint run  ./...

# DOCKER
dc-postgre:
	docker-compose -f postgre.yaml up --build

dc-test-run:
	 docker-compose --env-file docker.test.env up --build

dc-down:
	docker-compose down

#
docker_compose_restart: dc-down dc-test-run

create_db_for_test: recreate_database add_data_to_database

run_test:
	docker-compose -f postgre.yaml up --build -d
	psql postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/postgres -f ./database/migrations/todo-database.sql
	psql postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/postgres -f ./database/seeds/data.sql
	docker-compose -f postgre.yaml down



.PHONY: race_test build