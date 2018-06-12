FROM golang:alpine

ADD . $GOPATH/src/aliddns

WORKDIR $GOPATH/src/aliddns

ENV PORT="8888"

EXPOSE ${PORT}

# RUN echo http://mirrors.aliyun.com/alpine/v3.7/main > /etc/apk/repositories && \
#     echo http://mirrors.aliyun.com/alpine/v3.7/community >> /etc/apk/repositories

RUN apk add --no-cache --virtual .build-deps git && \
    cd ${GOPATH}/src/aliddns && \
    go get -v -u -d ./... && \
    go build -ldflags '-s -w' && \
    go install && \
    apk del .build-deps

CMD aliddns run --port=${PORT}
