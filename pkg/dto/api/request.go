package api

import (
	"reflect"

	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/pkg/x/stringsx"
)

func NewRequestParams(el interface{}, renderType string) (interface{}, error) {
	builder := NewFormBuilder()

	renderType = stringsx.DefaultString(renderType, response.RenderTypeForm)
	switch renderType {
	case response.RenderTypeTable:
		return builder.BuildTableConfig(reflect.TypeOf(el))
	case response.RenderTypeForm:
		return builder.BuildFormConfig(reflect.TypeOf(el), renderType)
	default:
		return builder.BuildFormConfig(reflect.TypeOf(el), renderType)
	}
}
