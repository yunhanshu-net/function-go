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
		case WidgetListInput, WidgetList:
			return newListInputWidget(info)
		case WidgetFormInput, WidgetForm:
			return newFormWidget(info)
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
		case WidgetListInput, WidgetList:
			return newListInputWidget(info)
		case WidgetFormInput, WidgetForm:
			return newFormWidget(info)
		}
		return NewInputWidget(info)
	}

	return NewInputWidget(info)
}
