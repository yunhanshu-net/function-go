package api

import (
	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/pkg/x/stringsx"
)

func NewRequestParams(el interface{}, renderType string) (interface{}, error) {
	renderType = stringsx.DefaultString(renderType, response.RenderTypeForm)
	switch renderType {
	case response.RenderTypeTable:
		return NewTableRequestParams(el)
	case response.RenderTypeForm:
		return NewFormRequestParams(el, renderType)
	default:
		return NewFormRequestParams(el, renderType)
	}
}
