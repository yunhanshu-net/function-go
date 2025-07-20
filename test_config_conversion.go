package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yunhanshu-net/function-go/runner"
)

// ConfigDemoConfig 配置演示结构体
type ConfigDemoConfig struct {
	AppName       string `json:"app_name"`
	Environment   string `json:"environment"`
	Version       string `json:"version"`
	APIEndpoint   string `json:"api_endpoint"`
	SecretKey     string `json:"secret_key"`
	DebugMode     bool   `json:"debug_mode"`
	EnableLogging bool   `json:"enable_logging"`
	EnableCache   bool   `json:"enable_cache"`
	MaxRetries    int    `json:"max_retries"`
	Timeout       int    `json:"timeout"`
	BatchSize     int    `json:"batch_size"`
}

func main() {
	fmt.Println("=== 配置类型转换测试 ===")

	// 创建上下文
	ctx := runner.NewContext(context.Background(), "POST", "/cmp/config_demo")

	// 获取配置
	config := ctx.GetConfig()
	if config == nil {
		log.Fatal("配置为空")
	}

	fmt.Printf("配置类型: %T\n", config)
	fmt.Printf("配置内容: %+v\n", config)

	// 尝试类型断言
	if configStruct, ok := config.(ConfigDemoConfig); ok {
		fmt.Println("\n=== 配置转换成功 ===")
		fmt.Printf("应用名称: %s\n", configStruct.AppName)
		fmt.Printf("环境: %s\n", configStruct.Environment)
		fmt.Printf("版本: %s\n", configStruct.Version)
		fmt.Printf("API端点: %s\n", configStruct.APIEndpoint)
		fmt.Printf("调试模式: %t\n", configStruct.DebugMode)
		fmt.Printf("启用日志: %t\n", configStruct.EnableLogging)
		fmt.Printf("启用缓存: %t\n", configStruct.EnableCache)
		fmt.Printf("最大重试: %d\n", configStruct.MaxRetries)
		fmt.Printf("超时: %d\n", configStruct.Timeout)
		fmt.Printf("批处理大小: %d\n", configStruct.BatchSize)
	} else {
		fmt.Println("\n=== 配置转换失败 ===")
		fmt.Printf("期望类型: ConfigDemoConfig\n")
		fmt.Printf("实际类型: %T\n", config)
	}

	fmt.Println("\n=== 测试完成 ===")
} 