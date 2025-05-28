default: build

VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X github.com/NodyHub/reggidump/cmd.buildVersion=$(VERSION) -X github.com/NodyHub/reggidump/cmd.buildCommit=$(COMMIT) -X github.com/NodyHub/reggidump/cmd.buildDate=$(DATE)"

build:
	@go build $(LDFLAGS) -o reggidump .

install: build
	@mv reggidump $(GOPATH)/bin/reggidump

clean:
	@go clean
	@rm reggidump

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

test_coverage_view:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

test_coverage_html:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o=coverage.html

all: build install