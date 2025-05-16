# Makefile for Chromatic Curator

# Variables
APP_NAME := curator
DOCKER_IMAGE := $(APP_NAME):latest
BUILD_DIR := .
MAIN_PATH := ./cmd/...

# Go build flags
GO_BUILD_FLAGS := -v

# Docker build flags
DOCKER_BUILD_FLAGS := --no-cache

.PHONY: clean build docker-build

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@if exist $(BUILD_DIR)\$(APP_NAME) del /F /Q $(BUILD_DIR)\$(APP_NAME)
	@echo "Clean complete"

# Build the Go application
build:
	@echo "Building $(APP_NAME)..."
	go build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)\$(APP_NAME) $(MAIN_PATH)
	@echo "Build complete"

# Build Docker image
docker-build:
	@echo "Building Docker image $(DOCKER_IMAGE)..."
	docker build $(DOCKER_BUILD_FLAGS) -t $(DOCKER_IMAGE) .
	@echo "Docker build complete"

# Default target
all: clean build

# Help target
help:
	@echo "Available targets:"
	@echo "  clean        - Remove build artifacts"
	@echo "  build        - Build the Go application"
	@echo "  docker-build - Build Docker image"
	@echo "  all          - Run clean and build"
	@echo "  help         - Show this help message"
