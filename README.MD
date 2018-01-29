# Ali DDNS

调用阿里云提供的API，刷新域名解析，将其解析到当前IP地址。

本来想用php写的，但是php版本的sdk写的太烂了（大家可以感受一下：[github/aliyun-stupid-php-sdk](https://github.com/aliyun/aliyun-openapi-php-sdk)），没有composer，连注释也没有，不想用，故换成golang。

## 使用前提

IP必须是公网IP，如果被NAT成了内网IP，需要打电话去申请成公网IP。

内网IP是没有办法用这个的。

## 原理

![原理流程](https://user-images.githubusercontent.com/15375753/35424004-e7fd1948-028b-11e8-946d-03ec504ed3ff.png)

## 配置说明

解释下 `.env.example` 文件里的几个变量：

- `ACCESS_KEY` 和 `ACCESS_KEY_SECRET` 就是访问阿里云API的密钥，必填。在这里可有看到 https://ak-console.aliyun.com

- `REGION` 是区域ID，可以为空，默认是 `cn-hangzhou`

- `DOMAIN_NAME` 表示要解析的域名。必填，如 `darthpestilane.com`

- `RR` 表示主机记录。选填，默认为 `@`

- `AUTO_REFRESH` 是否开启自动刷新模式。选填，填 `true` 或 `false`

- `INTERVAL_MINUTES` 如果开启了自动刷新，则为必填，表示刷新的间隔时间（分钟），务必是大于0的整数

## 使用方式

首先复制 .env.example 文件到 .env，然后按照配置说明修改里面的变量。

如果有 `docker` ，先build镜像（如果还没有build过的话）:

```sh
docker build -t xxx/ddns:t123 .
```

然后run容器:

```sh
docker run --name ali-ddns xxx/ddns:t123
```

或者直接用仓库里的 `docker-compose.yml` 文件来创建容器:

```sh
docker-compose up
```

之后就可以定时去start容器来刷新解析:

```sh
docker start -i ali-ddns # 如果加了 -i 参数，则会输出日志
```

---

如果么有 `docker`，直接 `go run main.go` 即可。
