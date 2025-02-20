#!/bin/bash

# 获取自定义域名
read -p "请输入域名（例如：example.com）：" domain

# 创建证书配置文件
cat > cert.conf << EOF
[req]
default_bits = 2048
prompt = no
default_md = sha256
x509_extensions = v3_req
distinguished_name = dn

[dn]
C = CN
ST = Beijing
L = Beijing
O = Development
OU = Testing Department
CN = $domain

[v3_req]
subjectAltName = @alt_names
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment

[alt_names]
DNS.1 = $domain
DNS.2 = *.$domain
IP.1 = 127.0.0.1
IP.2 = ::1
EOF

# 生成私钥和证书
openssl req \
    -new \
    -newkey rsa:2048 \
    -sha256 \
    -days 3650 \
    -nodes \
    -x509 \
    -keyout server.key \
    -out server.crt \
    -config cert.conf

# 验证证书内容
openssl x509 -in server.crt -text -noout

cat server.key server.crt > server.pem
# 设置适当的权限
chmod 600 server.key
chmod 644 server.crt

# 清理配置文件
rm cert.conf

echo "证书生成完成！"
echo "- 私钥文件: server.key"
echo "- 证书文件: server.crt"
echo "- pem 文件: server.pem"
