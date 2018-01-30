FROM golang:1.9.2-alpine

COPY . $GOPATH/src/aliddns

WORKDIR $GOPATH/src/aliddns

RUN cd $GOPATH/src/aliddns && go build && go install

ENTRYPOINT ["aliddns"]
