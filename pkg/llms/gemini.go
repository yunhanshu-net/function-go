package llms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GeminiClient Gemini客户端实现
type GeminiClient struct {
	APIKey  string
	BaseURL string
	Options *ClientOptions
}

// GeminiRequest Gemini API请求结构
type GeminiRequest struct {
	Contents         []GeminiContent         `json:"contents"`
	GenerationConfig *GeminiGenerationConfig `json:"generationConfig,omitempty"`
	SafetySettings   []GeminiSafetySetting   `json:"safetySettings,omitempty"`
}

// GeminiContent 内容结构
type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

// GeminiPart 内容部分
type GeminiPart struct {
	Text string `json:"text"`
}

// GeminiGenerationConfig 生成配置
type GeminiGenerationConfig struct {
	Temperature     float64  `json:"temperature,omitempty"`
	TopP            float64  `json:"topP,omitempty"`
	TopK            int      `json:"topK,omitempty"`
	MaxOutputTokens int      `json:"maxOutputTokens,omitempty"`
	CandidateCount  int      `json:"candidateCount,omitempty"`
	StopSequences   []string `json:"stopSequences,omitempty"`
}

// GeminiSafetySetting 安全设置
type GeminiSafetySetting struct {
	Category  string `json:"category"`
	Threshold string `json:"threshold"`
}

// GeminiResponse Gemini API响应结构
type GeminiResponse struct {
	Candidates     []GeminiCandidate     `json:"candidates"`
	PromptFeedback *GeminiPromptFeedback `json:"promptFeedback,omitempty"`
	UsageMetadata  *GeminiUsageMetadata  `json:"usageMetadata,omitempty"`
}

// GeminiCandidate 候选响应
type GeminiCandidate struct {
	Content       GeminiContent        `json:"content"`
	FinishReason  string               `json:"finishReason"`
	Index         int                  `json:"index"`
	SafetyRatings []GeminiSafetyRating `json:"safetyRatings,omitempty"`
}

// GeminiPromptFeedback 提示反馈
type GeminiPromptFeedback struct {
	SafetyRatings []GeminiSafetyRating `json:"safetyRatings"`
}

// GeminiSafetyRating 安全评级
type GeminiSafetyRating struct {
	Category    string `json:"category"`
	Probability string `json:"probability"`
}

// GeminiUsageMetadata 使用元数据
type GeminiUsageMetadata struct {
	PromptTokenCount     int `json:"promptTokenCount"`
	CandidatesTokenCount int `json:"candidatesTokenCount"`
	TotalTokenCount      int `json:"totalTokenCount"`
}

// NewGeminiClient 创建Gemini客户端（保持向后兼容）
func NewGeminiClient(apiKey string) *GeminiClient {
	return NewGeminiClientWithOptions(apiKey, DefaultClientOptions())
}

// NewGeminiClientWithOptions 创建带配置的Gemini客户端
func NewGeminiClientWithOptions(apiKey string, options *ClientOptions) *GeminiClient {
	if options == nil {
		options = DefaultClientOptions()
	}

	baseURL := options.BaseURL
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com/v1beta/models"
	}

	return &GeminiClient{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Options: options,
	}
}

// Chat 实现LLMClient接口的Chat方法
func (c *GeminiClient) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// 直接使用原始上下文，超时由HTTP客户端控制

	// 转换为Gemini API请求格式
	geminiReq := &GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: req.Messages[len(req.Messages)-1].Content}, // 使用最后一条用户消息
				},
			},
		},
		GenerationConfig: &GeminiGenerationConfig{
			Temperature:     req.Temperature,
			MaxOutputTokens: req.MaxTokens,
		},
	}

	// 设置默认值
	if geminiReq.GenerationConfig.Temperature == 0 {
		geminiReq.GenerationConfig.Temperature = 0.7 // Gemini推荐值
	}
	if geminiReq.GenerationConfig.MaxOutputTokens == 0 {
		geminiReq.GenerationConfig.MaxOutputTokens = 1024 // 默认值
	}

	// 序列化请求
	jsonData, err := json.Marshal(geminiReq)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	// 构建完整的API URL
	apiURL := fmt.Sprintf("%s/gemini-2.0-flash-exp:generateContent?key=%s", c.BaseURL, c.APIKey)

	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")

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
	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查是否有候选响应
	if len(geminiResp.Candidates) == 0 {
		return &ChatResponse{
			Error: "API响应中没有候选内容",
		}, nil
	}

	// 获取第一个候选响应的文本内容
	content := ""
	if len(geminiResp.Candidates[0].Content.Parts) > 0 {
		content = geminiResp.Candidates[0].Content.Parts[0].Text
	}

	// 构建使用统计
	var usage *Usage
	if geminiResp.UsageMetadata != nil {
		usage = &Usage{
			PromptTokens:     geminiResp.UsageMetadata.PromptTokenCount,
			CompletionTokens: geminiResp.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      geminiResp.UsageMetadata.TotalTokenCount,
		}
	}

	chatResp := &ChatResponse{
		Content: content,
		Usage:   usage,
	}

	return chatResp, nil
}

// GetModelName 获取模型名称
func (c *GeminiClient) GetModelName() string {
	return "gemini-2.0-flash-exp"
}

// GetProvider 获取提供商名称
func (c *GeminiClient) GetProvider() string {
	return "Gemini"
}

// GetSupportedModels 获取支持的模型列表
func (c *GeminiClient) GetSupportedModels() []string {
	return []string{
		"gemini-2.0-flash-exp",
		"gemini-2.0-flash",
		"gemini-2.0-pro",
		"gemini-1.5-flash",
		"gemini-1.5-pro",
		"gemini-1.5-flash-001",
		"gemini-1.5-pro-001",
	}
}

// GetPricingInfo 获取价格信息
func (c *GeminiClient) GetPricingInfo() map[string]interface{} {
	return map[string]interface{}{
		"model":          "gemini-2.0-flash-exp",
		"context_length": "1M",
		"input_price":    "免费额度充足",
		"output_price":   "免费额度充足",
		"note":           "Gemini提供大量免费额度，适合开发和测试使用",
	}
}
