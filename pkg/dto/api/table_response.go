package api

import (
	"fmt"
	"reflect"
)

// NewTableResponseParams 创建表格响应参数，返回与Request一致的结构
func NewTableResponseParams(el interface{}) (*TableConfig, error) {
	rspType := reflect.TypeOf(el)
	if rspType.Kind() == reflect.Pointer {
		rspType = rspType.Elem()
	}

	if rspType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("输出参数仅支持Struct类型")
	}

	// 查找Items字段
	var itemsFieldType reflect.Type
	for i := 0; i < rspType.NumField(); i++ {
		field := rspType.Field(i)
		if field.Name == "Items" {
			itemsFieldType = field.Type
			break
		}
	}

	if itemsFieldType == nil {
		return nil, fmt.Errorf("not found items field")
	}

	// 使用新的FormBuilder构建表格配置
	builder := NewFormBuilder()
	tableConfig, err := builder.BuildTableConfig(itemsFieldType)
	if err != nil {
		return nil, err
	}

	// 直接返回TableConfig，与Request保持一致
	return tableConfig, nil
}
