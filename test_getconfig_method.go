package main

import (
	"fmt"
	"github.com/yunhanshu-net/function-go/runner"
)

// TestConfig 测试配置结构体
type TestConfig struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// TestFunction 测试函数
func TestFunction(ctx *runner.Context, req interface{}, resp interface{}) error {
	// 使用 GetConfig() 方法获取配置
	config := ctx.GetConfig()
	if config != nil {
		if testConfig, ok := config.(*TestConfig); ok {
			fmt.Printf("获取到配置: Name=%s, Value=%d\n", testConfig.Name, testConfig.Value)
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
	Tags: []string{"测试", "GetConfig"},
	AutoUpdateConfig: &runner.AutoUpdateConfig{
		ConfigStruct: &TestConfig{Name: "default", Value: 0},
	},
	EnglishName:  "test_getconfig",
	ChineseName:  "GetConfig测试",
	ApiDesc:      "测试GetConfig方法",
	FunctionType: runner.FunctionTypeDynamic,
}

func testGetConfigMethod() {
	fmt.Println("GetConfig方法测试:")
	fmt.Println("==================")

	// 初始化配置管理器
	configManager := runner.GetConfigManager()
	localStorage := runner.NewLocalFileStorage("./test_configs")
	configManager.SetStorage(localStorage)

	// 创建测试上下文
	ctx := runner.NewContext(nil, "GET", "/test/getconfig")

	// 测试 GetConfig() 方法
	fmt.Println("1. 测试 GetConfig() 方法")
	config := ctx.GetConfig()
	if config != nil {
		fmt.Println("✓ GetConfig() 方法正常工作")
	} else {
		fmt.Println("✗ GetConfig() 方法返回nil")
	}

	// 测试配置键生成
	fmt.Println("\n2. 测试配置键生成")
	configKey := ctx.generateConfigKey()
	fmt.Printf("生成的配置键: %s\n", configKey)
}

func init() {
	// 注册测试函数
	runner.Get("/test/getconfig", TestFunction, TestFunctionOption)

	// 运行测试
	testGetConfigMethod()

	fmt.Println("GetConfig方法测试完成")
} 