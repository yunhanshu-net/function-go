#!/bin/bash

echo "=== 测试文件哈希上传功能 ==="

# 设置七牛云上传配置
export UPLOAD_PROVIDER=qiniu
export UPLOAD_BUCKET=geeleo
export UPLOAD_ACCESS_KEY=ehF_E4x_EyO_wSN_nwqExyhXPe5hGl5Xjo89_cZ6
export UPLOAD_SECRET_KEY=FjfIpqUevEcVx9bQxdgiuX9Di-CUOrKFkR88CZAj
export DOWNLOAD_DOMAIN=http://cdn.geeleo.com

echo "1. 当前上传配置："
echo "   Provider: $UPLOAD_PROVIDER"
echo "   Bucket: $UPLOAD_BUCKET"
echo "   Download Domain: $DOWNLOAD_DOMAIN"
echo ""

# 检查测试数据文件是否存在
TEST_FILE="function-go/test_data/file_hash_request.json"
if [ ! -f "$TEST_FILE" ]; then
    echo "错误：测试数据文件不存在: $TEST_FILE"
    exit 1
fi

echo "2. 测试数据文件内容："
cat "$TEST_FILE"
echo ""
echo ""

echo "3. 开始执行文件哈希计算..."
echo "   这将测试以下功能："
echo "   - 从URL下载文件"
echo "   - 计算文件哈希值"
echo "   - 生成文本和JSON报告"
echo "   - 上传报告文件到云存储"
echo ""

# 这里应该调用你的应用程序
# 示例命令（需要根据实际情况调整）：
echo "执行命令："
echo "go run main.go run --file=$TEST_FILE --trace_id=test-$(date +%s)"
echo ""

echo "注意：请确保你的应用程序已经编译并且可以运行"
echo "预期结果："
echo "- 成功下载测试文件"
echo "- 计算MD5哈希值"
echo "- 生成 hash_report_md5.txt 和 hash_report_md5.json 两个报告文件"
echo "- 上传报告文件到七牛云"
echo "- 返回包含云存储下载链接的响应"

echo ""
echo "=== 测试准备完成 ===" 