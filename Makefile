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

all: clean dep test compile

compile:
	@echo "Compiling Shark SDK..."

	@echo "Building Shark SDK for current OS and Platform..."
	@GOARCH=amd64 GOOS=darwin go build ${BIN_FLAGS} -o ${BUILD_DIR}/bin/Darwin/amd64/${BINARY_NAME} ./cmd/bin/sdk
	@echo "[DONE]: Darwin AMD64"
	@GOARCH=amd64 GOOS=linux go build ${BIN_FLAGS} -o ${BUILD_DIR}/bin/Linux/amd64/${BINARY_NAME} ./cmd/bin/sdk
	@echo "[DONE]: Linux AMD64"
	@GOARCH=amd64 GOOS=windows go build ${BIN_FLAGS} -o ${BUILD_DIR}/bin/Windows/amd64/${BINARY_NAME}.exe ./cmd/bin/sdk
	@echo "[DONE]: Windows AMD64"
	@GOARCH=arm64 GOOS=darwin go build ${BIN_FLAGS} -o ${BUILD_DIR}/bin/Darwin/arm64/${BINARY_NAME} ./cmd/bin/sdk
	@echo "[DONE]: Darwin ARM64"

	@echo "Compiling the Shark Compiler..."
	@GOARCH=amd64 GOOS=darwin go build ${BIN_FLAGS} -o ${BUILD_DIR}/bin/Darwin/amd64/${BINARY_NAME}c ./cmd/bin/compiler
	@echo "[DONE]: Darwin AMD64"
	@GOARCH=amd64 GOOS=linux go build ${BIN_FLAGS} -o ${BUILD_DIR}/bin/Linux/amd64/${BINARY_NAME}c ./cmd/bin/compiler
	@echo "[DONE]: Linux AMD64"
	@GOARCH=amd64 GOOS=windows go build ${BIN_FLAGS} -o ${BUILD_DIR}/bin/Windows/amd64/${BINARY_NAME}c.exe ./cmd/bin/compiler
	@echo "[DONE]: Windows AMD64"
	@GOARCH=arm64 GOOS=darwin go build ${BIN_FLAGS} -o ${BUILD_DIR}/bin/Darwin/arm64/${BINARY_NAME}c ./cmd/bin/compiler
	@echo "[DONE]: Darwin ARM64"

	@echo "Compiling the Shark Virtual Machine..."
	@GOARCH=amd64 GOOS=darwin go build ${BIN_FLAGS} -o ${BUILD_DIR}/bin/Darwin/amd64/orpvm ./cmd/bin/vm
	@echo "[DONE]: Darwin AMD64"
	@GOARCH=amd64 GOOS=linux go build ${BIN_FLAGS} -o ${BUILD_DIR}/bin/Linux/amd64/orpvm ./cmd/bin/vm
	@echo "[DONE]: Linux AMD64"
	@GOARCH=amd64 GOOS=windows go build ${BIN_FLAGS} -o ${BUILD_DIR}/bin/Windows/amd64/orpvm.exe ./cmd/bin/vm
	@echo "[DONE]: Windows AMD64"
	@GOARCH=arm64 GOOS=darwin go build ${BIN_FLAGS} -o ${BUILD_DIR}/bin/Darwin/arm64/orpvm ./cmd/bin/vm
	@echo "[DONE]: Darwin ARM64"

	@echo "Compiling Shark C bindings for current OS and Platform.."
	@go build ${LIB_FLAGS} -o ${BUILD_DIR}/lib/${DETECTED_OS}/${MACHINE}/${LIB_REFIX}${BINARY_NAME}.${LIB_EXT} ./cmd/lib
	@echo "[DONE]: Lib compiled"

run:
	@go run ./cmd/bin/main.go

clean:
	@echo "Cleaning build directory..."
	@rm -rf ${BUILD_DIR}
	@echo "[DONE]: Cleaned build directory"
	@echo "Cleaning modcache..."
	@go clean -modcache -i -r
	@echo "[DONE]: Cleaned modcache"

test:
	@echo "Running tests..."
	@go test -coverprofile=coverage.out -cover -v ./...
	@echo "[DONE]: Tests completed, coverage report generated in coverage.out"

dep:
	@echo "Downloading dependencies..."
	@go mod download
	@echo "[DONE]: Dependencies downloaded"
	@echo "Tidying dependencies..."
	@go mod tidy
	@echo "[DONE]: Dependencies tidied"

check:
	@echo "Detected OS:   ${DETECTED_OS}"
	@echo "Detected ARCH: ${MACHINE}"

coverage:
	@go tool cover -html=coverage.out
