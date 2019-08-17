ROOT_PATH = $(abspath . )
VERSION = $(shell git describe --tags)
PKG_NAME = $(shell basename `git rev-parse --show-toplevel`)
BUILD_DIR = $(ROOT_PATH)/.build

.PHONY: deps
deps:
	go mod tidy
	go get ./...

.PHONY: build
build: deps
	mkdir -p $(BUILD_DIR)
	GO111MODULE=on go build -v -o $(BUILD_DIR)/$(PKG_NAME) -ldflags "-X main.appVer=$(VERSION)" ./$(PKG_NAME)/

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)