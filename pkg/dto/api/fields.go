package api

import (
	"fmt"
	"reflect"
)

// GetFields 获取字段信息，使用新的多标签解析器
func GetFields(el interface{}) ([]*FieldInfo, error) {
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
		formConfig, err = builder.BuildFormConfig(rspType, "form")
		if err != nil {
			return nil, err
		}
	} else {
		// 处理切片类型，获取元素类型
		elemType := rspType.Elem()
		if elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
		}
		formConfig, err = builder.BuildFormConfig(elemType, "form")
		if err != nil {
			return nil, err
		}
	}

	return formConfig.Fields, nil
}
