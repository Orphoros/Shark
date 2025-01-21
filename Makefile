BINARY_NAME=shark
BUILD_DIR=build
BUILD_VERSION=0.0.1
INSTRUCTION_FAMILY_CODENAME=Onos1

ifeq ($(OS),Windows_NT) 
    DETECTED_OS := Windows
	 ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
        MACHINE += AMD64
    endif
    ifeq ($(PROCESSOR_ARCHITECTURE),x86)
        MACHINE += IA32
    endif
else
    DETECTED_OS := $(shell sh -c 'uname 2>/dev/null || echo Unknown')
	MACHINE := $(shell sh -c 'uname -m 2>/dev/null || echo Unknown')
endif

ifeq ($(DETECTED_OS),Darwin)
	LIB_EXT=dylib
	LIB_REFIX=lib
endif
ifeq ($(DETECTED_OS),Linux)
	LIB_EXT=so
	LIB_REFIX=lib
endif
ifeq ($(DETECTED_OS),Windows)
	LIB_EXT=dll
	LIB_REFIX=
endif

BUILD_NUMBER=$(shell date +%y)$(shell date +%j)
BIN_FLAGS=-trimpath -gcflags=all="-l=10 -C" -ldflags="-s -w -X 'main.Version=${BUILD_VERSION}' -X 'main.Build=$(BUILD_NUMBER)' -X 'main.Codename=$(INSTRUCTION_FAMILY_CODENAME)'"
LIB_FLAGS=-ldflags="-s -w" -trimpath -gcflags=all="-l -C" -buildmode=c-shared

GOOS = linux darwin windows
GOARCH = amd64 arm64

.PHONY: all build clean test dep check coverage lint serve-coverage help field-align sec-check bench-profile

##@ Commands

all: clean dep test build ## Run all commands

build: ## Build the Shark binaries
	@echo "Compiling Shark SDK..."
	@for os in $(GOOS); do \
		if [ "$$os" = "windows" ]; then \
			EXE_EXT=".exe"; \
		fi; \
		for arch in $(GOARCH); do \
			GOARCH=$$arch GOOS=$$os go build ${BIN_FLAGS} -o ${BUILD_DIR}/bin/$$os/$$arch/${BINARY_NAME}$$EXE_EXT ./cmd/bin/sdk; \
		done; \
	done
	@echo "[DONE]: Shark SDK compiled"

	@echo "Compiling the Shark Compiler..."
	@for os in $(GOOS); do \
		if [ "$$os" = "windows" ]; then \
			EXE_EXT=".exe"; \
		fi; \
		for arch in $(GOARCH); do \
			GOARCH=$$arch GOOS=$$os go build ${BIN_FLAGS} -o ${BUILD_DIR}/bin/$$os/$$arch/${BINARY_NAME}c$$EXE_EXT ./cmd/bin/compiler; \
		done; \
	done
	@echo "[DONE]: Shark Compiler compiled"

	@echo "Compiling the Nidum Virtual Machine..."
	@for os in $(GOOS); do \
		if [ "$$os" = "windows" ]; then \
			EXE_EXT=".exe"; \
		fi; \
		for arch in $(GOARCH); do \
			GOARCH=$$arch GOOS=$$os go build ${BIN_FLAGS} -o ${BUILD_DIR}/bin/$$os/$$arch/nidum$$EXE_EXT ./cmd/bin/vm; \
		done; \
	done
	@echo "[DONE]: Shark Virtual Machine compiled"

	@echo "Compiling the Shark language server..."
	@for os in $(GOOS); do \
		if [ "$$os" = "windows" ]; then \
			EXE_EXT=".exe"; \
		fi; \
		for arch in $(GOARCH); do \
			GOARCH=$$arch GOOS=$$os go build ${BIN_FLAGS} -o ${BUILD_DIR}/bin/$$os/$$arch/${BINARY_NAME}ls$$EXE_EXT ./cmd/bin/lsp; \
		done; \
	done
	@echo "[DONE]: Shark language server compiled"

	@echo "Compiling Shark C bindings for current OS and Platform.."
	@go build ${LIB_FLAGS} -o ${BUILD_DIR}/lib/${DETECTED_OS}/${MACHINE}/${LIB_REFIX}${BINARY_NAME}.${LIB_EXT} ./cmd/lib
	@echo "[DONE]: Lib compiled"

run-lsp: ## Run the language server
	@go run ./cmd/bin/lsp/main.go

clean: ## Remove development artifacts and clean build directory
	@echo "Cleaning build directory..."
	@rm -rf ${BUILD_DIR}
	@echo "[DONE]: Cleaned build directory"
	@echo "Cleaning modcache..."
	@go clean -modcache -i -r
	@echo "[DONE]: Cleaned modcache"

test: ## Run unit tests
	@echo "Running tests..."
	@mkdir -p ${BUILD_DIR}
	@go test -coverprofile=./${BUILD_DIR}/coverage.out -cover -v ./...
	@echo "[DONE]: Tests completed, coverage report generated in coverage.out"
	@echo "Validating race conditions..."
	@go test -race -short ./...
	@echo "[DONE]: Testing completed"

dep: ## Download and tidy dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@echo "[DONE]: Dependencies downloaded"
	@echo "Tidying dependencies..."
	@go mod tidy
	@echo "[DONE]: Dependencies tidied"

check: ## Check OS and ARCH
	@echo "Detected OS:   ${DETECTED_OS}"
	@echo "Detected ARCH: ${MACHINE}"

lint: ## Run linter
	@echo "Running linter..."
	@command -v golangci-lint >/dev/null 2>&1 || { echo >&2 "golangci-lint is required but not installed. Aborting."; exit 1; }
	@golangci-lint run --timeout 5m
	@echo "[DONE]: Linter completed"

coverage-report: test ## Serve coverage report in browser
	@echo "Serving coverage report..."
	@go tool cover -html=./${BUILD_DIR}/coverage.out  
	@go tool cover -html=./${BUILD_DIR}/coverage.out -o ./${BUILD_DIR}/coverage.html

field-align: ## Run field analysis
	@echo "Running field analysis..."
	@fieldalignment -fix ./vm
	@fieldalignment -fix ./types
	@fieldalignment -fix ./token
	@fieldalignment -fix ./serializer
	@fieldalignment -fix ./parser
	@fieldalignment -fix ./object
	@fieldalignment -fix ./lsp
	@fieldalignment -fix ./lexer
	@fieldalignment -fix ./internal
	@fieldalignment -fix ./exception
	@fieldalignment -fix ./emitter
	@fieldalignment -fix ./config
	@fieldalignment -fix ./compiler
	@fieldalignment -fix ./code
	@fieldalignment -fix ./cmd
	@fieldalignment -fix ./bytecode
	@fieldalignment -fix ./ast
	@echo "[DONE]: Field analysis completed"

sec-check: ## Automatic Static Code Security Analysis
	@echo "Static Code Checking for Security Vulnerabilities..."
	@gosec ./...
	@echo "[DONE]: Security Check completed"

bench-profile: ## Run benchmark profiling
	@echo "Running benchmark profiling..."
	@go test -cpuprofile cpu.prof -memprofile mem.prof -trace trace.out -run ^TestRecursiveFibonacci$ -v ./vm
	@echo "[DONE]: Benchmark profiling completed"

help:
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make <command> \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
