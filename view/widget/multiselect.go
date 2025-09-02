package widget

import (
	"strings"

	"github.com/yunhanshu-net/pkg/x/tagx"
)

// MultiSelectWidget 多选择器组件
type MultiSelectWidget struct {
	// 选项配置
	Options      []string `json:"options"`       // 选项列表，支持 value(label) 格式
	DefaultValue []string `json:"default_value"` // 默认选中值

	// 多选配置
	MultipleLimit int `json:"multiple_limit"` // 最多选择数量，0为不限制

	// 显示配置
	Placeholder  string `json:"placeholder"`   // 占位符文本
	CollapseTags bool   `json:"collapse_tags"` // 是否折叠显示已选择的标签

	// 创建配置
	AllowCreate bool `json:"allow_create"` // 是否允许创建新条目

	// 回调配置
	Callback string `json:"callback"` // 回调配置字符串

	// 拼接配置（当后端字段是 string 类型时，如何将多选值拼接为字符串）
	// 默认用英文逗号 ","，可通过标签 separator:"/" 覆盖
	Separator string `json:"separator"`

	//// 交互配置
	//Disabled bool `json:"disabled"` // 是否禁用
}

// newMultiSelectWidget 创建多选择器组件
func newMultiSelectWidget(info *tagx.RunnerFieldInfo) (Widget, error) {
	multiSelect := &MultiSelectWidget{
		MultipleLimit: 0,     // 默认不限制
		CollapseTags:  false, // 默认不折叠标签
		AllowCreate:   false, // 默认不允许创建
		Separator:     ",",   // 默认使用逗号分隔
	}

	if info.Tags == nil {
		return multiSelect, nil
	}

	tag := info.Tags

	// 解析选项：支持转义的选项字符串
	if optionsStr, ok := tag["options"]; ok && optionsStr != "" {
		multiSelect.Options = parseOptionsWithEscape(optionsStr)
	}

	// 解析默认值：支持转义的默认值字符串
	if defaultValue, ok := tag["default_value"]; ok && defaultValue != "" {
		multiSelect.DefaultValue = parseOptionsWithEscape(defaultValue)
	}

	// 设置选择数量限制
	if limit, ok := tag["multiple_limit"]; ok {
		if limitNum := parseInt(limit); limitNum > 0 {
			multiSelect.MultipleLimit = limitNum
		}
	}

	// 设置占位符
	if placeholder, ok := tag["placeholder"]; ok && placeholder != "" {
		multiSelect.Placeholder = strings.TrimSpace(placeholder)
	}

	// 设置是否折叠标签
	if collapseTags, ok := tag["collapse_tags"]; ok {
		multiSelect.CollapseTags = collapseTags == "true"
	}

	// 设置是否允许创建
	if allowCreate, ok := tag["allow_create"]; ok {
		multiSelect.AllowCreate = allowCreate == "true"
	}

	// 设置回调配置
	if callback, ok := tag["callback"]; ok && callback != "" {
		multiSelect.Callback = strings.TrimSpace(callback)
	}

	// 设置拼接分隔符（当后端字段为 string 类型时，前端提交按该分隔符拼接）
	if sep, ok := tag["separator"]; ok && sep != "" {
		multiSelect.Separator = sep
	}

	//// 设置是否禁用
	//if disabled, ok := tag["disabled"]; ok {
	//	multiSelect.Disabled = disabled == "true"
	//}

	return multiSelect, nil
}

func (w *MultiSelectWidget) GetWidgetType() string {
	return WidgetMultiSelect
}
