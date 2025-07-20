package main

import (
	"context"
	"fmt"
	"reflect"

	"github.com/yunhanshu-net/function-go/runner"
)

// ConfigDemoConfig 配置演示的配置结构体
type ConfigDemoConfig struct {
	AppName     string `json:"app_name"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
	EnableLogging bool `json:"enable_logging"`
	EnableCache   bool `json:"enable_cache"`
	DebugMode     bool `json:"debug_mode"`
	MaxRetries    int  `json:"max_retries"`
	Timeout       int  `json:"timeout"`
	BatchSize     int  `json:"batch_size"`
	ApiEndpoint   string `json:"api_endpoint"`
	SecretKey     string `json:"secret_key"`
}

func main() {
	fmt.Println("=== 配置结构体注册调试 ===")

	// 创建上下文
	ctx := runner.NewContext(context.Background(), "POST", "/cmp/config_demo")
	
	// 生成配置键
	configKey := ctx.generateConfigKey()
	fmt.Printf("配置键: %s\n", configKey)

	// 获取配置管理器
	configManager := runner.GetConfigManager()

	// 注册配置结构体
	configStruct := &ConfigDemoConfig{
		AppName:       "配置演示应用",
		Version:       "1.0.0",
		Environment:   "开发环境",
		EnableLogging: true,
		EnableCache:   true,
		DebugMode:     false,
		MaxRetries:    3,
		Timeout:       30,
		BatchSize:     100,
		ApiEndpoint:   "https://api.example.com",
		SecretKey:     "demo_secret_key",
	}
	
	configManager.RegisterConfigStruct(configKey, configStruct)
	fmt.Printf("已注册配置结构体类型: %v\n", reflect.TypeOf(configStruct))

	// 检查注册是否成功
	configManager.mutex.RLock()
	registeredType, exists := configManager.configStructs[configKey]
	configManager.mutex.RUnlock()
	
	if exists {
		fmt.Printf("✓ 配置结构体已注册: %v\n", registeredType)
	} else {
		fmt.Printf("✗ 配置结构体未注册\n")
	}

	// 测试配置数据转换
	testData := map[string]interface{}{
		"app_name":       "测试应用",
		"version":        "2.0.0",
		"environment":    "测试环境",
		"enable_logging": true,
		"enable_cache":   false,
		"debug_mode":     true,
		"max_retries":    5,
		"timeout":        60,
		"batch_size":     200,
		"api_endpoint":   "https://test-api.example.com",
		"secret_key":     "test_secret_key",
	}

	// 创建配置数据
	configData := &runner.ConfigData{
		Type: "json",
		Data: testData,
	}

	// 测试转换
	result := configManager.GetConfigStruct(ctx, configKey)
	if result != nil {
		fmt.Printf("✓ 配置转换成功: %T\n", result)
		if config, ok := result.(ConfigDemoConfig); ok {
			fmt.Printf("✓ 类型断言成功: AppName=%s, Version=%s\n", config.AppName, config.Version)
		} else {
			fmt.Printf("✗ 类型断言失败: 期望 ConfigDemoConfig，实际 %T\n", result)
		}
	} else {
		fmt.Printf("✗ 配置转换失败\n")
	}
} 