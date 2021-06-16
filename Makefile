.PHONY: all test clean vendor

buildTime = `date +%Y-%m-%dT%T%z`
gitCommit = `git rev-parse --short HEAD`
gitTag = `git --no-pager tag --points-at HEAD`

ldflags = -ldflags="-s -w -X main.buildTime=${buildTime} -X main.gitCommit=${gitCommit} -X main.gitTag=${gitTag}"
gcflags = -gcflags="-trimpath=${PWD}"
output = -o=aliddns

build:
	CGO_ENABLED=0 go build ${ldflags} ${gcflags} -v ${output}

lint:
	CGO_ENABLED=0 golangci-lint run --concurrency=2

tidy:
	go mod tidy

image = "darthminion/aliddns"

build-image:
	docker build -t ${image} .

push-image:
	docker push ${image}
