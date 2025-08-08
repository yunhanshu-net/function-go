package runner

import (
	"context"
	"fmt"
	"log"
)

// UploadExample 展示如何使用新的上传功能
func UploadExample() {
	// 1. 创建带有trace_id的context
	ctx := context.WithValue(context.Background(), "trace_id", "example-trace-123")

	// 2. 创建runner context（新格式）
	runnerCtx := NewContext(ctx, "POST", "api/upload_demo")

	// 3. 获取上传路径信息
	uploadPath := runnerCtx.GetUploadPath()
	fmt.Printf("📁 上传路径: %s\n", uploadPath)

	// 4. 从数据创建文件并自动上传
	testData := []byte(`{
		"message": "Hello, World!",
		"timestamp": "2025-06-28T14:30:00Z",
		"data": [1, 2, 3, 4, 5]
	}`)

	files, err := runnerCtx.CreateFilesFromData("demo.json", testData)
	if err != nil {
		log.Printf("❌ 创建文件失败: %v", err)
		return
	}

	// 5. 检查上传结果
	if len(files.GetFiles()) > 0 {
		file := files.GetFiles()[0]
		fmt.Printf("✅ 文件上传成功!\n")
		fmt.Printf("   📄 文件名: %s\n", file.Name)
		fmt.Printf("   📏 文件大小: %d bytes\n", file.Size)
		fmt.Printf("   🔗 访问URL: %s\n", file.URL)
		fmt.Printf("   📅 创建时间: %s\n", file.CreatedAt)
		fmt.Printf("   ✨ 状态: %s\n", file.Status)
	}

	// 6. 设置文件属性
	files.SetNote("演示文件上传功能")
	files.SetConfig("demo", true)
	files.SetMetadata("purpose", "example")

	fmt.Printf("📝 备注: %s\n", files.GetNote())
	fmt.Printf("📊 文件总数: %d\n", files.GetFileCount())
	fmt.Printf("📦 总大小: %d bytes\n", files.GetTotalSize())
}

// UploadWithLifecycleExample 展示生命周期管理
func UploadWithLifecycleExample() {
	ctx := context.WithValue(context.Background(), "trace_id", "lifecycle-example-456")
	runnerCtx := NewContext(ctx, "POST", "api/lifecycle_demo")

	fmt.Println("\n🔄 生命周期管理示例:")

	// 1. 临时文件（下载一次后删除）
	tempFiles := runnerCtx.NewTemporaryFiles()
	tempData := []byte("这是临时文件内容")
	err := tempFiles.AddFileFromData("temp_report.txt", tempData)
	if err != nil {
		log.Printf("❌ 创建临时文件失败: %v", err)
	} else {
		fmt.Printf("🗂️  临时文件: %s (下载1次后删除)\n", tempFiles.GetFiles()[0].Name)
	}

	// 2. 有效期文件（7天后过期）
	expiringFiles := runnerCtx.NewExpiringFiles()
	expiringData := []byte("这是有效期文件内容")
	err = expiringFiles.AddFileFromData("weekly_report.pdf", expiringData)
	if err != nil {
		log.Printf("❌ 创建有效期文件失败: %v", err)
	} else {
		fmt.Printf("📅 有效期文件: %s (7天后过期)\n", expiringFiles.GetFiles()[0].Name)
	}

	// 3. 永久文件（无限制）
	permanentFiles := runnerCtx.NewPermanentFiles()
	permanentData := []byte("这是永久保存的文件内容")
	err = permanentFiles.AddFileFromData("backup.zip", permanentData)
	if err != nil {
		log.Printf("❌ 创建永久文件失败: %v", err)
	} else {
		fmt.Printf("♾️  永久文件: %s (无限制)\n", permanentFiles.GetFiles()[0].Name)
	}
}

// PathFormatExample 展示路径格式
func PathFormatExample() {
	fmt.Println("\n📁 路径格式示例:")

	examples := []struct {
		router string
		method string
	}{
		{"api/upload", "POST"},
		{"customer/export", "GET"},
		{"report/generate", "POST"},
	}

	for _, example := range examples {
		ctx := context.WithValue(context.Background(), "trace_id", "path-example")
		runnerCtx := NewContext(ctx, example.method, example.router)

		uploadPath := runnerCtx.GetUploadPath()
		fmt.Printf("   %s %s -> %s\n",
			example.method, example.router, uploadPath)
	}
}

// 运行所有示例
func RunUploadExamples() {
	fmt.Println("🚀 文件上传功能示例")
	fmt.Println("====================================================")

	UploadExample()
	UploadWithLifecycleExample()
	PathFormatExample()

	fmt.Println("\n✨ 示例运行完成!")
}
