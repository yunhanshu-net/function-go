package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/yunhanshu-net/function-go/runner"
	"github.com/yunhanshu-net/function-go/pkg/dto/syscallback"
)

// TestConfig 测试配置结构体
type TestConfig struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Enabled bool   `json:"enabled"`
	Timeout int64  `json:"timeout"`
}

func main() {
	fmt.Println("=== 测试配置缓存设计 ===")

	// 创建上下文
	ctx := &runner.Context{}

	// 获取配置管理器
	configManager := runner.GetConfigManager()

	// 注册配置结构体
	configKey := "function.test.config.get"
	configManager.RegisterConfigStruct(configKey, &TestConfig{})

	fmt.Println("1. 测试首次获取配置（应该解析并缓存）")
	
	// 模拟配置数据
	configData := &syscallback.ConfigData{
		Type: "json",
		Data: `{"name":"张三","age":25,"enabled":true,"timeout":30000}`,
	}

	// 直接设置缓存（模拟从存储加载）
	configManager.UpdateConfig(ctx, configKey, configData)

	// 获取配置结构体
	config1 := configManager.GetConfigStruct(ctx, configKey)
	if config1 == nil {
		fmt.Println("❌ 首次获取配置失败")
		return
	}

	// 断言类型
	if testConfig, ok := config1.(*TestConfig); ok {
		fmt.Printf("✅ 首次获取成功: Name=%s, Age=%d, Enabled=%v, Timeout=%d\n", 
			testConfig.Name, testConfig.Age, testConfig.Enabled, testConfig.Timeout)
	} else {
		fmt.Println("❌ 类型断言失败")
		return
	}

	fmt.Println("\n2. 测试重复获取配置（应该直接返回缓存的指针）")
	
	// 再次获取，应该直接返回缓存的指针
	config2 := configManager.GetConfigStruct(ctx, configKey)
	if config2 == nil {
		fmt.Println("❌ 重复获取配置失败")
		return
	}

	// 验证是否是同一个指针
	if config1 == config2 {
		fmt.Println("✅ 重复获取返回相同指针（缓存生效）")
	} else {
		fmt.Println("❌ 重复获取返回不同指针（缓存未生效）")
	}

	fmt.Println("\n3. 测试配置无感知更新")
	
	// 更新配置
	newConfigData := &syscallback.ConfigData{
		Type: "json",
		Data: `{"name":"李四","age":30,"enabled":false,"timeout":60000}`,
	}

	err := configManager.UpdateConfig(ctx, configKey, newConfigData)
	if err != nil {
		fmt.Printf("❌ 更新配置失败: %v\n", err)
		return
	}

	// 获取更新后的配置（应该还是同一个指针）
	config3 := configManager.GetConfigStruct(ctx, configKey)
	if config3 == nil {
		fmt.Println("❌ 获取更新后配置失败")
		return
	}

	// 验证指针是否相同
	if config1 == config3 {
		fmt.Println("✅ 配置更新后指针相同（无感知更新生效）")
	} else {
		fmt.Println("❌ 配置更新后指针不同（无感知更新未生效）")
	}

	// 验证数据是否正确更新
	if testConfig, ok := config3.(*TestConfig); ok {
		fmt.Printf("✅ 数据已更新: Name=%s, Age=%d, Enabled=%v, Timeout=%d\n", 
			testConfig.Name, testConfig.Age, testConfig.Enabled, testConfig.Timeout)
	} else {
		fmt.Println("❌ 更新后类型断言失败")
	}

	fmt.Println("\n4. 测试并发安全性")
	
	// 模拟并发访问
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() { done <- true }()
			
			config := configManager.GetConfigStruct(ctx, configKey)
			if config != nil {
				if testConfig, ok := config.(*TestConfig); ok {
					fmt.Printf("  协程 %d: Name=%s\n", id, testConfig.Name)
				}
			}
		}(i)
	}

	// 等待所有协程完成
	for i := 0; i < 10; i++ {
		<-done
	}

	fmt.Println("✅ 并发访问测试完成")

	fmt.Println("\n5. 测试缓存统计")
	
	cacheSize := configManager.GetCacheSize()
	fmt.Printf("缓存大小: %d\n", cacheSize)

	// 清空缓存
	configManager.ClearCache()
	cacheSize = configManager.GetCacheSize()
	fmt.Printf("清空后缓存大小: %d\n", cacheSize)

	fmt.Println("\n=== 测试完成 ===")
} 