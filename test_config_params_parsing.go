package main

import (
	"encoding/json"
	"fmt"

	"github.com/yunhanshu-net/function-go/runner"
	"github.com/yunhanshu-net/function-go/pkg/dto/api"
)

// TestConfig 测试配置结构体
type TestConfig struct {
	ApiKey     string `json:"api_key" runner:"code:api_key;name:API密钥;desc:用于访问第三方服务的密钥" default_value:"sk-123456789"`
	Timeout    int    `json:"timeout" runner:"code:timeout;name:超时时间;desc:请求超时时间（秒）" default_value:"30"`
	MaxRetries int    `json:"max_retries" runner:"code:max_retries;name:最大重试次数;desc:请求失败时的最大重试次数" default_value:"3"`
	Debug      bool   `json:"debug" runner:"code:debug;name:调试模式;desc:是否启用调试模式" default_value:"false"`
}

func main() {
	fmt.Println("=== 测试ParamsConfig解析 ===")

	// 创建带初始值的配置结构体
	initialConfig := TestConfig{
		ApiKey:     "sk-123456789",
		Timeout:    30,
		MaxRetries: 3,
		Debug:      false,
	}

	fmt.Println("1. 测试直接解析配置结构体")
	
	// 直接使用api包解析配置结构体
	configParams, err := api.NewRequestParamsWithFunctionInfo(initialConfig, "form", nil)
	if err != nil {
		fmt.Printf("❌ 解析配置结构体失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 配置结构体解析成功，类型: %T\n", configParams)

	// 输出解析结果
	configJson, _ := json.MarshalIndent(configParams, "", "  ")
	fmt.Printf("解析结果:\n%s\n", string(configJson))

	fmt.Println("\n2. 测试在buildApiInfo中的解析")
	
	// 模拟FunctionOptions
	functionOptions := &runner.FunctionOptions{
		RenderType: "form",
		AutoUpdateConfig: &runner.AutoUpdateConfig{
			ConfigStruct: initialConfig,
		},
	}

	// 模拟routerInfo
	worker := &runner.routerInfo{
		Router:       "/api/test",
		Method:       "POST",
		FunctionInfo: functionOptions,
	}

	// 模拟Runner
	r := &runner.Runner{}

	// 构建API信息
	apiInfo, err := r.buildApiInfo(worker)
	if err != nil {
		fmt.Printf("❌ 构建API信息失败: %v\n", err)
		return
	}

	fmt.Printf("✅ API信息构建成功:\n")
	fmt.Printf("  Router: %s\n", apiInfo.Router)
	fmt.Printf("  Method: %s\n", apiInfo.Method)
	fmt.Printf("  ParamsConfig类型: %T\n", apiInfo.ParamsConfig)

	// 验证ParamsConfig是否正确解析
	if apiInfo.ParamsConfig != nil {
		paramsJson, _ := json.MarshalIndent(apiInfo.ParamsConfig, "", "  ")
		fmt.Printf("  ParamsConfig内容:\n%s\n", string(paramsJson))
	} else {
		fmt.Println("  ParamsConfig为空")
	}

	fmt.Println("\n3. 测试配置结构体标签解析")
	
	// 验证runner标签是否正确解析
	if formConfig, ok := configParams.(*api.FormConfig); ok {
		fmt.Println("字段解析结果:")
		for _, field := range formConfig.Fields {
			fmt.Printf("  - %s (%s): %s\n", field.Code, field.Name, field.Desc)
		}
	}

	fmt.Println("\n=== 测试完成 ===")
} 