package api

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/pkg/typex/files"
)

// NewTagSystemReq 新标签系统测试请求结构
type NewTagSystemReq struct {
	// 基础字符串字段
	Title string `json:"title" form:"title" runner:"code:title;name:标题" widget:"type:input;placeholder:请输入标题" data:"type:string;example:产品发布会;default_value:" validate:"required,max=200"`

	// 多行文本字段
	Description string `json:"description" form:"description" runner:"code:description;name:描述" widget:"type:input;mode:text_area;placeholder:请输入详细描述" data:"type:string;example:这是一个详细描述" validate:"max=1000"`

	// 选择器字段 - 简化options格式
	Category string `json:"category" form:"category" runner:"code:category;name:分类" widget:"type:select;placeholder:请选择分类;options:产品,营销,技术,运营" data:"type:string;default_value:产品" validate:"required"`

	// 多选字段 - []string类型
	Tags []string `json:"tags" form:"tags" runner:"code:tags;name:标签" widget:"type:multiselect;placeholder:选择相关标签;options:重要,紧急,公开,内部" data:"type:[]string;example:重要,公开" validate:"min=1"`

	// 数字字段
	Priority int `json:"priority" form:"priority" runner:"code:priority;name:优先级" widget:"type:slider" data:"type:number;default_value:5;example:8" validate:"min=1,max=10"`

	// 布尔字段
	IsPublic bool `json:"is_public" form:"is_public" runner:"code:is_public;name:是否公开" widget:"type:switch;true_label:公开;false_label:私有" data:"type:boolean;default_value:false"`

	// 带搜索功能的用户选择 - 字段级别回调
	AssignedUser string `json:"assigned_user" form:"assigned_user" runner:"code:assigned_user;name:负责人" widget:"type:select;placeholder:搜索并选择负责人" data:"type:string;source:api://users" callback:"OnInputFuzzy(delay:300,min:2)" validate:"required"`

	// 只读字段 - 新权限系统
	CreatedAt string `json:"created_at" form:"created_at" runner:"code:created_at;name:创建时间" widget:"type:input" data:"type:string" permission:"read"`

	// 仅创建时可编辑的字段
	ProjectType string `json:"project_type" form:"project_type" runner:"code:project_type;name:项目类型" widget:"type:select;options:内部项目,外部项目,合作项目" data:"type:string;default_value:内部项目" permission:"read,create" validate:"required"`

	// 带多个回调的复杂字段
	Budget float64 `json:"budget" form:"budget" runner:"code:budget;name:预算" widget:"type:input;placeholder:请输入预算金额" data:"type:number;example:10000.50" callback:"OnInputChange();OnBlur(validate:true,format:currency)" validate:"min=0"`

	// 文件上传字段 - 使用正确的files.Files指针类型
	Attachments *files.Files `json:"attachments" form:"attachments" runner:"code:attachments;name:附件" widget:"type:file_upload;accept:.pdf,.doc,.docx,.jpg,.png" data:"type:files" callback:"OnFileChange(max_size:10MB,max_count:5)"`
}

// NewTagSystemResp 新标签系统测试响应结构
type NewTagSystemResp struct {
	// 基础返回字段
	ID string `json:"id" runner:"code:id;name:任务ID" data:"type:string"`

	// 显示用户输入的数据
	Title       string   `json:"title" runner:"code:title;name:标题" data:"type:string"`
	Description string   `json:"description" runner:"code:description;name:描述" data:"type:string"`
	Category    string   `json:"category" runner:"code:category;name:分类" data:"type:string"`
	Tags        []string `json:"tags" runner:"code:tags;name:标签" data:"type:[]string"`
	Priority    int      `json:"priority" runner:"code:priority;name:优先级" data:"type:number"`
	IsPublic    bool     `json:"is_public" runner:"code:is_public;name:是否公开" data:"type:boolean"`

	// 系统生成字段
	Status      string `json:"status" runner:"code:status;name:状态" data:"type:string;example:已创建"`
	ProcessTime string `json:"process_time" runner:"code:process_time;name:处理时间" data:"type:string"`
}

// PermissionTestReq 权限控制测试请求
type PermissionTestReq struct {
	// 默认权限（全部权限）
	NormalField string `json:"normal_field" form:"normal_field" runner:"code:normal_field;name:普通字段" widget:"type:input" data:"type:string"`

	// 只读权限
	ReadOnlyField string `json:"readonly_field" form:"readonly_field" runner:"code:readonly_field;name:只读字段" widget:"type:input" data:"type:string" permission:"read"`

	// 读+创建权限
	CreateOnlyField string `json:"create_only_field" form:"create_only_field" runner:"code:create_only_field;name:仅创建字段" widget:"type:input" data:"type:string" permission:"read,create"`

	// 读+更新权限
	UpdateOnlyField string `json:"update_only_field" form:"update_only_field" runner:"code:update_only_field;name:仅更新字段" widget:"type:input" data:"type:string" permission:"read,update"`
}

// CallbackTestReq 回调功能测试请求
type CallbackTestReq struct {
	// 单个回调
	SearchField string `json:"search_field" form:"search_field" runner:"code:search_field;name:搜索字段" widget:"type:input;placeholder:输入搜索内容" data:"type:string" callback:"OnInputFuzzy(delay:300,min:2)"`

	// 多个回调
	ComplexField string `json:"complex_field" form:"complex_field" runner:"code:complex_field;name:复杂字段" widget:"type:input" data:"type:string" callback:"OnInputChange();OnBlur(validate:true);OnFocus(highlight:true)"`

	// 带参数的回调
	AmountField float64 `json:"amount_field" form:"amount_field" runner:"code:amount_field;name:金额字段" widget:"type:input" data:"type:float" callback:"OnValueChange(format:currency,precision:2,min:0)"`
}

// FileTestReq 文件类型测试请求
type FileTestReq struct {
	// 单文件上传 - 使用files.Files指针类型（参考旧版本）
	Avatar *files.Files `json:"avatar" form:"avatar" runner:"code:avatar;name:头像" widget:"type:file_upload;accept:.jpg,.png,.gif" data:"type:files" callback:"OnFileChange(max_size:2MB,max_count:1)"`

	// 多文件上传 - files.Files指针类型（参考旧版本file_manager.go）
	Documents *files.Files `json:"documents" form:"documents" runner:"code:documents;name:文档列表" widget:"type:file_upload;accept:.pdf,.doc,.docx" data:"type:files" callback:"OnFileChange(max_size:10MB,max_count:5)"`

	// 响应文件类型 - files.Writer
	ProcessedFiles files.Writer `json:"processed_files" form:"processed_files" runner:"code:processed_files;name:处理后的文件" widget:"type:file_display" data:"type:files"`
}

// TestNewTagSystemRequest 测试新标签系统的请求参数解析
func TestNewTagSystemRequest(t *testing.T) {
	fmt.Println("=== 测试新标签系统 - 请求参数 ===")

	params, err := NewRequestParams(&NewTagSystemReq{}, response.RenderTypeForm)
	if err != nil {
		t.Fatalf("解析请求参数失败: %v", err)
	}

	marshal, err := json.MarshalIndent(params, "", "  ")
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	fmt.Println("新标签系统请求参数JSON:")
	fmt.Println(string(marshal))

	// 尝试类型断言来验证字段
	if formConfig, ok := params.(*FormConfig); ok {
		if len(formConfig.Fields) == 0 {
			t.Error("字段列表为空")
		}
		fmt.Printf("解析到 %d 个字段\n", len(formConfig.Fields))
	} else {
		fmt.Printf("返回类型: %T\n", params)
	}
}

// TestNewTagSystemResponse 测试新标签系统的响应参数解析
func TestNewTagSystemResponse(t *testing.T) {
	fmt.Println("\n=== 测试新标签系统 - 响应参数 ===")

	params, err := NewResponseParams(&NewTagSystemResp{}, response.RenderTypeForm)
	if err != nil {
		t.Fatalf("解析响应参数失败: %v", err)
	}

	marshal, err := json.MarshalIndent(params, "", "  ")
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	fmt.Println("新标签系统响应参数JSON:")
	fmt.Println(string(marshal))

	// 尝试类型断言来验证字段
	if formConfig, ok := params.(*FormConfig); ok {
		fmt.Printf("解析到 %d 个字段\n", len(formConfig.Fields))
	} else {
		fmt.Printf("返回类型: %T\n", params)
	}
}

// TestPermissionControl 测试权限控制功能
func TestPermissionControl(t *testing.T) {
	fmt.Println("\n=== 测试权限控制 ===")

	params, err := NewRequestParams(&PermissionTestReq{}, response.RenderTypeForm)
	if err != nil {
		t.Fatalf("解析权限测试参数失败: %v", err)
	}

	marshal, err := json.MarshalIndent(params, "", "  ")
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	fmt.Println("权限控制测试JSON:")
	fmt.Println(string(marshal))

	// 验证权限字段
	if formConfig, ok := params.(*FormConfig); ok {
		for _, field := range formConfig.Fields {
			if field.Permission != nil {
				fmt.Printf("字段 %s 的权限: read=%v, update=%v, create=%v\n",
					field.Code, field.Permission.Read, field.Permission.Update, field.Permission.Create)
			} else {
				fmt.Printf("字段 %s 的权限: 默认权限（无限制）\n", field.Code)
			}
		}
	}
}

// TestCallbackFunctions 测试回调函数功能
func TestCallbackFunctions(t *testing.T) {
	fmt.Println("\n=== 测试回调函数 ===")

	params, err := NewRequestParams(&CallbackTestReq{}, response.RenderTypeForm)
	if err != nil {
		t.Fatalf("解析回调测试参数失败: %v", err)
	}

	marshal, err := json.MarshalIndent(params, "", "  ")
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	fmt.Println("回调函数测试JSON:")
	fmt.Println(string(marshal))

	// 验证回调字段
	if formConfig, ok := params.(*FormConfig); ok {
		for _, field := range formConfig.Fields {
			if len(field.Callbacks) > 0 {
				fmt.Printf("字段 %s 的回调函数:\n", field.Code)
				for _, callback := range field.Callbacks {
					fmt.Printf("  - %s: %v\n", callback.Event, callback.Params)
				}
			}
		}
	}
}

// TestFileTypes 测试文件类型字段
func TestFileTypes(t *testing.T) {
	fmt.Println("\n=== 测试文件类型 ===")

	params, err := NewRequestParams(&FileTestReq{}, response.RenderTypeForm)
	if err != nil {
		t.Fatalf("解析文件测试参数失败: %v", err)
	}

	marshal, err := json.MarshalIndent(params, "", "  ")
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	fmt.Println("文件类型测试JSON:")
	fmt.Println(string(marshal))

	// 验证文件类型字段
	if formConfig, ok := params.(*FormConfig); ok {
		for _, field := range formConfig.Fields {
			fmt.Printf("字段 %s 的类型: %s\n", field.Code, field.Data.Type)
			if field.Widget.Type == "file_upload" {
				fmt.Printf("  - Widget配置: %v\n", field.Widget.Config)
			}
		}
	}
}

// TestTableRequest 测试表格请求参数
func TestTableRequest(t *testing.T) {
	fmt.Println("\n=== 测试表格请求参数 ===")

	params, err := NewRequestParams(&NewTagSystemReq{}, response.RenderTypeTable)
	if err != nil {
		t.Fatalf("解析表格请求参数失败: %v", err)
	}

	marshal, err := json.MarshalIndent(params, "", "  ")
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	fmt.Println("表格请求参数JSON:")
	fmt.Println(string(marshal))
}

// TestComplexTypes 测试复合类型
func TestComplexTypes(t *testing.T) {
	fmt.Println("\n=== 测试复合类型 ===")

	type ComplexTypeReq struct {
		StringArray []string               `json:"string_array" runner:"code:string_array;name:字符串数组" data:"type:[]string"`
		NumberArray []int                  `json:"number_array" runner:"code:number_array;name:数字数组" data:"type:[]number"`
		ObjectField map[string]interface{} `json:"object_field" runner:"code:object_field;name:对象字段" data:"type:object"`
		FilesField  files.Files            `json:"files_field" runner:"code:files_field;name:文件字段" data:"type:files"`
	}

	params, err := NewRequestParams(&ComplexTypeReq{}, response.RenderTypeForm)
	if err != nil {
		t.Fatalf("解析复合类型参数失败: %v", err)
	}

	marshal, err := json.MarshalIndent(params, "", "  ")
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	fmt.Println("复合类型测试JSON:")
	fmt.Println(string(marshal))

	// 验证类型推断
	if formConfig, ok := params.(*FormConfig); ok {
		for _, field := range formConfig.Fields {
			fmt.Printf("字段 %s 的推断类型: %s\n", field.Code, field.Data.Type)
		}
	}
}

// TestValidationRules 测试验证规则
func TestValidationRules(t *testing.T) {
	fmt.Println("\n=== 测试验证规则 ===")

	type ValidationTestReq struct {
		RequiredField string  `json:"required_field" validate:"required"`
		EmailField    string  `json:"email_field" validate:"required,email"`
		NumberField   int     `json:"number_field" validate:"min=1,max=100"`
		LengthField   string  `json:"length_field" validate:"min=3,max=50"`
		FloatField    float64 `json:"float_field" validate:"min=0.1,max=999.99"`
	}

	params, err := NewRequestParams(&ValidationTestReq{}, response.RenderTypeForm)
	if err != nil {
		t.Fatalf("解析验证规则参数失败: %v", err)
	}

	marshal, err := json.MarshalIndent(params, "", "  ")
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	fmt.Println("验证规则测试JSON:")
	fmt.Println(string(marshal))

	// 验证验证规则
	if formConfig, ok := params.(*FormConfig); ok {
		for _, field := range formConfig.Fields {
			if field.Validation != "" {
				fmt.Printf("字段 %s 的验证规则: %s\n", field.Code, field.Validation)
			}
		}
	}
}

// TestSimpleTagParsing 测试简单的标签解析
func TestSimpleTagParsing(t *testing.T) {
	fmt.Println("\n=== 测试简单标签解析 ===")

	type SimpleReq struct {
		Name string `json:"name" form:"name" runner:"code:name;name:姓名" widget:"type:input;placeholder:请输入姓名" data:"type:string;example:张三" validate:"required"`
	}

	params, err := NewRequestParams(&SimpleReq{}, response.RenderTypeForm)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	marshal, err := json.MarshalIndent(params, "", "  ")
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	fmt.Println("简单标签解析结果:")
	fmt.Println(string(marshal))

	// 验证字段
	if formConfig, ok := params.(*FormConfig); ok {
		if len(formConfig.Fields) > 0 {
			field := formConfig.Fields[0]
			fmt.Printf("字段代码: %s\n", field.Code)
			fmt.Printf("字段名称: %s\n", field.Name)
			fmt.Printf("Widget类型: %s\n", field.Widget.Type)
			if placeholder, ok := field.Widget.Config["placeholder"]; ok {
				fmt.Printf("占位符: %s\n", placeholder)
			}
			fmt.Printf("数据类型: %s\n", field.Data.Type)
			fmt.Printf("示例值: %s\n", field.Data.Example)
			fmt.Printf("验证规则: %s\n", field.Validation)
		}
	}
}

// TestTableMode 测试table模式下的权限控制
func TestTableMode(t *testing.T) {
	fmt.Println("\n=== 测试Table模式权限控制 ===")

	params, err := NewRequestParams(&PermissionTestReq{}, response.RenderTypeTable)
	if err != nil {
		t.Fatalf("解析table模式失败: %v", err)
	}

	marshal, err := json.MarshalIndent(params, "", "  ")
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	fmt.Println("Table模式JSON:")
	fmt.Println(string(marshal))

	// 验证table模式下权限字段的处理
	if tableConfig, ok := params.(*TableConfig); ok {
		fmt.Printf("Table模式解析到 %d 个列\n", len(tableConfig.Columns))

		for _, column := range tableConfig.Columns {
			if column.Permission == nil {
				fmt.Printf("- 字段 %s: 无权限限制（前端按全部权限处理）\n", column.Code)
			} else {
				fmt.Printf("- 字段 %s: read=%v, update=%v, create=%v\n",
					column.Code, column.Permission.Read, column.Permission.Update, column.Permission.Create)
			}
		}

		// 验证具体的权限处理逻辑
		// 查找有permission标签的字段
		var readOnlyField *FieldInfo
		var normalField *FieldInfo

		for _, column := range tableConfig.Columns {
			if column.Code == "readonly_field" {
				readOnlyField = column
			}
			if column.Code == "normal_field" {
				normalField = column
			}
		}

		// 有permission标签的字段应该有Permission配置
		if readOnlyField != nil && readOnlyField.Permission == nil {
			t.Errorf("有permission标签的字段 %s 应该有Permission配置", readOnlyField.Code)
		}

		// 没有permission标签的字段应该返回null（表示无权限限制）
		if normalField != nil && normalField.Permission != nil {
			t.Errorf("没有permission标签的字段 %s 应该返回null，表示无权限限制", normalField.Code)
		}
	} else {
		t.Errorf("期望返回TableConfig，实际返回: %T", params)
	}
}

// TestComplexWidgetConfig 测试复杂的Widget配置解析
func TestComplexWidgetConfig(t *testing.T) {
	fmt.Println("\n=== 测试复杂Widget配置 ===")

	type ComplexWidgetReq struct {
		// 文件上传组件 - 支持多种配置
		Avatar *files.Files `json:"avatar" form:"avatar" runner:"code:avatar;name:头像" widget:"type:file_upload;accept:.jpg,.png,.gif;max_size:2MB;max_count:1;preview:true;drag_drop:true;placeholder:拖拽或点击上传头像" data:"type:files" callback:"OnFileChange(max_size:2MB,max_count:1)"`

		// 颜色选择器 - 支持透明度和预设颜色
		ThemeColor string `json:"theme_color" form:"theme_color" runner:"code:theme_color;name:主题色" widget:"type:color;format:rgba;show_alpha:true;predefine:#FF0000,#00FF00,#0000FF;placeholder:选择主题颜色" data:"type:string;default_value:rgba(64,158,255,0.8)"`

		// 日期时间选择器 - 支持范围和快捷选项
		EventTime string `json:"event_time" form:"event_time" runner:"code:event_time;name:活动时间" widget:"type:datetime;format:datetimerange;separator:至;start_placeholder:开始时间;end_placeholder:结束时间;default_time:12:00:00" data:"type:string"`

		// 滑块组件 - 支持范围选择
		PriceRange string `json:"price_range" form:"price_range" runner:"code:price_range;name:价格区间" widget:"type:slider;mode:range;min:0;max:10000;step:100" data:"type:string;default_value:1000,5000"`

		// 多选组件 - 支持搜索
		Tags []string `json:"tags" form:"tags" runner:"code:tags;name:标签" widget:"type:multiselect;placeholder:搜索并选择标签;options:重要,紧急,公开,内部,机密;filterable:true;multiple_limit:5" data:"type:[]string"`
	}

	params, err := NewRequestParams(&ComplexWidgetReq{}, response.RenderTypeForm)
	if err != nil {
		t.Fatalf("解析复杂Widget配置失败: %v", err)
	}

	marshal, err := json.MarshalIndent(params, "", "  ")
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	fmt.Println("复杂Widget配置JSON:")
	fmt.Println(string(marshal))

	// 验证各种组件的配置
	if formConfig, ok := params.(*FormConfig); ok {
		for _, field := range formConfig.Fields {
			fmt.Printf("\n字段 %s (%s) 的Widget配置:\n", field.Code, field.Widget.Type)
			for key, value := range field.Widget.Config {
				fmt.Printf("  - %s: %v\n", key, value)
			}
		}
	}
}
