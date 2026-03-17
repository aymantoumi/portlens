BINARY=portlens
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

.PHONY: build test lint clean install

build:
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(BINARY) .

build-all:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-linux-amd64 .
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY)-linux-arm64 .

test:
	go test ./... -v -race

lint:
	golangci-lint run ./...

clean:
	rm -f $(BINARY)
	rm -rf dist/

install:
	CGO_ENABLED=0 go install $(LDFLAGS) .
