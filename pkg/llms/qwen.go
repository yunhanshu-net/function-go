package llms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// QwenClient 千问客户端实现
type QwenClient struct {
	APIKey  string         `json:"api_key"`
	BaseURL string         `json:"base_url"`
	Options *ClientOptions `json:"options"`
}

// NewQwenClient 创建千问客户端（保持向后兼容）
func NewQwenClient(apiKey string) *QwenClient {
	return NewQwenClientWithOptions(apiKey, DefaultClientOptions())
}

// NewQwenClientWithOptions 创建带配置的千问客户端
func NewQwenClientWithOptions(apiKey string, options *ClientOptions) *QwenClient {
	// 如果没有提供options，使用默认配置
	if options == nil {
		options = DefaultClientOptions()
	}

	// 设置默认BaseURL
	baseURL := options.BaseURL
	if baseURL == "" {
		baseURL = "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"
	}

	// 🎯 不再在构造函数中创建固定的HTTP客户端
	// 而是在每次Chat请求时动态创建，以支持不同的超时时间

	return &QwenClient{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Options: options,
	}
}

// Chat 实现LLMClient接口
func (q *QwenClient) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// 直接使用原始上下文，超时由HTTP客户端控制

	// 构造千问API请求
	apiReq := map[string]interface{}{
		"model": req.Model,
		"input": map[string]interface{}{
			"messages": req.Messages,
		},
		"parameters": map[string]interface{}{
			"max_tokens":  req.MaxTokens,
			"temperature": req.Temperature,
		},
	}

	if req.Model == "" {
		apiReq["model"] = "qwen-turbo"
	}
	if req.MaxTokens == 0 {
		apiReq["max_tokens"] = 4000
	}
	if req.Temperature == 0 {
		apiReq["temperature"] = 0.7
	}

	// 🎯 动态创建HTTP客户端，支持请求级别的超时配置
	timeout := q.Options.Timeout // 默认使用客户端配置的超时时间
	if req.Timeout != nil && *req.Timeout > 0 {
		timeout = *req.Timeout // 如果请求中指定了超时时间，则使用请求的超时时间
	}

	httpClient := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			MaxIdleConns:       q.Options.MaxIdleConns,
			IdleConnTimeout:    q.Options.IdleConnTimeout,
			DisableCompression: true,
		},
	}

	// 发送HTTP请求
	jsonData, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", q.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+q.APIKey)

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var apiResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查错误
	if errMsg, exists := apiResp["message"]; exists {
		return &ChatResponse{
			Error: fmt.Sprintf("API错误: %v", errMsg),
		}, nil
	}

	// 提取回答内容
	output, ok := apiResp["output"].(map[string]interface{})
	if !ok {
		return &ChatResponse{
			Error: "响应格式错误：没有找到output",
		}, nil
	}

	choices, ok := output["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return &ChatResponse{
			Error: "响应格式错误：没有找到choices",
		}, nil
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return &ChatResponse{
			Error: "响应格式错误：choice格式错误",
		}, nil
	}

	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return &ChatResponse{
			Error: "响应格式错误：message格式错误",
		}, nil
	}

	content, ok := message["content"].(string)
	if !ok {
		return &ChatResponse{
			Error: "响应格式错误：content格式错误",
		}, nil
	}

	return &ChatResponse{
		Content: content,
	}, nil
}

// GetModelName 获取模型名称
func (q *QwenClient) GetModelName() string {
	return "qwen-turbo"
}

// GetProvider 获取提供商名称
func (q *QwenClient) GetProvider() string {
	return "千问"
}
