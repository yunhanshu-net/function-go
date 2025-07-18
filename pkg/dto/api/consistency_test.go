package api

import (
	"encoding/json"
	"fmt"
	"testing"
)

// 测试各种组件类型的一致性输出结构

// SimpleInputStruct 简单输入组件测试
type SimpleInputStruct struct {
	Title string `runner:"code:title;name:标题" widget:"type:input" data:"type:string"`
}

// SelectStruct 选择组件测试
type SelectStruct struct {
	Category string `runner:"code:category;name:分类" widget:"type:select;options:手机,笔记本,平板" data:"type:string"`
}

// SwitchStruct 开关组件测试
type SwitchStruct struct {
	Enabled bool `runner:"code:enabled;name:启用" widget:"type:switch" data:"type:boolean"`
}

// ListInputStruct 列表输入组件测试
type ListInputStruct struct {
	Items []struct {
		Name string `runner:"code:name;name:名称" widget:"type:input" data:"type:string"`
		Age  int    `runner:"code:age;name:年龄" widget:"type:input" data:"type:number"`
	} `runner:"code:items;name:项目列表" widget:"type:list_input;placeholder:添加项目" data:"type:[]struct"`
}

// FormStruct 表单组件测试
type FormStruct struct {
	Config struct {
		Title   string `runner:"code:title;name:标题" widget:"type:input" data:"type:string"`
		Content string `runner:"code:content;name:内容" widget:"type:input" data:"type:string"`
	} `runner:"code:config;name:配置" widget:"type:form_input;title:配置信息" data:"type:struct"`
}

// ComplexStruct 复杂嵌套结构测试
type ComplexStruct struct {
	Name     string `runner:"code:name;name:名称" widget:"type:input" data:"type:string"`
	Users    []struct {
		Username string `runner:"code:username;name:用户名" widget:"type:input" data:"type:string"`
		Email    string `runner:"code:email;name:邮箱" widget:"type:input" data:"type:string"`
		Profile  struct {
			Age      int    `runner:"code:age;name:年龄" widget:"type:input" data:"type:number"`
			Gender   string `runner:"code:gender;name:性别" widget:"type:select;options:男,女,其他" data:"type:string"`
			Location string `runner:"code:location;name:所在地" widget:"type:input" data:"type:string"`
		} `runner:"code:profile;name:用户资料" widget:"type:form_input;title:个人资料" data:"type:struct"`
	} `runner:"code:users;name:用户列表" widget:"type:list_input;placeholder:添加用户" data:"type:[]struct"`
}

func TestWidgetOutputConsistency(t *testing.T) {
	tests := []struct {
		name     string
		structType interface{}
		expected  string
	}{
		{
			name:       "简单输入组件",
			structType: SimpleInputStruct{},
			expected:   "input",
		},
		{
			name:       "选择组件",
			structType: SelectStruct{},
			expected:   "select",
		},
		{
			name:       "开关组件",
			structType: SwitchStruct{},
			expected:   "switch",
		},
		{
			name:       "列表输入组件",
			structType: ListInputStruct{},
			expected:   "list_display",
		},
		{
			name:       "表单组件",
			structType: FormStruct{},
			expected:   "form_display",
		},
		{
			name:       "复杂嵌套结构",
			structType: ComplexStruct{},
			expected:   "complex",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := NewUnifiedFormResponse(tt.structType, "form")
			if err != nil {
				t.Fatalf("NewUnifiedFormResponse failed: %v", err)
			}

			// 输出JSON结构
			jsonData, err := json.MarshalIndent(response, "", "  ")
			if err != nil {
				t.Fatalf("JSON marshal failed: %v", err)
			}

			fmt.Printf("\n=== %s ===\n", tt.name)
			fmt.Printf("输出结构:\n%s\n", string(jsonData))

			// 验证结构一致性
			validateWidgetStructure(t, response, tt.expected)
		})
	}
}

func validateWidgetStructure(t *testing.T, response *UnifiedAPIResponse, expectedType string) {
	for _, field := range response.Fields {
		// 验证基本字段存在
		if field.Widget.Type == "" {
			t.Errorf("字段 %s 缺少 widget.type", field.Code)
		}

		// 验证config字段存在（即使是空map）
		if field.Widget.Config == nil {
			t.Errorf("字段 %s 缺少 widget.config", field.Code)
		}

		// 对于复合类型，验证fields字段
		if field.Widget.Type == "list_display" || field.Widget.Type == "form_display" {
			if field.Widget.Fields == nil {
				t.Errorf("复合类型字段 %s 缺少 widget.fields", field.Code)
			}
			if len(field.Widget.Fields) == 0 {
				t.Errorf("复合类型字段 %s 的 widget.fields 为空", field.Code)
			}

			// 递归验证子字段
			for _, subField := range field.Widget.Fields {
				if subField.Widget.Type == "" {
					t.Errorf("子字段 %s 缺少 widget.type", subField.Code)
				}
				if subField.Widget.Config == nil {
					t.Errorf("子字段 %s 缺少 widget.config", subField.Code)
				}
			}
		} else {
			// 对于简单类型，fields应该为空或nil
			if len(field.Widget.Fields) > 0 {
				t.Errorf("简单类型字段 %s 不应该有 widget.fields", field.Code)
			}
		}
	}
}

func TestWidgetStructureAnalysis(t *testing.T) {
	fmt.Println("\n=== 组件结构分析 ===")

	// 测试简单组件
	simpleStruct := SimpleInputStruct{}
	response, err := NewUnifiedFormResponse(simpleStruct, "form")
	if err != nil {
		t.Fatalf("NewUnifiedFormResponse failed: %v", err)
	}

	fmt.Println("简单组件结构:")
	for _, field := range response.Fields {
		fmt.Printf("  - %s: type=%s, hasConfig=%t, hasFields=%t\n", 
			field.Code, 
			field.Widget.Type, 
			field.Widget.Config != nil, 
			len(field.Widget.Fields) > 0)
	}

	// 测试复合组件
	complexStruct := ComplexStruct{}
	response, err = NewUnifiedFormResponse(complexStruct, "form")
	if err != nil {
		t.Fatalf("NewUnifiedFormResponse failed: %v", err)
	}

	fmt.Println("\n复合组件结构:")
	for _, field := range response.Fields {
		fmt.Printf("  - %s: type=%s, hasConfig=%t, hasFields=%t\n", 
			field.Code, 
			field.Widget.Type, 
			field.Widget.Config != nil, 
			len(field.Widget.Fields) > 0)
		
		// 显示子字段信息
		for _, subField := range field.Widget.Fields {
			fmt.Printf("    └─ %s: type=%s, hasConfig=%t, hasFields=%t\n", 
				subField.Code, 
				subField.Widget.Type, 
				subField.Widget.Config != nil, 
				len(subField.Widget.Fields) > 0)
		}
	}
} 