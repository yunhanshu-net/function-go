package llms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yunhanshu-net/pkg/logger"
)

// GLMThinkingConfig GLM思考模式配置
type GLMThinkingConfig struct {
	Type string `json:"type"` // enabled 或 disabled
}

// GLMResponseFormat GLM响应格式配置
type GLMResponseFormat struct {
	Type string `json:"type"` // text 或 json_object
}

// GLMAPIRequest GLM API请求结构体
type GLMAPIRequest struct {
	Model          string             `json:"model"`
	Messages       []Message          `json:"messages"`
	MaxTokens      int                `json:"max_tokens,omitempty"`
	Temperature    float64            `json:"temperature,omitempty"`
	TopP           float64            `json:"top_p,omitempty"`
	DoSample       bool               `json:"do_sample,omitempty"`
	Stream         bool               `json:"stream,omitempty"`
	Thinking       *GLMThinkingConfig `json:"thinking,omitempty"`
	ResponseFormat *GLMResponseFormat `json:"response_format,omitempty"`
}

// GLMAPIResponse GLM API响应结构体
type GLMAPIResponse struct {
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
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage,omitempty"`
}

// GLMClient GLM客户端实现
type GLMClient struct {
	APIKey  string         `json:"api_key"`
	BaseURL string         `json:"base_url"`
	Options *ClientOptions `json:"options"`
	Model   string         `json:"model"` // 模型名称
}

// NewGLMClient 创建GLM客户端（保持向后兼容）
func NewGLMClient(apiKey string) *GLMClient {
	return NewGLMClientWithOptions(apiKey, DefaultClientOptions())
}

// NewGLMClientWithOptions 创建带配置的GLM客户端
func NewGLMClientWithOptions(apiKey string, options *ClientOptions) *GLMClient {
	if options == nil {
		options = DefaultClientOptions()
	}

	baseURL := options.BaseURL
	if baseURL == "" {
		baseURL = "https://open.bigmodel.cn/api/paas/v4/chat/completions"
	}

	return &GLMClient{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Options: options,
		Model:   "glm-4.5", // 设置默认模型为GLM-4.5
	}
}

// SetModel 设置模型名称
func (g *GLMClient) SetModel(model string) {
	g.Model = model
}

// GetModelName 获取模型名称
func (g *GLMClient) GetModelName() string {
	return g.Model
}

// GetProvider 获取提供商名称
func (g *GLMClient) GetProvider() string {
	return "GLM"
}

// Chat 实现LLMClient接口
func (g *GLMClient) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// 构造GLM API请求
	apiReq := GLMAPIRequest{
		Model:       req.Model,
		Messages:    req.Messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		Stream:      false, // 默认非流式
		Thinking: &GLMThinkingConfig{
			Type: "enabled", // 默认启用思考模式
		},
	}

	// 设置默认值
	if apiReq.Model == "" {
		apiReq.Model = g.Model
	}
	if apiReq.MaxTokens <= 0 {
		apiReq.MaxTokens = 4096
	}
	if apiReq.Temperature == 0 {
		apiReq.Temperature = 0.6
	}

	// 根据请求参数控制思考模式
	if req.UseThinking != nil {
		if *req.UseThinking {
			apiReq.Thinking.Type = "enabled"
		} else {
			apiReq.Thinking.Type = "disabled"
		}
	}

	// 动态创建HTTP客户端，支持请求级别的超时配置
	timeout := g.Options.Timeout // 默认使用客户端配置的超时时间
	if req.Timeout != nil && *req.Timeout > 0 {
		timeout = *req.Timeout // 如果请求中指定了超时时间，则使用请求的超时时间
	}

	httpClient := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			MaxIdleConns:       g.Options.MaxIdleConns,
			IdleConnTimeout:    g.Options.IdleConnTimeout,
			DisableCompression: true,
		},
	}

	// 发送HTTP请求
	jsonData, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", g.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+g.APIKey)

	// 设置自定义User-Agent
	if g.Options != nil && g.Options.UserAgent != "" {
		httpReq.Header.Set("User-Agent", g.Options.UserAgent)
	}

	// 启用日志记录
	if g.Options != nil && g.Options.EnableLogging {
		fmt.Printf("[GLM] 请求体: %s\n", string(jsonData))
		logger.Errorf(ctx, "[GLM] 发送请求到:%s 请求体: %s\n", g.BaseURL, string(jsonData))
	}

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var apiResp GLMAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 记录响应日志
	jsonData, err = json.Marshal(apiResp)
	if err != nil {
		logger.Errorf(ctx, "[GLM] body 序列化失败")
		return nil, err
	}
	logger.Infof(ctx, "[GLM] body : %s", string(jsonData))

	// 检查错误
	if apiResp.Error != nil {
		return nil, fmt.Errorf("GLM API错误: %s - %s", apiResp.Error.Code, apiResp.Error.Message)
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
			PromptTokens:     apiResp.Usage.PromptTokens,
			CompletionTokens: apiResp.Usage.CompletionTokens,
			TotalTokens:      apiResp.Usage.TotalTokens,
		}
	}

	// 启用日志记录
	if g.Options != nil && g.Options.EnableLogging {
		fmt.Printf("[GLM] 响应成功，内容长度: %d\n", len(content))
		logger.Infof(ctx, "[GLM] 响应成功，:%s 内容长度: %d\n", string(content), len(content))
	}

	return &ChatResponse{
		Content: content,
		Usage:   usage,
	}, nil
}

// ChatWithThinking 带思考模式的聊天（GLM特有功能）
func (g *GLMClient) ChatWithThinking(ctx context.Context, req *ChatRequest, enableThinking bool) (*ChatResponse, error) {
	// 构造GLM API请求
	apiReq := GLMAPIRequest{
		Model:       req.Model,
		Messages:    req.Messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		Stream:      false, // 默认非流式
		Thinking: &GLMThinkingConfig{
			Type: "enabled", // 默认启用思考模式
		},
	}

	// 设置默认值
	if apiReq.Model == "" {
		apiReq.Model = g.Model
	}
	if apiReq.MaxTokens <= 0 {
		apiReq.MaxTokens = 4096
	}
	if apiReq.Temperature == 0 {
		apiReq.Temperature = 0.6
	}

	// 根据参数控制思考模式
	if !enableThinking {
		apiReq.Thinking.Type = "disabled"
	}

	// 动态创建HTTP客户端，支持请求级别的超时配置
	timeout := g.Options.Timeout // 默认使用客户端配置的超时时间
	if req.Timeout != nil && *req.Timeout > 0 {
		timeout = *req.Timeout // 如果请求中指定了超时时间，则使用请求的超时时间
	}

	httpClient := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			MaxIdleConns:       g.Options.MaxIdleConns,
			IdleConnTimeout:    g.Options.IdleConnTimeout,
			DisableCompression: true,
		},
	}

	// 发送HTTP请求
	jsonData, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", g.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+g.APIKey)

	// 设置自定义User-Agent
	if g.Options != nil && g.Options.UserAgent != "" {
		httpReq.Header.Set("User-Agent", g.Options.UserAgent)
	}

	// 启用日志记录
	if g.Options != nil && g.Options.EnableLogging {
		fmt.Printf("[GLM] 思考模式请求体: %s\n", string(jsonData))
		logger.Errorf(ctx, "[GLM] 发送思考模式请求到:%s 请求体: %s\n", g.BaseURL, string(jsonData))
	}

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var apiResp GLMAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 记录响应日志
	jsonData, err = json.Marshal(apiResp)
	if err != nil {
		logger.Errorf(ctx, "[GLM] 思考模式body 序列化失败")
		return nil, err
	}
	logger.Infof(ctx, "[GLM] 思考模式body : %s", string(jsonData))

	// 检查错误
	if apiResp.Error != nil {
		return nil, fmt.Errorf("GLM API错误: %s - %s", apiResp.Error.Code, apiResp.Error.Message)
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
			PromptTokens:     apiResp.Usage.PromptTokens,
			CompletionTokens: apiResp.Usage.CompletionTokens,
			TotalTokens:      apiResp.Usage.TotalTokens,
		}
	}

	// 启用日志记录
	if g.Options != nil && g.Options.EnableLogging {
		fmt.Printf("[GLM] 思考模式响应成功，内容长度: %d\n", len(content))
		logger.Infof(ctx, "[GLM] 思考模式响应成功，:%s 内容长度: %d\n", string(content), len(content))
	}

	return &ChatResponse{
		Content: content,
		Usage:   usage,
	}, nil
}

// GetSupportedModels 获取支持的模型列表
func (g *GLMClient) GetSupportedModels() []string {
	return []string{
		"glm-4.5",       // 最强大的推理模型，3550亿参数
		"glm-4.5-air",   // 高性价比轻量级强性能
		"glm-4.5-x",     // 高性能强推理极速响应
		"glm-4.5-airx",  // 轻量级强性能极速响应
		"glm-4.5-flash", // 免费高效多功能
	}
}

// IsThinkingEnabled 检查当前模型是否支持思考模式
func (g *GLMClient) IsThinkingEnabled() bool {
	// GLM-4.5系列都支持思考模式
	supportedModels := g.GetSupportedModels()
	for _, model := range supportedModels {
		if g.Model == model {
			return true
		}
	}
	return false
}
