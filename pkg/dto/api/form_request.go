package api

import (
	"encoding/json"
	"strings"

	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/function-go/view/widget"
	"github.com/yunhanshu-net/pkg/x/tagx"
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

// ToFormRequestParamInfo 将 FieldInfo 转换为 FormRequestParamInfo（向下兼容）
func (f *FieldInfo) ToFormRequestParamInfo() *FormRequestParamInfo {
	// 转换回调为字符串格式
	var callbackStrings []string
	for _, callback := range f.Callbacks {
		callbackStrings = append(callbackStrings, callback.Event)
	}
	callbacksStr := strings.Join(callbackStrings, ",")

	// 转换权限为旧格式
	var show, hidden string
	if f.Permission != nil {
		var showParts []string
		if f.Permission.Create {
			showParts = append(showParts, "create")
		}
		if f.Permission.Update {
			showParts = append(showParts, "update")
		}
		if f.Permission.Read {
			showParts = append(showParts, "list")
		}
		show = strings.Join(showParts, ",")

		// 如果没有任何权限，设置为隐藏
		if !f.Permission.Read && !f.Permission.Update && !f.Permission.Create {
			hidden = "all"
		}
	}

	// 处理默认值
	var defaultValue interface{}
	if f.Data.DefaultValue != "" {
		defaultValue = f.Data.DefaultValue
	}

	return &FormRequestParamInfo{
		Code:         f.Code,
		Name:         f.Name,
		Desc:         f.Desc,
		Required:     strings.Contains(f.Validation, "required"),
		DefaultValue: defaultValue,
		Callbacks:    callbacksStr,
		Validates:    f.Validation,
		WidgetConfig: f.Widget.Config,
		WidgetType:   f.Widget.Type,
		ValueType:    f.Data.Type,
		Example:      f.Data.Example,
		Show:         show,
		Hidden:       hidden,
	}
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

//func NewFormRequestParams(el interface{}, renderType string) (*FormRequestParams, error) {
//	renderType = stringsx.DefaultString(renderType, response.RenderTypeForm)
//	typeOf := reflect.TypeOf(el)
//	if typeOf.Kind() == reflect.Pointer {
//		typeOf = typeOf.Elem()
//	}
//	if typeOf.Kind() != reflect.Struct {
//		return nil, fmt.Errorf("输入参数仅支持Struct类型")
//	}
//	reqFields, err := tagx.ParseStructFieldsTypeOf(typeOf, "runner")
//	if err != nil {
//		return nil, err
//	}
//
//	var searchCond []string
//	//	判断不同数据类型form,table,echarts,bi,3D ....
//	children := make([]*FormRequestParamInfo, 0, len(reqFields))
//	for _, field := range reqFields {
//		if field.IsSearchCond() {
//			searchCond = append(searchCond, field.GetCode())
//			continue
//		}
//		info, err := newFormRequestParamInfo(field)
//		if err != nil {
//			return nil, err
//		}
//		children = append(children, info)
//	}
//
//	return &FormRequestParams{
//		SearchCondList: strings.Join(searchCond, ","),
//		RenderType:     renderType,
//		Children:       children}, nil
//}
