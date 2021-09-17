# v2ray-helper

V2ray一键部署服务，支持对V2ray常用代理协议(TCP/KCP/WebSocket/HTTP2)的一键部署与服务伪装，可视化UI面板配置

## 特性

- 可视化配置面板
- 客户端配置管理
- 服务器一键部署
- VMess链接生成
- 站点伪装(一键部署Cloudreve，暂未支持其他，有需要可以提交issue)
- HTTPS证书管理(自动申请、续期)

![alt 配置列表](https://github.com/Luna-CY/v2ray-helper/raw/master/resources/image/v2ray-helper-1.png)
![alt 服务器部署](https://github.com/Luna-CY/v2ray-helper/raw/master/resources/image/v2ray-helper-2.png)

## 安装教程

可以下载已打包的二进制包或手动构建

`mkdir /usr/local/v2ray-helper && tar zxf v2ray-helper-x.x.x.tgz -C /usr/local/v2ray-helper`

## 手动构建

克隆项目到本地，执行./pack.sh进行打包，将打包后的`v2ray-helper-x.x.x.tgz`上传到服务器

## 使用说明

首次运行时将会在执行程序文件所在目录自动创建运行所需文件，或者通过参数`-home-dir`指定运行时根目录

- `/usr/local/v2ray-helper/v2ray-helper` 直接运行
- `/usr/local/v2ray-helper/v2ray-helper -home-dir /usr/local/v2ray-helper/v2ray-helper` 直接运行并设置`-home-dir`

支持通过参数`-install`或`-install-with-https`参数安装为系统服务

- `/usr/local/v2ray-helper/v2ray-helper -install` 安装为系统服务
- `/usr/local/v2ray-helper/v2ray-helper -install -home-dir /usr/local/v2ray-helper/v2ray-helper` 安装为系统服务并设置`-home-dir`
- `/usr/local/v2ray-helper/v2ray-helper -install-with-https -https-host your.host` 安装为系统服务且启用TLS，此方式安装时必须提供域名

安装为服务后可通过`systemd`进行管理

- 启动：`service v2ray-helper start`
- 停止：`service v2ray-helper stop`
- 重启：`service v2ray-helper restart`
