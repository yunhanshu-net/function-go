package main

import (
	"fmt"
	"strings"

	"github.com/yunhanshu-net/function-go/pkg/dto/usercall"
)

// generateConfigKey 生成配置键（复制自 default.go）
func generateConfigKey(router, method string) string {
	// 将路由中的路径分隔符替换为点号
	routerKey := strings.ReplaceAll(strings.Trim(router, "/"), "/", ".")
	// 去除前后多余的点号
	routerKey = strings.Trim(routerKey, ".")

	// 生成配置键格式: function.{router}.{method}
	return fmt.Sprintf("function.%s.%s", routerKey, strings.ToLower(method))
}

func main() {
	// 测试用例
	testCases := []struct {
		router string
		method string
	}{
		{"/widgets/add", "POST"},
		{"/widgets/product_list", "GET"},
		{"/cmp/config_demo", "POST"},
		{"/api/test", "PUT"},
		{"/", "GET"},
	}

	fmt.Println("=== 配置键生成一致性测试 ===")
	
	for i, tc := range testCases {
		fmt.Printf("\n测试用例 %d:\n", i+1)
		fmt.Printf("  路由: %s\n", tc.router)
		fmt.Printf("  方法: %s\n", tc.method)
		
		// 使用 default.go 中的函数
		key1 := generateConfigKey(tc.router, tc.method)
		fmt.Printf("  default.go 生成: %s\n", key1)
		
		// 使用 usercall 中的方法
		updateReq := &usercall.UpdateConfigReq{
			Router: tc.router,
			Method: tc.method,
		}
		key2 := updateReq.GenerateConfigKey()
		fmt.Printf("  usercall 生成:  %s\n", key2)
		
		// 检查一致性
		if key1 == key2 {
			fmt.Printf("  ✅ 一致\n")
		} else {
			fmt.Printf("  ❌ 不一致\n")
		}
	}
	
	fmt.Println("\n=== 测试完成 ===")
} 