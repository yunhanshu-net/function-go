#!/bin/bash

# HTTP上传配置
export UPLOAD_PROVIDER=http
export UPLOAD_DOMAIN=https://upload.example.com/api/upload
export DOWNLOAD_DOMAIN=https://cdn.example.com

# 启动应用
echo "启动应用，使用HTTP上传配置..."
echo "Provider: $UPLOAD_PROVIDER"
echo "Upload Domain: $UPLOAD_DOMAIN"
echo "Download Domain: $DOWNLOAD_DOMAIN"

# 这里替换为你的实际启动命令
# ./your-app --runner_id=your-runner-id 