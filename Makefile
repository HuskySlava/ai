# Makefile

BINARY_NAME=ai
CONFIG_NAME=config.yaml
BIN_DIR=./bin

# Default OS/ARCH
GOOS ?= darwin
GOARCH ?= arm64

.PHONY: all build build-macos build-linux build-windows move clean

all: build move

build:
	@echo "Building $(BINARY_NAME) for $(GOOS)/$(GOARCH)..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BINARY_NAME) ./cmd/ai

build-macos:
	$(MAKE) build GOOS=darwin GOARCH=arm64

build-linux:
	$(MAKE) build GOOS=linux GOARCH=amd64

build-windows:
	$(MAKE) build GOOS=windows GOARCH=amd64
	mv -f $(BINARY_NAME) $(BINARY_NAME).exe

move:
	@echo "Moving binary to $(BIN_DIR)..."
	mkdir -p $(BIN_DIR)
	if [ -f $(BINARY_NAME) ]; then mv -f $(BINARY_NAME) $(BIN_DIR)/; fi
	if [ -f $(BINARY_NAME).exe ]; then mv -f $(BINARY_NAME).exe $(BIN_DIR)/; fi
	@echo "Moving config to $(BIN_DIR)..."
	cp -f $(CONFIG_NAME) $(BIN_DIR)/

clean:
	@echo "Cleaning..."
	rm -rf $(BIN_DIR)/*