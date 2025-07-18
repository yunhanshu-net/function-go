package api

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/yunhanshu-net/function-go/pkg/dto/response"
)

// 简单的测试结构体
type SimpleListTest struct {
	Name string `json:"name" runner:"code:name;name:名称" widget:"type:input" data:"type:string"`
	Age  int    `json:"age" runner:"code:age;name:年龄" widget:"type:input" data:"type:number"`
}

type SimpleFormTest struct {
	Title   string `json:"title" runner:"code:title;name:标题" widget:"type:input" data:"type:string"`
	Content string `json:"content" runner:"code:content;name:内容" widget:"type:input" data:"type:string"`
}

// 测试递归解析的结构体
type RecursiveTest struct {
	// 基础字段
	Title string `json:"title" runner:"code:title;name:标题" widget:"type:input" data:"type:string"`

	// 列表输入组件
	Items []SimpleListTest `json:"items" runner:"code:items;name:项目列表" widget:"type:list_input;placeholder:添加项目" data:"type:[]struct"`

	// 嵌套表单组件
	Config SimpleFormTest `json:"config" runner:"code:config;name:配置" widget:"type:form_input;title:配置信息" data:"type:struct"`
}

// TestSimpleRecursive 测试简单的递归解析
func TestSimpleRecursive(t *testing.T) {
	fmt.Println("=== 测试简单递归解析 ===")

	builder := NewFormBuilder()
	config, err := builder.BuildFormConfig(reflect.TypeOf(RecursiveTest{}), response.RenderTypeForm)
	if err != nil {
		t.Fatalf("构建配置失败: %v", err)
	}

	// 验证字段
	fmt.Println("=== 字段验证 ===")
	for _, field := range config.Fields {
		fmt.Printf("字段: %s (%s)\n", field.Name, field.Code)
		fmt.Printf("  数据类型: %s\n", field.Data.Type)
		fmt.Printf("  组件类型: %s\n", field.Widget.Type)

		// 检查list_input组件的子字段
		if field.Widget.Type == "list_input" {
			if fields, ok := field.Widget.Config["fields"]; ok {
				if subFields, ok := fields.([]*FieldInfo); ok {
					fmt.Printf("  list_input子字段数量: %d\n", len(subFields))
					for _, subField := range subFields {
						fmt.Printf("    - %s (%s): %s\n", subField.Name, subField.Code, subField.Widget.Type)
					}
				}
			}
		}

		// 检查form_input组件的子字段
		if field.Widget.Type == "form_input" {
			if fields, ok := field.Widget.Config["fields"]; ok {
				if subFields, ok := fields.([]*FieldInfo); ok {
					fmt.Printf("  form子字段数量: %d\n", len(subFields))
					for _, subField := range subFields {
						fmt.Printf("    - %s (%s): %s\n", subField.Name, subField.Code, subField.Widget.Type)
					}
				}
			}
		}
		fmt.Println()
	}

	// 输出JSON配置
	marshal, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	fmt.Println("=== 完整配置JSON ===")
	fmt.Println(string(marshal))
}

// 测试嵌套的复杂结构
type NestedComplexTest struct {
	// 基础字段
	Name string `json:"name" runner:"code:name;name:名称" widget:"type:input" data:"type:string"`

	// 列表输入：包含嵌套表单
	Users []UserWithProfile `json:"users" runner:"code:users;name:用户列表" widget:"type:list_input;placeholder:添加用户" data:"type:[]struct"`
}

type UserWithProfile struct {
	Username string `json:"username" runner:"code:username;name:用户名" widget:"type:input" data:"type:string"`
	Email    string `json:"email" runner:"code:email;name:邮箱" widget:"type:input" data:"type:string"`

	// 嵌套表单：用户资料
	Profile UserProfile `json:"profile" runner:"code:profile;name:用户资料" widget:"type:form;title:个人资料" data:"type:struct"`
}

type UserProfile struct {
	Age      int    `json:"age" runner:"code:age;name:年龄" widget:"type:input" data:"type:number"`
	Gender   string `json:"gender" runner:"code:gender;name:性别" widget:"type:select;options:男,女,其他" data:"type:string"`
	Location string `json:"location" runner:"code:location;name:所在地" widget:"type:input" data:"type:string"`
}

// TestNestedComplex 测试嵌套的复杂结构
func TestNestedComplex(t *testing.T) {
	fmt.Println("\n=== 测试嵌套复杂结构 ===")

	builder := NewFormBuilder()
	config, err := builder.BuildFormConfig(reflect.TypeOf(NestedComplexTest{}), response.RenderTypeForm)
	if err != nil {
		t.Fatalf("构建配置失败: %v", err)
	}

	// 验证嵌套结构
	fmt.Println("=== 嵌套结构验证 ===")
	for _, field := range config.Fields {
		if field.Widget.Type == "list_input" {
			fmt.Printf("找到list_input组件: %s\n", field.Code)
			
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
			}
		}
	}

	// 输出JSON配置
	marshal, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	fmt.Println("\n=== 完整配置JSON ===")
	fmt.Println(string(marshal))
} 