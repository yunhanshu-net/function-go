#!/bin/bash

# 七牛云上传配置
export UPLOAD_PROVIDER=qiniu
export UPLOAD_BUCKET=geeleo
export UPLOAD_ACCESS_KEY=ehF_E4x_EyO_wSN_nwqExyhXPe5hGl5Xjo89_cZ6
export UPLOAD_SECRET_KEY=FjfIpqUevEcVx9bQxdgiuX9Di-CUOrKFkR88CZAj
export DOWNLOAD_DOMAIN=http://cdn.geeleo.com

# 启动应用
echo "启动应用，使用七牛云上传配置..."
echo "Provider: $UPLOAD_PROVIDER"
echo "Bucket: $UPLOAD_BUCKET"
echo "Download Domain: $DOWNLOAD_DOMAIN"

# 这里替换为你的实际启动命令
# ./your-app --runner_id=your-runner-id 