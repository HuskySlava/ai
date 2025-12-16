# Makefile

# Binary name
BINARY_NAME=ai

# Config name
CONFIG_NAME=config.yaml

# Output directory
BIN_DIR=./bin

# Go build flags
GOOS=darwin
GOARCH=arm64

# Default target
all: build move

# Build the Go binary for macOS
build:
	@echo "Building $(BINARY_NAME) for macOS..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BINARY_NAME) ./cmd/ai


# Build for MacOS
build-macos:
	@echo "Building $(BINARY_NAME) for MacOS/arm64..."
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME) ./cmd/

# Build for Linux
build-linux:
	@echo "Building $(BINARY_NAME) for linux/amd64..."
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME) ./cmd/

# Build for Windows
build-windows:
	@echo "Building $(BINARY_NAME).exe for windows/amd64..."
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME).exe ./cmd/

# Move binary to target folder
move:
	@echo "Moving binary to $(BIN_DIR)..."
	mkdir -p $(BIN_DIR)
	mv -f $(BINARY_NAME) $(BIN_DIR)/
	@echo "Moving config to $(BIN_DIR)..."
	cp -f $(CONFIG_NAME) $(BIN_DIR)/

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BIN_DIR)/$(BINARY_NAME)
