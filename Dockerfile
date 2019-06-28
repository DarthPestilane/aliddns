FROM golang:alpine as builder

ADD . /aliddns-src

RUN cp /etc/apk/repositories /etc/apk/repositories.backup && \
    sed -i -E "s|http://.+/alpine|http://mirrors\.aliyun\.com/alpine|" /etc/apk/repositories && \
    apk add --no-cache --virtual .build-deps git && \
    cd /aliddns-src && \
    go build -gcflags=-trimpath=${GOPATH} -asmflags=-trimpath=${GOPATH} -ldflags '-s -w' -o aliddns && \
    cp aliddns /

FROM alpine:latest

COPY --from=builder /aliddns /
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip

ENV PORT=8888

EXPOSE ${PORT}

CMD ["/aliddns", "run"]
