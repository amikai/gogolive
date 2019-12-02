EXECUTABLE := gogolive

CC := gcc-9
GO ?= go
PACKAGES ?= $(shell $(GO) list ./...)
GOFILES := $(shell find . -name "*.go" -type f)

.PHONY: test
test: 
	@$(GO) test -v -cover -coverprofile coverage.txt $(PACKAGES) && echo "\n==>\033[32m Ok\033[m\n" || exit 1
	@$(GO) clean -testcache ./...

.PHONY: build
build: $(EXECUTABLE)

.PHONY: $(EXECUTABLE)
$(EXECUTABLE): $(GOFILES)
	$(GO) mod tidy
	$(GO) build -o $@

.PHONY: clean
clean:
	@$(GO) clean -modcache -x -i ./...
	@$(RM) coverage.txt