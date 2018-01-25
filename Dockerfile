FROM golang:1.9.2-alpine

COPY . $GOPATH/src/ddns

WORKDIR $GOPATH/src/ddns

ENTRYPOINT [ "go", "run", "main.go" ]
