package llms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// KimiClient Kimi客户端实现
type KimiClient struct {
	APIKey  string
	BaseURL string
	Options *ClientOptions
}

// KimiRequest Kimi API请求结构
type KimiRequest struct {
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
	ResponseFormat   interface{} `json:"response_format,omitempty"`
}

// KimiResponse Kimi API响应结构
type KimiResponse struct {
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

// NewKimiClient 创建Kimi客户端（保持向后兼容）
func NewKimiClient(apiKey string) *KimiClient {
	return NewKimiClientWithOptions(apiKey, DefaultClientOptions())
}

// NewKimiClientWithOptions 创建带配置的Kimi客户端
func NewKimiClientWithOptions(apiKey string, options *ClientOptions) *KimiClient {
	if options == nil {
		options = DefaultClientOptions()
	}

	baseURL := options.BaseURL
	if baseURL == "" {
		baseURL = "https://api.moonshot.cn/v1/chat/completions"
	}

	return &KimiClient{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Options: options,
	}
}

// Chat 实现LLMClient接口的Chat方法
func (c *KimiClient) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// 直接使用原始上下文，超时由HTTP客户端控制

	// 转换为Kimi API请求格式
	kimiReq := &KimiRequest{
		Model:       "kimi-k2-0711-preview", // 默认使用kimi-k2-0711-preview模型
		Messages:    req.Messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		Stream:      false,
	}

	// 如果用户指定了模型，使用用户指定的模型
	if req.Model != "" {
		kimiReq.Model = req.Model
	}

	// 设置默认值
	if kimiReq.Temperature == 0 {
		kimiReq.Temperature = 0.6 // Kimi推荐值
	}
	if kimiReq.MaxTokens == 0 {
		kimiReq.MaxTokens = 1024 // 默认值
	}

	// 序列化请求
	jsonData, err := json.Marshal(kimiReq)
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

	// 调试日志：记录超时设置
	fmt.Printf("🔍 Kimi超时调试: 客户端默认超时=%v, 请求超时=%v, 最终使用超时=%v\n",
		c.Options.Timeout, req.Timeout, timeout)

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
	var kimiResp KimiResponse
	if err := json.Unmarshal(body, &kimiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查API错误
	if kimiResp.Error != nil {
		return &ChatResponse{
			Error: fmt.Sprintf("API错误: %s", kimiResp.Error.Message),
		}, nil
	}

	// 转换为通用响应格式
	if len(kimiResp.Choices) == 0 {
		return &ChatResponse{
			Error: "API响应中没有choices",
		}, nil
	}

	chatResp := &ChatResponse{
		Content: kimiResp.Choices[0].Message.Content,
		Usage: &Usage{
			PromptTokens:     kimiResp.Usage.PromptTokens,
			CompletionTokens: kimiResp.Usage.CompletionTokens,
			TotalTokens:      kimiResp.Usage.TotalTokens,
		},
	}

	return chatResp, nil
}

// GetModelName 获取模型名称
func (c *KimiClient) GetModelName() string {
	return "kimi-k2-0711-preview"
}

// GetProvider 获取提供商名称
func (c *KimiClient) GetProvider() string {
	return "Kimi"
}

// ChatStream 实现流式聊天接口
func (c *KimiClient) ChatStream(ctx context.Context, req *ChatRequest) (<-chan *StreamChunk, error) {
	// 创建流式响应通道
	chunkChan := make(chan *StreamChunk, 1)

	// 在goroutine中处理
	go func() {
		defer close(chunkChan)
		chunkChan <- &StreamChunk{
			Error: "Kimi 客户端暂不支持流式响应，请使用 Chat 方法",
			Done:  true,
		}
	}()

	return chunkChan, nil
}

// GetSupportedModels 获取支持的模型列表
func (c *KimiClient) GetSupportedModels() []string {
	return []string{
		"kimi-k2-0711-preview",
		"moonshot-v1-8k",
		"moonshot-v1-32k",
		"moonshot-v1-128k",
		"moonshot-v1-auto",
		"kimi-latest",
		"moonshot-v1-8k-vision-preview",
		"moonshot-v1-32k-vision-preview",
		"moonshot-v1-128k-vision-preview",
		"kimi-thinking-preview",
	}
}

// GetPricingInfo 获取价格信息
func (c *KimiClient) GetPricingInfo() map[string]interface{} {
	return map[string]interface{}{
		"model":          "kimi-k2-0711-preview",
		"context_length": "32K",
		"input_price":    "0.0001元/1K tokens",
		"output_price":   "0.0002元/1K tokens",
		"note":           "价格信息仅供参考，请以官方最新定价为准",
	}
}
