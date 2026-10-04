.PHONY: all build-local clean test deps version

BINARY_NAME := penit
DIST_DIR := dist

ifeq ($(origin VERSION), environment)
else
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
endif

BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

LDFLAGS := -ldflags="-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

all: build-local

build-local:
	go build $(LDFLAGS) -o $(BINARY_NAME) .

install: build-local
	@mkdir -p $(HOME)/.local/bin
	mv $(BINARY_NAME) $(HOME)/.local/bin/

test:
	go test ./...

clean:
	rm -rf $(DIST_DIR)
	rm -f $(BINARY_NAME)

deps:
	go mod tidy && go mod download

version:
	@echo "Version: $(VERSION) | Build: $(BUILD_TIME) | Commit: $(GIT_COMMIT)"
