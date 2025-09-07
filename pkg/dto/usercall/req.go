package usercall

import (
	"encoding/json"
	"fmt"
	"github.com/yunhanshu-net/function-go/view/widget/types"
	"strings"

	"github.com/yunhanshu-net/pkg/x/jsonx"
)

// GenerateConfigKey 统一的配置key生成函数
func GenerateConfigKey(router, method string) string {
	// 将路由中的路径分隔符替换为点号
	routerKey := strings.ReplaceAll(strings.Trim(router, "/"), "/", ".")
	// 去除前后多余的点号
	routerKey = strings.Trim(routerKey, ".")
	// 使用大写 method
	return fmt.Sprintf("function.%s.%s", routerKey, strings.ToUpper(method))
}

// NoData 空数据结构
type NoData struct{}

// ApiInfoRequest API信息请求
type ApiInfoRequest struct {
	Router string `json:"router" form:"router"` // API路由路径
	Method string `json:"method" form:"method"` // HTTP方法（GET/POST）
}

type OnPageLoadReq struct {
}

type OnPageLoadMessage struct {
	Type    string `json:"type"`    //warn/error
	Title   string `json:"title"`   //标题
	Content string `json:"content"` //内容
}
type OnPageLoadResp struct {
	Request interface{} `json:"request"` //会初始化前端的表单参数

	Message    *OnPageLoadMessage `json:"message"`
	DisableRun bool               `json:"disable_run"` //禁止此次运行
	AutoRun    bool               `json:"auto_run"`    //是否自动运行
}

type OnApiCreatedReq struct {
	//Method string `json:"method"`
	//Router string `json:"router"`
}

type OnApiUpdatedReq struct {
	Method string `json:"method"`
	Router string `json:"router"`
}

type BeforeApiDeleteReq struct {
	Method string `json:"method"`
	Router string `json:"router"`
}

type AfterApiDeletedReq struct {
	Method string `json:"method"`
	Router string `json:"router"`
}

type BeforeRunnerCloseReq struct {
}

type AfterRunnerCloseReq struct {
}

type Change struct {
	Method string `json:"method"`
	Router string `json:"router"`
	Type   string `json:"type"`
}

func (c *Change) String() string {
	return fmt.Sprintf(`{"method": "%s", "router": "%s","type","%s"}`, c.Method, c.Router, c.Type)
}

type OnVersionChangeReq struct {
	Change []Change `json:"change"`
}

type OnInputFuzzyReq struct {
	Code      string      `json:"code"`  //回调的这个字段的key
	Value     interface{} `json:"value"` //用户输入的值
	Request   interface{} `json:"request"`
	InputType string      `json:"input_type"` //by_filed_value/by_filed_values
	ValueType string      `json:"value_type"` //
	keywork   string
}

func (r *OnInputFuzzyReq) GetFiledValues() interface{} {

	switch r.ValueType {
	case types.ValueFloat, types.ValueFloats:
		switch r.Value.(type) {
		case []interface{}:
			var floats []float64
			jsonx.Convert(r.Value, &floats)
			return floats
		}
	case types.ValueNumber, types.ValueNumbers:
		switch r.Value.(type) {
		case []interface{}:
			var ints []int
			jsonx.Convert(r.Value, &ints)
			return ints
		}

	case types.ValueString, types.ValueStrings:
		switch r.Value.(type) {
		case []interface{}:
			var strs []string
			jsonx.Convert(r.Value, &strs)
			return strs
		}
	}

	return r.Value
}
func (r *OnInputFuzzyReq) GetFiledValue() interface{} {

	if r.ValueType == types.ValueNumber {
		return int(r.Value.(float64))
	}
	return r.Value
}

//	func (r *OnInputFuzzyReq) IsByFieldValue() bool {
//		return r.InputType == "by_filed_value"
//	}
func (r *OnInputFuzzyReq) IsByFiledValues() bool {
	return r.InputType == "by_field_values"
}
func (r *OnInputFuzzyReq) IsByFiledValue() bool {
	return r.InputType == "by_field_value"
}

func (r *OnInputFuzzyReq) Keyword() string {
	if r.keywork == "" {
		r.keywork = fmt.Sprintf("%s", r.Value)
	}
	return r.keywork
}

//func (r *OnInputFuzzyReq) ByID() int {
//	//i, err := strconv.ParseInt(fmt.Sprintf("%s", r.Value), 10, 64)
//	//if err != nil {
//	//	return 0
//	//}
//	v, ok := r.Value.(int)
//	if !ok {
//		return v
//	}
//
//	return 0
//}

//func (r *OnInputFuzzyReq) ByIDS() []int {
//	ints, ok := r.Value.([]int)
//	if !ok {
//		return []int{0}
//	}
//	return ints
//}

func (c *OnInputFuzzyReq) DecodeBy(el interface{}) error {
	err := jsonx.Convert(c.Request, el)
	if err != nil {
		return err
	}
	return nil
}

type OnInputValidateReq struct {
	Code    string      `json:"code"`
	Value   interface{} `json:"value"`
	Request interface{} `json:"request"`
}

func (c *OnInputValidateReq) DecodeBy(el interface{}) error {
	err := jsonx.Convert(c.Request, el)
	if err != nil {
		return err
	}
	return nil
}

type OnTableDeleteRowsReq struct {
	Ids []int `json:"ids"`
}

type OnTableAddRowsReq struct {
	Rows interface{} `json:"rows"`
}

type OnTableAddRowsResp struct {
}

func (r *OnTableAddRowsReq) DecodeBy(el interface{}) error {
	err := jsonx.Convert(r.Rows, el)
	if err != nil {
		return err
	}
	return nil
}

type OnTableSearchReq struct {
	Cond map[string]string `json:"cond"`
}
type Request struct {
	Method string      `json:"method"`
	Router string      `json:"router"`
	Type   string      `json:"type"`
	Body   interface{} `json:"body"`
}

func (c *Request) DecodeData(el interface{}) error {
	marshal, err := json.Marshal(c.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, &el)
	if err != nil {
		return err
	}
	return nil
}

type Response struct {
	Request  interface{} `json:"request"`
	Response interface{} `json:"response"`
}

type InputFuzzyItem struct {
	Value       interface{}            `json:"value"`
	Label       string                 `json:"label"`
	Icon        string                 `json:"icon"`
	DisplayInfo map[string]interface{} `json:"display_info"`
}

type OnInputFuzzyResp struct {
	//只有在结构体数组或者切片下的select和multiselect组件才会有聚合计算的功能，场景例如收银，我一个[]Orders
	//下面有ProductId，然后每个产品虽然选择产品id，但是DisplayInfo里返回了价格，这时候我想价格求和来计算，statistics"价格":"sum"即可

	Statistics map[string]interface{} `json:"statistics"`
	Values     []*InputFuzzyItem      `json:"values"`
}

type OnInputValidateResp struct {
	ErrorMsg string `json:"error_msg"`
}

type OnTableDeleteRowsResp struct {
}

type OnTableUpdateRowsResp struct {
}

type OnTableSearchResp struct {
}

// DryRunCase 抽象接口 - 所有危险操作都应该实现此接口
type DryRunCase interface {
	// Type 返回操作类型
	Type() string

	// Map 返回操作详情
	Map() map[string]interface{}

	// Metadata 返回元数据
	Metadata() map[string]interface{}
}

// OnDryRunReq DryRun 请求结构体
type OnDryRunReq struct {
	Body interface{} `json:"body"` // 原始请求体
}

// DecodeBody 解码请求体到指定类型
func (r *OnDryRunReq) DecodeBody(el interface{}) error {
	return jsonx.Convert(r.Body, el)
}

// OnDryRunResp DryRun 响应结构体
type OnDryRunResp struct {
	Valid   bool         `json:"valid"`   // 是否有效
	Cases   []DryRunCase `json:"cases"`   // DryRun 案例列表
	Message string       `json:"message"` // 提示信息
}

// UpdateConfigReq 配置更新请求
type UpdateConfigReq struct {
	Router     string                 `json:"router"`      // 路由路径
	Method     string                 `json:"method"`      // HTTP方法
	ConfigData map[string]interface{} `json:"config_data"` // 配置数据
}

// ToConfigData 转换为ConfigData结构
func (req *UpdateConfigReq) ToConfigData() *ConfigData {
	return &ConfigData{
		Type: "json",
		Data: req.ConfigData,
	}
}

// GenerateConfigKey 生成配置键
func (req *UpdateConfigReq) GenerateConfigKey() string {
	return GenerateConfigKey(req.Router, req.Method)
}

// GetConfigReq 配置获取请求
type GetConfigReq struct {
	Router string `json:"router"` // 路由路径
	Method string `json:"method"` // HTTP方法
}

// GenerateConfigKey 生成配置键
func (req *GetConfigReq) GenerateConfigKey() string {
	return GenerateConfigKey(req.Router, req.Method)
}

// UpdateConfigResp 配置更新响应
type UpdateConfigResp struct {
	Success bool   `json:"success"` // 是否成功
	Message string `json:"message"` // 响应消息
	Error   string `json:"error"`   // 错误信息
}

// GetConfigResp 配置获取响应
type GetConfigResp struct {
	Success bool        `json:"success"` // 是否成功
	Config  *ConfigData `json:"config"`  // 配置数据
	Error   string      `json:"error"`   // 错误信息
}

// ConfigData 配置数据结构
type ConfigData struct {
	Type string      `json:"type,omitempty"` // 配置类型：json, yaml, toml, xml 等（可选，默认为json）
	Data interface{} `json:"data"`           // 配置数据（直接存储，避免双重序列化）
}
