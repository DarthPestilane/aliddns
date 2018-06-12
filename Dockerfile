FROM golang:alpine

ADD . $GOPATH/src/aliddns

WORKDIR $GOPATH/src/aliddns

ENV PORT="8888"

EXPOSE ${PORT}

RUN apk add --no-cache git && cd $GOPATH/src/aliddns && go get -v -u -d ./... && go build -ldflags '-s -w' && go install

CMD aliddns run --port=${PORT}
