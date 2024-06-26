PKG_NAME    := github.com/ArtalkJS/Artalk
BIN_NAME	:= ./bin/artalk
VERSION     ?= $(shell git describe --tags --abbrev=0 --match 'v*')
COMMIT_HASH ?= $(shell git rev-parse --short HEAD)

HAS_RICHGO  := $(shell which richgo)
GOTEST      ?= $(if $(HAS_RICHGO), richgo test, go test)
ARGS        ?= server

all: install build

install:
	go mod tidy

build:
	env CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC="x86_64-linux-musl-gcc" CXX="x86_64-linux-musl-g++" \
	go build \
    	-ldflags "-s -w -X $(PKG_NAME)/internal/config.Version=$(VERSION) \
        -X $(PKG_NAME)/internal/config.CommitHash=$(COMMIT_HASH)" \
        -o $(BIN_NAME)

build-frontend:
	./scripts/build-frontend.sh

run: all
	$(BIN_NAME) $(ARGS)

build-debug:
	@echo "Building Artalk $(VERSION) for debugging..."
	@env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 CC="x86_64-linux-musl-gcc" CXX="x86_64-linux-musl-g++" \
	 go build \
		-ldflags "-X $(PKG_NAME)/internal/config.Version=$(VERSION) \
		  -X $(PKG_NAME)/internal/config.CommitHash=$(COMMIT_HASH)" \
		-gcflags "all=-N -l" \
		-o $(BIN_NAME)

dev: build-debug
	$(BIN_NAME) $(ARGS)

test:
	$(GOTEST) -timeout 20m $(or $(TEST_PATHS), ./...)

test-coverage:
	$(GOTEST) -cover $(or $(TEST_PATHS), ./...)

test-coverage-html:
	$(GOTEST) -v -coverprofile=coverage.out $(or $(TEST_PATHS), ./...)
	go tool cover -html=coverage.out

update-i18n:
	go generate ./internal/i18n

update-conf-docs:
	go run ./internal/config/meta/gen -f ./docs/docs/guide/env.md

update-swagger:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g server/server.go --output ./docs/swagger --requiredByDefault
	pnpm -r swagger:build-http-client

docker-build:
	./scripts/docker-build.sh

docker-push:
	./scripts/docker-build.sh --push

test-frontend-e2e:
	./scripts/frontend-e2e-test.sh $(if $(REPORT), --show-report)

.PHONY: all install build build-frontend build-debug \
	dev test test-coverage test-coverage-html update-i18n \
	docker-build docker-push test-frontend-e2e;

push: build-frontend build docker-build docker-push