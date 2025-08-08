# 📁 文件上传路径规范指南

## 🎯 概述

本指南介绍了新的文件上传路径规范功能，该功能确保每个函数的文件上传都遵循统一的路径格式，便于文件管理和避免冲突。

## 📋 路径格式规范

### 基本格式
```
/{租户}/{应用}/{函数路由}/{HTTP方法}/output/{日期}/{文件名}
```

### 示例
```
/user123/myapp/api/upload/POST/output/20250628/document_1751092194257657000.pdf
/company/crm/customer/export/GET/output/20250628/report_1751092194257658000.csv
/admin/dashboard/report/generate/POST/output/20250628/chart_1751092194257659000.png
```

## 🚀 使用方法

### 1. 获取上传路径

```go
// 创建runner context
ctx := context.WithValue(context.Background(), "trace_id", "your-trace-id")
runnerCtx := NewContext(ctx, "POST", "api/upload")

// 获取规范的上传路径
uploadPath := runnerCtx.GetUploadPath()
fmt.Printf("上传路径: %s\n", uploadPath)
// 输出: /user/app/api/upload/POST/output/20250628
```

### 2. 直接上传文件

```go
// 从数据创建文件并自动上传
data := []byte("文件内容")
files, err := runnerCtx.CreateFilesFromData("document.txt", data)
if err != nil {
    log.Printf("上传失败: %v", err)
    return
}

// 检查上传结果
if len(files.GetFiles()) > 0 {
    file := files.GetFiles()[0]
    fmt.Printf("文件URL: %s\n", file.URL)
}
```

### 3. 从本地路径上传

```go
// 从本地文件上传
files, err := runnerCtx.CreateFilesFromPath("/path/to/local/file.pdf")
if err != nil {
    log.Printf("上传失败: %v", err)
    return
}
```

## 🔄 生命周期管理

### 临时文件（下载一次后删除）
```go
tempFiles := runnerCtx.NewTemporaryFiles()
err := tempFiles.AddFileFromData("temp_report.txt", data)
```

### 有效期文件（7天后过期）
```go
expiringFiles := runnerCtx.NewExpiringFiles()
err := expiringFiles.AddFileFromData("weekly_report.pdf", data)
```

### 永久文件（无限制）
```go
permanentFiles := runnerCtx.NewPermanentFiles()
err := permanentFiles.AddFileFromData("backup.zip", data)
```

## ⚙️ 配置说明

### 环境变量配置
```bash
# 上传提供商
UPLOAD_PROVIDER=qiniu

# 七牛云配置
UPLOAD_BUCKET=your-bucket
UPLOAD_ACCESS_KEY=your-access-key
UPLOAD_SECRET_KEY=your-secret-key
DOWNLOAD_DOMAIN=https://your-domain.com

# 或者直接使用上传Token
UPLOAD_TOKEN=your-upload-token
```

### 支持的上传提供商
- ✅ **qiniu** - 七牛云对象存储
- 🚧 **aliyun** - 阿里云OSS（待实现）
- 🚧 **aws** - AWS S3（待实现）
- ✅ **http** - HTTP multipart上传

## 🔧 Context方法

### NewContext 方法
统一的参数格式：

```go
// 标准格式
runnerCtx := NewContext(ctx, method, router)
```

注意：用户、应用名称、版本等信息会自动从 `env` 包中获取。

### 获取信息的方法
```go
// 获取上传路径
uploadPath := runnerCtx.GetUploadPath()

// 获取FunctionMsg对象
functionMsg := runnerCtx.GetFunctionMsg()

// 获取用户信息
user := runnerCtx.GetUsername()
```

### Files创建方法
```go
// 基础方法
files := runnerCtx.NewFiles(input)

// 生命周期方法
tempFiles := runnerCtx.NewTemporaryFiles()
expiringFiles := runnerCtx.NewExpiringFiles()
permanentFiles := runnerCtx.NewPermanentFiles()

// 便捷方法
files, err := runnerCtx.CreateFilesFromData(filename, data)
files, err := runnerCtx.CreateFilesFromPath(localPath)
```

## 📊 实际上传示例

### 七牛云上传路径示例
```
原始文件名: document.pdf
生成的Key: user123/myapp/api/upload/POST/output/20250628/document_1751092194257657000.pdf
访问URL: https://cdn.example.com/user123/myapp/api/upload/POST/output/20250628/document_1751092194257657000.pdf
```

### 路径组成部分
- **租户**: `user123` - 来自 FunctionMsg.User
- **应用**: `myapp` - 来自 FunctionMsg.Runner  
- **函数路由**: `api/upload` - 来自 FunctionMsg.Router
- **HTTP方法**: `POST` - 来自 FunctionMsg.Method
- **输出标识**: `output` - 固定值，区分输入和输出文件
- **日期**: `20250628` - 当前日期（YYYYMMDD格式）
- **文件名**: `document_1751092194257657000.pdf` - 原名+时间戳+扩展名

## 🔍 调试和测试

### 运行测试
```bash
cd function-go/runner
go test -v -run TestContext
```

### 查看上传配置
运行时会自动打印上传配置信息：
```
[上传配置] Provider: qiniu
[上传配置] Bucket: your-bucket
[上传配置] DownloadDomain: https://your-domain.com
[上传配置] AccessKey: your-key***
```

## 📝 最佳实践

1. **使用规范路径**: 始终通过 `runnerCtx.GetUploadPath()` 获取上传路径
2. **设置生命周期**: 根据文件用途选择合适的生命周期策略
3. **错误处理**: 妥善处理上传错误，提供用户友好的错误信息
4. **文件命名**: 使用有意义的文件名，系统会自动添加时间戳避免冲突
5. **配置管理**: 通过环境变量管理上传配置，避免硬编码

## 🚨 注意事项

1. **权限配置**: 确保上传服务的访问密钥有足够权限
2. **文件大小**: 注意上传服务的文件大小限制
3. **网络稳定**: 大文件上传时注意网络稳定性
4. **存储成本**: 合理设置文件生命周期以控制存储成本

## 🔗 相关文件

- `context.go` - Context实现和上传方法
- `upload_config.go` - 上传配置管理
- `qiniu_uploader.go` - 七牛云上传器实现
- `pkg/trace/function_msg.go` - FunctionMsg和路径生成
- `pkg/typex/files/` - Files类型和文件管理

---

📚 更多信息请参考相关源码和测试文件。 