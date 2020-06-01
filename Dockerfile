FROM golang:alpine as builder

ADD . /aliddns-src

ENV GOPROXY='https://goproxy.cn/,direct'
ENV GOSUMDB=off

RUN cp /etc/apk/repositories /etc/apk/repositories.backup && \
    sed -i -E "s|http://.+/alpine|http://mirrors\.aliyun\.com/alpine|" /etc/apk/repositories && \
    apk add --no-cache git make && \
    cd /aliddns-src && \
    make build && \
    cp aliddns /

FROM alpine:latest

COPY --from=builder /aliddns /
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip

ENV PORT=8888

EXPOSE ${PORT}

CMD ["/aliddns", "run"]
