# Ali DDNS

![Aliddns](https://github.com/DarthPestilane/aliddns/workflows/Aliddns/badge.svg)

调用阿里云提供的 API，刷新域名解析，将其解析到当前 IP 地址。

本来想用 php 写的，但是 php 版本的 sdk 写的太烂了(大家可以感受一下：[github/aliyun-stupid-php-sdk](https://github.com/aliyun/aliyun-openapi-php-sdk))，没有 composer，连注释也没有，不想用，故换成 golang。

## 使用前提

IP 必须是公网 IP，如果被 NAT 成了内网 IP，需要打电话去申请成公网 IP。内网 IP 是没有办法用这个的。

## 使用方式

### 1. 以服务 (server) 的方式运行

#### 使用前提

此方式需要一台外网可访问的服务器(vps)。

#### 配置说明

解释下 `.env.example` 文件里的几个变量：

- `ACCESS_KEY` 和 `ACCESS_KEY_SECRET` 就是访问阿里云 API 的密钥，**必填**。在这里可以看到 https://ak-console.aliyun.com
- `REGION` 是区域 ID，可以为空，默认是 `cn-hangzhou`
- `PORT` 表示服务监听的端口，默认是`8888`

#### 部署

镜像: [darthminion/aliddns](https://hub.docker.com/r/darthminion/aliddns)

- docker run

`docker run -p 8888:8888 -e ACCESS_KEY="key" -e ACCESS_KEY_SECRET="secret" -e REGION="cn-hangzhou" darthminion/aliddns`

- docker-compose

```yml
version: "3"

services:
  aliddns-server:
    image: darthminion/aliddns
    ports:
      - 8888:8888
    environment:
      - ACCESS_KEY=key
      - ACCESS_KEY_SECRET=secret
      - REGION=cn-hangzhou
```

#### 请求调用

请求方式: `GET host:port/?domain_name=example.com&rr=@`

Example: `curl 'host:port/?domain_name=example.com&rr=@'`

Query String 说明：

- `domain-name`: 要解析的域名，**必填**。_注意_: 这里的域名是你的主域名，一般来说是 abc.com 这种形式， **而不是** 完整的 blog.abc.com。
- `rr`: 解析的 RR 值，选填，默认是`@`。也就是子域名部分，比如 'blog'，则最终解析的完整域名是 blog.abc.com。

### 2. 以命令行 (cli) 方式运行

#### 运行

```sh
$ aliddns sync 域名 --ip='current ip' --rr='@' --access-key='xxx' --access-secret='xxx' --region='xx'
```

说明:

- Arguments:

  - `域名`: **必填**。_注意_: 这里的域名是你的主域名，一般来说是 abc.com 这种形式， **而不是** 完整的 blog.abc.com

- Flags:
  - `ip`: 是要解析到域名的 IP 地址。为空则会从 `https://api.ip.sb/geoip/` 获取当前设备的公网 IP，会有一定的耗时。
  - `rr`: 解析的 RR 值，选填，默认是`@`。也就是子域名部分，比如 'blog'，则最终解析的完整域名是 blog.abc.com。
  - `access-key` `access-secret`: 访问阿里云 API 的密钥，**必填**。在这里可以看到 https://ak-console.aliyun.com。默认是读取 `ACCESS_KEY` 和 `ACCESS_KEY_SECRET` 两个环境变量
  - `region`: 是区域 ID。默认是读取 `REGION` 环境变量，如果为空，则默认为 `cn-hangzhou`。

可以通过 `aliddns sync --help` 来查看使用方法。

Example:

```sh
# setup env...

$ aliddns 'abc.com' --rr='blog' # 将域名 blog.abc.com 解析到本地的公网 IP
```
