# 文件哈希计算器演示

## 功能概述

文件哈希计算器是一个强大的工具，可以：
- 从URL下载文件
- 计算文件的MD5、SHA1或SHA256哈希值
- 生成详细的哈希报告（文本和JSON格式）
- 自动上传报告文件到云存储
- 返回云存储下载链接

## API 路由

```
POST /files/file_hash
```

## 请求参数

```json
{
  "files": [
    {
      "name": "example.html",
      "url": "https://www.example.com/",
      "content_type": "text/html"
    }
  ],
  "algorithm": "md5"
}
```

### 参数说明

- `files`: 文件列表，支持多文件批量处理
  - `name`: 文件名
  - `url`: 文件下载URL
  - `content_type`: 文件类型（可选）
- `algorithm`: 哈希算法，支持 `md5`、`sha1`、`sha256`

## 响应结果

```json
{
  "success": true,
  "message": "文件哈希计算成功，共处理2个文件，生成了文本和JSON两种格式的报告",
  "files": [
    {
      "name": "hash_report_md5.txt",
      "url": "http://cdn.geeleo.com/hash_report_md5.txt",
      "size": 1024,
      "content_type": "text/plain"
    },
    {
      "name": "hash_report_md5.json",
      "url": "http://cdn.geeleo.com/hash_report_md5.json",
      "size": 512,
      "content_type": "application/json"
    }
  ],
  "hashes": [
    "d41d8cd98f00b204e9800998ecf8427e",
    "098f6bcd4621d373cade4e832627b4f6"
  ]
}
```

## 生成的报告文件

### 文本格式报告 (hash_report_md5.txt)

```
文件哈希计算报告 (md5)
===========================================

文件 1:
  文件名: example.html
  文件大小: 1256 字节
  文件类型: text/html
  哈希算法: md5
  哈希值: d41d8cd98f00b204e9800998ecf8427e

文件 2:
  文件名: httpbin.json
  文件大小: 429 字节
  文件类型: application/json
  哈希算法: md5
  哈希值: 098f6bcd4621d373cade4e832627b4f6

报告生成时间: 2024-01-15 14:30:25
总计处理文件: 2 个
```

### JSON格式报告 (hash_report_md5.json)

```json
{
  "algorithm": "md5",
  "total_files": 2,
  "generated_at": "2024-01-15 14:30:25",
  "results": [
    {
      "file_name": "example.html",
      "file_size": 1256,
      "algorithm": "md5",
      "hash_value": "d41d8cd98f00b204e9800998ecf8427e",
      "content_type": "text/html"
    },
    {
      "file_name": "httpbin.json",
      "file_size": 429,
      "algorithm": "md5",
      "hash_value": "098f6bcd4621d373cade4e832627b4f6",
      "content_type": "application/json"
    }
  ]
}
```

## 测试方法

### 1. 使用测试脚本

```bash
./function-go/scripts/test_files_hash.sh
```

### 2. 使用测试数据

测试数据文件位于：`function-go/test_data/files_hash_request.json`

### 3. 环境配置

确保设置了以下环境变量：

```bash
export UPLOAD_PROVIDER=qiniu
export UPLOAD_BUCKET=geeleo
export UPLOAD_ACCESS_KEY=your_access_key
export UPLOAD_SECRET_KEY=your_secret_key
export DOWNLOAD_DOMAIN=http://cdn.geeleo.com
```

## 工作流程

1. **文件下载**: 从提供的URL下载文件到临时目录
2. **哈希计算**: 使用指定算法计算每个文件的哈希值
3. **报告生成**: 创建文本和JSON两种格式的详细报告
4. **文件上传**: 自动上传报告文件到云存储
5. **响应返回**: 返回包含云存储下载链接的响应

## 特性亮点

- ✅ **多文件支持**: 支持批量处理多个文件
- ✅ **多算法支持**: MD5、SHA1、SHA256
- ✅ **双格式报告**: 文本和JSON两种格式
- ✅ **自动上传**: 报告文件自动上传到云存储
- ✅ **详细信息**: 包含文件大小、类型等元数据
- ✅ **错误处理**: 完善的错误处理和日志记录

## 使用场景

- 文件完整性验证
- 批量文件哈希计算
- 文件变更检测
- 安全审计
- 数据备份验证 