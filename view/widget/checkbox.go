package widget

import (
	"strings"

	"github.com/yunhanshu-net/pkg/x/tagx"
)

// CheckboxWidget 多选框组件
type CheckboxWidget struct {
	// 选项配置
	Options      []string `json:"options"`       // 选项列表，支持 value(label) 格式
	DefaultValue []string `json:"default_value"` // 默认选中值

	// 多选配置
	MultipleLimit int `json:"multiple_limit"` // 最多选择数量，0为不限制

	// 布局配置
	//Direction string `json:"direction"` // 排列方向：horizontal/vertical，默认vertical
	//Columns   int    `json:"columns"`   // 列数（grid布局），仅当direction为vertical时有效

	// 交互配置
	ShowSelectAll bool `json:"show_select_all"` // 是否显示全选/反选按钮
	//Disabled      bool `json:"disabled"`        // 是否禁用整个组件
}

// newCheckboxWidget 创建多选框组件
func newCheckboxWidget(info *tagx.RunnerFieldInfo) (Widget, error) {
	checkbox := &CheckboxWidget{
		//Direction:     "vertical", // 默认垂直排列
		MultipleLimit: 0, // 默认不限制
		//Columns:       1,          // 默认1列
	}

	if info.Tags == nil {
		return checkbox, nil
	}

	tag := info.Tags

	// 解析选项：直接存储原始字符串，前端负责解析 value(label) 格式
	if optionsStr, ok := tag["options"]; ok && optionsStr != "" {
		// 按逗号分割选项
		parts := strings.Split(optionsStr, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				checkbox.Options = append(checkbox.Options, part)
			}
		}
	}

	// 解析默认值：支持逗号分隔的多个值
	if defaultValue, ok := tag["default_value"]; ok && defaultValue != "" {
		checkbox.DefaultValue = strings.Split(defaultValue, ",")
		// 去除空格
		for i, v := range checkbox.DefaultValue {
			checkbox.DefaultValue[i] = strings.TrimSpace(v)
		}
	}

	//// 设置排列方向
	//if direction, ok := tag["direction"]; ok {
	//	if direction == "horizontal" || direction == "vertical" {
	//		checkbox.Direction = direction
	//	}
	//}

	//// 设置列数
	//if columns, ok := tag["columns"]; ok {
	//	if cols := parseInt(columns); cols > 0 {
	//		checkbox.Columns = cols
	//	}
	//}

	// 设置选择数量限制
	if limit, ok := tag["multiple_limit"]; ok {
		if limitNum := parseInt(limit); limitNum > 0 {
			checkbox.MultipleLimit = limitNum
		}
	}

	// 设置是否显示全选按钮
	if showSelectAll, ok := tag["show_select_all"]; ok {
		checkbox.ShowSelectAll = showSelectAll == "true"
	}
	//
	//// 设置是否禁用
	//if disabled, ok := tag["disabled"]; ok {
	//	checkbox.Disabled = disabled == "true"
	//}

	return checkbox, nil
}

func (w *CheckboxWidget) GetWidgetType() string {
	return WidgetCheckbox
}
