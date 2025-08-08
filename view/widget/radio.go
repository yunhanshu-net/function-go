package widget

import (
	"strings"

	"github.com/yunhanshu-net/pkg/x/tagx"
)

// RadioWidget 单选框组件
type RadioWidget struct {
	// 选项配置
	Options      []string `json:"options"`       // 选项列表，支持 value(label) 格式
	DefaultValue string   `json:"default_value"` // 默认值

	// 布局配置
	//Direction string `json:"direction"` // 排列方向：horizontal/vertical，默认vertical

	// 交互配置
	//Disabled bool `json:"disabled"` // 是否禁用整个组件
}

// newRadioWidget 创建单选框组件
func newRadioWidget(info *tagx.RunnerFieldInfo) (Widget, error) {
	radio := &RadioWidget{
		//Direction: "vertical", // 默认垂直排列
	}

	if info.Tags == nil {
		return radio, nil
	}

	tag := info.Tags

	// 解析选项：支持转义的选项字符串
	if optionsStr, ok := tag["options"]; ok && optionsStr != "" {
		radio.Options = parseOptionsWithEscape(optionsStr)
	}

	// 解析默认值
	if defaultValue, ok := tag["default_value"]; ok {
		radio.DefaultValue = strings.TrimSpace(defaultValue)
	}

	//// 设置排列方向
	//if direction, ok := tag["direction"]; ok {
	//	if direction == "horizontal" || direction == "vertical" {
	//		radio.Direction = direction
	//	}
	//}
	//
	//// 设置是否禁用
	//if disabled, ok := tag["disabled"]; ok {
	//	radio.Disabled = disabled == "true"
	//}

	return radio, nil
}

func (w *RadioWidget) GetWidgetType() string {
	return WidgetRadio
}
