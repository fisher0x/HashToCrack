.PHONY: all windows linux darwin clean

BINARY_NAME=HashToCrack
VERSION=1.0.0
BUILD_PATH=./cmd/HashToCrack

all: windows linux darwin

windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(BINARY_NAME)_windows_amd64.exe $(BUILD_PATH)
	GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o bin/$(BINARY_NAME)_windows_386.exe $(BUILD_PATH)

linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(BINARY_NAME)_linux_amd64 $(BUILD_PATH)
	GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o bin/$(BINARY_NAME)_linux_386 $(BUILD_PATH)
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o bin/$(BINARY_NAME)_linux_arm64 $(BUILD_PATH)
	GOOS=linux GOARCH=arm go build -ldflags="-s -w" -o bin/$(BINARY_NAME)_linux_arm $(BUILD_PATH)

darwin:
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(BINARY_NAME)_darwin_amd64 $(BUILD_PATH)
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bin/$(BINARY_NAME)_darwin_arm64 $(BUILD_PATH)

clean:
	@echo "Cleaning..."
	rm -rf bin/

build:
	@echo "Building for current platform..."
	go build -ldflags="-s -w" -o $(BINARY_NAME) $(BUILD_PATH)

test:
	go test -v ./...

install:
	go install $(BUILD_PATH)
