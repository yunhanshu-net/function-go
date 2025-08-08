package widget

import (
	"strings"

	"github.com/yunhanshu-net/pkg/x/tagx"
)

// FormWidget 表单组件
type FormWidget struct {
	// 基础配置
	Title string `json:"title,omitempty"` // 表单标题
	
	// 组件配置
	Config interface{} `json:"config,omitempty"` // 组件配置
	
	// 子字段配置 - 使用interface{}避免循环导入
	Fields interface{} `json:"fields,omitempty"` // 子字段定义
}

// newFormWidget 创建表单组件
func newFormWidget(info *tagx.RunnerFieldInfo) (Widget, error) {
	form := &FormWidget{}

	if info.Tags == nil {
		return form, nil
	}

	tag := info.Tags

	// 设置标题
	if title, ok := tag["title"]; ok && title != "" {
		form.Title = strings.TrimSpace(title)
	}

	return form, nil
}

func (w *FormWidget) GetValueType() string {
	return TypeStruct
}

func (w *FormWidget) GetWidgetType() string {
	return WidgetFormInput
}

// SetConfig 设置组件配置
func (w *FormWidget) SetConfig(config interface{}) {
	w.Config = config
}

// GetConfig 获取组件配置
func (w *FormWidget) GetConfig() interface{} {
	return w.Config
}

// SetFields 设置子字段配置
func (w *FormWidget) SetFields(fields interface{}) {
	w.Fields = fields
}

// GetFields 获取子字段配置
func (w *FormWidget) GetFields() interface{} {
	return w.Fields
} 