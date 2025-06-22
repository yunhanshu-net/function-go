package api

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/yunhanshu-net/pkg/typex/files"
)

// DebugTestProduct 测试产品结构体 - 使用新的标签系统
type DebugTestProduct struct {
	ID       int          `json:"id" form:"id" runner:"code:id;name:产品ID" data:"type:number;example:1" permission:"read"`
	Name     string       `json:"name" form:"name" runner:"code:name;name:产品名称" widget:"type:input;placeholder:请输入产品名称" data:"type:string;example:iPhone 15" validate:"required"`
	Category string       `json:"category" form:"category" runner:"code:category;name:产品分类" widget:"type:select;options:电子产品,服装,食品,图书" data:"type:string;default_value:电子产品" validate:"required"`
	Price    float64      `json:"price" form:"price" runner:"code:price;name:产品价格" widget:"type:input;prefix:¥;precision:2" data:"type:float;example:999.99" validate:"required;min:0"`
	Status   bool         `json:"status" form:"status" runner:"code:status;name:产品状态" widget:"type:switch;true_label:启用;false_label:禁用" data:"type:boolean;default_value:true" validate:"required"`
	Image    *files.Files `json:"image" form:"image" runner:"code:image;name:产品图片" widget:"type:file_upload;accept:.jpg,.png;max_size:2MB" data:"type:file"`
}

func TestNewTagSystemDebug(t *testing.T) {
	fmt.Println("=== 测试新的标签系统 ===")

	// 测试Form模式
	fmt.Println("\n--- 测试Form模式 ---")
	formResponse, err := NewUnifiedFormResponse(DebugTestProduct{}, "form")
	if err != nil {
		t.Fatalf("NewUnifiedFormResponse failed: %v", err)
	}

	// 打印结果
	formJson, _ := json.MarshalIndent(formResponse, "", "  ")
	fmt.Printf("Form响应结果:\n%s\n", string(formJson))

	// 测试Table模式
	fmt.Println("\n--- 测试Table模式 ---")
	tableResponse, err := NewUnifiedTableResponse([]DebugTestProduct{})
	if err != nil {
		t.Fatalf("NewUnifiedTableResponse failed: %v", err)
	}

	// 打印结果
	tableJson, _ := json.MarshalIndent(tableResponse, "", "  ")
	fmt.Printf("Table响应结果:\n%s\n", string(tableJson))

	// 测试新的RequestParams接口
	fmt.Println("\n--- 测试NewRequestParams接口 ---")
	requestParams, err := NewRequestParams(DebugTestProduct{}, "form")
	if err != nil {
		t.Fatalf("NewRequestParams failed: %v", err)
	}

	requestJson, _ := json.MarshalIndent(requestParams, "", "  ")
	fmt.Printf("Request响应结果:\n%s\n", string(requestJson))

	// 验证字段数量
	if len(formResponse.Fields) == 0 {
		t.Error("Form响应中没有字段")
	}

	if len(tableResponse.Columns) == 0 {
		t.Error("Table响应中没有列")
	}

	fmt.Println("\n=== 测试完成 ===")
}

func TestOldInterfaceCompatibility(t *testing.T) {
	fmt.Println("=== 测试旧接口兼容性 ===")

	// 测试旧的NewFormRequestParams接口
	fmt.Println("\n--- 测试旧的NewFormRequestParams接口 ---")
	oldFormParams, err := NewFormRequestParams(DebugTestProduct{}, "form")
	if err != nil {
		t.Fatalf("NewFormRequestParams failed: %v", err)
	}

	oldFormJson, _ := json.MarshalIndent(oldFormParams, "", "  ")
	fmt.Printf("旧Form接口响应结果:\n%s\n", string(oldFormJson))

	// 验证旧接口是否还在返回旧格式
	if len(oldFormParams.Children) == 0 {
		t.Error("旧Form接口没有返回children字段")
	}

	fmt.Println("\n=== 旧接口测试完成 ===")
}

func TestDebugErrorHandling(t *testing.T) {
	fmt.Println("=== 测试错误处理 ===")

	// 测试非结构体类型
	_, err := NewRequestParams("not a struct", "form")
	if err == nil {
		t.Error("应该返回错误，因为输入不是结构体")
	}
	fmt.Printf("预期错误: %v\n", err)

	// 测试空结构体
	type EmptyStruct struct{}
	emptyResponse, err := NewRequestParams(EmptyStruct{}, "form")
	if err != nil {
		t.Fatalf("空结构体应该成功处理: %v", err)
	}

	emptyJson, _ := json.MarshalIndent(emptyResponse, "", "  ")
	fmt.Printf("空结构体响应:\n%s\n", string(emptyJson))

	fmt.Println("\n=== 错误处理测试完成 ===")
}
