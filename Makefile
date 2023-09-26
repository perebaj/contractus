export GO_VERSION=1.21.1
export GOLANGCI_LINT_VERSION=v1.54.0

base_image=registry.heroku.com/contractus/web
version:=$(shell git rev-parse --short HEAD)
image:=$(base_image):latest
devrunopts:=--no-deps
container=contractus
devrun:=docker-compose run --rm $(devrunopts) $(container)

## Build the service
.PHONY: build
build:
	go build ./cmd/contractus

## Build service image
.PHONY: image
image:
	docker build . \
	--build-arg GO_VERSION=$(GO_VERSION) \
	-t $(image)

## Publish the service image
.PHONY: image/publish
image/publish: image
	docker push $(image)

## Run tests, if testcase=<testcase> only run that testcase
.PHONY: test
test:
	@echo "Running tests..."
	go test -run="$(testcase)" -cover ./...

## Run integration tests
.PHONY: integration-test
integration-test:
	go test -tags integration -run="$(testcase)" ./...

## Run lint
.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run ./... -v
	go run golang.org/x/vuln/cmd/govulncheck ./...

## Create a new migration, use make migration/new name=<migration_name>
.PHONY: migration/new
migration/new: param-name
	@echo "Creating new migration..."
	go run github.com/golang-migrate/migrate/v4/cmd/migrate \
		create \
		-dir ./postgres/migrations \
		-ext sql \
		-seq \
		$(name)

## Run migrations
.PHONY: migration/up
migration/up: param-CONTRACTUS_POSTGRES_URL
	@echo "Running migration up..."
	go run --tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate \
		-path=./postgres/migrations \
		-database=$(CONTRACTUS_POSTGRES_URL) \
		-verbose \
		up

## Step down a single migration
.PHONY: migration/down
migration/down: param-CONTRACTUS_POSTGRES_URL
	@echo "Running migration down..."
	go run --tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate \
		-path=./postgres/migrations \
		-database=$(CONTRACTUS_POSTGRES_URL) \
		-verbose \
		down 1

## Drop the database
.PHONY: migration/drop
migration/drop: param-CONTRACTUS_POSTGRES_URL
	@echo "Running migration drop..."
	go run --tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate \
		-path=./postgres/migrations \
		-database=$(CONTRACTUS_POSTGRES_URL) \
		-verbose \
		drop

## Clean containers, images and volumes
.PHONY: dev/clean
dev/clean:
	@echo "Cleaning containers, images and volumes..."
	@docker-compose down --rmi all --volumes --remove-orphans

## Build dev image service
.PHONY: dev/image
dev/image:
	docker-compose build

## Start containers, additionaly you can provide rebuild=true to force rebuild
.PHONY: dev/start
dev/start:
	@echo "Starting development server..."
	@if [ "$(rebuild)" = "true" ]; then \
		docker-compose up -d --build; \
	else \
		docker-compose up -d; \
	fi

## Stop containers
.PHONY: dev/stop
dev/stop:
	@echo "Stopping development server..."
	@docker-compose down

## Restart containers, if container=<name> is provided only it will be restarted
.PHONY: dev/restart
dev/restart: container=
dev/restart:
	@echo "Restarting development server..."
	@docker-compose restart $(container)

## Show logs, if container=<name> is provided logs for only that container will be shown
.PHONY: dev/logs
dev/logs:
	@echo "Showing logs..."
	@docker-compose logs -f $(container)

## Access the container
.PHONY: dev
dev:
	@$(devrun) bash

## run a make target inside the dev container.
dev/%:
	@$(devrun) make ${*}

## Display help for all targets
.PHONY: help
help:
	@awk '/^.PHONY: / { \
		msg = match(lastLine, /^## /); \
			if (msg) { \
				cmd = substr($$0, 9, 100); \
				msg = substr(lastLine, 4, 1000); \
				printf "  ${GREEN}%-30s${RESET} %s\n", cmd, msg; \
			} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

param-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Param \"$*\" is missing, use: make $(MAKECMDGOALS) $*=<value>"; \
		exit 1; \
	fi
