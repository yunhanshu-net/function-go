package runner

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/yunhanshu-net/function-go/pkg/dto/usercall"
	"github.com/yunhanshu-net/pkg/logger"
)

type FunctionType string
type OperateTableType string

const (
	FunctionTypeStatic  FunctionType = "static_function"  //无需参数，或者输入参数，但是结果永远恒定
	FunctionTypeDynamic FunctionType = "dynamic_function" //default，请求参数不可预测，响应参数不可预测（比如查询用户信息，用户的信息时时刻刻都可能在变化）输入zhangsan 现在可能他18岁，但是过了一会他就19岁了，所以结果不可预知
	FunctionTypePure    FunctionType = "pure_function"    //纯函数，例如1+1=2 or  sin(x) cos(x) 这种
)
const (
	OperateTableTypeAdd    OperateTableType = "add"    //增加数据
	OperateTableTypeUpdate OperateTableType = "update" //修改数据
	OperateTableTypeDelete OperateTableType = "delete" //删除数据
	OperateTableTypeGet    OperateTableType = "get"    //获取数据
)

type AutoCrud struct {
}

// AutoUpdateConfig 自动更新配置结构体

type FunctionOptions struct {
	AutoUpdateConfig *AutoUpdateConfig                  `json:"auto_update_config"` //这里是函数的配置，这里的
	Router           string                             `json:"router"`             //api的路由
	Method           string                             `json:"method"`             //api的method
	ApiDesc          string                             `json:"api_desc"`           //函数介绍
	IsPublicApi      bool                               `json:"is_public_api"`      //是否是公共api，默认false
	ChineseName      string                             `json:"chinese_name"`       //中文名称
	EnglishName      string                             `json:"english_name"`       //英文名称，需要符合go的文件名称规范和路由规范
	Classify         string                             `json:"classify"`           //分类
	Tags             []string                           `json:"tags"`               //tags
	Async            bool                               `json:"async"`              //是否异步，比较耗时的api，或者需要后台慢慢处理的api
	FunctionType     FunctionType                       `json:"function_type"`      //函数类型 默认：dynamic_function
	Timeout          int                                `json:"timeout"`            //超时时间，单位毫秒,0表示不超时
	RenderType       string                             `json:"widget"`             // 渲染类型	//form，table，echarts
	CreateTables     []interface{}                      `json:"create_tables"`      //创建该api时候会自动帮忙创建这个数据库表gorm的model列表
	UseTables        []interface{}                      `json:"use_tables"`         //这里需要记录这个函数用到的数据表，方便梳理引用关系
	OperateTables    map[interface{}][]OperateTableType `json:"-"`                  //用到了哪些表，对表进行了哪些操作方便梳理引用关系
	AutoRun          bool                               `json:"-"`                  //是否自动运行，默认false，如果为true，则在用户访问这个函数时候，会自动运行一次
	Request          interface{}                        `json:"-"`                  //这里是用户request请求的model，需要在相关字段打上runner 标签，runner:"-" 会忽略这个字段
	Response         interface{}                        `json:"-"`                  //这里是用户response的model，需要在相关字段打上runner 标签，runner:"-" 会忽略这个字段

	AutoCrudTable interface{} `json:"-"` //table函数列表接口，这里注册操作的表的model，内部可以自动实现表的更新和删除操作
	//用map的都是字段级别的回调，其他的都是接口级别回调

	//注意标注已经实现的可以被回调，没标注的，则不会被回调（正在实现中，如果生成代码，请生成已经实现的回调来）
	OnPageLoad OnPageLoad `json:"-"` //已经实现，优先级最高，先初始化表单参数，然后再判断是否有自动运行的回调，如果有，则执行

	OnApiCreated    OnApiCreated    `json:"-"` //已经实现
	OnApiUpdated    OnApiUpdated    `json:"-"`
	BeforeApiDelete BeforeApiDelete `json:"-"`
	AfterApiDeleted AfterApiDeleted `json:"-"`

	BeforeRunnerClose BeforeRunnerClose `json:"-"` // 运行器关闭前回调
	AfterRunnerClose  AfterRunnerClose  `json:"-"` // 运行器关闭后回调
	OnVersionChange   OnVersionChange   `json:"-"` // 版本变更回调

	OnTableDeleteRows OnTableDeleteRows `json:"-"` //已经实现
	OnTableUpdateRows OnTableUpdateRows `json:"-"` //已经实现
	OnTableAddRows    OnTableAddRows    `json:"-"` //已经实现
	OnTableSearch     OnTableSearch     `json:"-"`

	OnInputFuzzyMap    map[string]OnInputFuzzy    `json:"-"` //key是字段的code，字段级回调
	OnInputValidateMap map[string]OnInputValidate `json:"-"` //key是字段的code，字段级回调

	OnDryRun OnDryRun `json:"-"` // DryRun 回调，用于预览危险操作
}

func (f *FunctionOptions) defaultDeleteRows(ctx *Context, req *usercall.OnTableDeleteRowsReq) error {
	return ctx.MustGetOrInitDB().Model(f.AutoCrudTable).Delete("id in ?", req.Ids).Error
}

func (f *FunctionOptions) defaultUpdateRows(ctx *Context, req *usercall.OnTableUpdateRowsReq) error {
	for k, field := range req.Fields {
		switch field.(type) {
		case map[string]interface{}:
			marshal, err := json.Marshal(field)
			if err != nil {
				return err
			}
			req.Fields[k] = json.RawMessage(marshal)
		}
	}

	return ctx.MustGetOrInitDB().Model(f.AutoCrudTable).Where("id in ?", req.Ids).Updates(req.Fields).Error
}

func (f *FunctionOptions) defaultAddRows(ctx *Context, req *usercall.OnTableAddRowsReq) error {
	//  获取模型的类型
	modelType := reflect.TypeOf(f.AutoCrudTable)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	//  创建可寻址的slice
	sliceType := reflect.SliceOf(modelType)
	slicePtr := reflect.New(sliceType) // 创建slice指针
	slice := slicePtr.Elem()           // 获取slice值

	//  解码数据到slice中
	err := req.DecodeBy(slicePtr.Interface())
	if err != nil {
		logger.Errorf(ctx, "DecodeBy failed: %v", err)
		return err
	}

	logger.Infof(ctx, "defaultAddRows slice: %+v", slice.Interface())

	// 检查是否有数据
	if slice.Len() == 0 {
		return fmt.Errorf("没有要添加的数据")
	}

	//  执行数据库插入
	return ctx.MustGetOrInitDB().Create(slicePtr.Interface()).Error
}

// GetOnInputFuzzyMap 实现FunctionInfoProvider接口
func (f *FunctionOptions) GetOnInputFuzzyMap() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range f.OnInputFuzzyMap {
		result[k] = v
	}
	return result
}

// GetOnInputValidateMap 实现FunctionInfoProvider接口
func (f *FunctionOptions) GetOnInputValidateMap() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range f.OnInputValidateMap {
		result[k] = v
	}
	return result
}

// ==================== Option接口实现 ====================

// GetFunctionType 实现Option接口
func (f *FunctionOptions) GetFunctionType() FunctionType {
	return f.FunctionType
}

// GetRenderType 实现Option接口
func (f *FunctionOptions) GetRenderType() string {
	return f.RenderType
}

// GetBaseConfig 实现Option接口
func (f *FunctionOptions) GetBaseConfig() *BaseConfig {
	// 将FunctionOptions转换为BaseConfig
	return &BaseConfig{
		Router:           f.Router,
		Method:           f.Method,
		EnglishName:      f.EnglishName,
		ChineseName:      f.ChineseName,
		ApiDesc:          f.ApiDesc,
		Tags:             f.Tags,
		Classify:         f.Classify,
		Async:            f.Async,
		FunctionType:     f.FunctionType,
		Timeout:          f.Timeout,
		IsPublicApi:      f.IsPublicApi,
		Request:          f.Request,
		Response:         f.Response,
		CreateTables:     f.CreateTables,
		OperateTables:    f.OperateTables,
		AutoUpdateConfig: f.AutoUpdateConfig,
		AutoRun:          f.AutoRun,
		// Group字段在旧系统中不存在，设为nil
		Group: nil,
	}
}

// Validate 实现Option接口
func (f *FunctionOptions) Validate() error {
	if f.Router == "" {
		return fmt.Errorf("router is required")
	}
	if f.EnglishName == "" {
		return fmt.Errorf("english_name is required")
	}
	if f.Method == "" {
		return fmt.Errorf("method is required")
	}
	if f.FunctionType == "" {
		return fmt.Errorf("function_type is required")
	}
	if f.RenderType == "" {
		return fmt.Errorf("render_type is required")
	}
	
	// 表格函数特殊验证
	if f.RenderType == "table" {
		if f.AutoCrudTable == nil {
			return fmt.Errorf("auto_crud_table is required for table functions")
		}
	}
	
	return nil
}

// GetCreateTables 实现Option接口
func (f *FunctionOptions) GetCreateTables() []interface{} {
	return f.CreateTables
}

// GetAutoCrudTable 实现Option接口
func (f *FunctionOptions) GetAutoCrudTable() interface{} {
	return f.AutoCrudTable
}

// GetCallbacks 实现Option接口
func (f *FunctionOptions) GetCallbacks() map[string]interface{} {
	callbacks := make(map[string]interface{})
	
	// 生命周期回调
	if f.OnApiCreated != nil {
		callbacks["OnApiCreated"] = f.OnApiCreated
	}
	if f.OnApiUpdated != nil {
		callbacks["OnApiUpdated"] = f.OnApiUpdated
	}
	if f.BeforeApiDelete != nil {
		callbacks["BeforeApiDelete"] = f.BeforeApiDelete
	}
	if f.AfterApiDeleted != nil {
		callbacks["AfterApiDeleted"] = f.AfterApiDeleted
	}
	if f.BeforeRunnerClose != nil {
		callbacks["BeforeRunnerClose"] = f.BeforeRunnerClose
	}
	if f.AfterRunnerClose != nil {
		callbacks["AfterRunnerClose"] = f.AfterRunnerClose
	}
	if f.OnVersionChange != nil {
		callbacks["OnVersionChange"] = f.OnVersionChange
	}
	
	// 通用回调
	if f.OnPageLoad != nil {
		callbacks["OnPageLoad"] = f.OnPageLoad
	}
	
	// 组件级回调
	if f.OnInputFuzzyMap != nil {
		callbacks["OnInputFuzzyMap"] = f.OnInputFuzzyMap
	}
	if f.OnInputValidateMap != nil {
		callbacks["OnInputValidateMap"] = f.OnInputValidateMap
	}
	
	// 表单专用回调
	if f.OnDryRun != nil {
		callbacks["OnDryRun"] = f.OnDryRun
	}
	
	// 表格专用回调
	if f.OnTableDeleteRows != nil {
		callbacks["OnTableDeleteRows"] = f.OnTableDeleteRows
	}
	if f.OnTableUpdateRows != nil {
		callbacks["OnTableUpdateRows"] = f.OnTableUpdateRows
	}
	if f.OnTableAddRows != nil {
		callbacks["OnTableAddRows"] = f.OnTableAddRows
	}
	if f.OnTableSearch != nil {
		callbacks["OnTableSearch"] = f.OnTableSearch
	}
	
	return callbacks
}
