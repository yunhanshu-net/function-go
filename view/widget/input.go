package widget

import (
	"fmt"
	"github.com/yunhanshu-net/function-go/view/widget/types"
	"github.com/yunhanshu-net/pkg/x/stringsx"
	"github.com/yunhanshu-net/pkg/x/tagx"
)

// InputWidget 输入框组件
type InputWidget struct {
	//会把该字段渲染成什么组件，默认：input

	//line_text(默认),text_area
	Mode string `json:"mode"` //line_text / text_area / password
	//占位符（文本框提示信息）
	Placeholder string `json:"placeholder"`
	//默认值
	DefaultValue string `json:"default_value"`
}

// NewInputWidget 创建输入框组件
func NewInputWidget(info *tagx.RunnerFieldInfo) (*InputWidget, error) {
	valueType := info.GetValueType()
	if !types.IsValueType(valueType) {
		return nil, fmt.Errorf("不是合法的值类型：%s", valueType)
	}
	if info == nil {
		return nil, fmt.Errorf("nil")
	}
	if info.Tags == nil {
		info.Tags = map[string]string{}
	}

	input := &InputWidget{
		Mode:         stringsx.DefaultString(info.Tags["mode"], "line_text"),
		Placeholder:  info.Tags["placeholder"],
		DefaultValue: info.Tags["default_value"],
	}
	return input, nil
}

func (w *InputWidget) GetWidgetType() string {
	return WidgetInput
}
