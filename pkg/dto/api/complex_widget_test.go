package api

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/yunhanshu-net/function-go/pkg/dto/response"
)

// ComplexWidgetReq 复杂组件测试请求结构
type ComplexWidgetReq struct {
	// 基础字段
	SendMessage bool   `json:"send_message" runner:"code:send_message;name:是否发送消息" widget:"type:switch;true_label:是;false_label:否" data:"type:boolean;default_value:false"`
	Remark      string `json:"remark" runner:"code:remark;name:备注" widget:"type:input;placeholder:请输入备注" data:"type:string;default_value:"`

	// 列表输入组件 - 批量操作列表
	List []BatchOperation `json:"list" runner:"code:list;name:批量操作列表" widget:"type:list_input;placeholder:添加批量操作项目" data:"type:[]struct;default_value:[]"`
}

// BatchOperation 批量操作结构
type BatchOperation struct {
	Name    string   `json:"name" runner:"code:name;name:名称" widget:"type:input;placeholder:请输入名称" data:"type:string;default_value:" validate:"required"`
	Type    string   `json:"type" runner:"code:type;name:类型" widget:"type:select;options:user,product,task;placeholder:请选择类型" data:"type:string;default_value:" validate:"required"`
	Enabled bool     `json:"enabled" runner:"code:enabled;name:是否启用" widget:"type:switch;true_label:启用;false_label:禁用" data:"type:boolean;default_value:true"`
	Tags    []string `json:"tags" runner:"code:tags;name:标签" widget:"type:multiselect;options:前端,Vue,React,手机,苹果,新品,开发,重要,紧急;placeholder:选择标签" data:"type:[]string;default_value:[]"`

	// 嵌套表单组件 - 配置
	Config OperationConfig `json:"config" runner:"code:config;name:配置" widget:"type:form_input;title:操作配置" data:"type:struct"`
}

// OperationConfig 操作配置结构
type OperationConfig struct {
	AutoSave    bool   `json:"auto_save" runner:"code:auto_save;name:自动保存" widget:"type:switch;true_label:是;false_label:否" data:"type:boolean;default_value:false"`
	MaxRetries  int    `json:"max_retries" runner:"code:max_retries;name:最大重试次数" widget:"type:input;placeholder:请输入重试次数" data:"type:number;default_value:3" validate:"min:1,max:10"`
	ProcessMode string `json:"process_mode" runner:"code:process_mode;name:处理模式" widget:"type:select;options:sync,async,batch;placeholder:请选择处理模式" data:"type:string;default_value:sync" validate:"required"`
}

// TestComplexWidgets 测试复杂组件
func TestComplexWidgets(t *testing.T) {
	fmt.Println("\n=== 测试复杂组件 ===")

	// 构建表单配置
	builder := NewFormBuilder()
	config, err := builder.BuildFormConfig(reflect.TypeOf(ComplexWidgetReq{}), response.RenderTypeForm)
	if err != nil {
		t.Fatalf("构建表单配置失败: %v", err)
	}

	// 序列化为JSON
	marshal, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	fmt.Println("复杂组件配置JSON:")
	fmt.Println(string(marshal))

	// 验证字段
	fmt.Println("\n=== 字段验证 ===")
	for _, field := range config.Fields {
		fmt.Printf("字段: %s\n", field.Code)
		fmt.Printf("  名称: %s\n", field.Name)
		fmt.Printf("  数据类型: %s\n", field.Data.Type)
		fmt.Printf("  组件类型: %s\n", field.Widget.Type)
		
		// 检查是否有子字段
		if field.Widget.Type == "list_input" || field.Widget.Type == "form_input" {
			fmt.Printf("  子字段数量: %d\n", len(field.Widget.Config))
		}
		fmt.Println()
	}
}

// TestNestedStructTypes 测试嵌套结构体类型
func TestNestedStructTypes(t *testing.T) {
	fmt.Println("\n=== 测试嵌套结构体类型 ===")

	// 测试结构体类型
	type TestStruct struct {
		Name string `json:"name" runner:"code:name;name:名称" widget:"type:input" data:"type:string"`
		Age  int    `json:"age" runner:"code:age;name:年龄" widget:"type:input" data:"type:number"`
	}

	// 测试结构体数组类型
	type TestStructArray struct {
		Items []TestStruct `json:"items" runner:"code:items;name:项目列表" widget:"type:list_input" data:"type:[]struct"`
	}

	builder := NewFormBuilder()
	config, err := builder.BuildFormConfig(reflect.TypeOf(TestStructArray{}), response.RenderTypeForm)
	if err != nil {
		t.Fatalf("构建嵌套结构体配置失败: %v", err)
	}

	marshal, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	fmt.Println("嵌套结构体配置JSON:")
	fmt.Println(string(marshal))
}

// TestRecursiveFieldParsing 测试递归字段解析
func TestRecursiveFieldParsing(t *testing.T) {
	fmt.Println("\n=== 测试递归字段解析 ===")

	// 构建表单配置
	builder := NewFormBuilder()
	config, err := builder.BuildFormConfig(reflect.TypeOf(ComplexWidgetReq{}), response.RenderTypeForm)
	if err != nil {
		t.Fatalf("构建表单配置失败: %v", err)
	}

	// 验证list_input组件的子字段
	fmt.Println("=== 验证list_input组件的子字段 ===")
	for _, field := range config.Fields {
		if field.Widget.Type == "list_input" {
			fmt.Printf("找到list_input组件: %s\n", field.Code)
			
			// 检查是否有子字段配置
			if fields, ok := field.Widget.Config["fields"]; ok {
				if subFields, ok := fields.([]*FieldInfo); ok {
					fmt.Printf("  子字段数量: %d\n", len(subFields))
					for _, subField := range subFields {
						fmt.Printf("    - %s (%s): %s\n", subField.Name, subField.Code, subField.Widget.Type)
						
								// 检查嵌套的form_input组件
		if subField.Widget.Type == "form_input" {
							if subFields2, ok := subField.Widget.Config["fields"]; ok {
								if subSubFields, ok := subFields2.([]*FieldInfo); ok {
									fmt.Printf("      嵌套form子字段数量: %d\n", len(subSubFields))
									for _, subSubField := range subSubFields {
										fmt.Printf("        - %s (%s): %s\n", subSubField.Name, subSubField.Code, subSubField.Widget.Type)
									}
								}
							}
						}
					}
				}
			} else {
				fmt.Printf("  警告: 没有找到子字段配置\n")
			}
		}
	}

	// 序列化为JSON查看完整结构
	marshal, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	fmt.Println("\n=== 完整配置JSON ===")
	fmt.Println(string(marshal))
} 