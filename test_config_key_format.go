package main

import (
	"fmt"
	"github.com/yunhanshu-net/function-go/pkg/dto/syscallback"
)

// 测试配置键格式
func testConfigKeyFormat() {
	fmt.Println("配置键格式测试:")
	fmt.Println("==================")

	testCases := []struct {
		router string
		method string
	}{
		{"/api/users", "GET"},
		{"/api/users", "POST"},
		{"/api/users", "PUT"},
		{"/api/users", "DELETE"},
		{"/widgets/add", "POST"},
		{"/widgets/calculator", "GET"},
		{"/admin/settings", "POST"},
		{"/", "GET"},
		{"/api/v1/users", "GET"},
		{"/api/v1/users/profile", "PUT"},
		// 边界情况测试
		{"api/users", "GET"},      // 没有前导 /
		{"/api/users/", "POST"},   // 有尾随 /
		{"api/users/", "PUT"},     // 没有前导但有尾随 /
		{"", "GET"},               // 空路由
		{"/", "POST"},             // 只有 /
		{"/api//users", "GET"},    // 双斜杠
	}

	for _, tc := range testCases {
		// 测试 ConfigUpdateRequest
		updateReq := &syscallback.ConfigUpdateRequest{
			Router: tc.router,
			Method: tc.method,
		}
		updateKey := updateReq.GenerateConfigKey()
		
		// 测试 ConfigGetRequest
		getReq := &syscallback.ConfigGetRequest{
			Router: tc.router,
			Method: tc.method,
		}
		getKey := getReq.GenerateConfigKey()
		
		fmt.Printf("路由: %s, 方法: %s\n", tc.router, tc.method)
		fmt.Printf("更新配置键: %s\n", updateKey)
		fmt.Printf("获取配置键: %s\n", getKey)
		fmt.Println("---")
	}
}

func init() {
	testConfigKeyFormat()
} 