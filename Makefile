## simple makefile to log workflow
.PHONY: all test clean build bench

#GOFLAGS := $(GOFLAGS:-race -v)

all: build test
	@# done

build: clean
	@go get $(GOFLAGS) ./...

test: build
	@go test $(GOFLAGS) ./...

bench: test
	@go test $(GOFLAGS) -bench=. ./...

clean:
	@go clean $(GOFLAGS) -i ./...

## EOF
