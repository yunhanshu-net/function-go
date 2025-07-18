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

	//// 交互配置
	//Disabled bool `json:"disabled"` // 是否禁用
}

// newMultiSelectWidget 创建多选择器组件
func newMultiSelectWidget(info *tagx.RunnerFieldInfo) (Widget, error) {
	multiSelect := &MultiSelectWidget{
		MultipleLimit: 0,     // 默认不限制
		CollapseTags:  false, // 默认不折叠标签
		AllowCreate:   false, // 默认不允许创建
	}

	if info.Tags == nil {
		return multiSelect, nil
	}

	tag := info.Tags

	// 解析选项：直接存储原始字符串，前端负责解析 value(label) 格式
	if optionsStr, ok := tag["options"]; ok && optionsStr != "" {
		// 按逗号分割选项
		parts := strings.Split(optionsStr, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				multiSelect.Options = append(multiSelect.Options, part)
			}
		}
	}

	// 解析默认值：支持逗号分隔的多个值
	if defaultValue, ok := tag["default_value"]; ok && defaultValue != "" {
		parts := strings.Split(defaultValue, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				multiSelect.DefaultValue = append(multiSelect.DefaultValue, part)
			}
		}
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

	//// 设置是否禁用
	//if disabled, ok := tag["disabled"]; ok {
	//	multiSelect.Disabled = disabled == "true"
	//}

	return multiSelect, nil
}

func (w *MultiSelectWidget) GetWidgetType() string {
	return WidgetMultiSelect
}
