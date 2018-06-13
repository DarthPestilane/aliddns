FROM golang:alpine as builder

ADD . $GOPATH/src/aliddns

# RUN echo http://mirrors.aliyun.com/alpine/v3.7/main > /etc/apk/repositories && \
#     echo http://mirrors.aliyun.com/alpine/v3.7/community >> /etc/apk/repositories

RUN apk add --no-cache --virtual .build-deps git && \
    cd ${GOPATH}/src/aliddns && \
    go get -v -u -d ./... && \
    go build -gcflags=-trimpath=${GOPATH} -asmflags=-trimpath=${GOPATH} -ldflags '-s -w' -o aliddns.bin && \
    cp aliddns.bin /

FROM alpine:latest

COPY --from=builder /aliddns.bin /
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip

ENV PORT=8888

EXPOSE ${PORT}

CMD ["/aliddns.bin", "run"]
