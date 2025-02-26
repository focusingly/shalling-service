# Shalling Space
> 一个使用基于 golang 语言开发的轻量类 CMS 管理系统, 适合搭建个人博客, 简单的官网以及二次开发等, 使用了前后端分离的方式进行开发, 方便使用各自的技术栈进行开发, 也方便使用不同的技术栈实现 "前端/客户端", 同时, 如果你的配置较低(如轻量 vps)不适合部署 nginx, docker, redis 等服务, 也可以使用内置的 SPA 路由处理器 + 使用 Systemd 进行服务配置

## 主要的功能
- 提供了文章的管理; 分类管理; 标签管理管理, 以及基于分词 + 数据库特性实现的全文分页搜索(支持返回高亮位置和截取片段)
- 提供了文章多级评论的支持, 类似于 B 站 PC 的评论效果, 支持单独的评论管理
- 用户登录设备限制登录; 账户下线; 账户禁用功能
- 支持本地账户登录 + 第三方(Oauth2) 集成登录, 默认配置了 google 和 github
- 提供了站点的 UV 统计, 支持返回统计趋势
- 提供了接口限流管理, 手动添加封禁 IP 管理
- 基于基于分库的日志管理, 降低业务数据库的压力
- 公开展示的社交媒体管理
- 系统性能监控, 以及提供基于定时任务的资源高负载邮件预警
- 支持 md5 去重的文件上传, 以及图片自动转 webp 处理的本地文件存储功能; 以及添加 s3 文件直传集成, 默认集成了 cloudflare R2 的免费额度监控定时任务; 可在产生额外费用发送邮件预警
- 独立的 spa web 集成部署(使用 embed), 适合在缺乏 nginx 等环境下使用
- 集成友链管理
- 提供了动态菜单列表编辑功能, 可绑定站点文章或者外链为独立选项
- 提供了定时任务管理(基于 cron 表达式的处理); 提供了实用的预设的任务;
- 独立的邮件配置, 支持配置多个发件配置
- 基于中间件的客户端 IP 地址/属地提取; 基于中间件的客户端标识提取
- 基于中间件实现的统一自适应协商返回响应(json,xml, yaml, toml)
- 基于条件编译的 h2c 和 http2 默认集成; 以及自适应响应体压缩(br, gzip...) 返回, 对于站点和静态文件, 默认集成了强缓存/协商缓存的处理
- ...

## 安装部署
- 对于个人博客或者低配 vps, 默认使用 sqlite3 作为底层数据库, 在 Linux 平台上使用 systemd 进行部署, 可以参考底下这个配置:
  ```toml
  [Unit]
  Description=Your App
  After=network.target

  [Service]
  ExecStart=<your app path>
  WorkingDirectory=<your log path>
  Restart=on-failure
  User=nobody
  Group=nogroup

  [Install]
  WantedBy=multi-user.target
  ```
- 如果有其它需求, 如使用 docker 部署运行或者使用更加稳健的数据库: 如 PostgresSQL 或者 MySQL, 可以自行修改配置文件的连接配置选项, 并引入对应的驱动即可

## 开发环境
- 项目提供了配置文件的 schema 文件, 如过您使用 vscode 进行开发, 可以参考如下配置, 以获得配置文件的补全提示
    > 编辑 .vscode/settings.json 文件, 添加如下内容
    ```json
    "yaml.schemaStore.enable": true,
    "yaml.schemas": {
        ".vscode/config-yaml-schema.json": "config*.yml"
    }
    ```
