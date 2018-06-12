# Ali DDNS

调用阿里云提供的API，刷新域名解析，将其解析到当前IP地址。

本来想用php写的，但是php版本的sdk写的太烂了（大家可以感受一下：[github/aliyun-stupid-php-sdk](https://github.com/aliyun/aliyun-openapi-php-sdk)），没有composer，连注释也没有，不想用，故换成golang。

## 使用前提

IP必须是公网IP，如果被NAT成了内网IP，需要打电话去申请成公网IP。

内网IP是没有办法用这个的。

## 配置说明

解释下 `.env.example` 文件里的几个变量：

- `ACCESS_KEY` 和 `ACCESS_KEY_SECRET` 就是访问阿里云API的密钥，**必填**。在这里可以看到 https://ak-console.aliyun.com

- `REGION` 是区域ID，可以为空，默认是 `cn-hangzhou`

## 部署

`docker run -p 8888:8888 -e ACCESS_KEY="key" -e ACCESS_KEY_SECRET="secret" -e REGION="cn-hangzhou" darthminion/aliddns

- docker-compose

```yml
version: '3'

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

## 调用

然后发起请求 `GET host:port/?domain_name=example.com&rr=@`

Query String 说明：

- `--domain-name`: 要解析的域名，**必填**。

- `--rr`: 解析的RR值，选填，默认是`@`。
