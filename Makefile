.PHONY: build serve
VERSION := $(shell git describe --always |sed -e "s/^v//")

build:
	mkdir -p build
	go build $(GO_EXTRA_BUILD_ARGS) -ldflags "-s -w -X main.version=$(VERSION)" -o build/ccfiw cmd/ccfiw/main.go


# shortcuts for development

dev-requirements:
	go mod download

serve: build
	@echo "Starting CCFIW"
	./build/ccfiw
