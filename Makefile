COMMIT := $(shell git describe --dirty --always)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
REPO := github.com/dstreamcloud/metallb-dns

GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
GOFILES := $(wildcard *.go)

LDFLAGS=-ldflags "-X $(REPO)/internal/version.commit=$(COMMIT) -X $(REPO)/internal/version.branch=$(BRANCH)"
MAKEFLAGS += --silent

clean:
	@-rm $(GOBIN)/controller 2> /dev/null
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

build: install compile package

compile: $(GOFILES)
	@echo "  >  Building binary..."
	@go build -v -o $(GOBIN)/controller $(LDFLAGS) $(REPO)/controller

install:
	@echo "  >  Checking if there is any missing dependencies..."
	@go mod download

package: $(GOBIN)/controller
	@echo "  >  Building Docker..."
	@docker build -t dstreamcloud/metallb-dns:$(BRANCH)-$(COMMIT) -f Dockerfile ./

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "github.com/dstreamcloud/metallb-dns":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo