package llms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// QwenStreamResponse 千问流式响应结构体
type QwenStreamResponse struct {
	Output struct {
		Text         string `json:"text"`
		FinishReason string `json:"finish_reason,omitempty"`
	} `json:"output"`
	Usage *struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage,omitempty"`
	Error *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// QwenClient 千问客户端实现
type QwenClient struct {
	APIKey  string         `json:"api_key"`
	BaseURL string         `json:"base_url"`
	Options *ClientOptions `json:"options"`
}

// NewQwenClient 创建千问客户端（保持向后兼容）
func NewQwenClient(apiKey string) *QwenClient {
	// 如果传入的apiKey为空，尝试从环境变量获取
	if apiKey == "" {
		apiKey = os.Getenv("QIANWEN_API_KEY")
	}
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

// ChatStream 实现流式聊天接口
func (q *QwenClient) ChatStream(ctx context.Context, req *ChatRequest) (<-chan *StreamChunk, error) {
	// 创建流式响应通道
	chunkChan := make(chan *StreamChunk, 10) // 缓冲通道，避免阻塞

	// 在goroutine中处理流式请求
	go func() {
		defer close(chunkChan)

		// 构造千问API请求 - 修复格式
		modelName := req.Model
		if modelName == "" {
			modelName = "qwen-turbo"
		}

		maxTokens := req.MaxTokens
		if maxTokens <= 0 {
			maxTokens = 4000
		}

		temperature := req.Temperature
		if temperature == 0 {
			temperature = 0.7
		}

		apiReq := map[string]interface{}{
			"model": modelName,
			"input": map[string]interface{}{
				"messages": req.Messages,
			},
			"parameters": map[string]interface{}{
				"max_tokens":  maxTokens,
				"temperature": temperature,
				"stream":      true, // 启用流式
			},
		}

		// 动态创建HTTP客户端，支持请求级别的超时配置
		timeout := q.Options.Timeout
		if req.Timeout != nil && *req.Timeout > 0 {
			timeout = *req.Timeout
		}

		httpClient := &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConns:       q.Options.MaxIdleConns,
				IdleConnTimeout:    q.Options.IdleConnTimeout,
				DisableCompression: true,
			},
		}

		// 序列化请求
		jsonData, err := json.Marshal(apiReq)
		if err != nil {
			chunkChan <- &StreamChunk{
				Error: fmt.Sprintf("序列化请求失败: %v", err),
				Done:  true,
			}
			return
		}

		// 创建HTTP请求
		httpReq, err := http.NewRequestWithContext(ctx, "POST", q.BaseURL, bytes.NewBuffer(jsonData))
		if err != nil {
			chunkChan <- &StreamChunk{
				Error: fmt.Sprintf("创建HTTP请求失败: %v", err),
				Done:  true,
			}
			return
		}

		// 设置请求头
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+q.APIKey)
		if q.Options.UserAgent != "" {
			httpReq.Header.Set("User-Agent", q.Options.UserAgent)
		}

		// 发送请求
		resp, err := httpClient.Do(httpReq)
		if err != nil {
			chunkChan <- &StreamChunk{
				Error: fmt.Sprintf("HTTP请求失败: %v", err),
				Done:  true,
			}
			return
		}
		defer resp.Body.Close()

		// 检查HTTP状态码
		if resp.StatusCode != http.StatusOK {
			chunkChan <- &StreamChunk{
				Error: fmt.Sprintf("HTTP请求失败，状态码: %d", resp.StatusCode),
				Done:  true,
			}
			return
		}

		// 解析流式响应
		decoder := json.NewDecoder(resp.Body)
		var finalUsage *Usage

		for {
			var streamResp QwenStreamResponse
			if err := decoder.Decode(&streamResp); err != nil {
				if err.Error() == "EOF" {
					// 流结束，发送最终的使用统计
					chunkChan <- &StreamChunk{
						Usage: finalUsage,
						Done:  true,
					}
					break
				}
				chunkChan <- &StreamChunk{
					Error: fmt.Sprintf("解析流式响应失败: %v", err),
					Done:  true,
				}
				return
			}

			// 检查错误
			if streamResp.Error != nil {
				chunkChan <- &StreamChunk{
					Error: fmt.Sprintf("千问API错误: %s - %s", streamResp.Error.Code, streamResp.Error.Message),
					Done:  true,
				}
				return
			}

			// 发送内容片段
			if streamResp.Output.Text != "" {
				chunkChan <- &StreamChunk{
					Content: streamResp.Output.Text,
					Done:    false,
				}
			}

			// 检查是否完成
			if streamResp.Output.FinishReason != "" {
				// 保存使用统计
				if streamResp.Usage != nil {
					finalUsage = &Usage{
						PromptTokens:     streamResp.Usage.InputTokens,
						CompletionTokens: streamResp.Usage.OutputTokens,
						TotalTokens:      streamResp.Usage.TotalTokens,
					}
				}

				// 发送完成信号
				chunkChan <- &StreamChunk{
					Usage: finalUsage,
					Done:  true,
				}
				break
			}
		}
	}()

	return chunkChan, nil
}

// GetModelName 获取模型名称
func (q *QwenClient) GetModelName() string {
	return "qwen-turbo"
}

// GetProvider 获取提供商名称
func (q *QwenClient) GetProvider() string {
	return "千问"
}
