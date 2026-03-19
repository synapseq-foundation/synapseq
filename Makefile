# Makefile for building SynapSeq

# Binary information
BIN_NAME 	    := synapseq
BIN_DIR 	    := bin

# Go build metadata
VERSION 	    := $(shell cat VERSION)
COMMIT  	    := $(shell git rev-parse --short HEAD 2>/dev/null || echo $(shell echo ${GITHUB_SHA} | cut -c1-7))
DATE    	    := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

# Windows configuration
MAJOR_VERSION 			 := $(shell echo $(VERSION) | cut -d. -f1)
MINOR_VERSION 			 := $(shell echo $(VERSION) | cut -d. -f2)
PATCH_VERSION 			 := $(shell echo $(VERSION) | cut -d. -f3)
GO_VERSION_INFO_CMD 	 := github.com/josephspurrier/goversioninfo/cmd/goversioninfo@v1.5.0
GO_VERSION_INFO_CMD_ARGS := -company="SynapSeq Foundation <synapseq.org>" \
							-description="Text-Driven Audio Sequencer for Brainwave Entrainment" \
					  		-copyright="GPL v2" \
					  		-original-name="$(BIN_NAME).exe" \
							-product-name="SynapSeq" \
							-product-version="$(VERSION).0" \
					  		-comment="Main SynapSeq executable" \
							-icon="assets/synapseq.ico" \
							-ver-major=$(MAJOR_VERSION) -product-ver-major=$(MAJOR_VERSION) \
							-ver-minor=$(MINOR_VERSION) -product-ver-minor=$(MINOR_VERSION) \
							-ver-patch=$(PATCH_VERSION) -product-ver-patch=$(PATCH_VERSION) \
							-ver-build=0 -product-ver-build=0
# Go configuration
GO_METADATA     := -X github.com/synapseq-foundation/synapseq/v4/internal/info.VERSION=$(VERSION) \
				  -X github.com/synapseq-foundation/synapseq/v4/internal/info.BUILD_DATE=$(DATE) \
				  -X github.com/synapseq-foundation/synapseq/v4/internal/info.GIT_COMMIT=$(COMMIT)
GO_BUILD_FLAGS  := -ldflags="-s -w $(GO_METADATA)"
MAIN 		    := ./cmd/synapseq

.PHONY: all build clean test build-windows-amd64 build-windows-arm64 \
		build-linux-amd64 build-linux-arm64 \
		build-macos install

all: build

prepare:
	mkdir -p $(BIN_DIR)

# Windows resource file generation
windows-res-amd64:
	go run $(GO_VERSION_INFO_CMD) $(GO_VERSION_INFO_CMD_ARGS) -64 -o cmd/synapseq/synapseq.syso

windows-res-arm64:
	go run $(GO_VERSION_INFO_CMD) $(GO_VERSION_INFO_CMD_ARGS) -arm -o cmd/synapseq/synapseq.syso

test:
	go test -v ./...

build: prepare
	go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME) $(MAIN)

# Windows builds
build-windows-amd64: prepare windows-res-amd64
	GOOS=windows GOARCH=amd64 go build -tags=windows $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-windows-amd64.exe $(MAIN)

build-windows-arm64: prepare windows-res-arm64
	GOOS=windows GOARCH=arm64 go build -tags=windows $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-windows-arm64.exe $(MAIN)

# Linux builds
build-linux-amd64: prepare
	GOOS=linux GOARCH=amd64 go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-linux-amd64 $(MAIN)

build-linux-arm64: prepare
	GOOS=linux GOARCH=arm64 go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-linux-arm64 $(MAIN)

# macOS builds
build-macos: prepare
	GOOS=darwin GOARCH=arm64 go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-macos-arm64 $(MAIN)

# WASM build
build-wasm:
	GOOS=js GOARCH=wasm go build -tags=wasm $(GO_BUILD_FLAGS) -o wasm/$(BIN_NAME).wasm ./cmd/wasm
	cp $(shell go env GOROOT)/lib/wasm/wasm_exec.js wasm/wasm_exec.js

# POSIX installation
install:
	cp $(BIN_DIR)/$(BIN_NAME) /usr/local/bin/$(BIN_NAME)

# Clean build artifacts
clean:
	rm -rf $(BIN_DIR)
	rm -rf cmd/synapseq/*.syso
	rm -rf wasm/*.wasm wasm/wasm_exec.js