package widget

import (
	"strconv"
	"strings"

	"github.com/yunhanshu-net/pkg/x/tagx"
)

// SliderWidget 滑块组件
type SliderWidget struct {
	// 数值配置
	Min  float64 `json:"min"`  // 最小值
	Max  float64 `json:"max"`  // 最大值
	Step float64 `json:"step"` // 步长

	// 默认值配置
	DefaultValue interface{} `json:"default_value"` // 默认值，支持单值或范围值

	// 显示配置
	ShowMarks   bool   `json:"show_marks"`   // 是否显示刻度标记
	ShowTooltip bool   `json:"show_tooltip"` // 是否显示tooltip
	Unit        string `json:"unit"`         // 单位显示

	// 范围选择配置
	Range bool `json:"range"` // 是否为范围选择

	// 交互配置
	//Disabled bool `json:"disabled"` // 是否禁用
	Vertical bool `json:"vertical"` // 是否垂直显示
}

// newSliderWidget 创建滑块组件
func newSliderWidget(info *tagx.RunnerFieldInfo) (Widget, error) {
	slider := &SliderWidget{
		Min:         0,     // 默认最小值
		Max:         100,   // 默认最大值
		Step:        1,     // 默认步长
		ShowMarks:   false, // 默认不显示刻度
		ShowTooltip: true,  // 默认显示tooltip
		Range:       false, // 默认单值选择
		Vertical:    false, // 默认水平显示
	}

	if info.Tags == nil {
		return slider, nil
	}

	tag := info.Tags

	// 设置最小值
	if minStr, ok := tag["min"]; ok && minStr != "" {
		if min, err := strconv.ParseFloat(minStr, 64); err == nil {
			slider.Min = min
		}
	}

	// 设置最大值
	if maxStr, ok := tag["max"]; ok && maxStr != "" {
		if max, err := strconv.ParseFloat(maxStr, 64); err == nil {
			slider.Max = max
		}
	}

	// 设置步长
	if stepStr, ok := tag["step"]; ok && stepStr != "" {
		if step, err := strconv.ParseFloat(stepStr, 64); err == nil {
			slider.Step = step
		}
	}

	// 设置默认值
	if defaultValue, ok := tag["default_value"]; ok && defaultValue != "" {
		// 检查是否是范围值（包含逗号分隔）
		if strings.Contains(defaultValue, ",") {
			slider.Range = true
			slider.DefaultValue = defaultValue // 范围值保持字符串格式
		} else {
			// 单值，尝试解析为数字
			if val, err := strconv.ParseFloat(defaultValue, 64); err == nil {
				slider.DefaultValue = val
			} else {
				slider.DefaultValue = defaultValue
			}
		}
	}

	// 设置是否显示刻度
	if showMarks, ok := tag["show_marks"]; ok {
		slider.ShowMarks = showMarks == "true"
	}

	// 设置是否显示tooltip
	if showTooltip, ok := tag["show_tooltip"]; ok {
		slider.ShowTooltip = showTooltip == "true"
	}

	// 设置单位
	if unit, ok := tag["unit"]; ok && unit != "" {
		slider.Unit = strings.TrimSpace(unit)
	}

	// 设置是否为范围选择
	if rangeMode, ok := tag["range"]; ok {
		slider.Range = rangeMode == "true"
	}

	// 设置是否垂直显示
	if vertical, ok := tag["vertical"]; ok {
		slider.Vertical = vertical == "true"
	}

	//// 设置是否禁用
	//if disabled, ok := tag["disabled"]; ok {
	//	slider.Disabled = disabled == "true"
	//}

	return slider, nil
}

func (w *SliderWidget) GetWidgetType() string {
	return WidgetSlider
}
