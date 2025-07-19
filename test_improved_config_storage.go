package main

import (
	"encoding/json"
	"fmt"

	"github.com/yunhanshu-net/function-go/runner"
	"github.com/yunhanshu-net/function-go/pkg/dto/syscallback"
)

// TestConfig 测试配置结构体
type TestConfig struct {
	AppName     string `json:"app_name"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
	EnableLog   bool   `json:"enable_log"`
	Timeout     int    `json:"timeout"`
}

func main() {
	fmt.Println("=== 测试改进后的配置存储格式 ===")

	// 初始化配置管理器
	configManager := runner.GetConfigManager()
	localStorage := runner.NewLocalFileStorage("./test_configs")
	configManager.SetStorage(localStorage)

	// 创建测试上下文
	ctx := runner.NewContext(nil, "POST", "/test/improved")

	// 测试配置
	testConfig := TestConfig{
		AppName:     "改进测试应用",
		Version:     "2.0.0",
		Environment: "测试环境",
		EnableLog:   true,
		Timeout:     60,
	}

	configKey := "function.test.improved.post"

	fmt.Println("1. 注册配置结构体")
	configManager.RegisterConfigStruct(configKey, testConfig)

	fmt.Println("2. 创建改进后的配置数据")
	// 新的方式：直接存储配置对象
	configData := &syscallback.ConfigData{
		Type: "json",
		Data: testConfig, // 直接存储对象，不再双重序列化
	}

	fmt.Println("3. 写入配置")
	err := configManager.UpdateConfig(ctx, configKey, configData)
	if err != nil {
		fmt.Printf("❌ 写入配置失败: %v\n", err)
		return
	}
	fmt.Println("✅ 配置写入成功")

	fmt.Println("4. 读取配置")
	retrievedConfig := configManager.GetByKey(ctx, configKey)
	if retrievedConfig == nil {
		fmt.Println("❌ 读取配置失败")
		return
	}

	// 输出配置数据格式
	configJson, _ := json.MarshalIndent(retrievedConfig, "", "  ")
	fmt.Printf("存储的配置格式:\n%s\n", string(configJson))

	fmt.Println("5. 获取配置结构体")
	configStruct := configManager.GetConfigStruct(ctx, configKey)
	if configStruct == nil {
		fmt.Println("❌ 获取配置结构体失败")
		return
	}

	if testConfigRetrieved, ok := configStruct.(TestConfig); ok {
		fmt.Printf("✅ 成功获取配置结构体: AppName=%s, Version=%s, Environment=%s, EnableLog=%v, Timeout=%d\n",
			testConfigRetrieved.AppName, testConfigRetrieved.Version, testConfigRetrieved.Environment,
			testConfigRetrieved.EnableLog, testConfigRetrieved.Timeout)
	} else {
		fmt.Printf("❌ 类型断言失败，实际类型: %T\n", configStruct)
	}

	fmt.Println("\n=== 改进效果对比 ===")
	fmt.Println("旧格式: {\"type\":\"json\",\"data\":\"{\\\"app_name\\\":\\\"应用\\\",...}\"}")
	fmt.Println("新格式: {\"type\":\"json\",\"data\":{\"app_name\":\"应用\",...}}")
	fmt.Println("优势:")
	fmt.Println("- 避免双重序列化")
	fmt.Println("- 减少存储空间")
	fmt.Println("- 提高解析效率")
	fmt.Println("- 更好的可读性")
} 