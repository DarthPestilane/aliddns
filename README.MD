# Ali DDNS

调用阿里云提供的API，刷新域名解析，将其解析到当前IP地址。

本来想用php写的，但是php版本的sdk写的太烂了（大家可以感受一下:https://github.com/aliyun/aliyun-openapi-php-sdk）

没有composer，注释也没有，不想用，故换成golang。

## 使用之前

解释下 `.env.example` 文件里的几个变量：

- `ACCESS_KEY` 和 `ACCESS_KEY_SECRET` 就是访问阿里云API的密钥，必填。在这里可有看到 https://ak-console.aliyun.com

- `REGION` 是区域ID，可以为空，默认是 `cn-hangzhou`

- `DOMAIN_NAME` 表示要解析的域名。必填，如 `darthpestilane.com`

- `RR` 表示主机记录。选填，默认为 `@`

## 使用方式

首先复制 .env.example 文件到 .env，然后修改里面的变量。

如果有 `docker` ，先build镜像（如果还没有build过的话）:

```sh
docker build -t xxx/ddns:t123 .
```

然后run容器:

```sh
docker run --name aliddns xxx/ddns:t123
```

之后就可以定时去start容器来刷新解析

---

如果么有 `docker`，直接 `go run main.go` 即可。
