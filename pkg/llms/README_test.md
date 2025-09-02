# DeepSeek API Key 测试指南

## 🚀 快速开始

### 1. 运行单元测试

```bash
# 进入测试目录
cd function-go/pkg/llms

# 运行所有测试
go test -v

# 运行特定测试
go test -v -run TestDeepSeekChatBasic

# 运行性能测试
go test -v -bench=BenchmarkDeepSeekChat

# 跳过集成测试（快速测试）
go test -v -short
```

### 2. 使用测试脚本

```bash
# 给脚本添加执行权限
chmod +x run_deepseek_test.sh

# 运行测试脚本
./run_deepseek_test.sh
```

### 3. 测试你的API Key

你的API Key已经配置在测试文件中：
```go
const testAPIKey = "sk-1a8d81e59fc84205b42a3cf210ff49fc"
```

## 📋 测试内容

### 基础功能测试
- ✅ 客户端创建和配置
- ✅ 接口实现检查
- ✅ 基本聊天功能
- ✅ 系统提示功能
- ✅ 错误处理机制
- ✅ 超时处理
- ✅ 模型参数覆盖
- ✅ 默认值设置

### 高级测试
- 🔍 错误处理测试
- ⚡ 性能基准测试
- 🎯 集成测试
- 📊 Token使用统计

## 🎯 测试结果解读

### 成功情况
```
✅ AI回答: [AI的实际回答内容]
📊 Token使用: 输入X, 输出Y, 总计Z
```

### 常见错误
```
⚠️  API返回错误: [错误信息]
   这可能是API key无效或网络问题，请检查配置
```

### 网络错误
```
❌ 请求失败: [网络错误信息]
```

## 🔧 故障排除

### 1. API Key无效
- 检查API Key是否正确
- 确认API Key是否已激活
- 检查账户余额

### 2. 网络问题
- 检查网络连接
- 确认防火墙设置
- 尝试使用代理

### 3. 超时问题
- 增加超时时间
- 检查网络延迟
- 减少请求的token数量

## 📝 自定义测试

### 修改测试参数
```go
req := &ChatRequest{
    Messages: []Message{
        {Role: "user", Content: "你的问题"},
    },
    MaxTokens:   500,      // 增加输出长度
    Temperature: 0.5,      // 调整创造性
}
```

### 添加新的测试用例
```go
func TestCustomScenario(t *testing.T) {
    client := NewDeepSeekClient(testAPIKey)
    
    req := &ChatRequest{
        Messages: []Message{
            {Role: "user", Content: "自定义测试内容"},
        },
        MaxTokens: 200,
    }
    
    // 你的测试逻辑
}
```

## 🚀 性能优化

### 1. 并发测试
```bash
# 运行并发测试
go test -v -bench=BenchmarkDeepSeekChat -cpu=1,2,4
```

### 2. 内存分析
```bash
# 内存分析
go test -v -bench=BenchmarkDeepSeekChat -memprofile=mem.prof
```

### 3. CPU分析
```bash
# CPU分析
go test -v -bench=BenchmarkDeepSeekChat -cpuprofile=cpu.prof
```

## 📊 测试报告

测试完成后，你会看到详细的测试结果，包括：
- 测试通过/失败状态
- AI回答内容
- Token使用统计
- 错误信息（如果有）
- 性能指标

## 💡 提示

1. **首次测试**：建议先运行基本功能测试
2. **网络环境**：确保网络连接稳定
3. **API限制**：注意API调用频率限制
4. **错误处理**：测试会显示详细的错误信息
5. **性能测试**：性能测试可能需要较长时间

## 🔗 相关链接

- [DeepSeek API 文档](https://platform.deepseek.com/)
- [Go 测试文档](https://golang.org/pkg/testing/)
- [function-go 框架](https://github.com/yunhanshu-net/function-go)



