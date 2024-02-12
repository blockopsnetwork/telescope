
.DEFAULT_GOAL := all
PROJECT_NAME := "telescope"

.PHONY: all clean test run build build-docker run-docker test-docker clean-docker help

# Build the project
all: clean build test

# Remove all build artifacts
clean:
	rm telescope

# Run the tests
test: build
	go test -v ./...

# Build the project
build:
	go build -o telescope .