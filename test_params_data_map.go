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
	MaxRetries  int    `json:"max_retries"`
}

// TestFunction 测试函数
func TestFunction(ctx *runner.Context, req interface{}, resp interface{}) error {
	// 使用 GetConfig() 方法获取配置
	config := ctx.GetConfig()
	if config != nil {
		if testConfig, ok := config.(*TestConfig); ok {
			fmt.Printf("获取到配置: AppName=%s, Version=%s, Environment=%s, EnableLog=%v, Timeout=%d, MaxRetries=%d\n", 
				testConfig.AppName, testConfig.Version, testConfig.Environment, 
				testConfig.EnableLog, testConfig.Timeout, testConfig.MaxRetries)
		} else {
			fmt.Println("配置类型断言失败")
		}
	} else {
		fmt.Println("未获取到配置")
	}
	return nil
}

// 函数选项
var TestFunctionOption = &runner.FunctionOptions{
	Tags: []string{"测试", "ParamsData"},
	AutoUpdateConfig: &runner.AutoUpdateConfig{
		ConfigStruct: &TestConfig{
			AppName:     "测试应用",
			Version:     "1.0.0",
			Environment: "开发环境",
			EnableLog:   true,
			Timeout:     30,
			MaxRetries:  3,
		},
	},
	EnglishName:  "test_params_data",
	ChineseName:  "ParamsData测试",
	ApiDesc:      "测试ParamsData的map类型转换",
	FunctionType: runner.FunctionTypeDynamic,
}

func testParamsDataMap() {
	fmt.Println("=== 测试ParamsData的map类型转换 ===")

	// 初始化配置管理器
	configManager := runner.GetConfigManager()
	localStorage := runner.NewLocalFileStorage("./test_configs")
	configManager.SetStorage(localStorage)

	// 创建测试上下文
	ctx := runner.NewContext(nil, "POST", "/test/paramsdata")

	// 测试配置结构体
	testConfig := &TestConfig{
		AppName:     "测试应用",
		Version:     "1.0.0",
		Environment: "开发环境",
		EnableLog:   true,
		Timeout:     30,
		MaxRetries:  3,
	}

	fmt.Println("1. 测试structToMap函数")
	
	// 测试structToMap函数
	configMap, err := structToMap(testConfig)
	if err != nil {
		fmt.Printf("❌ structToMap失败: %v\n", err)
		return
	}

	// 输出转换结果
	configJson, _ := json.MarshalIndent(configMap, "", "  ")
	fmt.Printf("✅ 结构体转换为map成功:\n%s\n", string(configJson))

	// 验证map类型
	fmt.Printf("✅ map类型验证: %T\n", configMap)

	// 验证map内容
	fmt.Println("✅ map内容验证:")
	for key, value := range configMap {
		fmt.Printf("  %s: %v (%T)\n", key, value, value)
	}

	fmt.Println("\n2. 测试配置键生成")
	configKey := ctx.generateConfigKey()
	fmt.Printf("生成的配置键: %s\n", configKey)

	fmt.Println("\n3. 测试配置注册和获取")
	
	// 注册配置结构体
	configManager.RegisterConfigStruct(configKey, testConfig)
	
	// 创建配置数据
	configData := &syscallback.ConfigData{
		Type: "json",
		Data: testConfig,
	}
	
	// 写入配置
	err = configManager.UpdateConfig(ctx, configKey, configData)
	if err != nil {
		fmt.Printf("❌ 写入配置失败: %v\n", err)
		return
	}
	
	// 获取配置
	retrievedConfig := configManager.GetConfigStruct(ctx, configKey)
	if retrievedConfig != nil {
		if testConfigRetrieved, ok := retrievedConfig.(TestConfig); ok {
			fmt.Printf("✅ 成功获取配置: AppName=%s, Version=%s\n", 
				testConfigRetrieved.AppName, testConfigRetrieved.Version)
		} else {
			fmt.Printf("❌ 类型断言失败，实际类型: %T\n", retrievedConfig)
		}
	} else {
		fmt.Println("❌ 获取配置失败")
	}
}

// structToMap 将结构体转换为map[string]interface{}
func structToMap(obj interface{}) (map[string]interface{}, error) {
	// 先序列化为JSON
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("序列化结构体失败: %w", err)
	}
	
	// 再反序列化为map
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("反序列化为map失败: %w", err)
	}
	
	return result, nil
}

func init() {
	// 注册测试函数
	runner.Post("/test/paramsdata", TestFunction, TestFunctionOption)

	// 运行测试
	testParamsDataMap()

	fmt.Println("ParamsData map类型转换测试完成")
} 