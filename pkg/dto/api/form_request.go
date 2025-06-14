package api

import (
	"encoding/json"
	"fmt"
	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/function-go/view/widget"
	"github.com/yunhanshu-net/pkg/x/stringsx"
	"github.com/yunhanshu-net/pkg/x/tagx"
	"reflect"
	"strings"
)

type FormRequestParamInfo struct {
	//英文标识
	Code string `json:"code"`
	//中文名称
	Name string `json:"name"`

	//介绍
	Desc string `json:"desc"`
	//是否必填
	Required     bool        `json:"required"`      //是否必填
	DefaultValue interface{} `json:"default_value"` //默认值
	Callbacks    string      `json:"callbacks"`     //字段级别的回调，多个用逗号分隔，ps OnInputFuzzy,OnInputFocus 等等，（还没实现）
	Validates    string      `json:"validates"`     //验证规则
	WidgetConfig interface{} `json:"widget_config"` //这里是widget.Widget类型的接口，需要实现Widget接口，这里可以是每个不同组件的个性化属性
	WidgetType   string      `json:"widget_type"`   //组件类型
	ValueType    string      `json:"value_type"`    //type 类型
	Example      string      `json:"example"`       //示例值

	Show   string `json:"show"`   //是否展示 create,update,list (仅仅在这三个场景显示) 为空表示全部显示，如果想隐藏用hidden来控制
	Hidden string `json:"hidden"` //是否隐藏 all(全部隐藏)
}

type FormRequestParams struct {
	SearchCondList string `json:"search_cond_list"` //支持的查询条件
	//SearchCondBlickList map[string]string   `json:"search_cond_blick_list"` //禁止的查询条件
	RenderType string                  `json:"render_type"`
	Children   []*FormRequestParamInfo `json:"children"`
}

func (p *FormRequestParams) JSONRawMessage() (json.RawMessage, error) {
	marshal, err := json.Marshal(p)
	if err != nil {
		return json.RawMessage("{}"), err
	}
	return marshal, nil
}

func newFormRequestParamInfo(tag *tagx.RunnerFieldInfo) (*FormRequestParamInfo, error) {

	widgetIns, err := widget.NewWidget(tag, response.RenderTypeForm)
	if err != nil {
		return nil, err
	}
	param := &FormRequestParamInfo{
		Code:         tag.GetCode(),
		Name:         tag.GetName(),
		Desc:         tag.GetDesc(),
		Show:         tag.GetShow(),
		Hidden:       tag.GetHidden(),
		DefaultValue: tag.GetDefaultValue(),
		Required:     tag.GetRequired(),
		Validates:    tag.GetValidates(),
		Callbacks:    tag.GetCallbacks(),
		WidgetConfig: widgetIns,
		WidgetType:   widgetIns.GetWidgetType(),
		ValueType:    tag.GetValueType(),
		Example:      tag.GetExample(),
	}

	return param, nil
}

func NewFormRequestParams(el interface{}, renderType string) (*FormRequestParams, error) {
	renderType = stringsx.DefaultString(renderType, response.RenderTypeForm)
	typeOf := reflect.TypeOf(el)
	if typeOf.Kind() == reflect.Pointer {
		typeOf = typeOf.Elem()
	}
	if typeOf.Kind() != reflect.Struct {
		return nil, fmt.Errorf("输入参数仅支持Struct类型")
	}
	reqFields, err := tagx.ParseStructFieldsTypeOf(typeOf, "runner")
	if err != nil {
		return nil, err
	}

	var searchCond []string
	//	判断不同数据类型form,table,echarts,bi,3D ....
	children := make([]*FormRequestParamInfo, 0, len(reqFields))
	for _, field := range reqFields {
		if field.IsSearchCond() {
			searchCond = append(searchCond, field.GetCode())
			continue
		}
		info, err := newFormRequestParamInfo(field)
		if err != nil {
			return nil, err
		}
		children = append(children, info)
	}

	return &FormRequestParams{
		SearchCondList: strings.Join(searchCond, ","),
		RenderType:     renderType,
		Children:       children}, nil
}
