# v2ray-endpoint

## 介绍

V2ray节点管理面板，可以对v2ray的服务器节点进行管理

该工具能够方便的对v2ray服务器节点进行管理以及配置生成，免去配置存到云端或者在多个客户端之间重复配置的烦恼

该工具能够生成通用的vmess协议连接，可以在客户端中直接导入

### 安装教程

克隆项目到本地，执行./pack.sh进行打包，将打包后的`target.gz`上传到服务器

在服务器创建`/data/v2ray-subscription`目录

### 使用说明

`target.gz`文件目录说明

```
/ 根目录
-- /v2ray-subscription-server      服务器可执行二进制文件
-- /v2ray-subscription-migrate     数据库迁移工具
-- /config                         服务器配置文件目录
-- /migrations                     数据库迁移文件
-- /web                            Web服务的根目录，内含index.html以及所有的静态资源文件
```

服务器本地化配置支持在`config`目录下创建`db.local.config.yaml`与`main.local.config.yaml`配置文件，配置项与`db.prod.config.yaml`
、`main.prod.config.yaml`文件相同，同名配置可进行覆盖

## 附录

### Nginx配置示例

```
server {
  listen 80;
  listen [::]:80;
  server_name your.host;

  return 301 https://your.host$request_uri;
}

server {
  listen 443 ssl;
  listen [::]:443 ssl;
  server_name your.host;

  ssl_certificate /path/to/tls.pem;
  ssl_certificate_key /path/to/tls.key;

  ssl_protocols TLSv1 TLSv1.1 TLSv1.2;

  ssl_ciphers  HIGH:!aNULL:!MD5;
  ssl_prefer_server_ciphers  on;

  location /api {
    proxy_redirect off;
    proxy_pass http://127.0.0.1:8800;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
    proxy_set_header Host $http_host;

    # Show realip in v2ray access.log
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
  }

  location / {
    root /path/to/your application root/web;
    try_files $uri $uri/ /index.html;
  }
}
```

### Systemd配置示例

```
[Unit]
Description=V2ray Subscription
After=network.target
Wants=network.target

[Service]
WorkingDirectory=/path/to/your application root
Environment=GIN_MODE=release
ExecStart=/path/to/your application root/v2ray-subscription-server
Restart=on-abnormal
RestartSec=5s
KillMode=mixed

StandardOutput=null
StandardError=syslog

[Install]
WantedBy=multi-user.target
```

### 服务器发布脚本示例

```
#!/usr/bin/env bash

tar zxf /path/to/target.tgz -C /path/to/your application root
chown www-data:www-data -R /path/to/your application root

systemctl restart v2ray-subscription
```
