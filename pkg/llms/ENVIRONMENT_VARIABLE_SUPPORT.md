# 环境变量自动支持功能

## 🎯 功能概述

为了方便开发和使用，所有LLM客户端现在都支持自动从环境变量获取API密钥。当传入的`apiKey`参数为空字符串时，客户端会自动尝试从相应的环境变量中获取API密钥。

## 🔧 支持的环境变量

| 客户端 | 环境变量名 | 示例值 |
|--------|------------|--------|
| GLM | `GLM_API_KEY` | `259dba772d8c4d5babde26a2fa762cf7.LkwGJ6V4rPnl27QK` |
| DeepSeek | `DEEPSEEK_API_KEY` | `sk-1a8d81e59fc84205b42a3cf210ff49fc` |
| Qwen | `QIANWEN_API_KEY` | `sk-75ff015dbdec41f0804cf0203f20a387` |

## 📝 使用方法

### 1. 设置环境变量

```bash
# 在你的 ~/.zshrc 或 ~/.bashrc 中添加
export GLM_API_KEY="your_glm_api_key_here"
export DEEPSEEK_API_KEY="your_deepseek_api_key_here"
export QIANWEN_API_KEY="your_qwen_api_key_here"
```

### 2. 使用客户端

```go
// 方法1：直接传入API密钥（优先级最高）
client := NewGLMClient("your-api-key")

// 方法2：传入空字符串，自动从环境变量获取
client := NewGLMClient("")

// 方法3：不传参数（Go会自动传入空字符串）
client := NewGLMClient("")
```

## 🏆 优先级规则

1. **传入的API密钥** - 最高优先级
2. **环境变量** - 当传入参数为空时使用
3. **空字符串** - 当环境变量也未设置时

## ✅ 测试验证

### 环境变量回退测试

```bash
cd function-go/pkg/llms
go test -v -run TestEnvironmentVariableFallback
```

### 流式聊天测试

```bash
# 测试GLM流式功能
go test -v -run TestChatStream/GLM_Stream

# 测试所有流式功能
go test -v -run TestChatStream
```

## 🔍 实现细节

### 修改的函数

- `NewGLMClient(apiKey string)` - 支持从 `GLM_API_KEY` 环境变量获取
- `NewDeepSeekClient(apiKey string)` - 支持从 `DEEPSEEK_API_KEY` 环境变量获取  
- `NewQwenClient(apiKey string)` - 支持从 `QIANWEN_API_KEY` 环境变量获取

### 核心逻辑

```go
func NewGLMClient(apiKey string) *GLMClient {
    // 如果传入的apiKey为空，尝试从环境变量获取
    if apiKey == "" {
        apiKey = os.Getenv("GLM_API_KEY")
    }
    return NewGLMClientWithOptions(apiKey, DefaultClientOptions())
}
```

## 🚀 优势

1. **开发便利性** - 无需在代码中硬编码API密钥
2. **安全性** - API密钥存储在环境变量中，不会意外提交到代码仓库
3. **向后兼容** - 原有的显式传入API密钥的方式仍然有效
4. **灵活性** - 支持不同环境使用不同的API密钥

## 📋 注意事项

1. **环境变量名称** - 确保使用正确的环境变量名称
2. **API密钥格式** - 不同提供商的API密钥格式可能不同
3. **权限设置** - 确保环境变量有正确的读取权限
4. **错误处理** - 当环境变量未设置时，客户端会使用空字符串作为API密钥

## 🎉 测试结果

✅ **GLM流式功能** - 完全正常工作，支持实时流式响应  
✅ **环境变量回退** - 所有客户端都支持自动从环境变量获取API密钥  
✅ **向后兼容** - 原有的API密钥传入方式仍然有效  
✅ **优先级正确** - 传入的API密钥优先于环境变量  

现在你可以更方便地使用LLM客户端了！只需要设置好环境变量，就可以直接使用 `NewGLMClient("")` 来创建客户端。
