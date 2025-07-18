package api

import (
	"encoding/json"
	"strings"
)

// FieldInfo 字段信息结构体 - 重构后的核心数据结构
type FieldInfo struct {
	// 基础信息（来自runner标签）
	Code string `json:"code"` // 字段代码
	Name string `json:"name"` // 字段显示名称
	Desc string `json:"desc"` // 字段描述

	// Widget配置（来自widget标签）
	Widget WidgetConfig `json:"widget"`

	// 数据配置（来自data标签） - type字段迁移到这里
	Data DataConfig `json:"data"`

	// 权限配置（来自permission标签） - 简化为三权限（form模式下为null）
	Permission *PermissionConfig `json:"permission"`

	// 回调配置（来自callback标签） - 仅字段级别回调
	Callbacks []CallbackConfig `json:"callbacks"`

	// 验证配置（来自validate标签） - 简化为字符串
	Validation string `json:"validation"`

	// 搜索配置（来自search标签） - 可为nil表示不支持搜索
	Search *SearchConfig `json:"search"`
}

// WidgetConfig Widget配置 - 使用灵活的map结构
type WidgetConfig struct {
	Type   string                 `json:"type"`             // 组件类型：input, select, switch, checkbox等
	Config map[string]interface{} `json:"config"`           // 组件个性化配置
	Fields []*FieldInfo           `json:"fields,omitempty"` // 子字段配置（用于list_input和form组件）
}

// DataConfig 数据配置 - type字段迁移到这里
type DataConfig struct {
	Type         string `json:"type"`          // 数据类型：string, number, boolean, []string等
	Example      string `json:"example"`       // 示例值
	DefaultValue string `json:"default_value"` // 默认值
	Source       string `json:"source"`        // 数据源配置
	Format       string `json:"format"`        // 格式配置
}

// PermissionConfig 权限配置 - 简化为三权限
type PermissionConfig struct {
	Read   bool `json:"read"`   // 可读权限
	Update bool `json:"update"` // 可更新权限
	Create bool `json:"create"` // 可创建权限
}

// CallbackConfig 回调配置 - 仅字段级别回调
type CallbackConfig struct {
	Event  string            `json:"event"`  // 回调事件：OnInputFuzzy, OnBlur等
	Params map[string]string `json:"params"` // 回调参数
}

// SearchConfig 搜索配置 - 来自search标签
type SearchConfig struct {
	Operators []string `json:"operators"` // 支持的操作符：eq, like, in, gte, lte, gt, lt
}

// FunctionInfo 函数信息 - 用于配置函数级别的回调
type FunctionInfo struct {
	// 函数基本信息
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`

	// 函数级别的回调配置
	OnPageLoad *FunctionCallback `json:"on_page_load"` // 页面加载回调
	OnSubmit   *FunctionCallback `json:"on_submit"`    // 表单提交回调
	OnReset    *FunctionCallback `json:"on_reset"`     // 表单重置回调

	// 字段列表
	Fields []FieldInfo `json:"fields"`
}

// FunctionCallback 函数级别回调配置
type FunctionCallback struct {
	Enabled bool              `json:"enabled"` // 是否启用
	Params  map[string]string `json:"params"`  // 回调参数
}

// FormConfig 表单配置
type FormConfig struct {
	RenderType string       `json:"render_type"` // 渲染类型：form, table
	Fields     []*FieldInfo `json:"fields"`      // 字段列表
}

// TableConfig 表格配置
type TableConfig struct {
	RenderType string       `json:"render_type"` // 渲染类型
	Columns    []*FieldInfo `json:"columns"`     // 列配置
}

// JSONRawMessage 返回JSON原始消息
func (c *FormConfig) JSONRawMessage() (json.RawMessage, error) {
	marshal, err := json.Marshal(c)
	if err != nil {
		return json.RawMessage("{}"), err
	}
	return marshal, nil
}

func (c *TableConfig) JSONRawMessage() (json.RawMessage, error) {
	marshal, err := json.Marshal(c)
	if err != nil {
		return json.RawMessage("{}"), err
	}
	return marshal, nil
}

// MarshalJSON 自定义JSON序列化 - 移除权限特殊处理
func (f *FieldInfo) MarshalJSON() ([]byte, error) {
	type Alias FieldInfo
	// 直接序列化，不做权限的特殊处理
	// 权限为nil时会序列化为null，前端将其理解为无权限限制
	return json.Marshal((*Alias)(f))
}

// IsReadOnly 判断是否为只读字段
func (f *FieldInfo) IsReadOnly() bool {
	if f.Permission == nil {
		return false // form模式下没有权限控制
	}
	return f.Permission.Read && !f.Permission.Update && !f.Permission.Create
}

// IsCreateOnly 判断是否为仅创建字段
func (f *FieldInfo) IsCreateOnly() bool {
	if f.Permission == nil {
		return false // form模式下没有权限控制
	}
	return f.Permission.Create && !f.Permission.Update
}

// HasCallback 判断是否有指定的回调函数
func (f *FieldInfo) HasCallback(event string) bool {
	for _, callback := range f.Callbacks {
		if callback.Event == event {
			return true
		}
	}
	return false
}

// GetCallback 获取指定的回调函数配置
func (f *FieldInfo) GetCallback(event string) *CallbackConfig {
	for _, callback := range f.Callbacks {
		if callback.Event == event {
			return &callback
		}
	}
	return nil
}

// IsRequired 判断是否为必填字段
func (f *FieldInfo) IsRequired() bool {
	return strings.Contains(f.Validation, "required")
}

// IsSearchable 判断字段是否支持搜索
func (f *FieldInfo) IsSearchable() bool {
	return f.Search != nil && len(f.Search.Operators) > 0
}

// SupportsSearchOperator 判断字段是否支持指定的搜索操作符
func (f *FieldInfo) SupportsSearchOperator(operator string) bool {
	if f.Search == nil {
		return false
	}
	for _, op := range f.Search.Operators {
		if op == operator {
			return true
		}
	}
	return false
}

// GetSearchOperators 获取字段支持的搜索操作符列表
func (f *FieldInfo) GetSearchOperators() []string {
	if f.Search == nil {
		return []string{}
	}
	return f.Search.Operators
}
