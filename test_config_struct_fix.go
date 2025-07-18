package main

import (
	"encoding/json"
	"fmt"
	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/function-go/runner"
)

// TestConfig 测试配置
type TestConfig struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

// TestReq 测试请求
type TestReq struct {
	Input string `json:"input" form:"input" runner:"code:input;name:输入" widget:"type:input;placeholder:请输入内容" data:"type:string;default_value:hello" validate:"required"`
}

// TestResp 测试响应
type TestResp struct {
	Output string `json:"output" runner:"code:output;name:输出" data:"type:string;example:hello world"`
}

// TestFunction 测试函数
func TestFunction(ctx *runner.Context, req *TestReq, resp response.Response) error {
	configInterface := ctx.GetConfig()
	if configInterface != nil {
		if config, ok := configInterface.(*TestConfig); ok {
			return resp.Form(&TestResp{Output: fmt.Sprintf("配置值: %s, 计数: %d, 输入: %s", config.Value, config.Count, req.Input)}).Build()
		}
	}
	return resp.Form(&TestResp{Output: fmt.Sprintf("无配置, 输入: %s", req.Input)}).Build()
}

// 配置变更回调
func beforeConfigChange(ctx *runner.Context, oldConfig, newConfig interface{}) error {
	fmt.Printf("配置变更回调触发，路由: %s, 方法: %s\n", ctx.router, ctx.method)
	return nil
}

// 函数选项
var TestFunctionOption = &runner.FunctionOptions{
	Tags: []string{"测试", "配置结构体"},
	AutoUpdateConfig: &runner.AutoUpdateConfig{
		ConfigStruct:      &TestConfig{Value: "default", Count: 0},
		BeforeConfigChange: beforeConfigChange,
	},
	EnglishName:  "test_config_struct",
	ChineseName:  "配置结构体测试",
	ApiDesc:      "测试配置结构体注册修复",
	Request:      &TestReq{},
	Response:     &TestResp{},
	FunctionType: runner.FunctionTypeDynamic,
}

// 测试配置结构体注册
func testConfigStructRegistration() {
	fmt.Println("配置结构体注册测试:")
	fmt.Println("====================")

	// 初始化配置管理器
	configManager := runner.GetConfigManager()
	localStorage := runner.NewLocalFileStorage("./test_configs")
	configManager.SetStorage(localStorage)

	// 创建测试上下文
	ctx := runner.NewContext(nil, "POST", "/test/struct")

	configKey := "function.test.struct.POST"

	// 1. 测试注册配置结构体
	fmt.Println("1. 注册配置结构体")
	configManager.RegisterConfigStruct(configKey, &TestConfig{})
	fmt.Println("✓ 配置结构体注册成功")

	// 2. 测试配置更新
	fmt.Println("\n2. 更新配置")
	testConfig := &TestConfig{
		Value: "test-value",
		Count: 100,
	}

	configData, _ := json.Marshal(testConfig)
	configDataStruct := &syscallback.ConfigData{
		Type: "json",
		Data: string(configData),
	}

	err := configManager.UpdateConfig(ctx, configKey, configDataStruct)
	if err != nil {
		fmt.Printf("配置更新失败: %v\n", err)
		return
	}
	fmt.Println("✓ 配置更新成功")

	// 3. 测试获取配置结构体
	fmt.Println("\n3. 获取配置结构体")
	configInterface := configManager.GetConfigStruct(ctx, configKey)
	if configInterface != nil {
		if config, ok := configInterface.(*TestConfig); ok {
			fmt.Printf("✓ 成功获取配置结构体: Value=%s, Count=%d\n", config.Value, config.Count)
		} else {
			fmt.Println("✗ 配置结构体类型断言失败")
		}
	} else {
		fmt.Println("✗ 获取配置结构体失败")
	}

	// 4. 测试从上下文获取配置
	fmt.Println("\n4. 从上下文获取配置")
	ctxConfig := ctx.GetConfig()
	if ctxConfig != nil {
		if config, ok := ctxConfig.(*TestConfig); ok {
			fmt.Printf("✓ 从上下文获取配置成功: Value=%s, Count=%d\n", config.Value, config.Count)
		} else {
			fmt.Println("✗ 从上下文获取配置类型断言失败")
		}
	} else {
		fmt.Println("✗ 从上下文获取配置失败")
	}
}

func init() {
	// 注册测试函数
	runner.Post("/test/struct", TestFunction, TestFunctionOption)

	// 运行配置结构体测试
	testConfigStructRegistration()

	fmt.Println("配置结构体注册测试完成")
} 