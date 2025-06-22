# 文件上传配置说明

## 环境变量配置

系统支持通过环境变量配置文件上传参数，支持多种云存储提供商。

### 通用配置

| 环境变量 | 说明 | 默认值 | 必填 |
|---------|------|--------|------|
| `UPLOAD_PROVIDER` | 上传提供商 | `qiniu` | 否 |
| `DOWNLOAD_DOMAIN` | 下载域名 | `http://cdn.geeleo.com` | 否 |

### 七牛云配置 (provider=qiniu)

| 环境变量 | 说明 | 默认值 | 必填 |
|---------|------|--------|------|
| `UPLOAD_BUCKET` | 存储桶名称 | `geeleo` | 是 |
| `UPLOAD_ACCESS_KEY` | 访问密钥 | 已设置默认值 | 是 |
| `UPLOAD_SECRET_KEY` | 私钥 | 已设置默认值 | 是 |
| `UPLOAD_TOKEN` | 上传Token | 空 | 否 |

### HTTP上传配置 (provider=http)

| 环境变量 | 说明 | 默认值 | 必填 |
|---------|------|--------|------|
| `UPLOAD_DOMAIN` | 上传接口地址 | 空 | 是 |

## 配置示例

### 七牛云配置

```bash
export UPLOAD_PROVIDER=qiniu
export UPLOAD_BUCKET=my-bucket
export UPLOAD_ACCESS_KEY=your-access-key
export UPLOAD_SECRET_KEY=your-secret-key
export DOWNLOAD_DOMAIN=https://cdn.example.com
```

### HTTP上传配置

```bash
export UPLOAD_PROVIDER=http
export UPLOAD_DOMAIN=https://upload.example.com/api/upload
export DOWNLOAD_DOMAIN=https://cdn.example.com
```

### 使用上传Token（七牛云）

```bash
export UPLOAD_PROVIDER=qiniu
export UPLOAD_BUCKET=my-bucket
export UPLOAD_TOKEN=your-upload-token
export DOWNLOAD_DOMAIN=https://cdn.example.com
```

## 支持的提供商

- `qiniu` - 七牛云存储
- `http` - 通用HTTP multipart上传
- `aliyun` - 阿里云OSS（待实现）
- `aws` - AWS S3（待实现）

## 配置优先级

1. 环境变量
2. 默认值

## 注意事项

1. 七牛云配置中，如果设置了 `UPLOAD_TOKEN`，则会优先使用Token，否则使用AccessKey/SecretKey生成Token
2. HTTP上传需要确保上传接口支持multipart/form-data格式
3. 下载域名建议使用CDN域名以提高访问速度
4. 生产环境建议通过环境变量或配置文件管理敏感信息，避免硬编码 