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

type FunctionOptions struct {
	AutoUpdateConfig interface{}                        `json:"auto_update_config"` //这里是函数的配置，这里的
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
	OperateTables    map[interface{}][]OperateTableType `json:"operate_tables"`     //用到了哪些表，对表进行了哪些操作方便梳理引用关系
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

	BeforeRunnerClose BeforeRunnerClose `json:"-"`
	AfterRunnerClose  AfterRunnerClose  `json:"-"`
	OnVersionChange   OnVersionChange   `json:"-"`

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

// GetOnInputFuzzyMap 实现FunctionInfoInterface接口
func (f *FunctionOptions) GetOnInputFuzzyMap() map[string]interface{} {
	if f.OnInputFuzzyMap == nil {
		return nil
	}

	result := make(map[string]interface{})
	for key, value := range f.OnInputFuzzyMap {
		result[key] = value
	}
	return result
}

// GetOnInputValidateMap 实现FunctionInfoInterface接口
func (f *FunctionOptions) GetOnInputValidateMap() map[string]interface{} {
	if f.OnInputValidateMap == nil {
		return nil
	}

	result := make(map[string]interface{})
	for key, value := range f.OnInputValidateMap {
		result[key] = value
	}
	return result
}
