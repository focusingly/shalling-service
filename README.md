# Shalling Space

# 这是什么
一个使用基于 golang 语言开发的轻量 CMS 管理系统, 适合搭建个人的博客, 简单的官网以及二次开发等, 使用了前后端分离的方式进行开发, 方便使用各自的技术栈进行开发, 也方便使用不同的技术栈实现 "前端/客户端"

## 安装部署
- 对于个人博客或者低配 vps, 默认使用 sqlite3 作为底层数据库, 在 Linux 平台上使用 systemd 进行部署
- 如果有其它需求, 如使用 docker 部署运行或者使用更加稳健的数据库: 如 PostgresSQL 或者 MySQL, 可以自行修改配置文件的连接配置选项, 并引入对应的驱动即可, 项目提供了默认的 dockerfile 和 docker-compose.yml 参考文件

## 开发环境
- 项目提供了配置文件的 schema 文件, 如过您使用 vscode 进行开发, 可以参考如下配置, 以获得配置文件的补全提示
    > 编辑 .vscode/settings.json 文件, 添加如下内容
    ```json
    "yaml.schemaStore.enable": true,
    "yaml.schemas": {
        ".vscode/config-yaml-schema.json": "config*.yml"
    }
    ```
