MOBILESNIPER_MAIN_GO = ./cmd/mobilesniper/main.go

BUILD_PATH_RELEASE = build/release/mobilesniper
BUILD_PATH_DEBUG = build/debug/mobilesniper

GO_FLAGS_RELEASE = -a -buildmode=exe -ldflags "-s -w" -gcflags=all="-l -B"
GO_FLAGS_DEBUG = -a -race -cover -buildmode=default

GO_BUILD_LINUX = GOOS=linux GOARCH=amd64 go build
GO_BUILD_WINDOWS = GOOS=windows GOOARCH=amd64 go build

mobilesniper: clean
	@mkdir -p build/release build/debug
	$(GO_BUILD_LINUX) $(GO_FLAGS_RELEASE) $(MOBILESNIPER_MAIN_GO) -o $(BUILD_PATH_RELEASE)
	$(GO_BUILD_LINUX) $(GO_FLAGS_DEBUG) $(MOBILESNIPER_MAIN_GO) -o $(BUILD_PATH_DEBUG)
	$(GO_BUILD_WINDOWS) $(GO_FLAGS_RELEASE) $(MOBILESNIPER_MAIN_GO) -o $(BUILD_PATH_RELEASE).exe
	$(GO_BUILD_WINDOWS) $(GO_FLAGS_DEBUG) $(MOBILESNIPER_MAIN_GO) -o $(BUILD_PATH_DEBUG).exe

run:
	go run $(MOBILESNIPER_MAIN_GO)

clean:
	@rm -f $(BUILD_PATH_RELEASE)* >> /dev/null
	@rm -f $(BUILD_PATH_DEBUG)* >> /dev/null