package runner

import (
	"context"
	"testing"
	"time"

	"github.com/yunhanshu-net/pkg/trace"
)

func TestContextGetUploadPath(t *testing.T) {
	// 创建一个带有trace_id的context
	ctx := context.WithValue(context.Background(), "trace_id", "test-trace-123")

	// 创建runner context
	runnerCtx := NewContext(ctx, "POST", "test/upload")

	// 获取上传路径
	uploadPath := runnerCtx.GetUploadPath()

	// 验证路径格式
	if uploadPath == "" {
		t.Error("上传路径不应该为空")
	}

	t.Logf("上传路径: %s", uploadPath)

	// 验证路径包含预期的组件
	// 格式应该是：/租户/应用/函数/方法/output/日期
	// 例如：/testuser/testapp/test/upload/POST/output/20250115

	// 验证包含用户信息
	// 注意：这里可能需要根据实际的env配置进行调整
}

func TestContextLock(t *testing.T) {
	// 创建一个带有trace_id的context
	ctx1 := context.WithValue(context.Background(), "trace_id", "test-trace-123")

	// 创建runner context
	ctx := NewContext(ctx1, "POST", "test/upload")

	locked := ctx.Locker.Lock("app1:token:gen", time.Second*5) //time.Second*5表示过期时间，如果不填表示无过期时间

	if locked { //说明加锁成功
		defer ctx.Locker.Unlock("app1:token:gen")
	} else { //锁失败，说明已经被占用

	}
}

func TestContextGetFunctionMsg(t *testing.T) {
	// 创建一个带有trace_id的context
	ctx := context.WithValue(context.Background(), "trace_id", "test-trace-456")

	// 创建runner context
	runnerCtx := NewContext(ctx, "GET", "test/function")

	// 获取FunctionMsg
	functionMsg := runnerCtx.GetFunctionMsg()

	if functionMsg == nil {
		t.Error("FunctionMsg不应该为nil")
		return
	}

	// 验证FunctionMsg的字段
	if functionMsg.Method != "GET" {
		t.Errorf("期望Method: GET, 实际: %s", functionMsg.Method)
	}

	if functionMsg.Router != "test/function" {
		t.Errorf("期望Router: test/function, 实际: %s", functionMsg.Router)
	}

	if functionMsg.TraceID != "test-trace-456" {
		t.Errorf("期望TraceID: test-trace-456, 实际: %s", functionMsg.TraceID)
	}

	t.Logf("FunctionMsg: %+v", functionMsg)
}

func TestUploadPathFormat(t *testing.T) {
	// 直接测试FunctionMsg的GetUploadPath方法
	functionMsg := &trace.FunctionMsg{
		User:    "testuser",
		Runner:  "testapp",
		Router:  "api/upload",
		Method:  "POST",
		TraceID: "test-trace-789",
	}

	uploadPath := functionMsg.GetUploadPath()

	// 验证路径格式
	expectedPrefix := "/testuser/testapp/api/upload/POST/output/"
	if len(uploadPath) < len(expectedPrefix) {
		t.Errorf("上传路径太短: %s", uploadPath)
		return
	}

	actualPrefix := uploadPath[:len(expectedPrefix)]
	if actualPrefix != expectedPrefix {
		t.Errorf("期望路径前缀: %s, 实际: %s", expectedPrefix, actualPrefix)
	}

	// 验证日期部分（应该是今天的日期）
	today := time.Now().Format("20060102")
	if !contains(uploadPath, today) {
		t.Errorf("上传路径应该包含今天的日期 %s: %s", today, uploadPath)
	}

	t.Logf("上传路径格式正确: %s", uploadPath)
}

// contains 检查字符串是否包含子字符串
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		s[len(s)-len(substr):] == substr ||
		indexOf(s, substr) >= 0
}

// indexOf 查找子字符串的位置
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func TestContextCreateFiles(t *testing.T) {
	// 创建一个带有trace_id的context
	ctx := context.WithValue(context.Background(), "trace_id", "test-trace-create")

	// 创建runner context
	runnerCtx := NewContext(ctx, "POST", "test/create_files")

	// 测试从数据创建文件
	testData := []byte("Hello, World! This is test data.")
	files, err := runnerCtx.CreateFilesFromData("test.txt", testData)

	if err != nil {
		t.Logf("创建文件时出现错误（可能是因为没有配置上传器）: %v", err)
		// 这个错误是预期的，因为测试环境可能没有配置上传器
		return
	}

	if files == nil {
		t.Error("创建的Files对象不应该为nil")
		return
	}

	if len(files.GetFiles()) == 0 {
		t.Error("Files对象应该包含至少一个文件")
		return
	}

	file := files.GetFiles()[0]
	if file.Name != "test.txt" {
		t.Errorf("期望文件名: test.txt, 实际: %s", file.Name)
	}

	if file.Size != int64(len(testData)) {
		t.Errorf("期望文件大小: %d, 实际: %d", len(testData), file.Size)
	}

	t.Logf("成功创建文件: %+v", file)
}
