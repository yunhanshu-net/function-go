package api

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/pkg/typex/files"
)

// OldTagSystemReq 旧标签系统示例（所有配置在runner标签中）
type OldTagSystemReq struct {
	// 旧版本的标签格式 - 所有配置都在runner标签中
	Title string `json:"title" form:"title" 
			runner:"code:title;name:标题;type:string;widget:input;placeholder:请输入标题;example:产品发布会" 
			validate:"required,max=200"`

	Category string `json:"category" form:"category"
				runner:"code:category;name:分类;type:string;widget:select;placeholder:请选择分类;options:产品,营销,技术,运营;default_value:产品"
				validate:"required"`

	Files *files.Files `json:"files" form:"files"
		  runner:"code:files;name:文件列表;type:files;widget:file_upload;multiple:true;max_size:10240;accept:.pdf,.doc,.docx,.jpg,.png"
		  validate:"required"`

	IsActive bool `json:"is_active" form:"is_active"
			 runner:"code:is_active;name:是否启用;type:boolean;widget:switch;true_label:启用;false_label:禁用;default_value:true"`

	Priority int `json:"priority" form:"priority"
			 runner:"code:priority;name:优先级;type:number;widget:slider;min:1;max:10;default_value:5"`
}

// NewTagSystemReq 新标签系统示例（标签分离）
type NewTagSystemReq2 struct {
	// 新版本的标签格式 - 标签分离
	Title string `json:"title" form:"title"
			runner:"code:title;name:标题"
			widget:"type:input;placeholder:请输入标题"
			data:"type:string;example:产品发布会"
			validate:"required,max=200"`

	Category string `json:"category" form:"category"
				runner:"code:category;name:分类"
				widget:"type:select;placeholder:请选择分类;options:产品,营销,技术,运营"
				data:"type:string;default_value:产品"
				validate:"required"`

	Files *files.Files `json:"files" form:"files"
		  runner:"code:files;name:文件列表"
		  widget:"type:file_upload;accept:.pdf,.doc,.docx,.jpg,.png"
		  data:"type:files"
		  callback:"OnFileChange(max_size:10MB,max_count:5)"
		  validate:"required"`

	IsActive bool `json:"is_active" form:"is_active"
			 runner:"code:is_active;name:是否启用"
			 widget:"type:switch;true_label:启用;false_label:禁用"
			 data:"type:boolean;default_value:true"`

	Priority int `json:"priority" form:"priority"
			 runner:"code:priority;name:优先级"
			 widget:"type:slider"
			 data:"type:number;default_value:5;example:8"
			 validate:"min=1,max=10"`

	// 新标签系统的权限控制示例
	CreatedAt string `json:"created_at" form:"created_at"
			  runner:"code:created_at;name:创建时间"
			  widget:"type:input"
			  data:"type:string"
			  permission:"read"`

	// 新标签系统的回调示例
	SearchField string `json:"search_field" form:"search_field"
				runner:"code:search_field;name:搜索字段"
				widget:"type:input;placeholder:输入搜索内容"
				data:"type:string"
				callback:"OnInputFuzzy(delay:300,min:2)"`
}

// TestOldTagSystem 测试旧标签系统
func TestOldTagSystem(t *testing.T) {
	fmt.Println("=== 测试旧标签系统（所有配置在runner中） ===")

	params, err := NewRequestParams(&OldTagSystemReq{}, response.RenderTypeForm)
	if err != nil {
		t.Fatalf("解析旧标签系统失败: %v", err)
	}

	marshal, err := json.MarshalIndent(params, "", "  ")
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	fmt.Println("旧标签系统JSON:")
	fmt.Println(string(marshal))
}

// TestNewTagSystem 测试新标签系统
func TestNewTagSystem(t *testing.T) {
	fmt.Println("\n=== 测试新标签系统（标签分离） ===")

	params, err := NewRequestParams(&NewTagSystemReq2{}, response.RenderTypeForm)
	if err != nil {
		t.Fatalf("解析新标签系统失败: %v", err)
	}

	marshal, err := json.MarshalIndent(params, "", "  ")
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	fmt.Println("新标签系统JSON:")
	fmt.Println(string(marshal))
}

// TestTagSystemComparison 对比新旧标签系统
func TestTagSystemComparison(t *testing.T) {
	fmt.Println("\n=== 新旧标签系统对比测试 ===")

	// 测试旧系统
	oldParams, err := NewRequestParams(&OldTagSystemReq{}, response.RenderTypeForm)
	if err != nil {
		t.Fatalf("解析旧标签系统失败: %v", err)
	}

	// 测试新系统
	newParams, err := NewRequestParams(&NewTagSystemReq2{}, response.RenderTypeForm)
	if err != nil {
		t.Fatalf("解析新标签系统失败: %v", err)
	}

	// 对比字段数量
	var oldFieldCount, newFieldCount int

	if oldConfig, ok := oldParams.(*FormConfig); ok {
		oldFieldCount = len(oldConfig.Fields)
		fmt.Printf("旧系统字段数量: %d\n", oldFieldCount)
	}

	if newConfig, ok := newParams.(*FormConfig); ok {
		newFieldCount = len(newConfig.Fields)
		fmt.Printf("新系统字段数量: %d\n", newFieldCount)

		// 展示新系统的特色功能
		fmt.Println("\n新系统特色功能展示:")
		for _, field := range newConfig.Fields {
			// 权限控制
			if field.Permission != nil && (!field.Permission.Read || !field.Permission.Update || !field.Permission.Create) {
				fmt.Printf("- 字段 %s 有权限控制: read=%v, update=%v, create=%v\n",
					field.Code, field.Permission.Read, field.Permission.Update, field.Permission.Create)
			}

			// 回调函数
			if len(field.Callbacks) > 0 {
				fmt.Printf("- 字段 %s 有回调函数: ", field.Code)
				for _, callback := range field.Callbacks {
					fmt.Printf("%s ", callback.Event)
				}
				fmt.Println()
			}

			// 数据配置
			if field.Data.Type != "" || field.Data.Example != "" || field.Data.DefaultValue != "" {
				fmt.Printf("- 字段 %s 数据配置: type=%s, example=%s, default=%s\n",
					field.Code, field.Data.Type, field.Data.Example, field.Data.DefaultValue)
			}
		}
	}

	fmt.Printf("\n对比结果: 新系统支持更多功能，包括权限控制、回调分离、数据配置分离等\n")
}

// TestFileTypeRecognition 测试文件类型识别
func TestFileTypeRecognition(t *testing.T) {
	fmt.Println("\n=== 测试文件类型识别 ===")

	type FileTypeTest struct {
		SingleFile  *files.Files `json:"single_file" runner:"code:single_file;name:单文件" data:"type:files"`
		WriterFile  files.Writer `json:"writer_file" runner:"code:writer_file;name:写入文件" data:"type:files"`
		StringArray []string     `json:"string_array" runner:"code:string_array;name:字符串数组" data:"type:[]string"`
		NumberArray []int        `json:"number_array" runner:"code:number_array;name:数字数组" data:"type:[]number"`
	}

	params, err := NewRequestParams(&FileTypeTest{}, response.RenderTypeForm)
	if err != nil {
		t.Fatalf("解析文件类型测试失败: %v", err)
	}

	if formConfig, ok := params.(*FormConfig); ok {
		fmt.Println("文件类型识别结果:")
		for _, field := range formConfig.Fields {
			fmt.Printf("- 字段 %s: 类型=%s\n", field.Code, field.Data.Type)
		}
	}
}

// TestValidationRulesHandling 测试验证规则处理
func TestValidationRulesHandling(t *testing.T) {
	fmt.Println("\n=== 测试验证规则处理 ===")

	type ValidationTest struct {
		RequiredField string  `json:"required_field" validate:"required"`
		EmailField    string  `json:"email_field" validate:"required,email,max=100"`
		NumberField   int     `json:"number_field" validate:"min=1,max=1000"`
		FloatField    float64 `json:"float_field" validate:"min=0.1,max=999.99"`
	}

	params, err := NewRequestParams(&ValidationTest{}, response.RenderTypeForm)
	if err != nil {
		t.Fatalf("解析验证规则测试失败: %v", err)
	}

	if formConfig, ok := params.(*FormConfig); ok {
		fmt.Println("验证规则处理结果:")
		for _, field := range formConfig.Fields {
			if field.Validation != "" {
				fmt.Printf("- 字段 %s: 验证规则=%s\n", field.Code, field.Validation)
			}
		}
	}
}
