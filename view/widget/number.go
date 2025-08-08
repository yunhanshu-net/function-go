package widget

import (
	"strconv"
	"strings"

	"github.com/yunhanshu-net/pkg/x/tagx"
)

// NumberWidget 数字输入组件
type NumberWidget struct {
	// 数值配置
	Min  float64 `json:"min"`  // 最小值
	Max  float64 `json:"max"`  // 最大值
	Step float64 `json:"step"` // 步长

	// 默认值配置
	DefaultValue interface{} `json:"default_value"` // 默认值

	// 显示配置
	Placeholder string `json:"placeholder"` // 占位符
	Unit        string `json:"unit"`        // 单位显示
	Prefix      string `json:"prefix"`      // 前缀（如¥）
	Suffix      string `json:"suffix"`      // 后缀（如%）

	// 精度配置
	Precision int `json:"precision"` // 小数位数精度

	// 输入模式配置
	AllowNegative bool `json:"allow_negative"` // 是否允许负数
	AllowDecimal  bool `json:"allow_decimal"`  // 是否允许小数
}

// NewNumberWidget 创建数字输入组件
func NewNumberWidget(info *tagx.RunnerFieldInfo) (*NumberWidget, error) {
	number := &NumberWidget{
		Min:           0,      // 默认最小值
		Max:           999999, // 默认最大值
		Step:          1,      // 默认步长
		Precision:     0,      // 默认整数
		AllowNegative: false,  // 默认不允许负数
		AllowDecimal:  false,  // 默认不允许小数
	}

	if info.Tags == nil {
		info.Tags = make(map[string]string)
	}

	tag := info.Tags

	// 设置最小值
	if minStr, ok := tag["min"]; ok && minStr != "" {
		if min, err := strconv.ParseFloat(minStr, 64); err == nil {
			number.Min = min
		}
	}

	// 设置最大值
	if maxStr, ok := tag["max"]; ok && maxStr != "" {
		if max, err := strconv.ParseFloat(maxStr, 64); err == nil {
			number.Max = max
		}
	}

	// 设置步长
	if stepStr, ok := tag["step"]; ok && stepStr != "" {
		if step, err := strconv.ParseFloat(stepStr, 64); err == nil {
			number.Step = step
		}
	}

	// 设置默认值
	if defaultValue, ok := tag["default_value"]; ok && defaultValue != "" {
		if val, err := strconv.ParseFloat(defaultValue, 64); err == nil {
			number.DefaultValue = val
		} else {
			number.DefaultValue = defaultValue
		}
	}

	// 设置占位符
	if placeholder, ok := tag["placeholder"]; ok && placeholder != "" {
		number.Placeholder = strings.TrimSpace(placeholder)
	}

	// 设置单位
	if unit, ok := tag["unit"]; ok && unit != "" {
		number.Unit = strings.TrimSpace(unit)
	}

	// 设置前缀
	if prefix, ok := tag["prefix"]; ok && prefix != "" {
		number.Prefix = strings.TrimSpace(prefix)
	}

	// 设置后缀
	if suffix, ok := tag["suffix"]; ok && suffix != "" {
		number.Suffix = strings.TrimSpace(suffix)
	}

	// 设置精度
	if precisionStr, ok := tag["precision"]; ok && precisionStr != "" {
		if precision, err := strconv.Atoi(precisionStr); err == nil {
			number.Precision = precision
			number.AllowDecimal = precision > 0 // 有精度就允许小数
		}
	}

	// 设置是否允许负数
	if allowNegative, ok := tag["allow_negative"]; ok {
		number.AllowNegative = allowNegative == "true"
		if number.AllowNegative && number.Min >= 0 {
			number.Min = -999999 // 允许负数时设置默认最小值
		}
	}

	// 设置是否允许小数
	if allowDecimal, ok := tag["allow_decimal"]; ok {
		number.AllowDecimal = allowDecimal == "true"
		if number.AllowDecimal && number.Precision == 0 {
			number.Precision = 2 // 允许小数时设置默认精度
		}
	}

	return number, nil
}

func (w *NumberWidget) GetWidgetType() string {
	return WidgetNumber
}
