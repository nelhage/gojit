## simple makefile to log workflow
.PHONY: all test clean build

#GOFLAGS := $(GOFLAGS:-race -v)

all: build test
	@# done

build: clean
	@go get $(GOFLAGS) ./...

test: build
	@go test $(GOFLAGS) ./...

clean:
	@go clean $(GOFLAGS) -i ./...

## EOF
