package api

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/yunhanshu-net/function-go/pkg/dto/response"
)

// 使用示例：批量任务配置
type BatchTaskConfig struct {
	// 基础配置
	TaskName string `json:"task_name" runner:"code:task_name;name:任务名称" widget:"type:input;placeholder:请输入任务名称" data:"type:string" validate:"required"`
	Enabled  bool   `json:"enabled" runner:"code:enabled;name:是否启用" widget:"type:switch;true_label:启用;false_label:禁用" data:"type:boolean;default_value:true"`

	// 列表输入：批量操作列表
	Operations []Operation `json:"operations" runner:"code:operations;name:操作列表" widget:"type:list_input;placeholder:添加操作" data:"type:[]struct"`
}

// 单个操作配置
type Operation struct {
	Name        string        `json:"name" runner:"code:name;name:操作名称" widget:"type:input;placeholder:请输入操作名称" data:"type:string" validate:"required"`
	Type        string        `json:"type" runner:"code:type;name:操作类型" widget:"type:select;options:create,update,delete,query;placeholder:请选择操作类型" data:"type:string" validate:"required"`
	Description string        `json:"description" runner:"code:description;name:操作描述" widget:"type:input;placeholder:请输入操作描述" data:"type:string"`
	
	// 嵌套表单：操作参数配置
	Params OperationParams `json:"params" runner:"code:params;name:操作参数" widget:"type:form;title:参数配置" data:"type:struct"`
}

// 操作参数配置
type OperationParams struct {
	Timeout     int    `json:"timeout" runner:"code:timeout;name:超时时间" widget:"type:input;placeholder:请输入超时时间(秒)" data:"type:number;default_value:30"`
	RetryCount  int    `json:"retry_count" runner:"code:retry_count;name:重试次数" widget:"type:input;placeholder:请输入重试次数" data:"type:number;default_value:3"`
	Priority    string `json:"priority" runner:"code:priority;name:优先级" widget:"type:select;options:low,normal,high,critical;placeholder:请选择优先级" data:"type:string;default_value:normal"`
	Async       bool   `json:"async" runner:"code:async;name:异步执行" widget:"type:switch;true_label:是;false_label:否" data:"type:boolean;default_value:false"`
}

// ExampleComplexWidgets 展示复杂组件的使用
func ExampleComplexWidgets() {
	fmt.Println("=== 复杂组件使用示例 ===")

	// 构建表单配置
	builder := NewFormBuilder()
	config, err := builder.BuildFormConfig(reflect.TypeOf(BatchTaskConfig{}), response.RenderTypeForm)
	if err != nil {
		fmt.Printf("构建配置失败: %v\n", err)
		return
	}

	// 输出JSON配置
	marshal, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Printf("JSON序列化失败: %v\n", err)
		return
	}

	fmt.Println("生成的配置:")
	fmt.Println(string(marshal))

	// 分析字段结构
	fmt.Println("\n=== 字段分析 ===")
	for _, field := range config.Fields {
		fmt.Printf("字段: %s (%s)\n", field.Name, field.Code)
		fmt.Printf("  数据类型: %s\n", field.Data.Type)
		fmt.Printf("  组件类型: %s\n", field.Widget.Type)
		
		if field.Widget.Type == "list_input" {
			fmt.Printf("  -> 列表输入组件，支持动态添加/删除项目\n")
		} else if field.Widget.Type == "form" {
			fmt.Printf("  -> 嵌套表单组件，包含子字段配置\n")
		}
		fmt.Println()
	}
}

// 使用示例：用户管理配置
type UserManagementConfig struct {
	// 基础设置
	EnableUserMgmt bool   `json:"enable_user_mgmt" runner:"code:enable_user_mgmt;name:启用用户管理" widget:"type:switch;true_label:启用;false_label:禁用" data:"type:boolean;default_value:true"`
	AdminEmail     string `json:"admin_email" runner:"code:admin_email;name:管理员邮箱" widget:"type:input;placeholder:请输入管理员邮箱" data:"type:string" validate:"required,email"`

	// 列表输入：用户角色配置
	Roles []UserRole `json:"roles" runner:"code:roles;name:用户角色" widget:"type:list_input;placeholder:添加角色" data:"type:[]struct"`
}

// 用户角色配置
type UserRole struct {
	RoleName  string   `json:"role_name" runner:"code:role_name;name:角色名称" widget:"type:input;placeholder:请输入角色名称" data:"type:string" validate:"required"`
	RoleCode  string   `json:"role_code" runner:"code:role_code;name:角色代码" widget:"type:input;placeholder:请输入角色代码" data:"type:string" validate:"required"`
	IsActive  bool     `json:"is_active" runner:"code:is_active;name:是否激活" widget:"type:switch;true_label:激活;false_label:禁用" data:"type:boolean;default_value:true"`
	Tags      []string `json:"tags" runner:"code:tags;name:角色标签" widget:"type:multiselect;options:管理员,普通用户,访客,开发者,测试员;placeholder:选择角色标签" data:"type:[]string"`

	// 嵌套表单：权限配置
	Permissions RolePermissions `json:"permissions" runner:"code:permissions;name:权限配置" widget:"type:form;title:权限设置" data:"type:struct"`
}

// 角色权限配置
type RolePermissions struct {
	CanRead   bool `json:"can_read" runner:"code:can_read;name:读取权限" widget:"type:switch;true_label:允许;false_label:拒绝" data:"type:boolean;default_value:true"`
	CanWrite  bool `json:"can_write" runner:"code:can_write;name:写入权限" widget:"type:switch;true_label:允许;false_label:拒绝" data:"type:boolean;default_value:false"`
	CanDelete bool `json:"can_delete" runner:"code:can_delete;name:删除权限" widget:"type:switch;true_label:允许;false_label:拒绝" data:"type:boolean;default_value:false"`
	CanAdmin  bool `json:"can_admin" runner:"code:can_admin;name:管理权限" widget:"type:switch;true_label:允许;false_label:拒绝" data:"type:boolean;default_value:false"`
}

// ExampleUserManagement 展示用户管理配置示例
func ExampleUserManagement() {
	fmt.Println("\n=== 用户管理配置示例 ===")

	builder := NewFormBuilder()
	config, err := builder.BuildFormConfig(reflect.TypeOf(UserManagementConfig{}), response.RenderTypeForm)
	if err != nil {
		fmt.Printf("构建配置失败: %v\n", err)
		return
	}

	// 输出JSON配置
	marshal, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Printf("JSON序列化失败: %v\n", err)
		return
	}

	fmt.Println("用户管理配置:")
	fmt.Println(string(marshal))
} 