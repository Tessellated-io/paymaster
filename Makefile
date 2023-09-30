# Build / install

BINARY_NAME := paymaster
GOPATH ?= $(shell go env GOPATH)
INSTALL_PATH := $(GOPATH)/bin/$(BINARY_NAME)

.PHONY: install

install:
	go build -o $(BINARY_NAME) cmd/cmd.go
	mv $(BINARY_NAME) $(INSTALL_PATH)

# Protocol Buffers

check-proto-deps:
ifeq (,$(shell which protoc-gen-gogofaster))
	@go install github.com/gogo/protobuf/protoc-gen-gogofaster@latest
endif
.PHONY: check-proto-deps

check-proto-format-deps:
ifeq (,$(shell which clang-format))
	$(error "clang-format is required for Protobuf formatting. See instructions for your platform on how to install it.")
endif
.PHONY: check-proto-format-deps

proto-gen: check-proto-deps
	@echo "Generating Protobuf files"
	@buf generate
	@echo "Complete"
.PHONY: proto-gen
