package widget

import (
	"github.com/pkg/errors"
	"github.com/yunhanshu-net/pkg/x/tagx"
)

// TableWidget 表格组件
type TableWidget struct {
	// 组件类型，固定为table
	Widget string `json:"widget"`
	// 表格列定义
	Columns []TableColumn `json:"columns,omitempty"`
}

// TableColumn 表格列定义 - 统一使用FieldInfo结构
type TableColumn struct {
	Code         string      `json:"code"`
	Name         string      `json:"name"`
	ValueType    string      `json:"value_type"`
	WidgetType   string      `json:"widget_type"`   // 组件类型：input/select/switch/datetime等
	WidgetConfig interface{} `json:"widget_config"` // 组件配置对象
	// 统一使用FieldInfo结构，而不是AddFormConfig
	FieldConfig interface{} `json:"field_config"` // 完整的字段配置信息
}

func (w *TableWidget) GetValueType() string {
	return TypeArray
}

func (w *TableWidget) GetWidgetType() string {
	return w.Widget
}

func NewTable(info []*tagx.RunnerFieldInfo) (*TableWidget, error) {
	if info == nil {
		return nil, errors.New("NewTable info ==nil")
	}

	Columns := make([]TableColumn, 0)
	for _, v := range info {
		// 为每个字段创建对应的widget配置
		widgetConfig, err := NewWidget(v, "table") // 使用"table"渲染类型
		var widgetType string
		if err != nil {
			// 如果创建失败，使用默认的input组件
			widgetConfig, _ = NewInputWidget(v)
			widgetType = WidgetInput
		} else {
			widgetType = widgetConfig.GetWidgetType()
		}

		Columns = append(Columns, TableColumn{
			Code:         v.GetCode(),
			Name:         v.GetName(),
			ValueType:    v.GetValueType(),
			WidgetType:   widgetType,
			WidgetConfig: widgetConfig,
			FieldConfig:  v,
		})
	}

	return &TableWidget{Widget: WidgetTable, Columns: Columns}, nil
}
