package llms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ClaudeClient Claude客户端实现（通过2233代理）
type ClaudeClient struct {
	APIKey  string
	BaseURL string
	Options *ClientOptions
}

// ClaudeRequest Claude API请求结构（兼容OpenAI格式）
type ClaudeRequest struct {
	Model            string      `json:"model"`
	Messages         []Message   `json:"messages"`
	MaxTokens        int         `json:"max_tokens,omitempty"`
	Temperature      float64     `json:"temperature,omitempty"`
	TopP             float64     `json:"top_p,omitempty"`
	N                int         `json:"n,omitempty"`
	Stream           bool        `json:"stream,omitempty"`
	Stop             interface{} `json:"stop,omitempty"`
	PresencePenalty  float64     `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64     `json:"frequency_penalty,omitempty"`
}

// ClaudeResponse Claude API响应结构（兼容OpenAI格式）
type ClaudeResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// NewClaudeClient 创建Claude客户端（保持向后兼容）
func NewClaudeClient(apiKey string) *ClaudeClient {
	return NewClaudeClientWithOptions(apiKey, DefaultClientOptions())
}

// NewClaudeClientWithOptions 创建带配置的Claude客户端
func NewClaudeClientWithOptions(apiKey string, options *ClientOptions) *ClaudeClient {
	if options == nil {
		options = DefaultClientOptions()
	}

	baseURL := options.BaseURL
	if baseURL == "" {
		baseURL = "https://api.gptsapi.net/v1/chat/completions" // 2233代理API
	}

	return &ClaudeClient{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Options: options,
	}
}

// Chat 实现LLMClient接口的Chat方法
func (c *ClaudeClient) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// 直接使用原始上下文，超时由HTTP客户端控制

	// 转换为Claude API请求格式
	claudeReq := &ClaudeRequest{
		Model:       "claude-sonnet-4-20250514", // 默认使用claude-sonnet-4-20250514模型（性价比最高）
		Messages:    req.Messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		Stream:      false,
	}

	// 如果用户指定了模型，使用用户指定的模型
	if req.Model != "" {
		claudeReq.Model = req.Model
	}

	// 设置默认值
	if claudeReq.Temperature == 0 {
		claudeReq.Temperature = 0.7 // Claude推荐值
	}
	if claudeReq.MaxTokens == 0 {
		claudeReq.MaxTokens = 1024 // 默认值
	}

	// 序列化请求
	jsonData, err := json.Marshal(claudeReq)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

	// 🎯 动态创建HTTP客户端，支持请求级别的超时配置
	timeout := c.Options.Timeout // 默认使用客户端配置的超时时间
	if req.Timeout != nil && *req.Timeout > 0 {
		timeout = *req.Timeout // 如果请求中指定了超时时间，则使用请求的超时时间
	}

	httpClient := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			MaxIdleConns:       c.Options.MaxIdleConns,
			IdleConnTimeout:    c.Options.IdleConnTimeout,
			DisableCompression: true,
		},
	}

	// 发送请求
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var claudeResp ClaudeResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查API错误
	if claudeResp.Error != nil {
		return &ChatResponse{
			Error: fmt.Sprintf("API错误: %s", claudeResp.Error.Message),
		}, nil
	}

	// 转换为通用响应格式
	if len(claudeResp.Choices) == 0 {
		return &ChatResponse{
			Error: "API响应中没有choices",
		}, nil
	}

	chatResp := &ChatResponse{
		Content: claudeResp.Choices[0].Message.Content,
		Usage: &Usage{
			PromptTokens:     claudeResp.Usage.PromptTokens,
			CompletionTokens: claudeResp.Usage.CompletionTokens,
			TotalTokens:      claudeResp.Usage.TotalTokens,
		},
	}

	return chatResp, nil
}

// GetModelName 获取模型名称
func (c *ClaudeClient) GetModelName() string {
	return "claude-sonnet-4-20250514"
}

// GetProvider 获取提供商名称
func (c *ClaudeClient) GetProvider() string {
	return "Claude"
}

// ChatStream 实现流式聊天接口
func (c *ClaudeClient) ChatStream(ctx context.Context, req *ChatRequest) (<-chan *StreamChunk, error) {
	// 创建流式响应通道
	chunkChan := make(chan *StreamChunk, 1)

	// 在goroutine中处理
	go func() {
		defer close(chunkChan)
		chunkChan <- &StreamChunk{
			Error: "Claude 客户端暂不支持流式响应，请使用 Chat 方法",
			Done:  true,
		}
	}()

	return chunkChan, nil
}

// GetSupportedModels 获取支持的模型列表
func (c *ClaudeClient) GetSupportedModels() []string {
	return []string{
		"claude-sonnet-4-20250514",   // 推荐：性价比最高
		"claude-3-5-sonnet-20241022", // 备选：性能好但稍贵
		"claude-3-5-haiku-20241022",  // 备选：速度快但能力稍弱
		"claude-3-sonnet-20240229",   // 经典版本
		"claude-3-haiku-20240307",    // 轻量版本
		"claude-3-opus-20240229",     // 最强版本（最贵）
	}
}

// GetPricingInfo 获取价格信息
func (c *ClaudeClient) GetPricingInfo() map[string]interface{} {
	return map[string]interface{}{
		"model":          "claude-sonnet-4-20250514",
		"context_length": "200K",
		"input_price":    "$3.30 / 1M tokens",
		"output_price":   "$16.50 / 1M tokens",
		"note":           "性价比最高的Claude模型，推荐使用",
	}
}
