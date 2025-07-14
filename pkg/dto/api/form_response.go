package api

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/yunhanshu-net/function-go/pkg/dto/response"
)

type FormResponseParamInfo struct {
	//英文标识
	Code string `json:"code"`
	//中文名称
	Name string `json:"name"`
	//中文介绍
	Desc string `json:"desc"`
	//是否必填
	Required bool `json:"required"`

	Callbacks    string      `json:"callbacks"`
	Validates    string      `json:"validates"`
	WidgetConfig interface{} `json:"widget_config"` //这里是widget.Widget类型的接口
	WidgetType   string      `json:"widget_type"`
	ValueType    string      `json:"value_type"`
	Example      string      `json:"example"`
}

type FormResponseParams struct {
	RenderType string                   `json:"render_type"`
	Children   []*FormResponseParamInfo `json:"children"`
}

func (p *FormResponseParams) JSONRawMessage() (json.RawMessage, error) {
	marshal, err := json.Marshal(p)
	if err != nil {
		return json.RawMessage("{}"), err
	}
	return marshal, nil
}

// newFormResponseParamInfo 使用新的FieldInfo创建响应参数信息
func newFormResponseParamInfo(fieldInfo *FieldInfo) (*FormResponseParamInfo, error) {
	// 处理回调配置
	var callbacks []string
	for _, callback := range fieldInfo.Callbacks {
		callbacks = append(callbacks, callback.Event)
	}

	param := &FormResponseParamInfo{
		Code:         fieldInfo.Code,
		Name:         fieldInfo.Name,
		Desc:         fieldInfo.Desc,
		Required:     fieldInfo.IsRequired(),
		Validates:    fieldInfo.Validation,
		Callbacks:    strings.Join(callbacks, ";"),
		WidgetConfig: fieldInfo.Widget.Config,
		WidgetType:   fieldInfo.Widget.Type,
		ValueType:    fieldInfo.Data.Type,
		Example:      fieldInfo.Data.Example,
	}

	return param, nil
}

func NewFormResponseParams(el interface{}) (*FormResponseParams, error) {
	rspType := reflect.TypeOf(el)
	if rspType.Kind() == reflect.Pointer {
		rspType = rspType.Elem()
	}
	if rspType.Kind() != reflect.Struct && rspType.Kind() != reflect.Slice {
		return nil, fmt.Errorf("输出参数仅支持Struct和Slice类型")
	}

	// 使用新的FormBuilder构建表单配置
	builder := NewFormBuilder()
	var formConfig *FormConfig
	var err error

	if rspType.Kind() == reflect.Struct {
		formConfig, err = builder.BuildFormConfig(rspType, response.RenderTypeForm)
		if err != nil {
			return nil, err
		}
	} else {
		// 处理切片类型，获取元素类型
		elemType := rspType.Elem()
		if elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
		}
		formConfig, err = builder.BuildFormConfig(elemType, response.RenderTypeForm)
		if err != nil {
			return nil, err
		}
	}

	var children = make([]*FormResponseParamInfo, 0, len(formConfig.Fields))

	for _, fieldInfo := range formConfig.Fields {
		info, err := newFormResponseParamInfo(fieldInfo)
		if err != nil {
			return nil, err
		}
		children = append(children, info)
	}

	return &FormResponseParams{
		RenderType: response.RenderTypeForm,
		Children:   children,
	}, nil
}
