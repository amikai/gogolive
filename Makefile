EXECUTABLE := gogolive

export CC=$(shell brew --prefix llvm)/bin/clang
export CXX=$(shell brew --prefix llvm)/bin/clang++
GO=go

.PHONY: test
test: 
	@$(GO) test -v -cover -tags $(TAGS) -coverprofile=coverage.txt ./... && echo "\n==>\033[32m Ok\033[m\n" || exit 1

.PHONY: build
build: $(EXECUTABLE)

.PHONY: $(EXECUTABLE)
$(EXECUTABLE): $(GOFILES)
	$(GO) mod tidy
	$(GO) build -o $@

.PHONY: run
run: build
	./$(EXECUTABLE)

.PHONY: clean
clean:
	@$(GO) clean -modcache -x -i ./...
	@$(RM) coverage.txt
