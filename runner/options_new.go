package runner

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/yunhanshu-net/function-go/pkg/dto/usercall"
	"reflect"
)

// ==================== 函数组配置 ====================

// FunctionGroup 函数组配置（简化版）
type FunctionGroup struct {
	CnName string `json:"cn_name"` // 组名称，如 "JSON转换"
	EnName string `json:"en_name"` // 一般用文件名称命名
	// 后续按需要添加字段
	// Description string `json:"description"` // 组描述
	// Version     string `json:"version"`     // 版本
	// Author      string `json:"author"`      // 作者
	// Tags        []string `json:"tags"`      // 标签
	// IsAtomic    bool `json:"is_atomic"`    // 是否原子组
}

// ==================== 选项接口 ====================

// Option 选项接口
type Option interface {
	GetFunctionType() FunctionType
	GetRenderType() string
	GetBaseConfig() *BaseConfig
	Validate() error
	// 新增方法，用于获取创建表信息
	GetCreateTables() []interface{}
	// 新增方法，用于获取自动CRUD表信息
	GetAutoCrudTable() interface{}
	// 新增方法，用于获取回调信息
	GetCallbacks() map[string]interface{}
}

// ==================== 基础配置 ====================

// BaseConfig 基础配置（所有函数通用）
type BaseConfig struct {
	Router string `json:"router"` //api的路由
	Method string `json:"method"` //api的method

	// 名称配置
	EnglishName string   `json:"english_name" validate:"required"`
	ChineseName string   `json:"chinese_name" validate:"required"`
	ApiDesc     string   `json:"api_desc"`
	Tags        []string `json:"tags"`
	Classify    string   `json:"classify"`

	// 函数组配置
	Group *FunctionGroup `json:"group"` // 函数组配置

	// 执行配置
	Async        bool         `json:"async"`
	FunctionType FunctionType `json:"function_type"` // 函数类型：static/dynamic/pure
	Timeout      int          `json:"timeout"`

	// 权限配置
	IsPublicApi bool `json:"is_public_api"`

	// 请求响应
	Request  interface{} `json:"-"`
	Response interface{} `json:"-"`

	// 数据库配置
	CreateTables  []interface{}                      `json:"create_tables"`
	OperateTables map[interface{}][]OperateTableType `json:"-"`

	// 自动更新配置
	AutoUpdateConfig *AutoUpdateConfig `json:"auto_update_config"`

	// 自动运行
	AutoRun bool `json:"-"`
}

// ==================== 扁平化选项结构 ====================

// FormFunctionOptions 表单函数选项（扁平化设计）
type FormFunctionOptions struct {
	// 基础配置
	BaseConfig `json:",inline"`

	// 生命周期回调（所有函数通用）
	OnApiCreated      OnApiCreated      `json:"-"`
	OnApiUpdated      OnApiUpdated      `json:"-"`
	BeforeApiDelete   BeforeApiDelete   `json:"-"`
	AfterApiDeleted   AfterApiDeleted   `json:"-"`
	BeforeRunnerClose BeforeRunnerClose `json:"-"`
	AfterRunnerClose  AfterRunnerClose  `json:"-"`
	OnVersionChange   OnVersionChange   `json:"-"`

	// 通用回调
	OnPageLoad OnPageLoad `json:"-"`

	// 组件级回调
	OnInputFuzzyMap    map[string]OnInputFuzzy    `json:"-"`
	OnInputValidateMap map[string]OnInputValidate `json:"-"`

	// 表单专用回调
	OnDryRun OnDryRun `json:"-"`
}

// TableFunctionOptions 表格函数选项（扁平化设计）
type TableFunctionOptions struct {
	// 基础配置
	BaseConfig `json:",inline"`

	// 生命周期回调（所有函数通用）
	OnApiCreated      OnApiCreated      `json:"-"`
	OnApiUpdated      OnApiUpdated      `json:"-"`
	BeforeApiDelete   BeforeApiDelete   `json:"-"`
	AfterApiDeleted   AfterApiDeleted   `json:"-"`
	BeforeRunnerClose BeforeRunnerClose `json:"-"`
	AfterRunnerClose  AfterRunnerClose  `json:"-"`
	OnVersionChange   OnVersionChange   `json:"-"`

	// 通用回调
	OnPageLoad OnPageLoad `json:"-"`

	// 组件级回调
	OnInputFuzzyMap    map[string]OnInputFuzzy    `json:"-"`
	OnInputValidateMap map[string]OnInputValidate `json:"-"`

	// 表格专用回调
	OnTableDeleteRows OnTableDeleteRows `json:"-"`
	OnTableUpdateRows OnTableUpdateRows `json:"-"`
	OnTableAddRows    OnTableAddRows    `json:"-"`
	OnTableSearch     OnTableSearch     `json:"-"`

	// 表格特有配置
	AutoCrudTable interface{} `json:"-"`
}

// ==================== 接口实现 ====================

// FormFunctionOptions 实现
func (opt *FormFunctionOptions) GetFunctionType() FunctionType {
	return opt.FunctionType
}

func (opt *FormFunctionOptions) GetRenderType() string {
	return "form"
}

func (opt *FormFunctionOptions) GetBaseConfig() *BaseConfig {
	return &opt.BaseConfig
}

func (opt *FormFunctionOptions) Validate() error {
	if opt.EnglishName == "" {
		return errors.New("english_name is required")
	}
	return nil
}

func (opt *FormFunctionOptions) GetCreateTables() []interface{} {
	return opt.CreateTables
}

func (opt *FormFunctionOptions) GetAutoCrudTable() interface{} {
	return nil // 表单函数没有AutoCrudTable
}

// 实现 FunctionInfoInterface 接口
func (opt *FormFunctionOptions) GetOnInputFuzzyMap() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range opt.OnInputFuzzyMap {
		result[k] = v
	}
	return result
}

func (opt *FormFunctionOptions) GetOnInputValidateMap() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range opt.OnInputValidateMap {
		result[k] = v
	}
	return result
}

func (opt *FormFunctionOptions) GetCallbacks() map[string]interface{} {
	callbacks := make(map[string]interface{})

	// 生命周期回调
	if opt.OnApiCreated != nil {
		callbacks["OnApiCreated"] = opt.OnApiCreated
	}
	if opt.OnApiUpdated != nil {
		callbacks["OnApiUpdated"] = opt.OnApiUpdated
	}
	if opt.BeforeApiDelete != nil {
		callbacks["BeforeApiDelete"] = opt.BeforeApiDelete
	}
	if opt.AfterApiDeleted != nil {
		callbacks["AfterApiDeleted"] = opt.AfterApiDeleted
	}
	if opt.BeforeRunnerClose != nil {
		callbacks["BeforeRunnerClose"] = opt.BeforeRunnerClose
	}
	if opt.AfterRunnerClose != nil {
		callbacks["AfterRunnerClose"] = opt.AfterRunnerClose
	}
	if opt.OnVersionChange != nil {
		callbacks["OnVersionChange"] = opt.OnVersionChange
	}

	// 通用回调
	if opt.OnPageLoad != nil {
		callbacks["OnPageLoad"] = opt.OnPageLoad
	}

	// 组件级回调
	if opt.OnInputFuzzyMap != nil {
		callbacks["OnInputFuzzyMap"] = opt.OnInputFuzzyMap
	}
	if opt.OnInputValidateMap != nil {
		callbacks["OnInputValidateMap"] = opt.OnInputValidateMap
	}

	// 表单专用回调
	if opt.OnDryRun != nil {
		callbacks["OnDryRun"] = opt.OnDryRun
	}

	return callbacks
}

// TableFunctionOptions 实现
func (opt *TableFunctionOptions) GetFunctionType() FunctionType {
	return opt.FunctionType
}

func (opt *TableFunctionOptions) GetRenderType() string {
	return "table"
}

func (opt *TableFunctionOptions) GetBaseConfig() *BaseConfig {
	return &opt.BaseConfig
}

func (opt *TableFunctionOptions) Validate() error {
	if opt.EnglishName == "" {
		return errors.New("english_name is required")
	}
	if opt.AutoCrudTable == nil {
		return errors.New("auto_crud_table is required for table functions")
	}
	return nil
}

func (opt *TableFunctionOptions) GetCreateTables() []interface{} {
	return opt.CreateTables
}

func (opt *TableFunctionOptions) GetAutoCrudTable() interface{} {
	return opt.AutoCrudTable
}

func (opt *TableFunctionOptions) GetCallbacks() map[string]interface{} {
	callbacks := make(map[string]interface{})

	// 生命周期回调
	if opt.OnApiCreated != nil {
		callbacks["OnApiCreated"] = opt.OnApiCreated
	}
	if opt.OnApiUpdated != nil {
		callbacks["OnApiUpdated"] = opt.OnApiUpdated
	}
	if opt.BeforeApiDelete != nil {
		callbacks["BeforeApiDelete"] = opt.BeforeApiDelete
	}
	if opt.AfterApiDeleted != nil {
		callbacks["AfterApiDeleted"] = opt.AfterApiDeleted
	}
	if opt.BeforeRunnerClose != nil {
		callbacks["BeforeRunnerClose"] = opt.BeforeRunnerClose
	}
	if opt.AfterRunnerClose != nil {
		callbacks["AfterRunnerClose"] = opt.AfterRunnerClose
	}
	if opt.OnVersionChange != nil {
		callbacks["OnVersionChange"] = opt.OnVersionChange
	}

	// 通用回调
	if opt.OnPageLoad != nil {
		callbacks["OnPageLoad"] = opt.OnPageLoad
	}

	// 组件级回调
	if opt.OnInputFuzzyMap != nil {
		callbacks["OnInputFuzzyMap"] = opt.OnInputFuzzyMap
	}
	if opt.OnInputValidateMap != nil {
		callbacks["OnInputValidateMap"] = opt.OnInputValidateMap
	}

	// 表格专用回调
	if opt.OnTableDeleteRows != nil {
		callbacks["OnTableDeleteRows"] = opt.OnTableDeleteRows
	}
	if opt.OnTableUpdateRows != nil {
		callbacks["OnTableUpdateRows"] = opt.OnTableUpdateRows
	}
	if opt.OnTableAddRows != nil {
		callbacks["OnTableAddRows"] = opt.OnTableAddRows
	}
	if opt.OnTableSearch != nil {
		callbacks["OnTableSearch"] = opt.OnTableSearch
	}

	return callbacks
}

// 实现 FunctionInfoInterface 接口
func (opt *TableFunctionOptions) GetOnInputFuzzyMap() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range opt.OnInputFuzzyMap {
		result[k] = v
	}
	return result
}

func (opt *TableFunctionOptions) GetOnInputValidateMap() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range opt.OnInputValidateMap {
		result[k] = v
	}
	return result
}

func (opt *TableFunctionOptions) defaultDeleteRows(ctx *Context, req *usercall.OnTableDeleteRowsReq) error {
	return ctx.MustGetOrInitDB().Model(opt.AutoCrudTable).Delete("id in ?", req.Ids).Error
}

func (opt *TableFunctionOptions) defaultUpdateRows(ctx *Context, req *usercall.OnTableUpdateRowsReq) error {
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

	return ctx.MustGetOrInitDB().Model(opt.AutoCrudTable).Where("id in ?", req.Ids).Updates(req.Fields).Error
}

func (opt *TableFunctionOptions) defaultAddRows(ctx *Context, req *usercall.OnTableAddRowsReq) error {
	//  获取模型的类型
	modelType := reflect.TypeOf(opt.AutoCrudTable)
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
		ctx.Logger.Errorf("DecodeBy failed: %v", err)
		return err
	}

	ctx.Logger.Infof("defaultAddRows slice: %+v", slice.Interface())

	// 检查是否有数据
	if slice.Len() == 0 {
		return fmt.Errorf("没有要添加的数据")
	}

	//  执行数据库插入
	return ctx.MustGetOrInitDB().Create(slicePtr.Interface()).Error
}
