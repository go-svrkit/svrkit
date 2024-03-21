# Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
# See accompanying files LICENSE.txt

PWD = $(shell pwd)
GOBIN = $(PWD)/bin
GO?=go
PATH := $(GOBIN):$(PATH)

PROTOC_FLAGS = --go_opt=paths=source_relative --go_out=. \
	--go-vtproto_opt=paths=source_relative,features=marshal+unmarshal+size --go-vtproto_out=.

ALL_TEST_PKG=gopkg.in/svrkit.v1/...


.PHONY: clean test all

dep:
	$(GO) mod tidy
	$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	$(GO) install github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto@latest

testdata:
	@cd codec/testdata && protoc $(PROTOC_FLAGS) ./*.proto

vet:
	$(GO) vet ${ALL_TEST_PKG}

test:
	$(GO) test -v ${ALL_TEST_PKG} -cover -cpu=4
	#$(GO) test -v -bench ${ALL_TEST_PKG} -run ^Benchmark$ -benchmem

lint:
	cd src && golangci-lint run --timeout 10m ./... > golint.log

clean:
	$(GO) clean
