.PHONY: all test clean vendor

ldflags = -ldflags="-s -w"
gcflags = -gcflags="-trimpath=${PWD}"
output = -o=aliddns

build:
	CGO_ENABLED=0 go build ${ldflags} ${gcflags} -v ${output}

lint:
	CGO_ENABLED=0 golangci-lint run --concurrency=2

tidy:
	go mod tidy
