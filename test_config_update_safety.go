package main

import (
	"encoding/json"
	"fmt"
	"github.com/yunhanshu-net/function-go/pkg/dto/syscallback"
	"github.com/yunhanshu-net/function-go/runner"
)

// 测试配置更新安全性
func testConfigUpdateSafety() {
	fmt.Println("配置更新安全性测试:")
	fmt.Println("=====================")

	// 初始化配置管理器
	configManager := runner.GetConfigManager()
	localStorage := runner.NewLocalFileStorage("./test_configs")
	configManager.SetStorage(localStorage)

	// 创建测试上下文
	ctx := runner.NewContext(nil, "POST", "/test/safety")

	// 测试配置
	testConfig := map[string]interface{}{
		"value": "original-value",
		"count": 100,
	}

	// 序列化配置
	configData, _ := json.Marshal(testConfig)

	// 创建配置数据
	configDataStruct := &syscallback.ConfigData{
		Type: "json",
		Data: string(configData),
	}

	configKey := "function.test.safety.POST"

	// 1. 测试初始更新
	fmt.Println("1. 初始配置更新")
	err := configManager.UpdateConfig(ctx, configKey, configDataStruct)
	if err != nil {
		fmt.Printf("初始更新失败: %v\n", err)
		return
	}
	fmt.Println("初始更新成功")

	// 2. 验证配置已缓存
	fmt.Println("\n2. 验证配置已缓存")
	cachedConfig := configManager.GetByKey(ctx, configKey)
	if cachedConfig != nil {
		fmt.Printf("缓存配置: %s\n", cachedConfig.Data)
	} else {
		fmt.Println("配置未找到")
		return
	}

	// 3. 测试修改原始指针不会影响缓存
	fmt.Println("\n3. 测试修改原始指针")
	configDataStruct.Data = "modified-value"
	
	// 重新获取缓存中的配置
	cachedConfig2 := configManager.GetByKey(ctx, configKey)
	if cachedConfig2 != nil {
		fmt.Printf("缓存配置（修改后）: %s\n", cachedConfig2.Data)
		if cachedConfig2.Data != configDataStruct.Data {
			fmt.Println("✓ 原始指针修改不影响缓存")
		} else {
			fmt.Println("✗ 原始指针修改影响了缓存")
		}
	}

	// 4. 测试更新配置
	fmt.Println("\n4. 测试更新配置")
	newConfig := &syscallback.ConfigData{
		Type: "json",
		Data: "{\"value\":\"updated-value\",\"count\":200}",
	}
	
	err = configManager.UpdateConfig(ctx, configKey, newConfig)
	if err != nil {
		fmt.Printf("更新失败: %v\n", err)
		return
	}
	fmt.Println("更新成功")

	// 5. 验证更新后的配置
	fmt.Println("\n5. 验证更新后的配置")
	cachedConfig3 := configManager.GetByKey(ctx, configKey)
	if cachedConfig3 != nil {
		fmt.Printf("更新后缓存配置: %s\n", cachedConfig3.Data)
	} else {
		fmt.Println("更新后配置未找到")
	}

	// 6. 测试修改新配置指针
	fmt.Println("\n6. 测试修改新配置指针")
	newConfig.Data = "hacked-value"
	
	cachedConfig4 := configManager.GetByKey(ctx, configKey)
	if cachedConfig4 != nil {
		fmt.Printf("缓存配置（修改新指针后）: %s\n", cachedConfig4.Data)
		if cachedConfig4.Data != newConfig.Data {
			fmt.Println("✓ 新指针修改不影响缓存")
		} else {
			fmt.Println("✗ 新指针修改影响了缓存")
		}
	}
}

func init() {
	testConfigUpdateSafety()
} 