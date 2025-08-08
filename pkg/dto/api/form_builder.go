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
	parser       *tagx.MultiTagParser
	functionInfo FunctionInfoInterface // 添加FunctionInfo支持
}

// NewFormBuilder 创建表单构建器
func NewFormBuilder() *FormBuilder {
	return &FormBuilder{
		parser: tagx.NewMultiTagParser(),
	}
}

// NewFormBuilderWithFunctionInfo 创建支持FunctionInfo的表单构建器
func NewFormBuilderWithFunctionInfo(functionInfo FunctionInfoInterface) *FormBuilder {
	return &FormBuilder{
		parser:       tagx.NewMultiTagParser(),
		functionInfo: functionInfo,
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

	for _, field := range fields {
		// 构建字段信息
		fieldInfo := b.buildFieldInfo(field, renderType)
		formFields = append(formFields, fieldInfo)
	}

	return &FormConfig{
		RenderType: renderType,
		Fields:     formFields,
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

	// 设置搜索配置
	fieldInfo.Search = b.buildSearchConfig(field)

	return fieldInfo
}

func hasSubFields(tp string) bool {
	return tp == "list" || tp == "list_input" || tp == "form" || tp == "form_input"
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

		// 特殊处理：递归解析子字段
		if hasSubFields(config.Type) {
			// 递归解析子字段
			subFields, err := b.buildSubFields(field)
			if err == nil && len(subFields) > 0 {
				// 将子字段添加到Fields字段中
				config.Fields = subFields
			}
		}
	}

	return config
}

// buildSubFields 递归构建子字段配置
func (b *FormBuilder) buildSubFields(field *tagx.FieldConfig) ([]*FieldInfo, error) {
	var subFields []*FieldInfo

	// 根据字段类型确定要解析的结构体类型
	var structType reflect.Type

	switch field.FieldType.Kind() {
	case reflect.Slice:
		// 对于 []struct 类型，解析切片元素的类型
		elemType := field.FieldType.Elem()
		if elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
		}
		if elemType.Kind() == reflect.Struct {
			structType = elemType
		}
	case reflect.Struct:
		// 对于 struct 类型，直接使用
		structType = field.FieldType
	case reflect.Ptr:
		// 对于 *struct 类型，解引用
		if field.FieldType.Elem().Kind() == reflect.Struct {
			structType = field.FieldType.Elem()
		}
	}

	// 如果没有找到有效的结构体类型，返回空列表
	if structType == nil {
		return subFields, nil
	}

	// 解析结构体字段
	fields, err := b.parser.ParseStruct(structType)
	if err != nil {
		return subFields, err
	}

	// 构建子字段信息
	for _, subField := range fields {
		fieldInfo := b.buildFieldInfo(subField, "form") // 子字段默认使用form渲染类型
		subFields = append(subFields, fieldInfo)
	}

	return subFields, nil
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

	// 首先添加标签中定义的回调
	for _, callback := range field.Callbacks {
		configs = append(configs, CallbackConfig{
			Event:  callback.Event,
			Params: callback.Params,
		})
	}

	// 如果有FunctionInfo，自动注入字段级回调
	if b.functionInfo != nil {
		fieldCode := field.GetCode()

		// 检查OnInputFuzzyMap中是否有该字段的回调
		if fuzzyMap := b.functionInfo.GetOnInputFuzzyMap(); fuzzyMap != nil {
			if _, exists := fuzzyMap[fieldCode]; exists {
				// 检查是否已经有OnInputFuzzy回调，避免重复
				hasOnInputFuzzy := false
				for _, existing := range configs {
					if existing.Event == "OnInputFuzzy" {
						hasOnInputFuzzy = true
						break
					}
				}

				if !hasOnInputFuzzy {
					configs = append(configs, CallbackConfig{
						Event: "OnInputFuzzy",
						Params: map[string]string{
							"delay": "300", // 默认延迟300ms
							"min":   "2",   // 默认最少2个字符
						},
					})
				}
			}
		}

		// 检查OnInputValidateMap中是否有该字段的回调
		if validateMap := b.functionInfo.GetOnInputValidateMap(); validateMap != nil {
			if _, exists := validateMap[fieldCode]; exists {
				// 检查是否已经有OnInputValidate回调，避免重复
				hasOnInputValidate := false
				for _, existing := range configs {
					if existing.Event == "OnInputValidate" {
						hasOnInputValidate = true
						break
					}
				}

				if !hasOnInputValidate {
					configs = append(configs, CallbackConfig{
						Event: "OnInputValidate",
						Params: map[string]string{
							"trigger": "blur", // 默认失焦时触发
						},
					})
				}
			}
		}
	}

	return configs
}

// buildSearchConfig 构建搜索配置
func (b *FormBuilder) buildSearchConfig(field *tagx.FieldConfig) *SearchConfig {
	if field.Search == nil {
		return nil
	}

	return &SearchConfig{
		Operators: field.Search.Operators,
	}
}
