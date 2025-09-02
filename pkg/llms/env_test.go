package llms

import (
	"os"
	"testing"
)

// TestEnvironmentVariableSupport 测试环境变量支持
func TestEnvironmentVariableSupport(t *testing.T) {
	// 直接使用用户设置的环境变量，不手动设置

	t.Run("测试千问3 Coder环境变量", func(t *testing.T) {
		// 直接使用用户设置的环境变量
		client, err := NewQwen3CoderClientFromEnv()
		if err != nil {
			t.Fatalf("从环境变量创建千问3 Coder客户端失败: %v", err)
		}

		if client == nil {
			t.Fatal("客户端为空")
		}

		// 验证API Key
		qwenClient, ok := client.(*Qwen3CoderClient)
		if !ok {
			t.Fatal("客户端类型错误")
		}

		// 验证API Key不为空即可，不检查具体值
		if qwenClient.APIKey == "" {
			t.Error("API Key为空")
		}

		t.Logf("✅ 千问3 Coder环境变量测试通过，API Key: %s", qwenClient.APIKey[:20]+"...")
	})

	t.Run("测试千问环境变量", func(t *testing.T) {
		// 直接使用用户设置的环境变量
		client, err := NewQwenClientFromEnv()
		if err != nil {
			t.Fatalf("从环境变量创建千问客户端失败: %v", err)
		}

		if client == nil {
			t.Fatal("客户端为空")
		}

		// 验证API Key
		qwenClient, ok := client.(*QwenClient)
		if !ok {
			t.Fatal("客户端类型错误")
		}

		// 验证API Key不为空即可，不检查具体值
		if qwenClient.APIKey == "" {
			t.Error("API Key为空")
		}

		t.Logf("✅ 千问环境变量测试通过，API Key: %s", qwenClient.APIKey[:20]+"...")
	})

	t.Run("测试DeepSeek环境变量", func(t *testing.T) {
		// 直接使用用户设置的环境变量
		client, err := NewDeepSeekClientFromEnv()
		if err != nil {
			t.Fatalf("从环境变量创建DeepSeek客户端失败: %v", err)
		}

		if client == nil {
			t.Fatal("客户端为空")
		}

		// 验证API Key
		deepSeekClient, ok := client.(*DeepSeekClient)
		if !ok {
			t.Fatal("客户端类型错误")
		}

		// 验证API Key不为空即可，不检查具体值
		if deepSeekClient.APIKey == "" {
			t.Error("API Key为空")
		}

		t.Logf("✅ DeepSeek环境变量测试通过，API Key: %s", deepSeekClient.APIKey[:20]+"...")
	})

	t.Run("测试Kimi环境变量", func(t *testing.T) {
		// 直接使用用户设置的环境变量
		client, err := NewKimiClientFromEnv()
		if err != nil {
			t.Fatalf("从环境变量创建Kimi客户端失败: %v", err)
		}

		if client == nil {
			t.Fatal("客户端为空")
		}

		// 验证API Key
		kimiClient, ok := client.(*KimiClient)
		if !ok {
			t.Fatal("客户端类型错误")
		}

		// 验证API Key不为空即可，不检查具体值
		if kimiClient.APIKey == "" {
			t.Error("API Key为空")
		}

		t.Logf("✅ Kimi环境变量测试通过，API Key: %s", kimiClient.APIKey[:20]+"...")
	})

	t.Run("测试环境变量为空时的错误处理", func(t *testing.T) {
		// 保存原始环境变量
		originalQianwen := os.Getenv("QIANWEN_API_KEY")

		// 清除环境变量
		os.Unsetenv("QIANWEN_API_KEY")

		// 测试从环境变量创建客户端应该失败
		_, err := NewQwen3CoderClientFromEnv()
		if err == nil {
			t.Fatal("环境变量为空时应该返回错误")
		}

		// 恢复环境变量
		if originalQianwen != "" {
			os.Setenv("QIANWEN_API_KEY", originalQianwen)
		}

		t.Logf("✅ 环境变量为空时的错误处理测试通过: %v", err)
	})

	t.Run("测试通用环境变量创建函数", func(t *testing.T) {
		// 直接使用用户设置的环境变量
		client, err := NewLLMClientFromEnv(ProviderQwen3Coder)
		if err != nil {
			t.Fatalf("通用环境变量创建函数失败: %v", err)
		}

		if client == nil {
			t.Fatal("客户端为空")
		}

		t.Logf("✅ 通用环境变量创建函数测试通过")
	})

	t.Run("测试通用环境变量创建函数-Kimi", func(t *testing.T) {
		// 直接使用用户设置的环境变量
		client, err := NewLLMClientFromEnv(ProviderKimi)
		if err != nil {
			t.Fatalf("通用环境变量创建函数-Kimi失败: %v", err)
		}

		if client == nil {
			t.Fatal("客户端为空")
		}

		t.Logf("✅ 通用环境变量创建函数-Kimi测试通过")
	})

	t.Run("测试豆包环境变量", func(t *testing.T) {
		// 直接使用用户设置的环境变量
		client, err := NewDouBaoClientFromEnv()
		if err != nil {
			t.Fatalf("从环境变量创建豆包客户端失败: %v", err)
		}

		if client == nil {
			t.Fatal("客户端为空")
		}

		// 验证API Key
		douBaoClient, ok := client.(*DouBaoClient)
		if !ok {
			t.Fatal("客户端类型错误")
		}

		// 验证API Key不为空即可，不检查具体值
		if douBaoClient.APIKey == "" {
			t.Error("API Key为空")
		}

		t.Logf("✅ 豆包环境变量测试通过，API Key: %s", douBaoClient.APIKey[:20]+"...")
	})

	t.Run("测试通用环境变量创建函数-豆包", func(t *testing.T) {
		// 直接使用用户设置的环境变量
		client, err := NewLLMClientFromEnv(ProviderDouBao)
		if err != nil {
			t.Fatalf("通用环境变量创建函数-豆包失败: %v", err)
		}

		if client == nil {
			t.Fatal("客户端为空")
		}

		t.Logf("✅ 通用环境变量创建函数-豆包测试通过")
	})
}

// TestSameAPIKeyForQwen 测试千问和千问3 Coder使用同一个API Key
func TestSameAPIKeyForQwen(t *testing.T) {
	// 直接使用用户设置的环境变量，不手动设置

	t.Run("验证千问和千问3 Coder使用相同API Key", func(t *testing.T) {
		// 创建千问客户端
		qwenClient, err := NewQwenClientFromEnv()
		if err != nil {
			t.Fatalf("创建千问客户端失败: %v", err)
		}

		// 创建千问3 Coder客户端
		qwen3CoderClient, err := NewQwen3CoderClientFromEnv()
		if err != nil {
			t.Fatalf("创建千问3 Coder客户端失败: %v", err)
		}

		// 验证两个客户端使用相同的API Key
		qwenAPIKey := qwenClient.(*QwenClient).APIKey
		qwen3CoderAPIKey := qwen3CoderClient.(*Qwen3CoderClient).APIKey

		if qwenAPIKey != qwen3CoderAPIKey {
			t.Errorf("千问和千问3 Coder的API Key应该相同，千问: %s, 千问3 Coder: %s",
				qwenAPIKey, qwen3CoderAPIKey)
		}

		// 验证API Key不为空
		if qwenAPIKey == "" {
			t.Error("千问API Key为空")
		}
		if qwen3CoderAPIKey == "" {
			t.Error("千问3 Coder API Key为空")
		}

		t.Logf("✅ 千问和千问3 Coder使用相同API Key验证通过: %s", qwenAPIKey[:20]+"...")
	})
}
