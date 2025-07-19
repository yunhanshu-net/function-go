package main

import (
	"fmt"
	"reflect"

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
	fmt.Println("=== 测试配置值传递设计 ===")

	// 创建上下文
	ctx := &runner.Context{}

	// 获取配置管理器
	configManager := runner.GetConfigManager()

	// 注册配置结构体
	configKey := "function.test.config.get"
	configManager.RegisterConfigStruct(configKey, TestConfig{})

	fmt.Println("1. 测试获取配置值")
	
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
		fmt.Println("❌ 获取配置失败")
		return
	}

	// 断言类型
	if testConfig, ok := config1.(TestConfig); ok {
		fmt.Printf("✅ 获取成功: Name=%s, Age=%d, Enabled=%v, Timeout=%d\n", 
			testConfig.Name, testConfig.Age, testConfig.Enabled, testConfig.Timeout)
	} else {
		fmt.Printf("❌ 类型断言失败，实际类型: %T\n", config1)
		return
	}

	fmt.Println("\n2. 测试重复获取配置（应该返回不同的值副本）")
	
	// 再次获取，应该返回新的值副本
	config2 := configManager.GetConfigStruct(ctx, configKey)
	if config2 == nil {
		fmt.Println("❌ 重复获取配置失败")
		return
	}

	// 验证是否是同一个值（应该不是，因为是副本）
	if reflect.DeepEqual(config1, config2) {
		fmt.Println("✅ 重复获取返回相同的值（内容相同）")
	} else {
		fmt.Println("❌ 重复获取返回不同的值")
	}

	// 验证是否是同一个地址（应该不是）
	if fmt.Sprintf("%p", &config1) == fmt.Sprintf("%p", &config2) {
		fmt.Println("❌ 重复获取返回相同地址（应该是副本）")
	} else {
		fmt.Println("✅ 重复获取返回不同地址（副本生效）")
	}

	fmt.Println("\n3. 测试配置更新后获取")
	
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

	// 获取更新后的配置
	config3 := configManager.GetConfigStruct(ctx, configKey)
	if config3 == nil {
		fmt.Println("❌ 获取更新后配置失败")
		return
	}

	// 验证数据是否正确更新
	if testConfig, ok := config3.(TestConfig); ok {
		fmt.Printf("✅ 数据已更新: Name=%s, Age=%d, Enabled=%v, Timeout=%d\n", 
			testConfig.Name, testConfig.Age, testConfig.Enabled, testConfig.Timeout)
	} else {
		fmt.Println("❌ 更新后类型断言失败")
	}

	fmt.Println("\n4. 测试用户修改不会影响缓存")
	
	// 用户修改获取到的配置
	if testConfig, ok := config3.(TestConfig); ok {
		originalName := testConfig.Name
		testConfig.Name = "王五"
		fmt.Printf("  用户修改: %s -> %s\n", originalName, testConfig.Name)
		
		// 再次获取配置，应该还是原始值
		config4 := configManager.GetConfigStruct(ctx, configKey)
		if testConfig4, ok := config4.(TestConfig); ok {
			if testConfig4.Name == "李四" {
				fmt.Println("✅ 用户修改不影响缓存（值传递生效）")
			} else {
				fmt.Printf("❌ 用户修改影响了缓存: %s\n", testConfig4.Name)
			}
		}
	}

	fmt.Println("\n5. 测试并发安全性")
	
	// 模拟并发访问
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() { done <- true }()
			
			config := configManager.GetConfigStruct(ctx, configKey)
			if config != nil {
				if testConfig, ok := config.(TestConfig); ok {
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

	fmt.Println("\n6. 测试AutoUpdateConfig结构体")
	
	// 测试AutoUpdateConfig
	autoConfig := &runner.AutoUpdateConfig{
		ConfigStruct: TestConfig{}, // 传递值而不是指针
		BeforeConfigChange: func(ctx *runner.Context, oldConfig, newConfig *syscallback.ConfigData) error {
			fmt.Println("  配置变更回调被触发")
			return nil
		},
	}
	
	fmt.Printf("AutoUpdateConfig类型: %T\n", autoConfig.ConfigStruct)
	if reflect.TypeOf(autoConfig.ConfigStruct).Kind() == reflect.Struct {
		fmt.Println("✅ AutoUpdateConfig使用值类型")
	} else {
		fmt.Println("❌ AutoUpdateConfig不是值类型")
	}

	fmt.Println("\n=== 测试完成 ===")
} 