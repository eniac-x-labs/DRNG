APP_NAME := DRNG-service
CMD_DIR := .
CMD_FILE := main.go
BIN_DIR := bin
BIN_FILE := ./$(BIN_DIR)/$(APP_NAME)
PORT := 8080

SHELL := /bin/sh

GITCOMMIT := $(shell git rev-parse HEAD)
GITDATE := $(shell git show -s --format='%ct')

LDFLAGSSTRING +=-X main.GitCommit=$(GITCOMMIT)
LDFLAGSSTRING +=-X main.GitDate=$(GITDATE)
LDFLAGSSTRING +=-X main.GitVersion=$(GITVERSION)
LDFLAGS :=-ldflags "$(LDFLAGSSTRING)"

# 默认目标
#.DEFAULT_GOAL := run

build:
	@echo "Building $(APP_NAME)..."
	env GO111MODULE=on go build $(LDFLAGS) -o $(BIN_FILE) $(CMD_DIR)/$(CMD_FILE)
	@echo "Build complete."
.PHONY: build

run:
	@echo "Running $(APP_NAME)..."
	@go run $(CMD_DIR)/$(CMD_FILE)
.PHONY: run

run-bin: build
	@echo "Running $(APP_NAME)..."
	@$(BIN_FILE)
.PHONY: run-bin
