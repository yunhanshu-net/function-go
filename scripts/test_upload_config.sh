#!/bin/bash

echo "=== 测试上传配置 ==="

# 测试七牛云配置
echo "1. 测试七牛云配置..."
export UPLOAD_PROVIDER=qiniu
export UPLOAD_BUCKET=test-bucket
export UPLOAD_ACCESS_KEY=test-access-key
export UPLOAD_SECRET_KEY=test-secret-key
export DOWNLOAD_DOMAIN=https://cdn.example.com

echo "环境变量设置："
echo "UPLOAD_PROVIDER=$UPLOAD_PROVIDER"
echo "UPLOAD_BUCKET=$UPLOAD_BUCKET"
echo "UPLOAD_ACCESS_KEY=$UPLOAD_ACCESS_KEY"
echo "UPLOAD_SECRET_KEY=$UPLOAD_SECRET_KEY"
echo "DOWNLOAD_DOMAIN=$DOWNLOAD_DOMAIN"
echo ""

# 测试HTTP配置
echo "2. 测试HTTP配置..."
export UPLOAD_PROVIDER=http
export UPLOAD_DOMAIN=https://upload.example.com/api/upload
export DOWNLOAD_DOMAIN=https://cdn.example.com

echo "环境变量设置："
echo "UPLOAD_PROVIDER=$UPLOAD_PROVIDER"
echo "UPLOAD_DOMAIN=$UPLOAD_DOMAIN"
echo "DOWNLOAD_DOMAIN=$DOWNLOAD_DOMAIN"
echo ""

# 清理环境变量
echo "3. 测试默认配置（清理环境变量）..."
unset UPLOAD_PROVIDER
unset UPLOAD_BUCKET
unset UPLOAD_ACCESS_KEY
unset UPLOAD_SECRET_KEY
unset UPLOAD_DOMAIN
unset DOWNLOAD_DOMAIN
unset UPLOAD_TOKEN

echo "环境变量已清理，将使用默认配置"
echo ""

echo "=== 配置测试完成 ==="
echo "请运行你的应用来查看实际的配置加载情况" 