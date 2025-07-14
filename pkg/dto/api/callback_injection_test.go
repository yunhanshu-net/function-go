package api

import (
	"encoding/json"
	"fmt"
	"testing"
)

// MockFunctionInfo 模拟FunctionInfo用于测试
type MockFunctionInfo struct {
	OnInputFuzzyMap    map[string]interface{}
	OnInputValidateMap map[string]interface{}
}

func (m *MockFunctionInfo) GetOnInputFuzzyMap() map[string]interface{} {
	return m.OnInputFuzzyMap
}

func (m *MockFunctionInfo) GetOnInputValidateMap() map[string]interface{} {
	return m.OnInputValidateMap
}

// TestCallbackStruct 测试结构体
type TestCallbackStruct struct {
	// 有标签回调的字段
	Username string `json:"username" form:"username" runner:"code:username;name:用户名" widget:"type:input;placeholder:请输入用户名" data:"type:string" callback:"OnInputFuzzy(delay:500,min:3)" validate:"required"`

	// 没有标签回调但在FunctionInfo中有回调的字段
	CompanyName string `json:"company_name" form:"company_name" runner:"code:company_name;name:公司名称" widget:"type:input;placeholder:输入公司名称" data:"type:string"`
	UserInfo    string `json:"user_info" form:"user_info" runner:"code:user_info;name:用户信息" widget:"type:input;placeholder:输入用户信息" data:"type:string"`

	// 既有标签回调又在FunctionInfo中有回调的字段（应该不重复）
	Email string `json:"email" form:"email" runner:"code:email;name:邮箱" widget:"type:input;placeholder:请输入邮箱" data:"type:string" callback:"OnInputValidate(trigger:change)" validate:"required,email"`

	// 没有任何回调的字段
	Description string `json:"description" form:"description" runner:"code:description;name:描述" widget:"type:input;placeholder:请输入描述" data:"type:string"`
}

func TestCallbackInjection(t *testing.T) {
	fmt.Println("=== 测试回调自动注入功能 ===")

	// 创建模拟的FunctionInfo，定义字段级回调
	mockFunctionInfo := &MockFunctionInfo{
		OnInputFuzzyMap: map[string]interface{}{
			"company_name": func() {}, // 模拟回调函数
			"user_info":    func() {}, // 模拟回调函数
			"email":        func() {}, // 这个字段已有标签回调，应该不重复
		},
		OnInputValidateMap: map[string]interface{}{
			"user_info": func() {}, // 模拟回调函数
			"email":     func() {}, // 这个字段已有标签回调，应该不重复
		},
	}

	// 测试不带FunctionInfo的版本（原来的方式）
	fmt.Println("\n--- 测试不带FunctionInfo的版本 ---")
	paramsWithoutFunctionInfo, err := NewRequestParams(TestCallbackStruct{}, "form")
	if err != nil {
		t.Fatalf("NewRequestParams failed: %v", err)
	}

	withoutFunctionInfoJson, _ := json.MarshalIndent(paramsWithoutFunctionInfo, "", "  ")
	fmt.Printf("不带FunctionInfo的结果:\n%s\n", string(withoutFunctionInfoJson))

	// 测试带FunctionInfo的版本（新功能）
	fmt.Println("\n--- 测试带FunctionInfo的版本 ---")
	paramsWithFunctionInfo, err := NewRequestParamsWithFunctionInfo(TestCallbackStruct{}, "form", mockFunctionInfo)
	if err != nil {
		t.Fatalf("NewRequestParamsWithFunctionInfo failed: %v", err)
	}

	withFunctionInfoJson, _ := json.MarshalIndent(paramsWithFunctionInfo, "", "  ")
	fmt.Printf("带FunctionInfo的结果:\n%s\n", string(withFunctionInfoJson))

	// 验证结果
	if formConfig, ok := paramsWithFunctionInfo.(*FormConfig); ok {
		fmt.Println("\n--- 验证回调注入结果 ---")

		for _, field := range formConfig.Fields {
			fmt.Printf("字段 %s (%s):\n", field.Code, field.Name)

			if len(field.Callbacks) > 0 {
				fmt.Printf("  回调函数:\n")
				for _, callback := range field.Callbacks {
					fmt.Printf("    - %s: %v\n", callback.Event, callback.Params)
				}
			} else {
				fmt.Printf("  无回调函数\n")
			}

			// 验证预期的回调是否正确注入
			switch field.Code {
			case "username":
				// 应该只有标签中定义的OnInputFuzzy
				if len(field.Callbacks) != 1 || field.Callbacks[0].Event != "OnInputFuzzy" {
					t.Errorf("字段 username 的回调不正确，期望1个OnInputFuzzy，实际: %v", field.Callbacks)
				}

			case "company_name":
				// 应该有自动注入的OnInputFuzzy
				hasOnInputFuzzy := false
				for _, cb := range field.Callbacks {
					if cb.Event == "OnInputFuzzy" {
						hasOnInputFuzzy = true
						break
					}
				}
				if !hasOnInputFuzzy {
					t.Errorf("字段 company_name 应该有自动注入的OnInputFuzzy回调")
				}

			case "user_info":
				// 应该有自动注入的OnInputFuzzy和OnInputValidate
				hasOnInputFuzzy := false
				hasOnInputValidate := false
				for _, cb := range field.Callbacks {
					if cb.Event == "OnInputFuzzy" {
						hasOnInputFuzzy = true
					}
					if cb.Event == "OnInputValidate" {
						hasOnInputValidate = true
					}
				}
				if !hasOnInputFuzzy {
					t.Errorf("字段 user_info 应该有自动注入的OnInputFuzzy回调")
				}
				if !hasOnInputValidate {
					t.Errorf("字段 user_info 应该有自动注入的OnInputValidate回调")
				}

			case "email":
				// 应该有标签中定义的OnInputValidate，还会自动注入OnInputFuzzy
				validateCount := 0
				fuzzyCount := 0
				for _, cb := range field.Callbacks {
					if cb.Event == "OnInputValidate" {
						validateCount++
					}
					if cb.Event == "OnInputFuzzy" {
						fuzzyCount++
					}
				}
				if validateCount != 1 {
					t.Errorf("字段 email 应该只有1个OnInputValidate回调，实际有%d个", validateCount)
				}
				if fuzzyCount != 1 {
					t.Errorf("字段 email 应该有1个自动注入的OnInputFuzzy回调，实际有%d个", fuzzyCount)
				}

			case "description":
				// 应该没有任何回调
				if len(field.Callbacks) != 0 {
					t.Errorf("字段 description 不应该有任何回调，实际有: %v", field.Callbacks)
				}
			}
		}
	} else {
		t.Error("返回的不是FormConfig类型")
	}

	fmt.Println("\n=== 回调注入测试完成 ===")
}
