package widget

import (
	"fmt"
	"github.com/yunhanshu-net/function-go/view/widget/types"
	"github.com/yunhanshu-net/pkg/x/tagx"
)

// DateTimeWidget 日期时间选择器组件
type DateTimeWidget struct {
	// 组件类型，可选值：date(日期)、time(时间)、datetime(日期时间)、daterange(日期范围)、timerange(时间范围)
	Widget string `json:"widget"`
	// 数据类型，一般为time
	Type string `json:"type"`
	// 日期格式，如：yyyy-MM-dd、HH:mm:ss、yyyy-MM-dd HH:mm:ss
	Format string `json:"format,omitempty"`
	// 是否显示清除按钮
	Clearable bool `json:"clearable,omitempty"`
	// 占位符
	Placeholder string `json:"placeholder,omitempty"`
	// 默认值
	DefaultValue string `json:"default_value,omitempty"`
	// 最小日期/时间
	Min string `json:"min,omitempty"`
	// 最大日期/时间
	Max string `json:"max,omitempty"`
	// 是否禁用
	Disabled bool `json:"disabled,omitempty"`

	Readonly bool `json:"readonly"` // 是否只读

}

func (w *DateTimeWidget) GetValueType() string {
	return w.Type
}

func (w *DateTimeWidget) GetWidgetType() string {
	return w.Widget
}

// NewDateTimeWidget 创建时间组件
func NewDateTimeWidget(info *tagx.RunnerFieldInfo) (*DateTimeWidget, error) {
	valueType := info.GetValueType()
	if !types.IsValueType(valueType) {
		return nil, fmt.Errorf("不是合法的值类型：%s", valueType)
	}
	if info == nil {
		return nil, fmt.Errorf("<UNK>nil")
	}
	if info.Tags == nil {
		info.Tags = map[string]string{}
	}

	input := &DateTimeWidget{
		Placeholder:  info.Tags["placeholder"],
		DefaultValue: info.Tags["default_value"],
	}
	return input, nil
}
