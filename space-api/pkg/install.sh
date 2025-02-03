#!/bin/sh

echo "Prepare Install Shalling Blog"

dir=/etc/systemd/system
serviceFile=./shalling-blog.service
serviceName=shalling-blog.service

# 检查当前用户是否为 root
if [ "$(id -u)" -ne 0 ]; then
    echo "该脚本必须以 root 权限运行"
    exit 1
else
    echo "以 root 权限运行"
fi

install() {
    useradd shalling-blog
    mkdir -p /opt/shalling-blog && chown -R shalling-blog:shalling-blog
    cp ${serviceFile} ${dir} && systemctl daemon-reload && systemctl enable ${serviceName} && systemctl restart ${serviceName}
}

remove() {
    echo "停止服务"
    systemctl stop ${serviceName}
    echo "移除服务"
    systemctl disable ${serviceName}
    rm -f ${dir}/${serviceName}
    systemctl daemon-reload
    echo "移除服务完成"
}
