FROM golang:alpine

ADD . $GOPATH/src/aliddns

WORKDIR $GOPATH/src/aliddns

ENV PORT="8888"

EXPOSE ${PORT}

RUN cd $GOPATH/src/aliddns && go build -ldflags '-s -w' && go install

CMD aliddns run --port=${PORT}
