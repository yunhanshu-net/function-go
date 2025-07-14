package widget

import (
	"strings"

	"github.com/yunhanshu-net/pkg/x/tagx"
)

// DateTimeWidget 日期时间组件
type DateTimeWidget struct {
	// 格式配置
	Kind string `json:"kind"` // 具体格式：date/datetime/time/daterange/datetimerange/month/year/week

	Format string `json:"format"` // 具体格式：yyyy-MM-dd HH:mm:ss 或 yyyy-MM-dd HH:mm:ss
	// 占位符配置
	Placeholder      string `json:"placeholder"`       // 占位符文本
	StartPlaceholder string `json:"start_placeholder"` // 范围选择时开始日期占位符
	EndPlaceholder   string `json:"end_placeholder"`   // 范围选择时结束日期占位符

	// 默认值和限制
	DefaultValue string `json:"default_value"` // 默认值，支持today、now等特殊值
	DefaultTime  string `json:"default_time"`  // 选中日期后的默认具体时刻
	MinDate      string `json:"min_date"`      // 最小可选日期
	MaxDate      string `json:"max_date"`      // 最大可选日期

	// 范围选择配置
	Separator string `json:"separator"` // 日期范围分隔符，默认"至"

	// 快捷选项
	Shortcuts []string `json:"shortcuts"` // 快捷选项配置

	// 交互配置
	Disabled bool `json:"disabled"` // 是否禁用
}

// newDateTimeWidget 创建日期时间组件
func newDateTimeWidget(info *tagx.RunnerFieldInfo) (Widget, error) {
	dateTime := &DateTimeWidget{
		Kind:      "date",                // 默认日期格式
		Format:    "yyyy-MM-dd HH:mm:ss", // 默认格式
		Separator: "至",                   // 默认分隔符
	}

	if info.Tags == nil {
		return dateTime, nil
	}

	tag := info.Tags

	// 设置格式
	if kind, ok := tag["kind"]; ok && kind != "" {
		// 验证格式是否有效
		validFormats := map[string]bool{
			"date": true, "datetime": true, "time": true,
			"daterange": true, "datetimerange": true,
			"month": true, "year": true, "week": true,
		}
		if validFormats[kind] {
			dateTime.Kind = kind
		}
	}

	// 设置占位符
	if placeholder, ok := tag["placeholder"]; ok && placeholder != "" {
		dateTime.Placeholder = strings.TrimSpace(placeholder)
	}

	// 设置format
	if format, ok := tag["format"]; ok && format != "" {
		dateTime.Format = strings.TrimSpace(format)
	}

	if startPlaceholder, ok := tag["start_placeholder"]; ok && startPlaceholder != "" {
		dateTime.StartPlaceholder = strings.TrimSpace(startPlaceholder)
	}

	if endPlaceholder, ok := tag["end_placeholder"]; ok && endPlaceholder != "" {
		dateTime.EndPlaceholder = strings.TrimSpace(endPlaceholder)
	}

	// 设置默认值
	if defaultValue, ok := tag["default_value"]; ok && defaultValue != "" {
		dateTime.DefaultValue = strings.TrimSpace(defaultValue)
	}

	if defaultTime, ok := tag["default_time"]; ok && defaultTime != "" {
		dateTime.DefaultTime = strings.TrimSpace(defaultTime)
	}

	// 设置日期限制
	if minDate, ok := tag["min_date"]; ok && minDate != "" {
		dateTime.MinDate = strings.TrimSpace(minDate)
	}

	if maxDate, ok := tag["max_date"]; ok && maxDate != "" {
		dateTime.MaxDate = strings.TrimSpace(maxDate)
	}

	// 设置分隔符
	if separator, ok := tag["separator"]; ok && separator != "" {
		dateTime.Separator = strings.TrimSpace(separator)
	}

	// 解析快捷选项
	if shortcuts, ok := tag["shortcuts"]; ok && shortcuts != "" {
		// 按逗号分割快捷选项
		parts := strings.Split(shortcuts, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				dateTime.Shortcuts = append(dateTime.Shortcuts, part)
			}
		}
	}

	// 设置是否禁用
	if disabled, ok := tag["disabled"]; ok {
		dateTime.Disabled = disabled == "true"
	}

	return dateTime, nil
}

func (w *DateTimeWidget) GetWidgetType() string {
	return WidgetDateTime
}
