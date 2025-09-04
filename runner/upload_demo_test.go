package runner

import (
	"fmt"
	"github.com/yunhanshu-net/pkg/trace"
)

//func TestActualUploadDemo(t *testing.T) {
//	fmt.Println("🚀 实际上传测试")
//	fmt.Println("==================================================")
//
//	// 0. 设置测试环境变量
//	originalUser := env.User
//	originalName := env.Name
//	originalVersion := env.Version
//
//	env.User = "testuser"
//	env.Name = "testapp"
//	env.Version = "v1.0.0"
//
//	// 测试结束后恢复原值
//	defer func() {
//		env.User = originalUser
//		env.Name = originalName
//		env.Version = originalVersion
//	}()
//
//	// 1. 创建context
//	ctx := context.WithValue(context.Background(), "trace_id", "actual-upload-demo-12345")
//
//	// 2. 创建runner context
//	runnerCtx := NewContext(ctx, "POST", "api/actual_upload_test")
//
//	// 3. 显示上传路径信息
//	uploadPath := runnerCtx.GetUploadPath()
//	fmt.Printf("📁 规范上传路径: %s\n", uploadPath)
//
//	// 4. 获取FunctionMsg详细信息
//	functionMsg := runnerCtx.GetFunctionMsg()
//	if functionMsg != nil {
//		fmt.Printf("👤 用户: '%s'\n", functionMsg.User)
//		fmt.Printf("📱 应用: '%s'\n", functionMsg.Runner)
//		fmt.Printf("🔢 版本: '%s'\n", functionMsg.Version)
//		fmt.Printf("🌐 路由: %s\n", functionMsg.Router)
//		fmt.Printf("📋 方法: %s\n", functionMsg.Method)
//		fmt.Printf("🔍 TraceID: %s\n", functionMsg.TraceID)
//		fmt.Printf("☁️  上传提供商: %s\n", functionMsg.UploadConfig.Provider)
//		fmt.Printf("🪣 存储桶: %s\n", functionMsg.UploadConfig.Bucket)
//		fmt.Printf("🔗 下载域名: %s\n", functionMsg.UploadConfig.DownloadDomain)
//	}
//
//	fmt.Println("\n" + "====================================================")
//
//	// 5. 创建测试文件内容
//	testContent := `{
//	"test": "这是一个实际上传测试文件",
//	"timestamp": "2025-06-28T14:40:00Z",
//	"data": {
//		"numbers": [1, 2, 3, 4, 5],
//		"message": "Hello from actual upload test!",
//		"chinese": "中文测试内容",
//		"purpose": "验证上传路径是否符合规范"
//	},
//	"metadata": {
//		"author": "test-user",
//		"test_type": "actual_upload_verification",
//		"expected_path_format": "/user/app/api/actual_upload_test/POST/output/YYYYMMDD/filename"
//	}
//}`
//
//	// 6. 上传文件
//	fmt.Println("📤 开始实际上传文件...")
//	files, err := runnerCtx.CreateFilesFromData("actual_upload_demo.json", []byte(testContent))
//
//	if err != nil {
//		t.Logf("❌ 上传失败: %v", err)
//		return
//	}
//
//	// 7. 显示上传结果
//	if len(files.GetFiles()) > 0 {
//		file := files.GetFiles()[0]
//		fmt.Println("\n✅ 文件上传成功!")
//		fmt.Printf("📄 文件名: %s\n", file.Name)
//		fmt.Printf("📏 文件大小: %d bytes\n", file.Size)
//		fmt.Printf("🏷️  文件类型: %s\n", file.ContentType)
//		fmt.Printf("📅 创建时间: %s\n", file.CreatedAt)
//		fmt.Printf("🔄 更新时间: %s\n", file.UpdatedAt)
//		fmt.Printf("✨ 状态: %s\n", file.Status)
//		fmt.Printf("\n🔗 完整访问URL:\n%s\n", file.URL)
//
//		// 分析URL结构
//		fmt.Println("\n🔍 URL结构分析:")
//		analyzeURL(file.URL)
//
//		// 验证路径是否符合规范
//		fmt.Println("\n📋 路径规范验证:")
//		validateUploadPath(file.URL, functionMsg)
//
//	} else {
//		t.Error("❌ 没有文件被上传")
//	}
//
//	fmt.Println("\n" + "====================================================")
//	fmt.Println("✨ 实际上传测试完成!")
//}

// analyzeURL 分析URL结构
func analyzeURL(url string) {
	fmt.Printf("   完整URL: %s\n", url)

	// 查找域名部分
	protocolEnd := -1
	if idx := findSubstring(url, "://"); idx != -1 {
		protocolEnd = idx + 3
	}

	if protocolEnd != -1 {
		// 查找域名结束位置
		domainEnd := findSubstringFrom(url, "/", protocolEnd)
		if domainEnd != -1 {
			domain := url[:domainEnd]
			path := url[domainEnd:]
			fmt.Printf("   域名部分: %s\n", domain)
			fmt.Printf("   路径部分: %s\n", path)

			// 分析路径组成
			if len(path) > 1 {
				pathParts := splitString(path[1:], "/") // 去掉开头的 /
				fmt.Printf("   路径组成: %v\n", pathParts)

				if len(pathParts) >= 6 {
					fmt.Printf("     - 用户/租户: %s\n", pathParts[0])
					fmt.Printf("     - 应用名称: %s\n", pathParts[1])
					fmt.Printf("     - 函数路由: %s/%s\n", pathParts[2], pathParts[3])
					fmt.Printf("     - HTTP方法: %s\n", pathParts[4])
					fmt.Printf("     - 输出标识: %s\n", pathParts[5])
					if len(pathParts) >= 7 {
						fmt.Printf("     - 日期: %s\n", pathParts[6])
					}
					if len(pathParts) >= 8 {
						fmt.Printf("     - 文件名: %s\n", pathParts[7])
					}
				}
			}
		}
	}
}

// validateUploadPath 验证上传路径是否符合规范
func validateUploadPath(url string, functionMsg *trace.FunctionMsg) {
	expected := functionMsg.GetUploadPath()
	fmt.Printf("   期望路径格式: %s/filename\n", expected)

	// 检查URL是否包含期望的路径结构
	if containsSubstring(url, expected) {
		fmt.Printf("   ✅ 路径格式符合规范\n")
	} else {
		fmt.Printf("   ❌ 路径格式不符合规范\n")
		fmt.Printf("   实际包含: %s\n", extractPathFromURL(url))
	}
}

// 辅助函数
func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func findSubstringFrom(s, substr string, start int) int {
	for i := start; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func splitString(s, sep string) []string {
	if s == "" {
		return []string{}
	}

	var result []string
	start := 0

	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}

	// 添加最后一部分
	if start < len(s) {
		result = append(result, s[start:])
	}

	return result
}

func containsSubstring(s, substr string) bool {
	return findSubstring(s, substr) != -1
}

func extractPathFromURL(url string) string {
	protocolEnd := findSubstring(url, "://")
	if protocolEnd == -1 {
		return url
	}

	domainEnd := findSubstringFrom(url, "/", protocolEnd+3)
	if domainEnd == -1 {
		return ""
	}

	return url[domainEnd:]
}
