package api

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/yunhanshu-net/pkg/typex"
)

// 模拟callback_demo.go中的CallbackDemoReq结构体
type CallbackDemoReq struct {
	Username      string     `json:"username" form:"username" runner:"code:username;name:用户名" widget:"type:input;placeholder:请输入用户名" data:"type:string" validate:"required,min=2,max=20"`
	Email         string     `json:"email" form:"email" runner:"code:email;name:邮箱" widget:"type:input;placeholder:请输入邮箱地址" data:"type:string" validate:"required,email"`
	Department    string     `json:"department" form:"department" runner:"code:department;name:部门" widget:"type:select;placeholder:请选择部门;options:技术部,产品部,设计部,运营部,市场部,人事部,财务部" data:"type:string" validate:"required"`
	Position      string     `json:"position" form:"position" runner:"code:position;name:职位" widget:"type:select;placeholder:请选择职位;options:初级工程师,中级工程师,高级工程师,技术专家,产品经理,设计师,运营专员,系统管理员" data:"type:string" validate:"required"`
	Location      string     `json:"location" form:"location" runner:"code:location;name:工作地点" widget:"type:input;placeholder:请输入工作地点" data:"type:string" callback:"OnInputFuzzy(delay:300,min:1)" validate:"required"`
	CompanySearch string     `json:"company_search" form:"company_search" runner:"code:company_search;name:公司搜索" widget:"type:input;placeholder:搜索公司名称或行业" data:"type:string" callback:"OnInputFuzzy(delay:500,min:2)"`
	UserSearch    string     `json:"user_search" form:"user_search" runner:"code:user_search;name:用户搜索" widget:"type:input;placeholder:搜索用户名、邮箱或部门" data:"type:string" callback:"OnInputFuzzy(delay:300,min:2)"`
	WorkMode      string     `json:"work_mode" form:"work_mode" runner:"code:work_mode;name:工作模式" widget:"type:radio;options:远程,办公室,混合;direction:horizontal" data:"type:string;default_value:混合" validate:"required"`
	Priority      string     `json:"priority" form:"priority" runner:"code:priority;name:优先级" widget:"type:select;options:低,中,高,紧急" data:"type:string;default_value:中" validate:"required"`
	Skills        []string   `json:"skills" form:"skills" runner:"code:skills;name:技能标签" widget:"type:multiselect;placeholder:选择技能标签;options:Java,Python,Go,JavaScript,React,Vue,MySQL,Redis,Kubernetes,Docker;multiple_limit:8;collapse_tags:true" data:"type:[]string" validate:"required,min=1"`
	StartDate     typex.Time `json:"start_date" form:"start_date" runner:"code:start_date;name:开始日期" widget:"type:datetime;format:date;placeholder:请选择开始日期;min_date:$today" data:"type:string;default_value:$today;example:2025-01-15" validate:"required"`
	Notification  bool       `json:"notification" form:"notification" runner:"code:notification;name:消息通知" widget:"type:switch;true_label:开启;false_label:关闭" data:"type:boolean;default_value:true"`
	Remarks       string     `json:"remarks" form:"remarks" runner:"code:remarks;name:备注信息" widget:"type:input;mode:text_area;placeholder:请输入备注信息;max:500" data:"type:string"`
}

// MockCallbackDemoFunctionInfo 模拟callback_demo.go中的FunctionInfo
type MockCallbackDemoFunctionInfo struct{}

func (m *MockCallbackDemoFunctionInfo) GetOnInputFuzzyMap() map[string]interface{} {
	return map[string]interface{}{
		"company_search":   func() {}, // 对应callback_demo.go中的OnInputFuzzyMap["company_search"]
		"user_search":      func() {}, // 对应callback_demo.go中的OnInputFuzzyMap["user_search"]
		"location_suggest": func() {}, // 对应callback_demo.go中的OnInputFuzzyMap["location_suggest"]，但字段名是location
	}
}

func (m *MockCallbackDemoFunctionInfo) GetOnInputValidateMap() map[string]interface{} {
	return map[string]interface{}{
		"username": func() {}, // 对应callback_demo.go中的OnInputValidateMap["username"]
		"email":    func() {}, // 对应callback_demo.go中的OnInputValidateMap["email"]
	}
}

func TestRealCallbackDemo(t *testing.T) {
	fmt.Println("=== 测试真实的CallbackDemo功能 ===")

	// 创建模拟的FunctionInfo，对应callback_demo.go中的定义
	mockFunctionInfo := &MockCallbackDemoFunctionInfo{}

	// 测试带FunctionInfo的版本
	fmt.Println("\n--- 测试带FunctionInfo的版本（自动注入回调） ---")
	paramsWithFunctionInfo, err := NewRequestParamsWithFunctionInfo(CallbackDemoReq{}, "form", mockFunctionInfo)
	if err != nil {
		t.Fatalf("NewRequestParamsWithFunctionInfo failed: %v", err)
	}

	// 验证结果
	if formConfig, ok := paramsWithFunctionInfo.(*FormConfig); ok {
		fmt.Println("\n--- 回调注入效果验证 ---")

		for _, field := range formConfig.Fields {
			fmt.Printf("\n字段: %s (%s)\n", field.Code, field.Name)

			if len(field.Callbacks) > 0 {
				fmt.Printf("  回调函数:\n")
				for _, callback := range field.Callbacks {
					fmt.Printf("    - %s: %v\n", callback.Event, callback.Params)
				}
			} else {
				fmt.Printf("  无回调函数\n")
			}
		}

	} else {
		t.Error("返回的不是FormConfig类型")
	}

	// 输出最终的JSON结果
	finalJson, _ := json.MarshalIndent(paramsWithFunctionInfo, "", "  ")
	fmt.Printf("\n=== 最终的API参数配置 ===\n%s\n", string(finalJson))

	fmt.Println("\n=== 真实CallbackDemo测试完成 ===")
}
