ifneq (,$(wildcard .env))
	include .env
	export
endif

CONTAINER_NAME=localdev
APP_NAME ?= authorizer-service
IMAGE_NAME ?= $(REGISTRY)/$(APP_NAME)
VERSION := $(shell git describe --tags --abbrev=0 || echo "0.0.0")
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

build:
	go mod tidy
	go build -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)" -o ./bin/app/$(APP_NAME) ./cmd/api

run: build
	./bin/app/$(APP_NAME)

clean:
	rm -f ./bin/app/$(APP_NAME)
	
docker-build:
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		-t $(IMAGE_NAME):$(VERSION) .

docker-run: docker-build
	docker run --rm --name $(CONTAINER_NAME) $(IMAGE_NAME)

docker-push: docker-build
	docker tag $(IMAGE_NAME):$(VERSION) $(IMAGE_NAME):latest
	docker push $(IMAGE_NAME):$(VERSION)
	docker push $(IMAGE_NAME):latest

.PHONY: build test run docker-build docker-run docker-push