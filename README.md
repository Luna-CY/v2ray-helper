# v2ray-endpoint

#### 介绍

V2ray节点管理面板


#### 安装教程

克隆项目到本地，执行./pack.sh进行打包，将打包后的`target.gz`上传到服务器

在服务器创建`/data/v2ray-subscription`目录

#### 使用说明

`target.gz`文件目录说明

```
/ 根目录
-- /v2ray-subscription-server      服务器可执行二进制文件
-- /v2ray-subscription-migrate     数据库迁移工具
-- /config                         服务器配置文件目录
-- /migrations                     数据库迁移文件
-- /web                            Web服务的根目录，内含index.html以及所有的静态资源文件
```

服务器本地化配置支持在`config`目录下创建`db.local.config.yaml`与`main.local.config.yaml`配置文件，配置项与`db.prod.config.yaml`、`main.prod.config.yaml`文件相同，同名配置可进行覆盖
