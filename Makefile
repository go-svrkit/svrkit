
PWD = $(shell pwd)
GOBIN = $(PWD)/bin
GO?=go
PATH := $(GOBIN):$(PATH)

ALL_TEST_PKG=gopkg.in/svrkit.v1/...


.PHONY: clean test all

vet:
	$(GO) vet ${ALL_TEST_PKG}

test:
	$(GO) test -v ${ALL_TEST_PKG} -cover -cpu=4
	#$(GO) test -v -bench ${ALL_TEST_PKG} -run ^Benchmark$ -benchmem

lint:
	cd src && golangci-lint run --timeout 10m ./... > ../golint.log

clean:
	$(GO) clean
