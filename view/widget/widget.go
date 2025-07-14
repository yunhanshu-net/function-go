package widget

import (
	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/pkg/x/stringsx"
	"github.com/yunhanshu-net/pkg/x/tagx"
)

func NewWidget(info *tagx.RunnerFieldInfo, renderType string) (Widget, error) {
	if info.Tags == nil {
		info.Tags = make(map[string]string)
	}
	widgetType := stringsx.DefaultString(info.Tags["widget"], "input")

	switch renderType {
	case response.RenderTypeTable:
		switch widgetType {
		case WidgetInput:
			return NewInputWidget(info)
		case WidgetNumber:
			return NewNumberWidget(info)
		case WidgetCheckbox:
			return newCheckboxWidget(info)
		case WidgetRadio:
			return newRadioWidget(info)
		case WidgetSelect:
			return NewSelectWidget(info)
		case WidgetSwitch:
			return newSwitchWidget(info)
		case WidgetSlider:
			return newSliderWidget(info)
		case WidgetColor:
			return newColorWidget(info)
		case WidgetDateTime:
			return newDateTimeWidget(info)
		case WidgetMultiSelect:
			return newMultiSelectWidget(info)
		case WidgetTag:
			return newTagWidget(info)
		case WidgetFileUpload:
			return newFileUploadWidget(info)
		case WidgetFileDisplay:
			return newFileDisplayWidget(info)
		}
	case response.RenderTypeForm:
		switch widgetType {
		case WidgetInput:
			return NewInputWidget(info)
		case WidgetNumber:
			return NewNumberWidget(info)
		case WidgetCheckbox:
			return newCheckboxWidget(info)
		case WidgetRadio:
			return newRadioWidget(info)
		case WidgetSelect:
			return NewSelectWidget(info)
		case WidgetSwitch:
			return newSwitchWidget(info)
		case WidgetSlider:
			return newSliderWidget(info)
		case WidgetColor:
			return newColorWidget(info)
		case WidgetDateTime:
			return newDateTimeWidget(info)
		case WidgetMultiSelect:
			return newMultiSelectWidget(info)
		case WidgetTag:
			return newTagWidget(info)
		case WidgetFileUpload:
			return newFileUploadWidget(info)
		case WidgetFileDisplay:
			return newFileDisplayWidget(info)
		}
		return NewInputWidget(info)
	}

	return NewInputWidget(info)
}

// parseInt 简单的字符串转整数函数
func parseInt(s string) int {
	if s == "" {
		return 0
	}

	result := 0
	for _, r := range s {
		if r >= '0' && r <= '9' {
			result = result*10 + int(r-'0')
		} else {
			return 0 // 非数字字符，返回0
		}
	}
	return result
}
