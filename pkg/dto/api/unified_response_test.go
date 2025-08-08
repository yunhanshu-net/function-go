package api

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/yunhanshu-net/function-go/pkg/dto/response"

	"github.com/yunhanshu-net/pkg/query"

	"github.com/yunhanshu-net/pkg/typex"
)

// TestProduct 测试用的产品结构
type TestProduct struct {
	ID        int        `json:"id" runner:"code:id;name:产品ID" data:"type:number"`
	Name      string     `json:"name" runner:"code:name;name:产品名称" widget:"type:input;placeholder:请输入产品名称" data:"type:string" validate:"required"`
	Price     float64    `json:"price" runner:"code:price;name:产品价格" widget:"type:input;prefix:¥;precision:2" data:"type:float" validate:"required,min=0"`
	Status    bool       `json:"status" runner:"code:status;name:产品状态" widget:"type:switch;true_label:启用;false_label:禁用" data:"type:boolean;default_value:true" permission:"read,update"`
	CreatedAt typex.Time `json:"created_at" runner:"code:created_at;name:创建时间" widget:"type:datetime;format:datetime" data:"type:string;default_value:$now" permission:"read"`
	UpdatedAt typex.Time `json:"updated_at" runner:"code:updated_at;name:更新时间" widget:"type:datetime;format:datetime" data:"type:string" permission:"read"`
}

// TestProductListResp 测试用的产品列表响应
type TestProductListResp struct {
	Items []TestProduct `json:"items"`
}

// TestProductReq 测试用的产品请求
type TestProductReq struct {
	Name       string     `json:"name" form:"name" runner:"code:name;name:产品名称" widget:"type:input;placeholder:请输入产品名称" data:"type:string" validate:"required"`
	Price      float64    `json:"price" form:"price" runner:"code:price;name:产品价格" widget:"type:input;prefix:¥;precision:2" data:"type:float" validate:"required,min=0"`
	Status     bool       `json:"status" form:"status" runner:"code:status;name:产品状态" widget:"type:switch;true_label:启用;false_label:禁用" data:"type:boolean;default_value:true"`
	LaunchDate typex.Time `json:"launch_date" form:"launch_date" runner:"code:launch_date;name:上线日期" widget:"type:datetime;format:date;min_date:$today" data:"type:string;default_value:$today" validate:"required"`
}

func TestUnifiedFormResponse(t *testing.T) {
	// 测试表单响应
	response, err := NewUnifiedFormResponse(TestProductReq{}, "form")
	if err != nil {
		t.Fatalf("NewUnifiedFormResponse failed: %v", err)
	}

	// 验证基本结构
	if response.RenderType != "form" {
		t.Errorf("Expected render_type 'form', got '%s'", response.RenderType)
	}

	if len(response.Fields) == 0 {
		t.Error("Expected fields to be populated")
	}

	if response.Columns != nil {
		t.Error("Form response should not have columns")
	}

	// 验证字段结构
	for _, field := range response.Fields {
		if field.Code == "" {
			t.Error("Field code should not be empty")
		}
		if field.Name == "" {
			t.Error("Field name should not be empty")
		}
		if field.Widget.Type == "" {
			t.Error("Widget type should not be empty")
		}
		if field.Data.Type == "" {
			t.Error("Data type should not be empty")
		}
	}

	// 输出JSON查看结构
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}
	t.Logf("Form Response JSON:\n%s", string(jsonData))
}

func TestUnifiedTableResponse(t *testing.T) {
	// 测试表格响应
	response, err := NewUnifiedTableResponse(TestProductListResp{})
	if err != nil {
		t.Fatalf("NewUnifiedTableResponse failed: %v", err)
	}

	// 验证基本结构
	if response.RenderType != "table" {
		t.Errorf("Expected render_type 'table', got '%s'", response.RenderType)
	}

	if len(response.Columns) == 0 {
		t.Error("Expected columns to be populated")
	}

	if response.Fields != nil {
		t.Error("Table response should not have fields")
	}

	// 验证列结构
	for _, column := range response.Columns {
		if column.Code == "" {
			t.Error("Column code should not be empty")
		}
		if column.Name == "" {
			t.Error("Column name should not be empty")
		}
		if column.Widget.Type == "" {
			t.Error("Widget type should not be empty")
		}
		if column.Data.Type == "" {
			t.Error("Data type should not be empty")
		}

		// Table模式下，有permission标签的字段应该有权限配置，没有permission标签的字段可以为nil
		// 这里只检查有permission标签的字段
		if column.Code == "status" || column.Code == "created_at" || column.Code == "updated_at" {
			if column.Permission == nil {
				t.Errorf("Table column %s should have permission config", column.Code)
			}
		}
	}

	// 输出JSON查看结构
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}
	t.Logf("Table Response JSON:\n%s", string(jsonData))
}

func TestStructureComparison(t *testing.T) {
	// 对比新旧结构的差异
	t.Log("=== 结构对比测试 ===")

	// 新的统一结构
	newResponse, err := NewUnifiedTableResponse(TestProductListResp{})
	if err != nil {
		t.Fatalf("NewUnifiedTableResponse failed: %v", err)
	}

	// 旧的Table结构（模拟）
	oldResponse := map[string]interface{}{
		"widget": "table",
		"columns": []map[string]interface{}{
			{
				"code":        "updated_at",
				"name":        "更新时间",
				"value_type":  "object",
				"widget_type": "input",
				"widget_config": map[string]interface{}{
					"mode":          "line_text",
					"placeholder":   "",
					"default_value": "",
				},
				"add_form_config": map[string]interface{}{
					"code":          "updated_at",
					"desc":          "",
					"name":          "更新时间",
					"show":          "",
					"hidden":        "",
					"example":       "",
					"required":      false,
					"callbacks":     "",
					"validates":     "",
					"value_type":    "object",
					"widget_type":   "input",
					"default_value": "",
					"widget_config": map[string]interface{}{
						"mode":          "line_text",
						"placeholder":   "",
						"default_value": "",
					},
				},
			},
		},
	}

	// 输出对比
	newJSON, _ := json.MarshalIndent(newResponse, "", "  ")
	oldJSON, _ := json.MarshalIndent(oldResponse, "", "  ")

	t.Logf("新的统一结构:\n%s", string(newJSON))
	t.Logf("旧的Table结构:\n%s", string(oldJSON))

	// 验证新结构的优势
	if len(newResponse.Columns) > 0 {
		firstColumn := newResponse.Columns[0]

		// 1. 没有冗余的嵌套结构
		if firstColumn.Widget.Config != nil {
			t.Log("✅ Widget配置结构清晰")
		}

		// 2. 权限配置简化
		if firstColumn.Permission != nil {
			t.Log("✅ 权限配置简化为三权限")
		}

		// 3. 数据类型配置统一
		if firstColumn.Data.Type != "" {
			t.Log("✅ 数据类型配置统一")
		}
	}
}

func TestFieldInfoSerialization(t *testing.T) {
	// 测试FieldInfo的序列化
	field := &FieldInfo{
		Code: "test_field",
		Name: "测试字段",
		Desc: "这是一个测试字段",
		Widget: WidgetConfig{
			Type: "input",
			Config: map[string]interface{}{
				"placeholder": "请输入测试值",
				"maxlength":   100,
			},
		},
		Data: DataConfig{
			Type:         "string",
			Example:      "测试示例",
			DefaultValue: "默认值",
		},
		Permission: &PermissionConfig{
			Read:   true,
			Update: true,
			Create: false,
		},
		Callbacks: []CallbackConfig{
			{
				Event: "OnInputFuzzy",
				Params: map[string]string{
					"debounce": "300ms",
				},
			},
		},
		Validation: "required,min=1,max=100",
	}

	// 序列化
	jsonData, err := json.MarshalIndent(field, "", "  ")
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	t.Logf("FieldInfo序列化结果:\n%s", string(jsonData))

	// 反序列化验证
	var deserializedField FieldInfo
	err = json.Unmarshal(jsonData, &deserializedField)
	if err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	// 验证关键字段
	if deserializedField.Code != field.Code {
		t.Errorf("Code mismatch: expected %s, got %s", field.Code, deserializedField.Code)
	}
	if deserializedField.Widget.Type != field.Widget.Type {
		t.Errorf("Widget type mismatch: expected %s, got %s", field.Widget.Type, deserializedField.Widget.Type)
	}
}

type Product struct {
	ID          int        `json:"id" gorm:"primaryKey;autoIncrement" runner:"code:id;name:产品ID" data:"type:number"`
	Name        string     `json:"name" gorm:"column:name;comment:产品名称" runner:"code:name;name:产品名称" widget:"type:input;placeholder:请输入产品名称" data:"type:string" validate:"required"`
	Category    string     `json:"category" gorm:"column:category;comment:产品分类" runner:"code:category;name:产品分类" widget:"type:select;options:手机,笔记本,平板,耳机,其他" data:"type:string" validate:"required"`
	Price       float64    `json:"price" gorm:"column:price;comment:产品价格" runner:"code:price;name:产品价格" widget:"type:input;prefix:¥;precision:2" data:"type:float" validate:"required,min=0"`
	Stock       int        `json:"stock" gorm:"column:stock;comment:库存数量" runner:"code:stock;name:库存数量" widget:"type:input;suffix:件" data:"type:number" validate:"required,min=0"`
	Description string     `json:"description" gorm:"column:description;comment:产品描述" runner:"code:description;name:产品描述" widget:"type:input;mode:text_area;max_length:200" data:"type:string"`
	Status      bool       `json:"status" gorm:"column:status;comment:产品状态" runner:"code:status;name:产品状态" widget:"type:switch;true_label:启用;false_label:禁用" data:"type:boolean;default_value:true" validate:"required"`
	Tags        string     `json:"tags" gorm:"column:tags;comment:产品标签" runner:"code:tags;name:产品标签" widget:"type:tag;separator:,;color:auto;max_tags:5" data:"type:string"`
	CreatedBy   string     `json:"created_by" gorm:"column:created_by;comment:创建人" runner:"code:created_by;name:创建人" permission:"read"`
	CreatedAt   typex.Time `json:"created_at" gorm:"autoCreateTime" runner:"code:created_at;name:创建时间" widget:"type:datetime;format:datetime" data:"type:string;example:2025-01-15 10:30:00" permission:"read"`
	UpdatedAt   typex.Time `json:"updated_at" gorm:"autoUpdateTime" runner:"code:updated_at;name:更新时间" widget:"type:datetime;format:datetime" data:"type:string;example:2025-01-15 14:20:00" permission:"read"`
}

type ProductListReq struct {
	query.PageInfoReq
	Name        string  `json:"name" form:"name" runner:"code:name;name:产品名称" widget:"type:input;placeholder:按产品名称搜索" data:"type:string"`
	Category    string  `json:"category" form:"category" runner:"code:category;name:产品分类" widget:"type:select;options:手机,笔记本,平板,耳机,其他;placeholder:选择产品分类" data:"type:string"`
	Status      string  `json:"status" form:"status" runner:"code:status;name:产品状态" widget:"type:select;options:启用,禁用;placeholder:选择产品状态" data:"type:string;default_value:启用"`
	MinPrice    float64 `json:"min_price" form:"min_price" runner:"code:min_price;name:最低价格" widget:"type:input;placeholder:最低价格;prefix:¥" data:"type:float" validate:"min=0"`
	MaxPrice    float64 `json:"max_price" form:"max_price" runner:"code:max_price;name:最高价格" widget:"type:input;placeholder:最高价格;prefix:¥" data:"type:float" validate:"min=0"`
	InStock     bool    `json:"in_stock" form:"in_stock" runner:"code:in_stock;name:仅显示有库存" widget:"type:switch;true_label:是;false_label:否" data:"type:boolean;default_value:false"`
	SortBy      string  `json:"sort_by" form:"sort_by" runner:"code:sort_by;name:排序方式" widget:"type:select;options:创建时间-降序,创建时间-升序,价格-降序,价格-升序,库存-降序,库存-升序;placeholder:选择排序方式" data:"type:string;default_value:创建时间-降序"`
	TagKeywords string  `json:"tag_keywords" form:"tag_keywords" runner:"code:tag_keywords;name:标签关键词" widget:"type:input;placeholder:搜索标签关键词" data:"type:string"`
}

func TestName111(t *testing.T) {
	// 测试Request参数
	requestParams, err := NewRequestParams(&ProductListReq{}, response.RenderTypeTable)
	if err != nil {
		panic(err)
	}

	// 测试Table Response参数
	params, err := NewTableResponseParams(query.PaginatedTable[[]Product]{})
	if err != nil {
		panic(err)
	}

	// 输出Request参数的JSON结构
	requestJSON, err := json.MarshalIndent(requestParams, "", "  ")
	if err != nil {
		t.Fatalf("Request JSON marshal failed: %v", err)
	}
	t.Logf("=== Request参数结构 ===\n%s", string(requestJSON))

	// 输出Table Response参数的JSON结构
	responseJSON, err := json.MarshalIndent(params, "", "  ")
	if err != nil {
		t.Fatalf("Response JSON marshal failed: %v", err)
	}
	t.Logf("=== Table Response参数结构 ===\n%s", string(responseJSON))

	// 对比结构差异
	t.Log("=== 结构对比分析 ===")
	t.Logf("Request类型: %T", requestParams)
	t.Logf("Response类型: %T", params)
}

func TestSliceTypeHandling(t *testing.T) {
	// 测试直接传入切片类型
	t.Log("=== 测试直接传入切片类型 ===")

	// 1. 测试 []Product 切片类型
	productSliceType := reflect.TypeOf([]Product{})
	builder := NewFormBuilder()

	tableConfig, err := builder.BuildTableConfig(productSliceType)
	if err != nil {
		t.Fatalf("BuildTableConfig with slice type failed: %v", err)
	}

	if len(tableConfig.Columns) == 0 {
		t.Error("Expected columns to be populated for slice type")
	}

	t.Logf("成功处理切片类型，生成了 %d 个列配置", len(tableConfig.Columns))

	// 2. 测试 []*Product 指针切片类型
	productPtrSliceType := reflect.TypeOf([]*Product{})
	tableConfig2, err := builder.BuildTableConfig(productPtrSliceType)
	if err != nil {
		t.Fatalf("BuildTableConfig with pointer slice type failed: %v", err)
	}

	if len(tableConfig2.Columns) == 0 {
		t.Error("Expected columns to be populated for pointer slice type")
	}

	t.Logf("成功处理指针切片类型，生成了 %d 个列配置", len(tableConfig2.Columns))

	// 3. 测试 NewUnifiedTableResponse 直接传入切片
	response, err := NewUnifiedTableResponse([]Product{})
	if err != nil {
		t.Fatalf("NewUnifiedTableResponse with slice failed: %v", err)
	}

	if response.RenderType != "table" {
		t.Errorf("Expected render_type 'table', got '%s'", response.RenderType)
	}

	if len(response.Columns) == 0 {
		t.Error("Expected columns to be populated")
	}

	// 输出结果
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}
	t.Logf("切片类型响应结果:\n%s", string(jsonData))
}

func TestErrorHandling(t *testing.T) {
	// 测试错误处理
	t.Log("=== 测试错误处理 ===")

	builder := NewFormBuilder()

	// 1. 测试非结构体、非切片类型
	stringType := reflect.TypeOf("test")
	_, err := builder.BuildTableConfig(stringType)
	if err == nil {
		t.Error("Expected error for string type, but got nil")
	}
	t.Logf("正确捕获字符串类型错误: %v", err)

	// 2. 测试非结构体切片类型
	intSliceType := reflect.TypeOf([]int{})
	_, err = builder.BuildTableConfig(intSliceType)
	if err == nil {
		t.Error("Expected error for int slice type, but got nil")
	}
	t.Logf("正确捕获整数切片类型错误: %v", err)

	// 3. 测试 NewUnifiedTableResponse 的错误处理
	_, err = NewUnifiedTableResponse("not a struct or slice")
	if err == nil {
		t.Error("Expected error for string input, but got nil")
	}
	t.Logf("正确捕获字符串输入错误: %v", err)
}

// TestIgnoreFields 测试忽略字段功能
type TestIgnoreFieldsStruct struct {
	ID         int        `json:"id" runner:"code:id;name:产品ID" data:"type:number"`
	Name       string     `json:"name" runner:"code:name;name:产品名称" widget:"type:input" data:"type:string" validate:"required"`
	InternalID string     `json:"-" runner:"-"`        // 应该被忽略
	TempField  string     `runner:"-"`                 // 应该被忽略
	Password   string     `json:"password" runner:"-"` // 应该被忽略
	CreatedAt  typex.Time `json:"created_at" runner:"code:created_at;name:创建时间" widget:"type:datetime" data:"type:string" permission:"read"`
}

func TestIgnoreFieldsWithRunnerDash(t *testing.T) {
	t.Log("=== 测试 runner:\"-\" 字段忽略功能 ===")

	// 测试Form响应
	formResponse, err := NewUnifiedFormResponse(TestIgnoreFieldsStruct{}, "form")
	if err != nil {
		t.Fatalf("NewUnifiedFormResponse failed: %v", err)
	}

	// 验证忽略的字段不在结果中
	fieldCodes := make([]string, 0, len(formResponse.Fields))
	for _, field := range formResponse.Fields {
		fieldCodes = append(fieldCodes, field.Code)
	}

	// 检查应该存在的字段
	expectedFields := []string{"id", "name", "created_at"}
	for _, expected := range expectedFields {
		found := false
		for _, code := range fieldCodes {
			if code == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected field '%s' not found in form response", expected)
		}
	}

	// 检查应该被忽略的字段
	ignoredFields := []string{"InternalID", "TempField", "Password"}
	for _, ignored := range ignoredFields {
		for _, field := range formResponse.Fields {
			if field.Code == ignored || strings.Contains(strings.ToLower(field.Code), strings.ToLower(ignored)) {
				t.Errorf("Field '%s' should be ignored but found in form response", ignored)
			}
		}
	}

	t.Logf("Form响应包含 %d 个字段: %v", len(formResponse.Fields), fieldCodes)

	// 测试Table响应
	tableResponse, err := NewUnifiedTableResponse([]TestIgnoreFieldsStruct{})
	if err != nil {
		t.Fatalf("NewUnifiedTableResponse failed: %v", err)
	}

	// 验证表格列
	columnCodes := make([]string, 0, len(tableResponse.Columns))
	for _, column := range tableResponse.Columns {
		columnCodes = append(columnCodes, column.Code)
	}

	// 检查应该存在的列
	for _, expected := range expectedFields {
		found := false
		for _, code := range columnCodes {
			if code == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected column '%s' not found in table response", expected)
		}
	}

	// 检查应该被忽略的列
	for _, ignored := range ignoredFields {
		for _, column := range tableResponse.Columns {
			if column.Code == ignored || strings.Contains(strings.ToLower(column.Code), strings.ToLower(ignored)) {
				t.Errorf("Column '%s' should be ignored but found in table response", ignored)
			}
		}
	}

	t.Logf("Table响应包含 %d 个列: %v", len(tableResponse.Columns), columnCodes)

	// 输出结果验证
	formJSON, _ := json.MarshalIndent(formResponse, "", "  ")
	t.Logf("Form响应结果:\n%s", string(formJSON))
}

// TestComplexIgnoreFields 测试复杂的忽略字段场景
type TestComplexIgnoreFields struct {
	// 正常字段
	ID   int    `json:"id" runner:"code:id;name:ID" data:"type:number"`
	Name string `json:"name" runner:"code:name;name:名称" widget:"type:input" data:"type:string"`

	// 各种忽略场景
	InternalField1 string `runner:"-"`                                 // 仅runner忽略
	InternalField2 string `json:"internal2" runner:"-"`                // json+runner忽略
	InternalField3 string `json:"-" runner:"-"`                        // json和runner都忽略
	PasswordField  string `json:"password" form:"password" runner:"-"` // 有json和form但runner忽略
	TempData       []byte `runner:"-" validate:"required"`             // 有其他标签但runner忽略

	// 正常字段继续
	Status    bool       `json:"status" runner:"code:status;name:状态" widget:"type:switch;true_label:启用;false_label:禁用" data:"type:boolean;default_value:true"`
	CreatedAt typex.Time `json:"created_at" runner:"code:created_at;name:创建时间" widget:"type:datetime" data:"type:string" permission:"read"`
}

func TestComplexIgnoreFieldsScenarios(t *testing.T) {
	t.Log("=== 测试复杂的 runner:\"-\" 字段忽略场景 ===")

	// 测试Form响应
	formResponse, err := NewUnifiedFormResponse(TestComplexIgnoreFields{}, "form")
	if err != nil {
		t.Fatalf("NewUnifiedFormResponse failed: %v", err)
	}

	// 应该存在的字段
	expectedFields := []string{"id", "name", "status", "created_at"}
	// 应该被忽略的字段（不应该出现）
	ignoredFields := []string{"InternalField1", "InternalField2", "InternalField3", "PasswordField", "TempData"}

	fieldCodes := make([]string, 0, len(formResponse.Fields))
	for _, field := range formResponse.Fields {
		fieldCodes = append(fieldCodes, field.Code)
	}

	// 验证应该存在的字段
	for _, expected := range expectedFields {
		found := false
		for _, code := range fieldCodes {
			if code == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected field '%s' not found in form response", expected)
		}
	}

	// 验证应该被忽略的字段确实不存在
	for _, ignored := range ignoredFields {
		for _, field := range formResponse.Fields {
			if field.Code == ignored || strings.Contains(strings.ToLower(field.Code), strings.ToLower(ignored)) {
				t.Errorf("Field '%s' should be ignored but found in form response with code '%s'", ignored, field.Code)
			}
		}
	}

	t.Logf("✅ Form响应正确忽略了 %d 个字段，保留了 %d 个字段: %v", len(ignoredFields), len(formResponse.Fields), fieldCodes)

	// 测试Table响应
	tableResponse, err := NewUnifiedTableResponse([]TestComplexIgnoreFields{})
	if err != nil {
		t.Fatalf("NewUnifiedTableResponse failed: %v", err)
	}

	columnCodes := make([]string, 0, len(tableResponse.Columns))
	for _, column := range tableResponse.Columns {
		columnCodes = append(columnCodes, column.Code)
	}

	// 验证表格列
	for _, expected := range expectedFields {
		found := false
		for _, code := range columnCodes {
			if code == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected column '%s' not found in table response", expected)
		}
	}

	for _, ignored := range ignoredFields {
		for _, column := range tableResponse.Columns {
			if column.Code == ignored || strings.Contains(strings.ToLower(column.Code), strings.ToLower(ignored)) {
				t.Errorf("Column '%s' should be ignored but found in table response with code '%s'", ignored, column.Code)
			}
		}
	}

	t.Logf("✅ Table响应正确忽略了 %d 个字段，保留了 %d 个列: %v", len(ignoredFields), len(tableResponse.Columns), columnCodes)

	// 验证字段数量
	if len(formResponse.Fields) != len(expectedFields) {
		t.Errorf("Expected %d fields but got %d", len(expectedFields), len(formResponse.Fields))
	}
	if len(tableResponse.Columns) != len(expectedFields) {
		t.Errorf("Expected %d columns but got %d", len(expectedFields), len(tableResponse.Columns))
	}

	t.Log("✅ 所有 runner:\"-\" 字段忽略测试通过！")
}
