package main

import (
	"fmt"
	"github.com/yunhanshu-net/function-go/runner"
)

// TestFunction 测试函数
func TestFunction(ctx *runner.Context, req interface{}, resp interface{}) error {
	fmt.Println("=== Context 方法测试 ===")
	
	// 测试基础信息方法
	fmt.Printf("User: %s\n", ctx.User())
	fmt.Printf("Name: %s\n", ctx.Name())
	fmt.Printf("Version: %s\n", ctx.Version())
	fmt.Printf("Router: %s\n", ctx.Router())
	fmt.Printf("Method: %s\n", ctx.Method())
	fmt.Printf("Now: %s\n", ctx.Now().Format("2006-01-02 15:04:05"))
	
	// 测试配置方法
	fmt.Printf("Config Manager: %v\n", ctx.Config())
	config := ctx.GetConfig()
	if config != nil {
		fmt.Printf("Config: %v\n", config)
	} else {
		fmt.Println("Config: nil")
	}
	
	// 测试数据库方法
	db := ctx.DB()
	fmt.Printf("DB: %v\n", db)
	
	// 测试文件方法
	files := ctx.NewFiles([]string{})
	fmt.Printf("Files: %v\n", files)
	
	// 测试上传路径
	uploadPath := ctx.GetUploadPath()
	fmt.Printf("Upload Path: %s\n", uploadPath)
	
	// 测试函数消息
	functionMsg := ctx.GetFunctionMsg()
	if functionMsg != nil {
		fmt.Printf("Function Msg: User=%s, Runner=%s, Version=%s\n", 
			functionMsg.User, functionMsg.Runner, functionMsg.Version)
	}
	
	// 测试配置键生成
	configKey := ctx.generateConfigKey()
	fmt.Printf("Config Key: %s\n", configKey)
	
	return nil
}

// 函数选项
var TestFunctionOption = &runner.FunctionOptions{
	Tags: []string{"测试", "Context方法"},
	EnglishName:  "test_context_methods",
	ChineseName:  "Context方法测试",
	ApiDesc:      "测试所有Context方法",
	FunctionType: runner.FunctionTypeDynamic,
}

func testContextMethods() {
	fmt.Println("Context 方法完整性测试:")
	fmt.Println("========================")

	// 创建测试上下文
	ctx := runner.NewContext(nil, "POST", "/test/context/methods")

	// 测试所有方法
	fmt.Println("1. 基础信息方法")
	fmt.Printf("   User: %s\n", ctx.User())
	fmt.Printf("   Name: %s\n", ctx.Name())
	fmt.Printf("   Version: %s\n", ctx.Version())
	fmt.Printf("   Router: %s\n", ctx.Router())
	fmt.Printf("   Method: %s\n", ctx.Method())
	fmt.Printf("   Now: %s\n", ctx.Now().Format("2006-01-02 15:04:05"))

	fmt.Println("\n2. 配置方法")
	fmt.Printf("   Config Manager: %v\n", ctx.Config())
	fmt.Printf("   GetConfig: %v\n", ctx.GetConfig())
	fmt.Printf("   Config Key: %s\n", ctx.generateConfigKey())

	fmt.Println("\n3. 数据库方法")
	fmt.Printf("   DB: %v\n", ctx.DB())
	fmt.Printf("   GetDB: %v\n", ctx.GetDB())

	fmt.Println("\n4. 文件方法")
	fmt.Printf("   NewFiles: %v\n", ctx.NewFiles([]string{}))
	fmt.Printf("   Upload Path: %s\n", ctx.GetUploadPath())

	fmt.Println("\n5. 函数消息")
	functionMsg := ctx.GetFunctionMsg()
	if functionMsg != nil {
		fmt.Printf("   Function Msg: User=%s, Runner=%s, Version=%s\n", 
			functionMsg.User, functionMsg.Runner, functionMsg.Version)
	}

	fmt.Println("\n✓ 所有Context方法测试完成")
}

func init() {
	// 注册测试函数
	runner.Post("/test/context/methods", TestFunction, TestFunctionOption)

	// 运行测试
	testContextMethods()
} 