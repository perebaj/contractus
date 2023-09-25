SHELL = /bin/bash

# TOOLS VERSIONS
GO_VERSION=1.21.0
GOLANGCI_LINT_VERSION=v1.54.0

# configuration/aliases
version=$(shell git rev-parse --short HEAD)
base_image=registry.heroku.com/contractus/web
image=$(base_image):latest
devimage=contractus-dev
# To avoid downloading deps everytime it runs on containers
gopkg=$(devimage)-gopkg
gocache=$(devimage)-gocache
devrun=docker run $(devrunopts) --rm \
	-v `pwd`:/app \
	-v $(gopkg):/go/pkg \
	-v $(gocache):/root/.cache/go-build \
	$(devimage)

covreport ?= coverage.txt

all: lint test image

## run isolated tests
.PHONY: test
test:
	go test ./... -timeout 10s -race -shuffle on

## Format go code
.PHONY: fmt
fmt:
	goimports -w .

## builds the service
.PHONY: service
service:
	go build -o ./cmd/contractus/contractus ./cmd/contractus

## runs the service locally
.PHONY: run
run: service
	./cmd/contractus/contractus

## tidy up go modules
.PHONY: mod
mod:
	go mod tidy

## lint the whole project
.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run ./...

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


## generates coverage report
.PHONY: test/coverage
test/coverage:
	go test -count=1 -coverprofile=$(covreport) ./...

## generates coverage report and shows it on the browser locally
.PHONY: test/coverage/show
test/coverage/show: test/coverage
	go tool cover -html=$(covreport)

## Build the service image
.PHONY: image
image:
	docker build . \
		--build-arg GO_VERSION=$(GO_VERSION) \
		-t $(image)

## Build a production ready container image and run it locally for testing.
.PHONY: image/run
image/run: image
	docker run --rm -ti \
		-v $(gopkg):/go/pkg \
		$(image)

## Publish the service image
.PHONY: image/publish
image/publish: image
	docker push $(image)

## Releases to production
.PHONY: release
release: release_version=release-$(shell date '+%Y-%m-%d')-$(version)
release: release_image=$(base_image):$(release_version)
release:
	@echo "releasing from image: $(image)"
	@echo "release image:        $(release_image)"
	@echo "git tag:              $(release_version)"
	docker pull $(image)
	docker image tag $(image) $(release_image)
	docker push $(release_image)
	git tag -a $(release_version) -m "release to production: $(release_image)"
	git push origin $(release_version)

## Create the dev container image
.PHONY: dev/image
dev/image:
	docker build \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--build-arg GOLANGCI_LINT_VERSION=$(GOLANGCI_LINT_VERSION) \
		-t $(devimage) \
		-f Dockerfile.dev \
		.

## Create a shell inside the dev container
.PHONY: dev
dev: devrunopts=-ti
dev: dev/image
	$(devrun)

## run a make target inside the dev container.
dev/%: dev/image
	$(devrun) make ${*}

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
