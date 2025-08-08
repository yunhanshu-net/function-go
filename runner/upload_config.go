package runner

import (
	"fmt"
	"github.com/yunhanshu-net/function-go/env"
	"os"

	"github.com/yunhanshu-net/pkg/trace"
)

// getUploadConfig 获取上传配置
func getUploadConfig() trace.UploadConfig {
	// 优先从环境变量读取
	provider := getEnvOrDefault("UPLOAD_PROVIDER", "qiniu")
	uploadDomain := getEnvOrDefault("UPLOAD_DOMAIN", "")
	downloadDomain := getEnvOrDefault("DOWNLOAD_DOMAIN", "http://cdn.geeleo.com")
	uploadToken := getEnvOrDefault("UPLOAD_TOKEN", "")
	bucket := getEnvOrDefault("UPLOAD_BUCKET", "geeleo")
	accessKey := getEnvOrDefault("UPLOAD_ACCESS_KEY", "ehF_E4x_EyO_wSN_nwqExyhXPe5hGl5Xjo89_cZ6")
	secretKey := getEnvOrDefault("UPLOAD_SECRET_KEY", "FjfIpqUevEcVx9bQxdgiuX9Di-CUOrKFkR88CZAj")

	config := trace.UploadConfig{
		Provider:       provider,
		UploadDomain:   uploadDomain,
		DownloadDomain: downloadDomain,
		UploadToken:    uploadToken,
		Bucket:         bucket,
		AccessKey:      accessKey,
		SecretKey:      secretKey,
	}

	// 打印配置信息（隐藏敏感信息）
	printUploadConfig(config)

	return config
}

// printUploadConfig 打印上传配置（隐藏敏感信息）
func printUploadConfig(config trace.UploadConfig) {
	fmt.Printf("[上传配置] Provider: %s\n", config.Provider)
	fmt.Printf("[上传配置] Bucket: %s\n", config.Bucket)
	fmt.Printf("[上传配置] UploadDomain: %s\n", config.UploadDomain)
	fmt.Printf("[上传配置] DownloadDomain: %s\n", config.DownloadDomain)

	// 隐藏敏感信息
	if config.AccessKey != "" {
		fmt.Printf("[上传配置] AccessKey: %s***\n", config.AccessKey[:min(8, len(config.AccessKey))])
	}
	if config.SecretKey != "" {
		fmt.Printf("[上传配置] SecretKey: %s***\n", config.SecretKey[:min(8, len(config.SecretKey))])
	}
	if config.UploadToken != "" {
		fmt.Printf("[上传配置] UploadToken: %s***\n", config.UploadToken[:min(8, len(config.UploadToken))])
	}
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// createFunctionMsg 创建FunctionMsg
func createFunctionMsg(traceId string, method string, router string) *trace.FunctionMsg {
	return &trace.FunctionMsg{
		User:         env.User,
		Runner:       env.Name,
		Version:      env.Version,
		Method:       method,
		Router:       router,
		TraceID:      traceId,
		UploadConfig: getUploadConfig(),
	}
}

// validateUploadConfig 验证上传配置
func validateUploadConfig(config trace.UploadConfig) error {
	switch config.Provider {
	case "qiniu":
		if config.Bucket == "" {
			return fmt.Errorf("七牛云上传需要配置UPLOAD_BUCKET")
		}
		if config.AccessKey == "" || config.SecretKey == "" {
			if config.UploadToken == "" {
				return fmt.Errorf("七牛云上传需要配置UPLOAD_ACCESS_KEY/UPLOAD_SECRET_KEY或UPLOAD_TOKEN")
			}
		}
	case "http":
		if config.UploadDomain == "" {
			return fmt.Errorf("HTTP上传需要配置UPLOAD_DOMAIN")
		}
	case "aliyun":
		return fmt.Errorf("阿里云OSS上传器暂未实现")
	case "aws":
		return fmt.Errorf("AWS S3上传器暂未实现")
	default:
		return fmt.Errorf("不支持的上传提供商: %s", config.Provider)
	}
	return nil
}
