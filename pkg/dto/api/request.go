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

// NewRequestParamsWithFunctionInfo 支持传入FunctionInfo来自动注入回调信息的版本
func NewRequestParamsWithFunctionInfo(el interface{}, renderType string, functionInfo FunctionInfoInterface) (interface{}, error) {
	builder := NewFormBuilderWithFunctionInfo(functionInfo)

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

// FunctionInfoInterface 定义FunctionInfo的接口，避免循环依赖
type FunctionInfoInterface interface {
	GetOnInputFuzzyMap() map[string]interface{}
	GetOnInputValidateMap() map[string]interface{}
}

//type FunctionWidgetCallback interface {
//
//	// WidgetCallbacks callbackType（OnInputFuzzy，OnInputValidate）
//
//	// WidgetCallbacks 返回的是key=字段名，value暂时好像没用，后面再说吧
//	WidgetCallbacks(callbackType string) map[string]interface{}
//}
