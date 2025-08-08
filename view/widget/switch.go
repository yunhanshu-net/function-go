package widget

import (
	"strings"

	"github.com/yunhanshu-net/pkg/x/tagx"
)

// SwitchWidget 开关组件
type SwitchWidget struct {
	// 显示配置
	TrueLabel  string `json:"true_label"`  // true状态的显示文本，默认"开启"
	FalseLabel string `json:"false_label"` // false状态的显示文本，默认"关闭"

	// 默认值配置
	DefaultValue bool `json:"default_value"` // 默认值，支持true/false，默认为false

	// 交互配置
	Disabled bool `json:"disabled"` // 是否禁用
}

// newSwitchWidget 创建开关组件
func newSwitchWidget(info *tagx.RunnerFieldInfo) (Widget, error) {
	switchWidget := &SwitchWidget{
		TrueLabel:    "开启",  // 默认true状态文本
		FalseLabel:   "关闭",  // 默认false状态文本
		DefaultValue: false, // 默认为关闭状态
	}

	if info.Tags == nil {
		return switchWidget, nil
	}

	tag := info.Tags

	// 设置true状态文本
	if trueLabel, ok := tag["true_label"]; ok && trueLabel != "" {
		switchWidget.TrueLabel = strings.TrimSpace(trueLabel)
	}

	// 设置false状态文本
	if falseLabel, ok := tag["false_label"]; ok && falseLabel != "" {
		switchWidget.FalseLabel = strings.TrimSpace(falseLabel)
	}

	// 设置默认值
	if defaultValue, ok := tag["default_value"]; ok {
		switchWidget.DefaultValue = defaultValue == "true"
	}

	// 设置是否禁用
	if disabled, ok := tag["disabled"]; ok {
		switchWidget.Disabled = disabled == "true"
	}

	return switchWidget, nil
}

func (w *SwitchWidget) GetWidgetType() string {
	return WidgetSwitch
}
