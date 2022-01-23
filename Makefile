
# Versioning information
GIT_COMMIT := $(shell git rev-parse HEAD)
GIT_BRANCH := $(shell git name-rev --name-only HEAD | sed "s/~.*//")

## Gets the current tag name or commit SHA
VERSION ?= $(shell git describe --tags ${COMMIT} 2> /dev/null || echo "$(GIT_COMMIT)")

## Gets the -ldflags for the go build command, this lets us set the version number in the binary
ROOT := github.com/yukitsune/maestro
LD_FLAGS := -X '$(ROOT).Version=$(VERSION)'

## Whether the repo has uncommitted changes
GIT_DIRTY := false
ifneq ($(shell git status -s),)
	GIT_DIRTY=true
endif

# Docker stuff

## Common docker build args
DOCKER_BUILD_ARGS := \
	--build-arg GIT_COMMIT="$(GIT_COMMIT)" \
	--build-arg GIT_BRANCH="$(GIT_BRANCH)" \
	--build-arg GIT_DIRTY="$(GIT_DIRTY)" \
	--build-arg VERSION="$(VERSION)" \

## The base docker-compose command
PROJECT_NAME := maestro
DOCKER_COMPOSE_CMD := docker-compose --project-name $(PROJECT_NAME) --file ./deployments/docker-compose.yml --file ./deployments/docker-compose.development.yml --env-file ./configs/.env

# Commands

.DEFAULT_GOAL := help

.PHONY: help
help: ## Shows this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Builds all programs and places their binaries in the bin/ directory
	mkdir -p bin
	go build -ldflags="$(LD_FLAGS)" -o ./bin/  ./cmd/...

.PHONY: test
test: compose-fresh-detach ## Runs all tests
	go test ./...

.PHONY: clean
clean: ## Removes the bin/ directory
	rm -rf bin

.PHONY: build-container
build-container: ## Builds the maestro docker container
	docker build \
		-t maestro \
		-f build/package/maestro/Dockerfile \
		$(DOCKER_BUILD_ARGS) \
		.

.PHONY: compose
compose: ## Runs docker compose
	$(DOCKER_COMPOSE_CMD) up

.PHONY: compose-detach
compose-detach: ## Runs docker compose in detached mode
	$(DOCKER_COMPOSE_CMD) up --detach

.PHONY: compose-fresh
compose-fresh: ## Rebuilds the containers and forces a recreation
	$(DOCKER_COMPOSE_CMD) up --build --force-recreate

.PHONY: compose-fresh-detach
compose-fresh-detach: ## Rebuilds the containers and forces a recreation in detached mode
	$(DOCKER_COMPOSE_CMD) up --build --force-recreate --detach

.PHONY: compose-down
compose-down: ## Tears down the docker instances created by compose-up
	$(DOCKER_COMPOSE_CMD) down
