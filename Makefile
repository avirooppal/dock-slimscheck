.PHONY: build install test clean release

# Set variables
BINARY_NAME=dock-slimcheck
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_FLAGS=-ldflags "-X main.Version=$(VERSION)"
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

# Build the binary
build:
	go build $(BUILD_FLAGS) -o $(BINARY_NAME) .

# Install the binary
install: build
	mv $(BINARY_NAME) /usr/local/bin/

# Run tests
test:
	go test -v ./...

# Run the binary
run: build
	./$(BINARY_NAME)

# Clean up
clean:
	rm -f $(BINARY_NAME)
	rm -rf dist/

# Build for all platforms
release:
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o dist/$(BINARY_NAME)_linux_amd64 .
	GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o dist/$(BINARY_NAME)_linux_arm64 .
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o dist/$(BINARY_NAME)_darwin_amd64 .
	GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o dist/$(BINARY_NAME)_darwin_arm64 .
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o dist/$(BINARY_NAME)_windows_amd64.exe .