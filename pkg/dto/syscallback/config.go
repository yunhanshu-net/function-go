package syscallback

import (
	"fmt"
	"strings"
)

// ConfigData 配置数据结构
type ConfigData struct {
	Type string      `json:"type,omitempty"` // 配置类型：json, yaml, toml, xml 等（可选，默认为json）
	Data interface{} `json:"data"`           // 配置数据（直接存储，避免双重序列化）
}

// ConfigUpdateRequest 配置更新请求
type ConfigUpdateRequest struct {
	Router     string      `json:"router"`     // 路由路径
	Method     string      `json:"method"`     // HTTP方法
	ConfigData *ConfigData `json:"config_data"` // 完整的修改后配置
}

// GenerateConfigKey 生成配置键
func (req *ConfigUpdateRequest) GenerateConfigKey() string {
	// 处理路由路径，将 / 替换为 . 以安全地用作配置键
	safeRouter := strings.ReplaceAll(req.Router, "/", ".")
	// 移除前后的点
	safeRouter = strings.Trim(safeRouter, ".")
	return fmt.Sprintf("function.%s.%s", safeRouter, req.Method)
}

// ConfigUpdateResponse 配置更新响应
type ConfigUpdateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// ConfigGetRequest 配置获取请求
type ConfigGetRequest struct {
	Router string `json:"router"` // 路由路径
	Method string `json:"method"` // HTTP方法
}

// GenerateConfigKey 生成配置键
func (req *ConfigGetRequest) GenerateConfigKey() string {
	// 处理路由路径，将 / 替换为 . 以安全地用作配置键
	safeRouter := strings.ReplaceAll(req.Router, "/", ".")
	// 移除前后的点
	safeRouter = strings.Trim(safeRouter, ".")
	return fmt.Sprintf("function.%s.%s", safeRouter, req.Method)
}

// ConfigGetResponse 配置获取响应
type ConfigGetResponse struct {
	Success bool        `json:"success"`
	Config  *ConfigData `json:"config,omitempty"`
	Error   string      `json:"error,omitempty"`
}
