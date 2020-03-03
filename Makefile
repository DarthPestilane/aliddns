.PHONY: all test clean vendor

ldflags = -ldflags="-s -w"
gcflags = -gcflags="-trimpath=${PWD}"
output = -o=aliddns

build: # build
	CGO_ENABLED=0 go build ${ldflags} ${gcflags} -v ${output}

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${ldflags} ${gcflags} -v ${output}

lint:
	CGO_ENABLED=0 GOGC=15 golangci-lint run

lint-ci:
	CGO_ENABLED=0 GOGC=15 golangci-lint run -v

lint-fix:
	CGO_ENABLED=0 GOGC=15 golangci-lint run --fix

tidy:
	go mod tidy
