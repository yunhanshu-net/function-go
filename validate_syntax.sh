#!/bin/bash

echo "=== 验证Go语法 ==="

# 检查是否有go命令
if ! command -v go &> /dev/null; then
    echo "警告: 未找到go命令，无法进行语法检查"
    echo "请确保已安装Go环境"
    exit 0
fi

echo "检查循环导入问题..."

# 尝试编译api包
echo "编译 pkg/dto/api 包..."
if go build ./pkg/dto/api; then
    echo "✅ pkg/dto/api 包编译成功"
else
    echo "❌ pkg/dto/api 包编译失败"
    exit 1
fi

# 尝试编译widget包
echo "编译 view/widget 包..."
if go build ./view/widget; then
    echo "✅ view/widget 包编译成功"
else
    echo "❌ view/widget 包编译失败"
    exit 1
fi

echo "✅ 所有包编译成功，没有循环导入问题" 