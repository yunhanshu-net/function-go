package api

import (
	"fmt"
	"reflect"

	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/pkg/x/stringsx"
	"github.com/yunhanshu-net/pkg/x/tagx"
)

// FormBuilder 表单构建器
type FormBuilder struct {
	parser *tagx.MultiTagParser
}

// NewFormBuilder 创建表单构建器
func NewFormBuilder() *FormBuilder {
	return &FormBuilder{
		parser: tagx.NewMultiTagParser(),
	}
}

// BuildFormConfig 构建表单配置
func (b *FormBuilder) BuildFormConfig(structType reflect.Type, renderType string) (*FormConfig, error) {
	renderType = stringsx.DefaultString(renderType, response.RenderTypeForm)

	if structType.Kind() == reflect.Pointer {
		structType = structType.Elem()
	}
	if structType.Kind() != reflect.Struct {
		err := fmt.Errorf("输入参数仅支持Struct类型，当前类型: %s", structType.Kind())
		return nil, err
	}

	fields, err := b.parser.ParseStruct(structType)
	if err != nil {
		return nil, err
	}

	var formFields []*FieldInfo
	var searchConditions []string

	for _, field := range fields {
		// 判断是否为搜索条件
		if field.IsSearchCondition() {
			searchConditions = append(searchConditions, field.GetCode())
			continue
		}

		// 构建字段信息
		fieldInfo := b.buildFieldInfo(field, renderType)
		formFields = append(formFields, fieldInfo)
	}

	return &FormConfig{
		RenderType:       renderType,
		Fields:           formFields,
		SearchConditions: searchConditions,
	}, nil
}

// BuildTableConfig 构建表格配置
func (b *FormBuilder) BuildTableConfig(structType reflect.Type) (*TableConfig, error) {
	// 处理指针类型
	if structType.Kind() == reflect.Pointer {
		structType = structType.Elem()
	}

	// 处理切片类型 - 提取切片元素的类型
	if structType.Kind() == reflect.Slice {
		structType = structType.Elem()
		// 如果切片元素还是指针，继续解引用
		if structType.Kind() == reflect.Pointer {
			structType = structType.Elem()
		}
	}

	// 最终必须是结构体类型
	if structType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("输入参数必须是Struct类型或Struct切片类型，当前类型: %s", structType.Kind())
	}

	fields, err := b.parser.ParseStruct(structType)
	if err != nil {
		return nil, err
	}

	var columns []*FieldInfo
	for _, field := range fields {
		// 表格不需要搜索条件字段
		if field.IsSearchCondition() {
			continue
		}

		fieldInfo := b.buildFieldInfo(field, response.RenderTypeTable)
		columns = append(columns, fieldInfo)
	}

	return &TableConfig{RenderType: response.RenderTypeTable, Columns: columns}, nil
}

// buildFieldInfo 构建字段信息
func (b *FormBuilder) buildFieldInfo(field *tagx.FieldConfig, renderType string) *FieldInfo {
	fieldInfo := &FieldInfo{Code: field.GetCode(), Name: field.GetName()}

	// 设置描述
	if field.Runner != nil {
		fieldInfo.Desc = field.Runner.Desc
	}

	// 构建Widget配置
	fieldInfo.Widget = b.buildWidgetConfig(field)

	// 构建数据配置
	fieldInfo.Data = b.buildDataConfig(field)

	// 构建权限配置
	fieldInfo.Permission = b.buildPermissionConfig(field, renderType)

	// 构建回调配置
	fieldInfo.Callbacks = b.buildCallbackConfigs(field)

	// 设置验证配置
	fieldInfo.Validation = field.Validation

	return fieldInfo
}

// buildWidgetConfig 构建Widget配置
func (b *FormBuilder) buildWidgetConfig(field *tagx.FieldConfig) WidgetConfig {
	config := WidgetConfig{
		Type:   "input", // 默认为input
		Config: make(map[string]interface{}),
	}

	if field.Widget != nil {
		if field.Widget.Type != "" {
			config.Type = field.Widget.Type
		}

		// 直接使用解析器已经处理好的配置
		config.Config = field.Widget.Config

		// 如果Config为空，初始化为空map
		if config.Config == nil {
			config.Config = make(map[string]interface{})
		}
	}

	return config
}

// buildDataConfig 构建数据配置
func (b *FormBuilder) buildDataConfig(field *tagx.FieldConfig) DataConfig {
	config := DataConfig{
		Type: field.GetType(), // 获取字段类型
	}

	if field.Data != nil {
		config.Type = field.Data.Type // 优先使用data标签中的type
		config.Example = field.Data.Example
		config.DefaultValue = field.Data.DefaultValue
		config.Source = field.Data.Source
		config.Format = field.Data.Format
	}

	// 如果data标签中没有type，使用推断的类型
	if config.Type == "" {
		config.Type = field.GetType()
	}

	return config
}

// buildPermissionConfig 构建权限配置
func (b *FormBuilder) buildPermissionConfig(field *tagx.FieldConfig, renderType string) *PermissionConfig {
	// 只有在字段有permission标签时才返回权限配置
	if field.Permission != nil {
		return &PermissionConfig{
			Read:   field.Permission.Read,
			Update: field.Permission.Update,
			Create: field.Permission.Create,
		}
	}

	// 没有permission标签的字段，无论是Form还是Table模式都返回nil
	// 前端将把nil理解为没有权限限制（即全部权限）
	return nil
}

// buildCallbackConfigs 构建回调配置
func (b *FormBuilder) buildCallbackConfigs(field *tagx.FieldConfig) []CallbackConfig {
	var configs []CallbackConfig

	for _, callback := range field.Callbacks {
		configs = append(configs, CallbackConfig{
			Event:  callback.Event,
			Params: callback.Params,
		})
	}

	return configs
}
