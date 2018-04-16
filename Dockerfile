FROM golang:1.10.1-alpine

COPY . $GOPATH/src/aliddns

WORKDIR $GOPATH/src/aliddns

RUN cd $GOPATH/src/aliddns && go build -ldflags '-s -w' && go install

ENTRYPOINT ["aliddns"]
