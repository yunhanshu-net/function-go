package widget

import (
	"strings"

	"github.com/yunhanshu-net/pkg/x/tagx"
)

// ListInputWidget 列表输入组件
type ListInputWidget struct {
	// 基础配置
	Placeholder string `json:"placeholder,omitempty"` // 占位符文本
	
	// 组件配置
	Config interface{} `json:"config,omitempty"` // 组件配置
	
	// 子字段配置 - 使用interface{}避免循环导入
	Fields interface{} `json:"fields,omitempty"` // 子字段定义
}

// newListInputWidget 创建列表输入组件
func newListInputWidget(info *tagx.RunnerFieldInfo) (Widget, error) {
	listInput := &ListInputWidget{
		Placeholder: "添加项目",
	}

	if info.Tags == nil {
		return listInput, nil
	}

	tag := info.Tags

	// 设置占位符
	if placeholder, ok := tag["placeholder"]; ok && placeholder != "" {
		listInput.Placeholder = strings.TrimSpace(placeholder)
	}

	return listInput, nil
}

func (w *ListInputWidget) GetValueType() string {
	return TypeListStruct
}

func (w *ListInputWidget) GetWidgetType() string {
	return WidgetListInput
}

// SetConfig 设置组件配置
func (w *ListInputWidget) SetConfig(config interface{}) {
	w.Config = config
}

// GetConfig 获取组件配置
func (w *ListInputWidget) GetConfig() interface{} {
	return w.Config
}

// SetFields 设置子字段配置
func (w *ListInputWidget) SetFields(fields interface{}) {
	w.Fields = fields
}

// GetFields 获取子字段配置
func (w *ListInputWidget) GetFields() interface{} {
	return w.Fields
} 