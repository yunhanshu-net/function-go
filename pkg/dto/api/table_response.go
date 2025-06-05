package api

import (
	"fmt"
	"github.com/yunhanshu-net/function-go/view/widget"
	"github.com/yunhanshu-net/pkg/x/tagx"
)

type TableResponseParams struct {
	*widget.TableWidget
}

func NewTableResponseParams(el interface{}) (*TableResponseParams, error) {
	fields, err := GetFields(el)
	if err != nil {
		return nil, err
	}
	items := &tagx.RunnerFieldInfo{}
	for _, field := range fields {
		if field.Name == "Items" {
			items = field
			break
		}
	}
	if items == nil {
		return nil, fmt.Errorf("not found items field")
	}
	fields, err = GetFields(items.Value.Interface())
	if err != nil {
		return nil, err
	}

	//提取items字段
	table, err := widget.NewTable(fields)
	if err != nil {
		return nil, err
	}
	mp := make(map[string]interface{})
	for _, field := range fields {
		info, err := newFormRequestParamInfo(field)
		if err != nil {
			continue
		}
		mp[field.GetCode()] = info
	}
	for i, column := range table.Columns {
		v, ok := mp[column.Code]
		if !ok {
			continue
		}
		table.Columns[i].AddFormConfig = v
	}
	return &TableResponseParams{TableWidget: table}, nil
}
