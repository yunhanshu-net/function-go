package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yunhanshu-net/function-go/runner"
	"github.com/yunhanshu-net/function-go/pkg/dto/syscallback"
)

// TestConfig 测试配置结构体
type TestConfig struct {
	ApiKey    string `json:"api_key" default_value:"sk-123456789"`
	Timeout   int    `json:"timeout" default_value:"30"`
	MaxRetries int   `json:"max_retries" default_value:"3"`
	Debug     bool   `json:"debug" default_value:"false"`
}

func main() {
	fmt.Println("=== 测试配置初始化设计 ===")

	// 创建上下文
	ctx := &runner.Context{}

	// 获取配置管理器
	configManager := runner.GetConfigManager()

	// 模拟路由信息
	router := "/api/test"
	method := "POST"
	configKey := generateConfigKey(router, method)

	fmt.Println("1. 测试配置结构体注册")
	
	// 创建带初始值的配置结构体
	initialConfig := TestConfig{
		ApiKey:     "sk-123456789",
		Timeout:    30,
		MaxRetries: 3,
		Debug:      false,
	}

	// 注册配置结构体
	configManager.RegisterConfigStruct(configKey, initialConfig)
	fmt.Printf("✅ 配置结构体注册成功: %T\n", initialConfig)

	fmt.Println("\n2. 测试初始配置写入")
	
	// 将初始配置序列化为JSON
	configData, err := json.Marshal(initialConfig)
	if err != nil {
		fmt.Printf("❌ 序列化配置失败: %v\n", err)
		return
	}

	// 创建配置数据
	config := &syscallback.ConfigData{
		Type: "json",
		Data: string(configData),
	}

	// 写入配置
	err = configManager.UpdateConfig(ctx, configKey, config)
	if err != nil {
		fmt.Printf("❌ 写入配置失败: %v\n", err)
		return
	}
	fmt.Println("✅ 初始配置写入成功")

	fmt.Println("\n3. 测试配置获取")
	
	// 获取配置
	retrievedConfig := configManager.GetConfigStruct(ctx, configKey)
	if retrievedConfig == nil {
		fmt.Println("❌ 获取配置失败")
		return
	}

	// 类型断言
	if testConfig, ok := retrievedConfig.(TestConfig); ok {
		fmt.Printf("✅ 配置获取成功: ApiKey=%s, Timeout=%d, MaxRetries=%d, Debug=%v\n",
			testConfig.ApiKey, testConfig.Timeout, testConfig.MaxRetries, testConfig.Debug)
	} else {
		fmt.Printf("❌ 类型断言失败，实际类型: %T\n", retrievedConfig)
		return
	}

	fmt.Println("\n4. 测试API信息构建")
	
	// 模拟FunctionOptions
	functionOptions := &runner.FunctionOptions{
		AutoUpdateConfig: &runner.AutoUpdateConfig{
			ConfigStruct: initialConfig,
		},
	}

	// 模拟routerInfo
	worker := &runner.routerInfo{
		Router:       router,
		Method:       method,
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

	// 验证ParamsConfig是否正确设置
	if apiInfo.ParamsConfig != nil {
		if config, ok := apiInfo.ParamsConfig.(TestConfig); ok {
			fmt.Printf("  ParamsConfig值: %+v\n", config)
		} else {
			fmt.Printf("  ParamsConfig类型断言失败: %T\n", apiInfo.ParamsConfig)
		}
	} else {
		fmt.Println("  ParamsConfig为空")
	}

	fmt.Println("\n5. 测试配置更新")
	
	// 更新配置
	newConfigData := &syscallback.ConfigData{
		Type: "json",
		Data: `{"api_key":"sk-new-key","timeout":60,"max_retries":5,"debug":true}`,
	}

	err = configManager.UpdateConfig(ctx, configKey, newConfigData)
	if err != nil {
		fmt.Printf("❌ 更新配置失败: %v\n", err)
		return
	}

	// 获取更新后的配置
	updatedConfig := configManager.GetConfigStruct(ctx, configKey)
	if testConfig, ok := updatedConfig.(TestConfig); ok {
		fmt.Printf("✅ 配置更新成功: ApiKey=%s, Timeout=%d, MaxRetries=%d, Debug=%v\n",
			testConfig.ApiKey, testConfig.Timeout, testConfig.MaxRetries, testConfig.Debug)
	} else {
		fmt.Println("❌ 更新后类型断言失败")
	}

	fmt.Println("\n=== 测试完成 ===")
}

// generateConfigKey 生成配置键（复制自default.go）
func generateConfigKey(router, method string) string {
	// 将路由中的路径分隔符替换为点号
	routerKey := strings.ReplaceAll(strings.Trim(router, "/"), "/", ".")
	// 去除前后多余的点号
	routerKey = strings.Trim(routerKey, ".")
	
	// 生成配置键格式: function.{router}.{method}
	return fmt.Sprintf("function.%s.%s", routerKey, strings.ToLower(method))
} 