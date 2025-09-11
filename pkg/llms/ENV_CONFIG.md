# 🔐 LLM API Key 环境变量配置指南

## 📋 概述

本项目支持通过环境变量配置各种LLM的API Key，提高安全性和灵活性。

## 🌍 环境变量配置

### 1. 千问系列 (阿里云百炼)

**重要**: 千问和千问3 Coder使用**同一个API Key**！

```bash
# 设置千问/千问3 Coder API Key
export QIANWEN_API_KEY="sk-75ff015dbdec41f0804cf0203f20a387"
```

**支持的模型**:
- `qwen3-coder-plus` - 千问3 Coder Plus
- `qwen3-coder-flash` - 千问3 Coder Flash
- `qwen2.5-coder-*` - 千问2.5 Coder系列

### 2. DeepSeek

```bash
# 设置DeepSeek API Key
export DEEPSEEK_API_KEY="sk-1a8d81e59fc84205b42a3cf210ff49fc"
```

### 3. GLM (智谱AI)

```bash
# 设置GLM API Key
export GLM_API_KEY="your-glm-api-key"
```

**支持的模型**:
- `glm-4.5` - 最强大的推理模型，3550亿参数
- `glm-4.5-air` - 高性价比轻量级强性能
- `glm-4.5-x` - 高性能强推理极速响应
- `glm-4.5-airx` - 轻量级强性能极速响应
- `glm-4.5-flash` - 免费高效多功能

### 4. 其他LLM (可选)

```bash
# 豆包
export DOUBAO_API_KEY="your-doubao-api-key"

# Kimi
export KIMI_API_KEY="your-kimi-api-key"

# Claude
export CLAUDE_API_KEY="your-claude-api-key"

# Gemini
export GEMINI_API_KEY="your-gemini-api-key"
```

## 🚀 使用方法

### 方法1: 从环境变量自动创建 (推荐)

```go
package main

import "github.com/yunhanshu-net/function-go/pkg/llms"

func main() {
    // 自动从环境变量读取API Key
    client, err := llms.NewQwen3CoderClientFromEnv()
    if err != nil {
        panic(err)
    }
    
    // 使用客户端...
}

// GLM使用示例
func glmExample() {
    // 创建GLM客户端
    client, err := llms.NewGLMClientFromEnv()
    if err != nil {
        panic(err)
    }
    
    // 使用GLM特殊功能
    glmClient := client.(*llms.GLMClient)
    glmClient.SetModel("glm-4.5") // 设置模型
    
    // 使用思考模式
    resp, err := glmClient.ChatWithThinking(ctx, req, true)
}
```

### 方法2: 手动指定API Key

```go
// 直接指定API Key
client := llms.NewQwen3CoderClient("sk-75ff015dbdec41f0804cf0203f20a387")
```

### 方法3: 通用环境变量创建

```go
// 使用通用函数
client, err := llms.NewLLMClientFromEnv(llms.ProviderQwen3Coder)
if err != nil {
    panic(err)
}
```

## 🔧 环境变量设置方法

### Linux/macOS

```bash
# 临时设置 (当前会话有效)
export QIANWEN_API_KEY="sk-75ff015dbdec41f0804cf0203f20a387"
export DEEPSEEK_API_KEY="sk-1a8d81e59fc84205b42a3cf210ff49fc"
export GLM_API_KEY="your-glm-api-key"

# 永久设置 (添加到 ~/.bashrc 或 ~/.zshrc)
echo 'export QIANWEN_API_KEY="sk-75ff015dbdec41f0804cf0203f20a387"' >> ~/.bashrc
echo 'export DEEPSEEK_API_KEY="sk-1a8d81e59fc84205b42a3cf210ff49fc"' >> ~/.bashrc
echo 'export GLM_API_KEY="your-glm-api-key"' >> ~/.bashrc
source ~/.bashrc
```

### Windows

```cmd
# 临时设置
set QIANWEN_API_KEY=sk-75ff015dbdec41f0804cf0203f20a387
set DEEPSEEK_API_KEY=sk-1a8d81e59fc84205b42a3cf210ff49fc
set GLM_API_KEY=your-glm-api-key

# 永久设置 (系统环境变量)
# 控制面板 -> 系统 -> 高级系统设置 -> 环境变量
```

### Docker

```dockerfile
# Dockerfile
ENV QIANWEN_API_KEY=sk-75ff015dbdec41f0804cf0203f20a387
ENV DEEPSEEK_API_KEY=sk-1a8d81e59fc84205b42a3cf210ff49fc
ENV GLM_API_KEY=your-glm-api-key
```

```yaml
# docker-compose.yml
services:
  app:
    environment:
      - QIANWEN_API_KEY=sk-75ff015dbdec41f0804cf0203f20a387
      - DEEPSEEK_API_KEY=sk-1a8d81e59fc84205b42a3cf210ff49fc
      - GLM_API_KEY=your-glm-api-key
```

## 🧪 测试环境变量配置

```bash
# 测试环境变量是否设置成功
echo $QIANWEN_API_KEY
echo $DEEPSEEK_API_KEY
echo $GLM_API_KEY

# 运行环境变量测试
go test -v -run TestEnvironmentVariableSupport
go test -v -run TestSameAPIKeyForQwen
```

## ⚠️ 安全注意事项

1. **不要将API Key提交到代码仓库**
2. **使用环境变量或配置文件存储敏感信息**
3. **定期轮换API Key**
4. **限制API Key的权限范围**

## 🔍 故障排除

### 问题1: "未提供API Key且环境变量中未找到配置"

**原因**: 环境变量未设置或名称错误

**解决方案**:
```bash
# 检查环境变量
env | grep API_KEY

# 重新设置环境变量
export QIANWEN_API_KEY="sk-75ff015dbdec41f0804cf0203f20a387"
```

### 问题2: "API Key无效"

**原因**: API Key过期或权限不足

**解决方案**:
1. 检查API Key是否正确
2. 确认API Key是否有效
3. 检查账户余额和权限

### 问题3: 环境变量在IDE中不生效

**解决方案**:
1. 重启IDE
2. 在IDE中设置环境变量
3. 使用IDE的环境变量配置文件

## 📚 相关文档

- [千问3 Coder使用指南](QWEN3_CODER_GUIDE.md)
- [DeepSeek集成说明](README_test.md)
- [LLM客户端接口文档](interface.go)

## 🎯 最佳实践

1. **开发环境**: 使用环境变量
2. **生产环境**: 使用密钥管理服务
3. **CI/CD**: 使用CI/CD平台的环境变量功能
4. **容器化**: 使用Docker secrets或Kubernetes secrets

---

祝你使用愉快！🚀
