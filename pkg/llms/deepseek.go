package llms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yunhanshu-net/pkg/logger"
)

// DeepSeekAPIResponse DeepSeek API响应结构体
type DeepSeekAPIResponse struct {
	Error *struct {
		Code    string      `json:"code"`
		Message string      `json:"message"`
		Param   interface{} `json:"param"`
		Type    string      `json:"type"`
	} `json:"error,omitempty"`
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices,omitempty"`
	Usage *struct {
		PromptTokens     float64 `json:"prompt_tokens"`
		CompletionTokens float64 `json:"completion_tokens"`
		TotalTokens      float64 `json:"total_tokens"`
	} `json:"usage,omitempty"`
}

// DeepSeekClient DeepSeek客户端实现
type DeepSeekClient struct {
	APIKey  string         `json:"api_key"`
	BaseURL string         `json:"base_url"`
	Options *ClientOptions `json:"options"`
	Model   string         `json:"model"` // 🆕 添加模型名称字段
}

// NewDeepSeekClient 创建DeepSeek客户端（保持向后兼容）
func NewDeepSeekClient(apiKey string) *DeepSeekClient {
	return NewDeepSeekClientWithOptions(apiKey, DefaultClientOptions())
}

// NewDeepSeekClientWithOptions 创建带配置的DeepSeek客户端
func NewDeepSeekClientWithOptions(apiKey string, options *ClientOptions) *DeepSeekClient {
	if options == nil {
		options = DefaultClientOptions()
	}

	baseURL := options.BaseURL
	if baseURL == "" {
		baseURL = "https://api.deepseek.com/v1/chat/completions"
	}

	// 🎯 不再在构造函数中创建固定的HTTP客户端
	// 而是在每次Chat请求时动态创建，以支持不同的超时时间

	return &DeepSeekClient{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Options: options,
		Model:   "deepseek-reasoner", // 🆕 设置默认模型
	}
}

// 🆕 SetModel 设置模型名称
func (d *DeepSeekClient) SetModel(model string) {
	d.Model = model
}

// 🆕 GetModelName 获取模型名称（现在返回实际设置的模型）
func (d *DeepSeekClient) GetModelName() string {
	return d.Model
}

// GetProvider 获取提供商名称
func (d *DeepSeekClient) GetProvider() string {
	return "DeepSeek"
}

// Chat 实现LLMClient接口
func (d *DeepSeekClient) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// 直接使用原始上下文，超时由HTTP客户端控制

	// 构造DeepSeek API请求
	apiReq := map[string]interface{}{
		"model":       req.Model,
		"messages":    req.Messages,
		"max_tokens":  req.MaxTokens,
		"temperature": req.Temperature,
	}

	if req.Model == "" {
		apiReq["model"] = "deepseek-reasoner"
	}
	if req.MaxTokens <= 10 {
		apiReq["max_tokens"] = 4000
	}
	if req.Temperature == 0 {
		apiReq["temperature"] = 0.7
	}

	// 🎯 动态创建HTTP客户端，支持请求级别的超时配置
	timeout := d.Options.Timeout // 默认使用客户端配置的超时时间
	if req.Timeout != nil && *req.Timeout > 0 {
		timeout = *req.Timeout // 如果请求中指定了超时时间，则使用请求的超时时间
	}

	httpClient := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			MaxIdleConns:       d.Options.MaxIdleConns,
			IdleConnTimeout:    d.Options.IdleConnTimeout,
			DisableCompression: true,
		},
	}

	// 发送HTTP请求
	jsonData, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", d.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+d.APIKey)

	// 设置自定义User-Agent
	if d.Options != nil && d.Options.UserAgent != "" {
		httpReq.Header.Set("User-Agent", d.Options.UserAgent)
	}

	// 启用日志记录
	if d.Options != nil && d.Options.EnableLogging {
		fmt.Printf("[DeepSeek] 请求体: %s\n", string(jsonData))
		logger.Errorf(ctx, "[DeepSeek] 发送请求到:%s 请求体: %s\n", d.BaseURL, string(jsonData))
	}

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var apiResp DeepSeekAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 记录响应日志
	jsonData, err = json.Marshal(apiResp)
	if err != nil {
		logger.Errorf(ctx, "[DeepSeek] body 序列化失败")
		return nil, err
	}
	logger.Infof(ctx, "[DeepSeek] body : %s", string(jsonData))

	// 检查错误
	if apiResp.Error != nil {
		return nil, fmt.Errorf("DeepSeek API错误: %s - %s", apiResp.Error.Code, apiResp.Error.Message)
	}

	// 提取回答内容
	if len(apiResp.Choices) == 0 {
		return nil, fmt.Errorf("响应格式错误：没有找到choices")
	}

	content := apiResp.Choices[0].Message.Content
	if content == "" {
		return nil, fmt.Errorf("响应格式错误：content为空")
	}

	// 提取使用统计
	var usage *Usage
	if apiResp.Usage != nil {
		usage = &Usage{
			PromptTokens:     int(apiResp.Usage.PromptTokens),
			CompletionTokens: int(apiResp.Usage.CompletionTokens),
			TotalTokens:      int(apiResp.Usage.TotalTokens),
		}
	}

	// 启用日志记录
	if d.Options != nil && d.Options.EnableLogging {
		fmt.Printf("[DeepSeek] 响应成功，内容长度: %d\n", len(content))
		logger.Infof(ctx, "[DeepSeek] 响应成功，:%s 内容长度: %d\n", string(content), len(content))
	}

	return &ChatResponse{
		Content: content,
		Usage:   usage,
	}, nil
}
