# v2ray-helper

V2ray配置管理面板，可以对v2ray的客户端配置进行管理，同时该工具支持一键部署服务器

## 特性

- 客户端配置管理
- 服务器一键部署
- VMess链接生成
- 站点伪装（一键部署Cloudreve，暂未支持其他，有需要可以提交issue）
- 自动HTTPS证书申请

![alt 配置列表](https://github.com/Luna-CY/v2ray-helper/raw/master/resources/image/v2ray-helper-1.png)
![alt 配置列表](https://github.com/Luna-CY/v2ray-helper/raw/master/resources/image/v2ray-helper-2.png)

## 示例站点

[http://vh.example.luna.xin](http://vh.example.luna.xin)，测试站点的key为默认的中短横线"-"，测试站点的数据定期删除

## 安装教程

可以下载已打包的二进制包或手动构建

`mkdir /usr/local/v2ray-helper && tar zxf v2ray-helper-1.0.0.tgz -C /usr/local/v2ray-helper`

## 手动构建

克隆项目到本地，执行./pack.sh进行打包，将打包后的`v2ray-helper-x.x.x.tgz`上传到服务器

## 使用说明

首次运行时将会在执行程序文件所在目录自动创建运行所需文件，或者通过参数`-home-dir`指定运行时根目录

- `/usr/local/v2ray-helper/v2ray-helper` 直接运行
- `/usr/local/v2ray-helper/v2ray-helper -home-dir /usr/local/v2ray-helper/v2ray-helper` 直接运行并设置`-home-dir`

支持通过参数`-install`或`-install-and-enable`参数安装为系统服务

- `/usr/local/v2ray-helper/v2ray-helper -install` 安装为系统服务
- `/usr/local/v2ray-helper/v2ray-helper -install -home-dir /usr/local/v2ray-helper/v2ray-helper` 安装为系统服务并设置`-home-dir`
- `/usr/local/v2ray-helper/v2ray-helper -install-and-enable` 安装为系统服务并且开机启动

安装为服务后可通过`systemd`进行管理

- 启动：`service v2ray-helper start`
- 停止：`service v2ray-helper stop`
- 重启：`service v2ray-helper restart`
