### THANK U KUBE-VIP

ifneq (,$(wildcard ./.env))
    include .env
endif

SHELL := /bin/sh

TARGET := ano

GO_BUILD_ENV=CGO_ENABLED=0 GOOS=linux GOARCH=amd64
GO_FILES=$(shell go list ./... | grep -v /vendor/)
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

CSV ?= malstat.csv
DB ?= ${DATABASE}
REPOSITORY ?= reidaa
DOCKERFILE ?= build/Dockerfile
DOCKERTAG ?= latest

.PHONY: all build clean install uninstall check run deploy ansible

all: build

$(TARGET):
	$(GO_BUILD_ENV) go build -v -o $@ .

build: $(TARGET)
	@true
.PHONY: build

test:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go test -cover "-coverprofile=cover.out" -v $(TESTFLAGS) $(GO_FILES)
.PHONY: test

clean:
	rm -f $(TARGET)
.PHONY: clean

vet:
	go vet $(GO_FILES)
.PHONY: vet

fmt:
	go fmt $(GO_FILES)
.PHONY: fmt

re: clean build
.PHONY: re

scrap: re
	./$(TARGET) $@
.PHONY: scrap

help: re
	./$(TARGET) $@
.PHONY: help

serve: re
	./$(TARGET) $@
.PHONY: serve

version: re
	./$(TARGET) $@
.PHONY: version

lint: build
	golangci-lint run
.PHONY: lint

docker:
	docker build -t ${REPOSITORY}/${TARGET}:${DOCKERTAG} -f ${DOCKERFILE} .
.PHONY: docker

docker-build-debug:
	docker build --progress=plain --no-cache -t ${REPOSITORY}/${TARGET}:debug -f ${DOCKERFILE} .
.PHONY: docker-build-debug

docker-run: docker
	docker run -p 8080:8080 ${REPOSITORY}/${TARGET}:${DOCKERTAG} version
.PHONY: docker-run
