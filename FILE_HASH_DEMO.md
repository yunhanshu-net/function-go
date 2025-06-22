# 文件哈希计算器演示

## 功能概述

这个文件哈希计算器演示了完整的文件处理和上传流程：

1. **文件下载** - 从URL下载用户指定的文件
2. **哈希计算** - 支持MD5、SHA1、SHA256算法
3. **报告生成** - 生成文本和JSON两种格式的详细报告
4. **文件上传** - 自动上传报告文件到云存储
5. **结果返回** - 返回可下载的报告文件链接

## 测试流程

### 1. 准备测试环境

```bash
# 设置上传配置
export UPLOAD_PROVIDER=qiniu
export UPLOAD_BUCKET=geeleo
export UPLOAD_ACCESS_KEY=your-access-key
export UPLOAD_SECRET_KEY=your-secret-key
export DOWNLOAD_DOMAIN=http://cdn.geeleo.com
```

### 2. 运行测试

```bash
# 使用测试脚本
./function-go/scripts/test_file_hash.sh

# 或直接运行
go run main.go run --file=function-go/test_data/file_hash_request.json --trace_id=test-123
```

### 3. 预期结果

#### 输入
```json
{
  "files": [
    {
      "name": "test1.txt",
      "url": "https://httpbin.org/base64/SGVsbG8gV29ybGQh",
      "content_type": "text/plain"
    }
  ],
  "algorithm": "md5"
}
```

#### 输出
```json
{
  "success": true,
  "message": "文件哈希计算成功，共处理1个文件，生成了文本和JSON两种格式的报告",
  "files": [
    {
      "name": "hash_report_md5.txt",
      "url": "http://cdn.geeleo.com/uploads/2024/01/01/hash_report_md5_123456.txt",
      "content_type": "text/plain"
    },
    {
      "name": "hash_report_md5.json", 
      "url": "http://cdn.geeleo.com/uploads/2024/01/01/hash_report_md5_123456.json",
      "content_type": "application/json"
    }
  ],
  "hashes": ["5d41402abc4b2a76b9719d911017c592"]
}
```

## 生成的报告文件

### 文本报告 (hash_report_md5.txt)
```
文件哈希计算报告 (md5)
===========================================

文件 1:
  文件名: test1.txt
  文件大小: 12 字节
  文件类型: text/plain
  哈希算法: md5
  哈希值: 5d41402abc4b2a76b9719d911017c592

报告生成时间: 2024-01-01 12:00:00
总计处理文件: 1 个
```

### JSON报告 (hash_report_md5.json)
```json
{
  "algorithm": "md5",
  "total_files": 1,
  "generated_at": "2024-01-01 12:00:00",
  "results": [
    {
      "file_name": "test1.txt",
      "file_size": 12,
      "algorithm": "md5",
      "hash_value": "5d41402abc4b2a76b9719d911017c592",
      "content_type": "text/plain"
    }
  ]
}
```

## 技术实现亮点

### 1. 自动文件上传
- 在响应构建时自动触发上传
- 支持多种云存储提供商
- 透明的错误处理

### 2. 多格式报告
- 人类可读的文本格式
- 机器可解析的JSON格式
- 详细的文件元信息

### 3. 配置化设计
- 环境变量配置
- 多种哈希算法支持
- 灵活的云存储配置

### 4. 完整的错误处理
- 文件下载失败处理
- 哈希计算错误处理
- 上传失败回滚机制

## 扩展可能性

1. **支持更多哈希算法** - SHA512、Blake2等
2. **批量文件处理** - 压缩包解压和批量处理
3. **文件完整性验证** - 对比已知哈希值
4. **病毒扫描集成** - 文件安全检查
5. **文件格式转换** - 自动格式转换功能

## 使用场景

- **文件完整性验证** - 确保文件传输完整性
- **重复文件检测** - 通过哈希值识别重复文件
- **数据备份验证** - 验证备份文件的完整性
- **安全审计** - 文件变更检测和记录 