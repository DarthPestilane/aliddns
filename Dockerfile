FROM golang:alpine as builder

ADD . /aliddns-src

# RUN echo http://mirrors.aliyun.com/alpine/v3.7/main > /etc/apk/repositories && \
#     echo http://mirrors.aliyun.com/alpine/v3.7/community >> /etc/apk/repositories

RUN apk add --no-cache --virtual .build-deps git && \
    cd aliddns-src && \
    go build -gcflags=-trimpath=${GOPATH} -asmflags=-trimpath=${GOPATH} -ldflags '-s -w' -o aliddns && \
    cp aliddns /

FROM alpine:latest

COPY --from=builder /aliddns /
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip

ENV PORT=8888

EXPOSE ${PORT}

CMD ["/aliddns", "run"]
