package api

import (
	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/pkg/x/stringsx"
)

func NewResponseParams(el interface{}, renderType string) (interface{}, error) {
	renderType = stringsx.DefaultString(renderType, response.RenderTypeForm)
	switch renderType {
	case response.RenderTypeForm:
		return NewFormResponseParams(el)
	case response.RenderTypeTable:
		return NewTableResponseParams(el)
	default:
		return NewFormResponseParams(el)
	}
}
