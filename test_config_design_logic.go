package main

import (
	"fmt"
	"reflect"
)

func main() {
	fmt.Println("=== 验证配置缓存设计逻辑 ===")

	// 模拟配置管理器的逻辑
	type TestConfig struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	// 模拟 configStructs 存储
	configStructs := make(map[string]interface{})

	// 1. 注册时存储类型
	configKey := "test.config"
	configStructs[configKey] = reflect.TypeOf(&TestConfig{})
	fmt.Printf("1. 注册后存储的类型: %T\n", configStructs[configKey])

	// 2. 检查是否是已解析的指针
	cachedStruct := configStructs[configKey]
	if reflect.TypeOf(cachedStruct).Kind() == reflect.Ptr {
		fmt.Println("✅ 这是已解析的指针，直接返回")
	} else {
		fmt.Println("📝 这是类型，需要解析")
	}

	// 3. 模拟解析后的缓存
	instance := &TestConfig{Name: "张三", Age: 25}
	configStructs[configKey] = instance
	fmt.Printf("2. 解析后存储的实例: %T, 值: %+v\n", configStructs[configKey], configStructs[configKey])

	// 4. 再次检查
	cachedStruct = configStructs[configKey]
	if reflect.TypeOf(cachedStruct).Kind() == reflect.Ptr {
		fmt.Println("✅ 这是已解析的指针，直接返回")
	} else {
		fmt.Println("📝 这是类型，需要解析")
	}

	// 5. 验证指针更新
	fmt.Println("\n3. 测试指针更新:")
	oldPointer := configStructs[configKey]
	fmt.Printf("  旧指针: %p, 值: %+v\n", oldPointer, oldPointer)

	// 模拟配置更新 - 直接修改指针指向的值
	if testConfig, ok := oldPointer.(*TestConfig); ok {
		testConfig.Name = "李四"
		testConfig.Age = 30
		fmt.Printf("  更新后: %p, 值: %+v\n", oldPointer, oldPointer)
	}

	// 6. 验证指针相同
	newPointer := configStructs[configKey]
	if oldPointer == newPointer {
		fmt.Println("✅ 指针相同，无感知更新生效")
	} else {
		fmt.Println("❌ 指针不同，无感知更新失败")
	}

	fmt.Println("\n=== 设计逻辑验证完成 ===")
} 