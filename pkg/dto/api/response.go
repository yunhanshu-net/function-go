package api

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/pkg/x/stringsx"
)

// UnifiedAPIResponse 统一的API响应结构
type UnifiedAPIResponse struct {
	RenderType string       `json:"render_type"` // 渲染类型：form, table
	Fields     []*FieldInfo `json:"fields"`      // 字段列表（form模式）
	Columns    []*FieldInfo `json:"columns"`     // 列配置（table模式）
}

// JSONRawMessage 返回JSON原始消息
func (r *UnifiedAPIResponse) JSONRawMessage() (json.RawMessage, error) {
	marshal, err := json.Marshal(r)
	if err != nil {
		return json.RawMessage("{}"), err
	}
	return marshal, nil
}

// NewUnifiedFormResponse 创建统一的表单响应
func NewUnifiedFormResponse(el interface{}, renderType string) (*UnifiedAPIResponse, error) {
	builder := NewFormBuilder()
	typeOf := reflect.TypeOf(el)

	formConfig, err := builder.BuildFormConfig(typeOf, renderType)
	if err != nil {
		return nil, err
	}

	return &UnifiedAPIResponse{
		RenderType: formConfig.RenderType,
		Fields:     formConfig.Fields,
	}, nil
}

// NewUnifiedTableResponse 创建统一的表格响应
func NewUnifiedTableResponse(el interface{}, functionInfo FunctionInfoInterface) (*UnifiedAPIResponse, error) {
	typeOf := reflect.TypeOf(el)

	// 处理指针类型
	if typeOf.Kind() == reflect.Pointer {
		typeOf = typeOf.Elem()
	}

	var itemsType reflect.Type

	// 如果直接传入的是切片类型，直接使用
	if typeOf.Kind() == reflect.Slice {
		itemsType = typeOf
	} else if typeOf.Kind() == reflect.Struct {
		// 如果是结构体，查找Items字段
		for i := 0; i < typeOf.NumField(); i++ {
			field := typeOf.Field(i)
			if field.Name == "Items" {
				itemsType = field.Type
				break
			}
		}

		if itemsType == nil {
			return nil, fmt.Errorf("not found items field in struct")
		}
	} else {
		return nil, fmt.Errorf("输入参数必须是包含Items字段的结构体或切片类型，当前类型: %s", typeOf.Kind())
	}

	builder := NewFormBuilder()
	builder.functionInfo = functionInfo
	tableConfig, err := builder.BuildTableConfig(itemsType)
	if err != nil {
		return nil, err
	}

	return &UnifiedAPIResponse{
		RenderType: tableConfig.RenderType,
		Columns:    tableConfig.Columns,
	}, nil
}

func NewResponseParams(el interface{}, renderType string, functionInfo FunctionInfoInterface) (interface{}, error) {
	renderType = stringsx.DefaultString(renderType, response.RenderTypeForm)
	switch renderType {
	case response.RenderTypeForm:
		return NewUnifiedFormResponse(el, renderType)
	case response.RenderTypeTable:
		return NewUnifiedTableResponse(el, functionInfo)
	default:
		return NewUnifiedFormResponse(el, renderType)
	}
}
