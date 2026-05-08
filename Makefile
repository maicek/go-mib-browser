
BINARY_NAME := dist/mib-browser
SRC_DIR     := ./cmd/mib-browser

BUILD_FLAGS := -a -trimpath -ldflags="-s -w -X main.embedPath=dist" -gcflags="all=-l=4" -tags=netgo



build:
	go build $(BUILD_FLAGS) -o $(BINARY_NAME) $(SRC_DIR)

build-windows-x86:
	CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-windows-amd64.exe $(SRC_DIR)

#build-windows-arm64:
#	CGO_ENABLED=1 CC=aarch64-w64-mingw32-gcc CXX=aarch64-w64-mingw32-g++ GOOS=windows GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-windows-arm64.exe $(SRC_DIR)

build-linux:
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-linux-amd64 $(SRC_DIR)


clean:
	rm -rf dist

build-all: clean build-windows-x86 build-linux