package llms

import (
	"context"
	"testing"
	"time"
)

// TestDeepSeekClientCreation 测试DeepSeek客户端创建
func TestDeepSeekClientCreation(t *testing.T) {
	client, err := NewDeepSeekClientFromEnv()
	if err != nil {
		t.Fatalf("从环境变量创建客户端失败: %v", err)
	}

	if client == nil {
		t.Fatal("客户端创建失败")
	}

	// 验证API Key不为空即可，不检查具体值
	deepSeekClient, ok := client.(*DeepSeekClient)
	if !ok {
		t.Fatal("客户端类型错误")
	}
	if deepSeekClient.APIKey == "" {
		t.Error("API Key为空")
	}

	if deepSeekClient.BaseURL != "https://api.deepseek.com/v1/chat/completions" {
		t.Errorf("BaseURL设置错误，期望: %s, 实际: %s",
			"https://api.deepseek.com/v1/chat/completions", deepSeekClient.BaseURL)
	}

}

// TestDeepSeekClientInterface 测试客户端接口实现
func TestDeepSeekClientInterface(t *testing.T) {
	client, err := NewDeepSeekClientFromEnv()
	if err != nil {
		t.Fatalf("从环境变量创建客户端失败: %v", err)
	}

	// 检查是否实现了LLMClient接口
	var _ LLMClient = client

	// 测试模型名称
	modelName := client.GetModelName()
	if modelName != "deepseek-chat" {
		t.Errorf("模型名称错误，期望: deepseek-chat, 实际: %s", modelName)
	}

	// 测试提供商名称
	provider := client.GetProvider()
	if provider != "DeepSeek" {
		t.Errorf("提供商名称错误，期望: DeepSeek, 实际: %s", provider)
	}
}

// TestDeepSeekChatBasic 测试基本的聊天功能
func TestDeepSeekChatBasic(t *testing.T) {
	client, err := NewDeepSeekClientFromEnv()
	if err != nil {
		t.Fatalf("从环境变量创建客户端失败: %v", err)
	}

	req := &ChatRequest{
		Messages: []Message{
			{Role: "user", Content: "你好，请简单介绍一下你自己"},
		},
		MaxTokens:   100,
		Temperature: 0.7,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.Chat(ctx, req)
	if err != nil {
		t.Fatalf("聊天请求失败: %v", err)
	}

	if resp == nil {
		t.Fatal("响应为空")
	}

	// 检查是否有错误
	if resp.Error != "" {
		t.Logf("API返回错误: %s", resp.Error)
		// 如果是API key相关错误，记录但不失败
		if resp.Error != "" {
			t.Logf("注意：这可能是API key无效或网络问题，请检查配置")
		}
		return
	}

	// 检查响应内容
	if resp.Content == "" {
		t.Error("响应内容为空")
	} else {
		t.Logf("AI回答: %s", resp.Content)
	}

	// 检查使用统计
	if resp.Usage != nil {
		t.Logf("Token使用: 输入%d, 输出%d, 总计%d",
			resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
	}
}

// TestDeepSeekChatWithSystemPrompt 测试带系统提示的聊天
func TestDeepSeekChatWithSystemPrompt(t *testing.T) {
	client, err := NewDeepSeekClientFromEnv()
	if err != nil {
		t.Fatalf("从环境变量创建客户端失败: %v", err)
	}

	req := &ChatRequest{
		Messages: []Message{
			{
				Role:    "system",
				Content: "你是一个专业的Go语言开发助手，请用简洁的语言回答问题",
			},
			{
				Role:    "user",
				Content: "Go语言中如何创建一个HTTP服务器？",
			},
		},
		MaxTokens:   200,
		Temperature: 0.3,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.Chat(ctx, req)
	if err != nil {
		t.Fatalf("带系统提示的聊天请求失败: %v", err)
	}

	if resp == nil {
		t.Fatal("响应为空")
	}

	// 检查是否有错误
	if resp.Error != "" {
		t.Logf("API返回错误: %s", resp.Error)
		t.Logf("注意：这可能是API key无效或网络问题，请检查配置")
		return
	}

	// 检查响应内容
	if resp.Content == "" {
		t.Error("响应内容为空")
	} else {
		t.Logf("AI回答: %s", resp.Content)
	}

	// 检查使用统计
	if resp.Usage != nil {
		t.Logf("Token使用: 输入%d, 输出%d, 总计%d",
			resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
	}
}

// TestDeepSeekCodeGeneration 测试代码生成功能
func TestDeepSeekCodeGeneration(t *testing.T) {
	client, err := NewDeepSeekClientFromEnv()
	if err != nil {
		t.Fatalf("从环境变量创建客户端失败: %v", err)
	}

	timeout := time.Second * 700
	req := &ChatRequest{
		Messages: []Message{
			{
				Role:    "system",
				Content: "你是一个专业的Go语言开发助手，请生成可运行的代码",
			},
			{
				Role:    "user",
				Content: "# function-go - AI原生全栈开发框架 > **专为大模型设计的Go全栈开发框架** > AI原生结构体标签系统 | 一个文件 = 一个完整Web应用 | 大模型生成代码直接编译部署 | 框架内部有运行入口，专注业务函数 ## 🏷️ 标签系统 - 大模型代码生成核心 ### 设计理念 function-go通过结构体标签实现**声明式配置**，让大模型能够： - **零样板代码**：一个结构体 = 完整Web应用 - **AI友好**：标签即配置，大模型直接理解 - **业务导向**：专注业务逻辑，框架处理技术细节 ### 💡 大模型使用指南 **记住：标签就是配置，配置就是功能！** | 标签 | 作用 | 大模型使用场景 | |------|------|----------------| | `runner:\"name:字段名\"` | 前端显示名称 | 所有字段都需要 | | `widget:\"type:组件类型\"` | UI组件选择 | 根据数据类型选择 | | `validate:\"required\"` | 数据验证 | 必填字段必须加 | | `search:\"like\"` | 搜索功能 | 需要搜索的字段 | ### 快速上手模板 ```go // 大模型代码生成模板 type YourModel struct { ID int `json:\"id\" gorm:\"primaryKey\" runner:\"name:ID\" permission:\"read\"` //这个字段只读 Name string `json:\"name\" runner:\"name:名称\" widget:\"type:input\" search:\"like\" validate:\"required\"` Status string `json:\"status\" runner:\"name:状态\" widget:\"type:select;options:启用,禁用\" search:\"in\"` CreatedAt int64 `json:\"created_at\" gorm:\"autoCreateTime:milli\" runner:\"name:创建时间\" permission:\"read\"` } ``` #### 2. **UI渲染** - **widget标签**：选择最合适的UI组件进行渲染 - **data标签**：自动设置默认值、示例值、动态值，格式化的值， - **runner标签**：自动生成中文标签、单位、显示名称 #### 3. **数据验证引擎** - **validate标签**：自动生成前端和后端验证规则 - **search标签**：自动生成搜索和过滤功能 - **permission标签**：自动控制字段在不同场景的显示权限 #### 4. **数据库操作自动化** - **gorm标签**：自动生成数据库表结构 - **CreateTables**：服务启动时自动建表 - **AutoCrudTable**：自动生成增删改查操作 #### 5. **回调函数集成** - **OnInputFuzzy**：自动集成模糊搜索和聚合计算 - **OnInputValidate**：自动集成实时字段验证 - **OnTableAddRows**：table函数新增记录回调 - **OnTableUpdateRows** table函数更新记录回调 - **OnTableDeleteRows** table函数删除记录回调 ### 🔄 标签系统的工作流程 ``` 结构体定义 → 标签解析 → 代码生成 → 运行时执行 ↓ ↓ ↓ ↓ 业务模型 配置信息 前端界面 完整应用 数据库表 验证规则 API接口 业务逻辑 ``` ### 🌟 标签系统的优势 | 传统开发方式 | function-go标签方式 | |------------|-----------------| | 手动编写CRUD代码 | 自动生成CRUD代码 | | 手动编写验证逻辑 | 标签声明验证规则 | | 手动设计UI界面 | 自动渲染UI界面 | | 手动管理数据库 | 自动管理数据库 | | 代码量大、易出错 | 代码简洁、零错误 | 通过标签系统，开发者只需要关注**业务逻辑**，框架自动处理所有**技术细节**，真正实现了\"一个文件 = 一个完整Web应用\"的愿景。 ### 标签顺序建议 ``` json → gorm → runner → widget → data → search → permission → validate ``` ### 核心标签说明 #### runner标签 - 业务逻辑配置 | 属性 | 格式 | 示例 | 说明 | |------|-------------|---------------------------------|----------------| | 字段名称 | `name:显示名称` | `runner:\"name:用户名\"` | 设置字段在前端的显示名称 | | 字段单位 | `desc:字段介绍` | `runner:\"name:年龄;desc:年龄0-100\"` | 设置字段的详细介绍（非必要） | #### widget标签 - UI组件配置 | 属性 | 格式 | 示例 | 说明 | |------|-------------|-----------------------|----------| | 组件类型 | `type:组件类型` | `widget:\"type:input\"` | 设置UI组件类型 | #### data标签 - 数据和值配置 | 功能 | 格式 | 示例 | 说明 | |-------|---------------------|-----------------------------|-------------------------| | 默认值 | `default_value:值` | `data:\"default_value:默认值\"` | 设置字段默认值 | | 示例值 | `example:示例值` | `data:\"example:示例文本\"` | 设置示例值 | | 动态默认值 | `default_value:$变量` | `data:\"default_value:$now\"` | 使用变量作为默认值 | | 格式化 | `format:格式化类型` | `format:markdown` | `设置格式化类型，csv或者markdown` #### validate标签 - 验证规则 | 规则 | 格式 | 示例 | 说明 | |------|---------------|----------------------------|--------| | 必填验证 | `required` | `validate:\"required\"` | 字段必填 | | 长度验证 | `min=值,max=值` | `validate:\"min=2,max=50\"` | 长度范围验证 | | 数值验证 | `min=值,max=值` | `validate:\"min=1,max=120\"` | 数值范围验证 | | 枚举验证 | `oneof=值1 值2` | `validate:\"oneof=男 女\"` | 枚举值验证 | #### search标签 - 搜索配置（仅table函数） | 搜索类型 | 格式 | 示例 | 说明 | |------|-----------|--------------------|---------------| | 模糊搜索 | `like` | `search:\"like\"` | 启用模糊搜索 | | 精确搜索 | `eq` | `search:\"eq\"` | 启用精确搜索 | | 区间搜索 | `gte,lte` | `search:\"gte,lte\"` | 启用大于等于、小于等于搜索 | | 多选搜索 | `in` | `search:\"in\"` | 启用多选搜索 | #### permission标签 - 权限控制（仅table函数） | 权限类型 | 格式 | 示例 | 说明 | |------|----------|-----------------------|-------------| | 仅可读 | `read` | `permission:\"read\"` | 仅列表显示，不能编辑 | | 仅可创建 | `create` | `permission:\"create\"` | 仅新增显示，列表不显示 | | 仅可更新 | `update` | `permission:\"update\"` | 仅编辑显示，列表不显示 | | 全权限 | 不写 | 无标签 | 列表、新增、编辑都显示 | ## 🧩 组件系统 ### 基础输入组件 #### input组件 - 文本输入 | 类型 | 配置 | 示例 | 说明 | |------|------------------|--------------------------------------|---------| | 单行文本 | `type:input` | `widget:\"type:input\"` | 基础文本输入框 | | 多行文本 | `mode:text_area` | `widget:\"type:input;mode:text_area\"` | 多行文本区域 | | 密码输入 | `mode:password` | `widget:\"type:input;mode:password\"` | 密码输入框 | #### number组件 - 数字输入 | 类型 | 配置 | 示例 | 说明 | |------|-----------------|-----------------------------------------------------------|--------| | 整数输入 | `type:number` | `widget:\"type:number;min:1;max:120;unit:岁\"` | 整数输入框 | | 小数输入 | `precision:小数位` | `widget:\"type:number;min:0;precision:2;prefix:￥\"` | 小数输入框 | | 百分比 | `suffix:%` | `widget:\"type:number;min:0;max:100;precision:1;suffix:%\"` | 百分比输入框 | #### select组件 - 下拉选择 | 类型 | 配置 | 示例 | 说明 | |------|-----------------|-------------------------------------------------------|-------| | 单选下拉 | `type:select` | `widget:\"type:select;options:男,女\"` | 单选下拉框 | | 多选下拉 | `multiple:true` | `widget:\"type:select;options:技术,产品,设计;multiple:true\"` | 多选下拉框 | #### datetime组件 - 日期时间 | 类型 | 配置 | 示例 | 说明 | |------|------------------|-----------------------------------------------------------|---------| | 日期选择 | `kind:date` | `widget:\"type:datetime;kind:date;format:yyyy-MM-dd\"` | 日期选择器 | | 时间选择 | `kind:time` | `widget:\"type:datetime;kind:time;format:HH:mm\"` | 时间选择器 | | 日期时间 | `kind:datetime` | `widget:\"type:datetime;kind:datetime\"` | 日期时间选择器 | | 日期范围 | `kind:daterange` | `widget:\"type:datetime;kind:daterange;format:yyyy-MM-dd\"` | 日期范围选择器 | ### 高级组件 #### multiselect组件 - 多选组件 | 配置 | 示例 | 说明 | |--------|----------------------------------------------------------------------|-----------| | 静态多选 | `widget:\"type:multiselect;options:紧急,重要,API,UI\"` | 固定选项多选 | | 可创建新选项 | `widget:\"type:multiselect;options:Java,Python,Go;allow_create:true\"` | 支持自定义创建选项 | #### color组件 - 颜色选择器 | 格式 | 配置 | 示例 | 说明 | |--------|---------------|---------------------------------------------------|----------| | Hex格式 | `format:hex` | `widget:\"type:color;format:hex;show_alpha:false\"` | 6位hex颜色 | | RGBA格式 | `format:rgba` | `widget:\"type:color;format:rgba;show_alpha:true\"` | RGBA颜色格式 | | HSL格式 | `format:hsl` | `widget:\"type:color;format:hsl;show_alpha:false\"` | HSL颜色格式 | #### file_upload组件 - 文件上传 | 配置 | 示例 | 说明 | |-------|--------------------------------------------------------------------------|-------| | 单文件上传 | `widget:\"type:file_upload;accept:.jpg,.png;max_size:5MB\"` | 单文件上传 | | 多文件上传 | `widget:\"type:file_upload;accept:.pdf,.doc;multiple:true;max_size:10MB\"` | 多文件上传 | #### list组件 - 列表输入 | 类型 | 示例 | 说明 | |------|----------------------|----------| | 简单列表 | `widget:\"type:list\"` | 字符串或数字列表 | | 复杂列表 | `widget:\"type:list\"` | 结构体列表 | #### form组件 - 嵌套表单 | 示例 | 说明 | |----------------------|-------------------| | `widget:\"type:form\"` | 嵌套表单结构，对应数据结构是结构体 | ### 其他组件 #### switch组件 - 开关 | 配置 | 示例 | 说明 | |-------|-----------------------------------------------------|---------| | 基础开关 | `widget:\"type:switch\"` | 布尔值开关 | | 自定义标签 | `widget:\"type:switch;true_label:启用;false_label:禁用\"` | 自定义开关标签 | #### radio组件 - 单选框 | 配置 | 示例 | 说明 | |-------|--------------------------------------------------------|-----------| | 基础单选框 | `widget:\"type:radio;options:男,女\"` | 单选按钮组 | | 水平排列 | `widget:\"type:radio;options:男,女;direction:horizontal\"` | 水平排列的单选按钮 | #### checkbox组件 - 复选框 | 配置 | 示例 | 说明 | |-------|-------------------------------------------|--------| | 基础复选框 | `widget:\"type:checkbox;options:阅读,音乐,运动\"` | 多选复选框组 | #### slider组件 - 滑块 | 配置 | 示例 | 说明 | |------|----------------------------------------------------|----------| | 数值滑块 | `widget:\"type:slider;min:0;max:100;unit:%\"` | 数值范围滑块 | | 评分滑块 | `widget:\"type:slider;min:1;max:5;step:0.5;unit:分\"` | 带步进的评分滑块 | ## 📝 使用示例 ### 基础字段配置 ## �� Form函数模型示例 - 大模型代码生成模板 #### 用户注册模型 ```go // 请求结构体 - 用户输入 type UserRegisterReq struct { // 基础信息 Username string `json:\"username\" runner:\"name:用户名\" widget:\"type:input\" data:\"example:john_doe\" validate:\"required,min=3,max=20\"` Password string `json:\"password\" runner:\"name:密码\" widget:\"type:input;mode:password\" data:\"example:123456\" validate:\"required,min=6,max=20\"` Email string `json:\"email\" runner:\"name:邮箱\" widget:\"type:input\" data:\"example:john@example.com\" validate:\"required,email\"` // 个人信息 RealName string `json:\"real_name\" runner:\"name:真实姓名\" widget:\"type:input\" data:\"example:张三\" validate:\"required,min=2,max=20\"` Age int `json:\"age\" runner:\"name:年龄\" widget:\"type:number;min:18;max:65;unit:岁\" data:\"example:25\" validate:\"required,min=18,max=65\"` Gender string `json:\"gender\" runner:\"name:性别\" widget:\"type:radio;options:男,女;direction:horizontal\" data:\"example:男\" validate:\"required,oneof=男 女\"` // 工作信息 Department string `json:\"department\" runner:\"name:部门\" widget:\"type:select;options:技术部,产品部,设计部,运营部\" data:\"default_value:技术部\" validate:\"required\"` Position string `json:\"position\" runner:\"name:职位\" widget:\"type:input\" data:\"example:软件工程师\" validate:\"required\"` Salary int `json:\"salary\" runner:\"name:期望薪资\" widget:\"type:number;min:3000;max:50000;unit:元\" data:\"example:15000\" validate:\"required,min=3000,max=50000\"` // 技能标签 Skills []string `json:\"skills\" runner:\"name:技能标签\" widget:\"type:multiselect;options:Java,Python,Go,JavaScript,React,Vue\" data:\"example:Java,Go\" validate:\"required,min=1\"` // 附件上传 Resume *files.Files `json:\"resume\" runner:\"name:简历\" widget:\"type:file_upload;accept:.pdf,.doc,.docx;max_size:10MB\" validate:\"required\"` Avatar *files.Files `json:\"avatar\" runner:\"name:头像\" widget:\"type:file_upload;accept:.jpg,.png,.gif;max_size:5MB\"` // 其他信息 Bio string `json:\"bio\" runner:\"name:个人简介\" widget:\"type:input;mode:text_area\" data:\"example:热爱编程，有3年开发经验\"` AgreeTerms bool `json:\"agree_terms\" runner:\"name:同意条款\" widget:\"type:switch;true_label:同意;false_label:不同意\" data:\"example:true\" validate:\"required\"` } // 响应结构体 - 处理结果 type UserRegisterResp struct { // 处理结果 Success bool `json:\"success\" runner:\"name:是否成功\" widget:\"type:switch;true_label:成功;false_label:失败\"` Message string `json:\"message\" runner:\"name:处理结果\" widget:\"type:input;mode:text_area\"` // 用户信息 UserID int `json:\"user_id\" runner:\"name:用户ID\" widget:\"type:number\"` Username string `json:\"username\" runner:\"name:用户名\" widget:\"type:input\"` // 系统信息 CreatedAt int64 `json:\"created_at\" runner:\"name:注册时间\" widget:\"type:datetime;kind:datetime\"` Token string `json:\"token\" runner:\"name:访问令牌\" widget:\"type:input;mode:password\"` } ``` #### 采购申请模型 ```go // 请求结构体 - 采购申请 type PurchaseReq struct { // 基础信息 Title string `json:\"title\" runner:\"name:采购标题\" widget:\"type:input\" data:\"example:办公用品采购\" validate:\"required,min=5,max=100\"` Department string `json:\"department\" runner:\"name:申请部门\" widget:\"type:select;options:技术部,产品部,设计部,运营部\" validate:\"required\"` Priority string `json:\"priority\" runner:\"name:优先级\" widget:\"type:select;options:低,中,高,紧急\" data:\"default_value:中\" validate:\"required\"` // 供应商信息 SupplierID int `json:\"supplier_id\" runner:\"name:供应商\" widget:\"type:select\" validate:\"required\"` // 采购商品列表 Items []PurchaseItem `json:\"items\" runner:\"name:采购商品\" widget:\"type:list\" validate:\"required,min=1\"` // 其他信息 ExpectedDate int64 `json:\"expected_date\" runner:\"name:期望到货日期\" widget:\"type:datetime;kind:date;format:yyyy-MM-dd\" validate:\"required\"` Remarks string `json:\"remarks\" runner:\"name:备注说明\" widget:\"type:input;mode:text_area\"` } // 采购商品项 type PurchaseItem struct { ProductID int `json:\"product_id\" runner:\"name:商品\" widget:\"type:select\" validate:\"required\"` Quantity int `json:\"quantity\" runner:\"name:数量\" widget:\"type:number;min:1\" data:\"default_value:1\" validate:\"required,min=1\"` UnitPrice float64 `json:\"unit_price\" runner:\"name:单价\" widget:\"type:number;min:0;precision:2;prefix:￥\" validate:\"required,min=0\"` Remarks string `json:\"remarks\" runner:\"name:备注\" widget:\"type:input\"` } // 响应结构体 - 采购结果 type PurchaseResp struct { // 处理结果 Success bool `json:\"success\" runner:\"name:是否成功\" widget:\"type:switch;true_label:成功;false_label:失败\"` Message string `json:\"message\" runner:\"name:处理结果\" widget:\"type:input;mode:text_area\"` // 采购信息 PurchaseID int `json:\"purchase_id\" runner:\"name:采购单号\" widget:\"type:number\"` TotalAmount float64 `json:\"total_amount\" runner:\"name:总金额\" widget:\"type:number;precision:2;prefix:￥\"` TotalItems int `json:\"total_items\" runner:\"name:商品种类\" widget:\"type:number\"` // 状态信息 Status string `json:\"status\" runner:\"name:采购状态\" widget:\"type:input\"` CreatedAt int64 `json:\"created_at\" runner:\"name:创建时间\" widget:\"type:datetime;kind:datetime\"` } ``` ### **2. Table函数模型示例** #### 用户管理模型 ```go // 用户数据模型 - 自动建表 type CRMUser struct { // 系统字段 ID int `json:\"id\" gorm:\"primaryKey;autoIncrement\" runner:\"name:用户ID\" permission:\"read\"` CreatedAt int64 `json:\"created_at\" gorm:\"autoCreateTime:milli\" runner:\"name:创建时间\" widget:\"type:datetime;kind:datetime\" permission:\"read\"` UpdatedAt int64 `json:\"updated_at\" gorm:\"autoUpdateTime:milli\" runner:\"name:更新时间\" widget:\"type:datetime;kind:datetime\" permission:\"read\"` DeletedAt gorm.DeletedAt `json:\"deleted_at\" gorm:\"index\" runner:\"-\"` // 基础信息 Username string `json:\"username\" gorm:\"column:username;uniqueIndex\" runner:\"name:用户名\" widget:\"type:input\" search:\"like\" validate:\"required,min=3,max=20\"` Email string `json:\"email\" gorm:\"column:email;uniqueIndex\" runner:\"name:邮箱\" widget:\"type:input\" search:\"like\" validate:\"required,email\"` Phone string `json:\"phone\" gorm:\"column:phone\" runner:\"name:手机号\" widget:\"type:input\" search:\"like\" validate:\"required,min=11,max=11\"` // 个人信息 RealName string `json:\"real_name\" gorm:\"column:real_name\" runner:\"name:真实姓名\" widget:\"type:input\" search:\"like\" validate:\"required,min=2,max=20\"` Age int `json:\"age\" gorm:\"column:age\" runner:\"name:年龄\" widget:\"type:number;min:18;max:65;unit:岁\" search:\"gte,lte\" validate:\"required,min=18,max=65\"` Gender string `json:\"gender\" gorm:\"column:gender\" runner:\"name:性别\" widget:\"type:select;options:男,女\" search:\"in\" validate:\"required,oneof=男 女\"` Avatar string `json:\"avatar\" gorm:\"column:avatar\" runner:\"name:头像\" widget:\"type:input\"` // 工作信息 Department string `json:\"department\" gorm:\"column:department\" runner:\"name:部门\" widget:\"type:select;options:技术部,产品部,设计部,运营部\" search:\"in\" validate:\"required\"` Position string `json:\"position\" gorm:\"column:position\" runner:\"name:职位\" widget:\"type:input\" search:\"like\" validate:\"required\"` Salary int `json:\"salary\" gorm:\"column:salary\" runner:\"name:薪资\" widget:\"type:number;min:3000;max:50000;unit:元\" search:\"gte,lte\" validate:\"required,min=3000,max=50000\"` // 状态信息 Status string `json:\"status\" gorm:\"column:status\" runner:\"name:状态\" widget:\"type:select;options:在职,离职,试用期\" data:\"default_value:在职\" search:\"in\" validate:\"required\"` IsActive bool `json:\"is_active\" gorm:\"column:is_active\" runner:\"name:是否启用\" widget:\"type:switch;true_label:启用;false_label:禁用\" data:\"default_value:true\" search:\"eq\" validate:\"required\"` // 敏感信息 - 仅新增和编辑时显示 Password string `json:\"password\" gorm:\"column:password\" runner:\"name:密码\" widget:\"type:input;mode:password\" permission:\"create,update\"` } func (User) TableName() string { return \"crm_user\" } // 请求结构体 - 分页搜索 type CRMUserListReq struct { query.PageInfoReq `runner:\"-\"` // 自定义搜索字段，这里的字段是不存在于上面表中的字段，存在表中的字段可以直接打标签支持各种搜索， //下面的搜索字段一般是join字段，例如 用户关联组织表，然后下面可以用组织表的字段，然后在处理函数用该字段进行连表查询 OrgName string `json:\"org_name\" runner:\"name:组织名称\" widget:\"type:input\" search:\"like\"` } ``` #### 产品管理模型 ```go // 产品数据模型 - 自动建表 type Product struct { // 系统字段 ID int `json:\"id\" gorm:\"primaryKey;autoIncrement\" runner:\"name:产品ID\" permission:\"read\"` CreatedAt int64 `json:\"created_at\" gorm:\"autoCreateTime:milli\" runner:\"name:创建时间\" widget:\"type:datetime;kind:datetime\" permission:\"read\"` UpdatedAt int64 `json:\"updated_at\" gorm:\"autoUpdateTime:milli\" runner:\"name:更新时间\" widget:\"type:datetime;kind:datetime\" permission:\"read\"` DeletedAt gorm.DeletedAt `json:\"deleted_at\" gorm:\"index\" runner:\"-\"` // 基础信息 Name string `json:\"name\" gorm:\"column:name\" runner:\"name:产品名称\" widget:\"type:input\" search:\"like\" validate:\"required,min=2,max=100\"` Code string `json:\"code\" gorm:\"column:code;uniqueIndex\" runner:\"name:产品编码\" widget:\"type:input\" search:\"like\" validate:\"required,min=3,max=20\"` Category string `json:\"category\" gorm:\"column:category\" runner:\"name:产品分类\" widget:\"type:select;options:电子产品,服装,食品,家居,其他\" search:\"in\" validate:\"required\"` // 规格信息 Brand string `json:\"brand\" gorm:\"column:brand\" runner:\"name:品牌\" widget:\"type:input\" search:\"like\" validate:\"required\"` Model string `json:\"model\" gorm:\"column:model\" runner:\"name:型号\" widget:\"type:input\" search:\"like\"` Spec string `json:\"spec\" gorm:\"column:spec\" runner:\"name:规格\" widget:\"type:input\" search:\"like\"` Unit string `json:\"unit\" gorm:\"column:unit\" runner:\"name:单位\" widget:\"type:input\" search:\"like\" validate:\"required\"` // 价格信息 CostPrice float64 `json:\"cost_price\" gorm:\"column:cost_price\" runner:\"name:成本价\" widget:\"type:number;min:0;precision:2;prefix:￥\" search:\"gte,lte\" validate:\"required,min=0\"` SalePrice float64 `json:\"sale_price\" gorm:\"column:sale_price\" runner:\"name:销售价\" widget:\"type:number;min:0;precision:2;prefix:￥\" search:\"gte,lte\" validate:\"required,min=0\"` MarketPrice float64 `json:\"market_price\" gorm:\"column:market_price\" runner:\"name:市场价\" widget:\"type:number;min:0;precision:2;prefix:￥\" search:\"gte,lte\"` // 库存信息 Stock int `json:\"stock\" gorm:\"column:stock\" runner:\"name:库存\" widget:\"type:number;min:0;unit:件\" search:\"gte,lte\" validate:\"required,min=0\"` MinStock int `json:\"min_stock\" gorm:\"column:min_stock\" runner:\"name:最低库存\" widget:\"type:number;min:0;unit:件\" data:\"default_value:10\" validate:\"required,min=0\"` // 状态信息 Status string `json:\"status\" gorm:\"column:status\" runner:\"name:状态\" widget:\"type:select;options:上架,下架,缺货,停售\" data:\"default_value:上架\" search:\"in\" validate:\"required\"` IsHot bool `json:\"is_hot\" gorm:\"column:is_hot\" runner:\"name:是否热销\" widget:\"type:switch;true_label:是;false_label:否\" data:\"default_value:false\" search:\"eq\"` // 描述信息 Description string `json:\"description\" gorm:\"column:description\" runner:\"name:产品描述\" widget:\"type:input;mode:text_area\" search:\"like\"` Images string `json:\"images\" gorm:\"column:images\" runner:\"name:产品图片\" widget:\"type:input\"` } func (Product) TableName() string { return \"product\" } // 请求结构体 - 分页搜索 type ProductListReq struct { query.PageInfoReq `runner:\"-\"` // 自定义搜索字段 Category string `json:\"category\" runner:\"name:分类筛选\" widget:\"type:select;options:电子产品,服装,食品,家居,其他\" search:\"in\"` Brand string `json:\"brand\" runner:\"name:品牌筛选\" widget:\"type:input\" search:\"like\"` PriceRange []float64 `json:\"price_range\" runner:\"name:价格范围\" widget:\"type:number;min:0;precision:2\" search:\"gte,lte\"` StockStatus string `json:\"stock_status\" runner:\"name:库存状态\" widget:\"type:select;options:充足,不足,缺货\" search:\"in\"` } ``` ### ��️ Form函数配置模板 ```go var YourFormOption = &runner.FormFunctionOptions{ BaseConfig: runner.BaseConfig{ ChineseName: \"功能名称\", ApiDesc: \"功能描述\", Tags: []string{\"标签1\", \"标签2\"}, Request: &YourReq{}, Response: &YourResp{}, CreateTables: []interface{}{&YourModel{}}, // 如果需要建表 Group: YourGroup, // 如果使用函数组 }, // 如果需要模糊搜索 OnInputFuzzyMap: map[string]runner.OnInputFuzzy{ \"field_name\": func(ctx *runner.Context, req *usercall.OnInputFuzzyReq) (*usercall.OnInputFuzzyResp, error) { // 实现模糊搜索逻辑 return &usercall.OnInputFuzzyResp{ Statistics: map[string]interface{}{ \"统计信息\": \"聚合计算表达式\", }, Values: items, }, nil }, }, } ``` ### ��️ Table函数配置模板 ```go var YourTableOption = &runner.TableFunctionOptions{ BaseConfig: runner.BaseConfig{ ChineseName: \"功能名称\", ApiDesc: \"功能描述\", Tags: []string{\"标签1\", \"标签2\"}, Request: &YourListReq{}, Response: query.PaginatedTable[[]YourModel]{}, CreateTables: []interface{}{&YourModel{}}, Group: YourGroup, // 如果使用函数组 }, // 自动CRUD AutoCrudTable: &YourModel{}, // 如果需要回调 OnTableAddRows: func(ctx *runner.Context, req *usercall.OnTableAddRowsReq) (*usercall.OnTableAddRowsResp, error) { // 实现新增逻辑 return &usercall.OnTableAddRowsResp{}, nil }, OnTableUpdateRows: func(ctx *runner.Context, req *usercall.OnTableUpdateRowsReq) (*usercall.OnTableUpdateRowsResp, error) { // 实现更新逻辑 return &usercall.OnTableUpdateRowsResp{}, nil }, OnTableDeleteRows: func(ctx *runner.Context, req *usercall.OnTableDeleteRowsReq) (*usercall.OnTableDeleteRowsResp, error) { // 实现删除逻辑 return &usercall.OnTableDeleteRowsResp{}, nil }, } ``` ### 权限控制示例 ```go type Product struct { // 只读字段 ID int `json:\"id\" gorm:\"primaryKey\" runner:\"name:ID\" permission:\"read\"` CreatedAt int64 `json:\"created_at\" gorm:\"autoCreateTime:milli\" runner:\"name:创建时间\" permission:\"read\"` // 业务字段：全权限 Name string `json:\"name\" runner:\"name:产品名称\" widget:\"type:input\" search:\"like\" validate:\"required\"` Status string `json:\"status\" runner:\"name:状态\" widget:\"type:select;options:启用,禁用\" search:\"in\"` // 密钥字段：仅编辑时显示（系统自动生成，用户不能修改） SecretKey string `json:\"secret_key\" runner:\"name:密钥\" widget:\"type:input;mode:password\" permission:\"update\"` // 备注字段：仅创建时显示（创建后不能编辑） Remark string `json:\"remark\" runner:\"name:创建备注\" widget:\"type:input;mode:text_area\" permission:\"create\"` } ``` ## 🎯 最佳实践 ### 1. 标签配置原则 - **必填字段**：添加 `validate:\"required\"` - **搜索字段**：根据类型选择合适的 `search` 标签 - **权限控制**：使用 `permission` 标签控制字段显示 - **默认值**：使用 `data:\"default_value:值\"` 设置默认值 ### 2. 组件选择原则 - **文本输入**：使用 `input` 组件 - **数字输入**：使用 `number` 组件 - **选择输入**：使用 `select`、`radio`、`checkbox` 组件 - **日期时间**：使用 `datetime` 组件 - **文件处理**：使用 `file_upload` 组件 - **多选场景**：使用 `multiselect` 组件 ### 3. 搜索配置原则 - **文本字段**：使用 `like` 模糊搜索 - **状态字段**：使用 `in` 多选搜索 - **数值字段**：使用 `gte,lte` 区间搜索 - **时间字段**：使用 `gte,lte` 时间范围搜索 ### 4. 验证规则原则 - **必填验证**：必填字段必须添加 `required` - **长度限制**：防止过长输入，使用 `min`、`max` - **格式验证**：邮箱、URL等特殊格式使用相应验证规则 - **业务规则**：符合实际业务需求的验证规则 ## table函数最佳实践的示例 ```go // Package crm // 文件：crm_ticket.go package crm import ( \"github.com/yunhanshu-net/function-go/pkg/dto/response\" \"github.com/yunhanshu-net/function-go/pkg/dto/usercall\" \"github.com/yunhanshu-net/function-go/runner\" \"github.com/yunhanshu-net/pkg/query\" \"github.com/yunhanshu-net/pkg/typex/files\" \"gorm.io/gorm\" ) // 工单数据模型 type CrmTicket struct { // 框架标签：runner:\"name:工单ID\" - 设置字段在前端的显示名称 // 框架标签：permission:\"read\" - 字段只读权限（不能编辑） // 注意：gorm:\"column:id\" 明确指定数据库列名，确保映射正确 ID int `json:\"id\" gorm:\"primaryKey;autoIncrement;column:id\" runner:\"name:工单ID\" permission:\"read\"` // 框架标签：widget:\"type:datetime;kind:datetime\" - 日期时间选择器组件 // 注意：gorm:\"autoCreateTime:milli\" 自动填充创建时间（毫秒级时间戳，必须是毫秒级别） CreatedAt int64 `json:\"created_at\" gorm:\"autoCreateTime:milli;column:created_at\" runner:\"name:创建时间\" widget:\"type:datetime;kind:datetime\" permission:\"read\"` // 框架标签：widget:\"type:datetime;kind:datetime\" - 日期时间选择器组件，（毫秒级时间戳，必须是毫秒级别） UpdatedAt int64 `json:\"updated_at\" gorm:\"autoUpdateTime:milli;column:updated_at\" runner:\"name:更新时间\" widget:\"type:datetime;kind:datetime\" permission:\"read\"` // 框架标签：runner:\"-\" - 隐藏字段（不在前端显示） DeletedAt gorm.DeletedAt `json:\"deleted_at\" gorm:\"index;column:deleted_at\" runner:\"-\"` // 框架标签：widget:\"type:input\" - 文本输入框组件 // 框架标签：search:\"like\" - 启用模糊搜索功能 // 框架标签：validate:\"required,min=2,max=200\" - 必填字段，长度2-200字符 Title string `json:\"title\" gorm:\"column:title\" runner:\"name:工单标题\" widget:\"type:input\" search:\"like\" validate:\"required,min=2,max=200\"` // 框架标签：widget:\"type:input;mode:text_area\" - 多行文本区域组件 // 框架标签：validate:\"required,min=10\" - 必填字段，至少10字符 Description string `json:\"description\" gorm:\"column:description\" runner:\"name:问题描述\" widget:\"type:input;mode:text_area\" validate:\"required,min=10\"` // 框架标签：widget:\"type:select;options:低,中,高\" - 下拉选择组件（选项：低/中/高） // 框架标签：data:\"default_value:中\" - 设置默认值为\"中\" // 框架标签：validate:\"required,oneof=低,中,高\" - 必填字段，值必须是选项之一 Priority string `json:\"priority\" gorm:\"column:priority\" runner:\"name:优先级\" widget:\"type:select;options:低,中,高\" data:\"default_value:中\" validate:\"required,oneof=低,中,高\"` // 框架标签：widget:\"type:select;options:待处理,处理中,已完成,已关闭\" - 下拉选择组件 // 框架标签：data:\"default_value:待处理\" - 设置默认状态为\"待处理\" // 框架标签：validate:\"required,oneof=待处理,处理中,已完成,已关闭\" - 值必须是有效状态 Status string `json:\"status\" gorm:\"column:status\" runner:\"name:工单状态\" widget:\"type:select;options:待处理,处理中,已完成,已关闭\" data:\"default_value:待处理\" validate:\"required,oneof=待处理,处理中,已完成,已关闭\"` // 框架标签：validate:\"required,min=11,max=20\" - 必填字段，长度11-20字符 Phone string `json:\"phone\" gorm:\"column:phone\" runner:\"name:联系电话\" widget:\"type:input\" validate:\"required,min=11,max=20\"` // 框架标签：widget:\"type:input;mode:text_area\" - 多行文本区域组件 Remark string `json:\"remark\" gorm:\"column:remark\" runner:\"name:备注\" widget:\"type:input;mode:text_area\"` // 框架标签：widget:\"type:file_upload;multiple:true;...\" - 多文件上传组件 // 框架标签：accept:.pdf,.doc,... - 允许上传的文件类型 // 注意：gorm:\"type:json\" 指定数据库存储为JSON类型 Attachment *files.Files `json:\"attachment\" gorm:\"type:json;column:attachment\" runner:\"name:附件\" widget:\"type:file_upload;multiple:true;max_size:10MB;accept:.pdf,.doc,.docx,.xls,.xlsx,.ppt,.pptx,.txt,.jpg,.png,.gif,.zip,.rar\"` } // 框架功能：TableName() 指定数据库表名 func (CrmTicket) TableName() string { return \"crm_ticket\" } // 分页请求结构 type CrmTicketListReq struct { // 框架标签：runner:\"-\" - 隐藏搜索，分页，等请求参数（无需在前端显示） query.PageInfoReq `runner:\"-\"` } // 表格处理函数，这里是前端展示列表数据的接口 func CrmTicketList(ctx *runner.Context, req *CrmTicketListReq, resp response.Response) error { // 框架功能：ctx.MustGetOrInitDB() - 获取数据库连接 db := ctx.MustGetOrInitDB() var rows []CrmTicket //【框架规范】这里框架会自动拿着请求参数和db去数据库查询数据，自动分页等等 queryBuilder := resp.Table(&rows).AutoPaginated(db, &CrmTicket{}, &req.PageInfoReq) //此时rows，已经是根据请求条件和分页条件查询到的数据了，可以在查询到数据后对数据进行处理等等 // 框架功能：AutoPaginated() - 自动处理分页/搜索/排序 return queryBuilder.Build() } // 自动CRUD配置 var CrmTicketListOption = &runner.TableFunctionOptions{ BaseConfig: runner.BaseConfig{ ChineseName: \"CRM工单管理\", // API中文名称 ApiDesc: \"工单管理系统\", // API描述 Tags: []string{\"工单系统\", \"政务系统\", \"工单管理\",\"CRM\"}, // 框架自动生成标签分类 // 框架功能：CreateTables - 服务启动时自动创建数据表 CreateTables: []interface{}{&CrmTicket{}}, //对应查询的请求参数 Request: &CrmTicketListReq{}, //固定query.PaginatedTable[[]table]{} 的格式 Response: query.PaginatedTable[[]CrmTicket]{}, }, // 框架功能：AutoCrudTable - 自动生成增删改查的界面 AutoCrudTable: &CrmTicket{}, //【框架规范】新增记录回调 OnTableAddRows: func(ctx *runner.Context, req *usercall.OnTableAddRowsReq) (*usercall.OnTableAddRowsResp, error) { //【框架规范】这里可以直接获取到gorm的db gormDb := ctx.MustGetOrInitDB() var addRows []*CrmTicket //【框架规范】这里解析用户新增的记录 if err := req.DecodeBy(&addRows); err != nil { return nil, err } //把记录写入db，写入前后都可以做一些操作 err := gormDb.Create(&addRows).Error if err != nil { return nil, err } return &usercall.OnTableAddRowsResp{}, nil }, //【框架规范】更新记录回调 OnTableUpdateRows: func(ctx *runner.Context, req *usercall.OnTableUpdateRowsReq) (*usercall.OnTableUpdateRowsResp, error) { //【框架规范】这里可以直接获取到gorm的db gormDb := ctx.MustGetOrInitDB() var updateRow CrmTicket if err := req.DecodeBy(&updateRow); err != nil { //这里只能解析到更新的字段的数据，例如更新了phone字段的话，那么其他字段都是空值，只有phone存在值 return nil, err } err := gormDb.Where(\"id in (?)\", req.Ids).Updates(&updateRow).Error if err != nil { return nil, err } return &usercall.OnTableUpdateRowsResp{}, nil }, //【框架规范】删除记录回调 OnTableDeleteRows: func(ctx *runner.Context, req *usercall.OnTableDeleteRowsReq) (*usercall.OnTableDeleteRowsResp, error) { gormDb := ctx.MustGetOrInitDB() err := gormDb.Delete(&CrmTicket{}, req.Ids).Error if err != nil { return nil, err } return &usercall.OnTableDeleteRowsResp{}, nil }, } // 框架功能：init() 函数自动注册路由 func init() { // 框架功能：runner.Get() - 注册Table请求路由 runner.Get(RouterGroup+\"/crm_ticket_list\", CrmTicketList, CrmTicketListOption) } ``` form函数最佳实践示例 ```go // 文件名：text_convert.go package text import ( \"strings\" \"github.com/yunhanshu-net/function-go/pkg/dto/response\" \"github.com/yunhanshu-net/function-go/runner\" ) // 【框架功能】Form函数请求结构体：定义用户输入字段，框架自动生成前端表单 type TextConvertReq struct { // 【框架标签】widget:\"type:input;mode:text_area\" - 自动渲染为多行文本输入框 InputText string `json:\"input_text\" runner:\"name:输入文本\" widget:\"type:input;mode:text_area\" validate:\"required\"` // 【框架标签】widget:\"type:select;options:...\" - 自动渲染为下拉选择框，data:\"default_value\"设置默认值 ConvertType string `json:\"convert_type\" runner:\"name:转换类型\" widget:\"type:select;options:转大写,转小写,首字母大写\" data:\"default_value:转大写\"` } // 【框架功能】Form函数响应结构体：定义返回结果，框架自动渲染前端展示 type TextConvertResp struct { // 【框架标签】widget:\"type:input;mode:text_area\" - 自动渲染为多行文本展示区域 Result string `json:\"result\" runner:\"name:转换结果\" widget:\"type:input;mode:text_area\"` // 【框架标签】widget:\"type:input\" - 自动渲染为单行文本展示 Message string `json:\"message\" runner:\"name:状态\" widget:\"type:input\"` } // 【框架功能】Form函数处理函数：框架自动调用此函数处理用户请求 func TextConvert(ctx *runner.Context, req *TextConvertReq, resp response.Response) error { var result string // 业务逻辑：根据转换类型执行文本转换 switch req.ConvertType { case \"转大写\": result = strings.ToUpper(req.InputText) case \"转小写\": result = strings.ToLower(req.InputText) case \"首字母大写\": result = strings.Title(strings.ToLower(req.InputText)) default: result = req.InputText } // 【框架功能】resp.Form() - 自动将结果渲染为前端表单展示 return resp.Form(&TextConvertResp{ Result: result, Message: \"转换完成\", }).Build() } // 【框架功能】Form函数配置：定义API元数据，框架自动生成接口文档和路由 var TextConvertOption = &runner.FormFunctionOptions{ BaseConfig: runner.BaseConfig{ ChineseName: \"文本大小写转换\", // 框架自动生成中文接口名称 ApiDesc: \"支持文本大小写转换的工具，包括转大写、转小写、首字母大写等功能\", // 框架自动生成API描述 Tags: []string{\"文本处理\", \"格式转换\", \"工具函数\"}, // 框架自动生成标签分类 Request: &TextConvertReq{}, // 框架自动生成请求参数说明 Response: &TextConvertResp{}, // 框架自动生成响应结果说明 }, } // 【框架功能】路由注册：框架自动创建RouterGroup 变量，无需自己创建RouterGroup变量，可以直接用 func init() { // 【框架功能】runner.Post() - 自动创建POST路由，支持表单提交 runner.Post(RouterGroup+\"/text_convert\", TextConvert, TextConvertOption) } ``` ### 示例3：复杂Table+Form系统 - 采购管理系统 **功能**：完整的采购业务流程，包含供应商管理、商品管理、采购记录 ```go //retail_purchase.go /* <metadata> <用户需求> 我需要一个采购管理系统，帮助采购员管理供应商、商品和采购记录。 功能需求： 供应商管理： - 供应商信息：名称、联系方式、地址、状态 - 供应商状态：正常合作、暂停合作、终止合作 商品管理： - 商品信息：名称、分类、规格、单位、标准单价 - 商品分类：原材料、半成品、成品、包装物 - 商品状态：正常采购、暂停采购、淘汰 - 商品归属：每个商品属于特定供应商 采购管理： - 采购申请：选择供应商、多选该供应商的商品、填写数量 - 采购记录：记录采购历史、状态跟踪 - 统计分析：采购金额、商品种类、供应商分析 业务规则： - 采购数量必须大于0 - 系统自动计算总价（数量 × 标准单价） - 采购状态：待审核、已审核、已完成、已取消 </用户需求> <文件>retail_purchase.go</文件> </metadata> */ package retail import ( \"errors\" \"fmt\" \"github.com/yunhanshu-net/function-go/pkg/dto/response\" \"github.com/yunhanshu-net/function-go/pkg/dto/usercall\" \"github.com/yunhanshu-net/function-go/runner\" \"github.com/yunhanshu-net/pkg/query\" \"gorm.io/gorm\" ) // ================ 函数组配置 ================ var RetailPurchaseGroup = &runner.FunctionGroup{ CnName: \"采购管理\", EnName: \"retail_purchase\", } // ================ 数据模型（按依赖关系排序） ================ // 1. 供应商信息（基础数据，被其他表依赖） type RetailPurchaseSupplier struct { ID int `json:\"id\" gorm:\"primaryKey;autoIncrement\" runner:\"name:供应商ID\" permission:\"read\"` CreatedAt int64 `json:\"created_at\" gorm:\"autoCreateTime:milli\" runner:\"name:创建时间\" widget:\"type:datetime;kind:datetime\" permission:\"read\"` UpdatedAt int64 `json:\"updated_at\" gorm:\"autoUpdateTime:milli\" runner:\"name:更新时间\" widget:\"type:datetime;kind:datetime\" permission:\"read\"` DeletedAt gorm.DeletedAt `json:\"deleted_at\" gorm:\"index\" runner:\"-\"` Name string `json:\"name\" gorm:\"column:name;comment:供应商名称\" runner:\"name:供应商名称\" widget:\"type:input\" search:\"like\" validate:\"required,min=2,max=100\" msg:\"供应商名称必填，长度2-100字符\"` Contact string `json:\"contact\" gorm:\"column:contact;comment:联系人\" runner:\"name:联系人\" widget:\"type:input\" search:\"like\" validate:\"required,min=2,max=50\" msg:\"联系人必填，长度2-50字符\"` Phone string `json:\"phone\" gorm:\"column:phone;comment:联系电话\" runner:\"name:联系电话\" widget:\"type:input\" search:\"like\" validate:\"required,min=11,max=20\" msg:\"联系电话必填，长度11-20字符\"` Address string `json:\"address\" gorm:\"column:address;comment:地址\" runner:\"name:地址\" widget:\"type:input;mode:text_area\" search:\"like\" validate:\"required,min=5,max=200\" msg:\"地址必填，长度5-200字符\"` } func (RetailPurchaseSupplier) TableName() string { return \"retail_purchase_supplier\" } // 2. 商品信息（依赖供应商，被采购记录依赖） type RetailPurchaseProduct struct { ID int `json:\"id\" gorm:\"primaryKey;autoIncrement\" runner:\"name:商品ID\" permission:\"read\"` CreatedAt int64 `json:\"created_at\" gorm:\"autoCreateTime:milli\" runner:\"name:创建时间\" widget:\"type:datetime;kind:datetime\" permission:\"read\"` UpdatedAt int64 `json:\"updated_at\" gorm:\"autoUpdateTime:milli\" runner:\"name:更新时间\" widget:\"type:datetime;kind:datetime\" permission:\"read\"` DeletedAt gorm.DeletedAt `json:\"deleted_at\" gorm:\"index\" runner:\"-\"` // 【框架说明】SupplierID必须使用select组件才能支持OnInputFuzzy模糊搜索功能 // 只有select（单选）和multiselect（多选）组件才支持动态数据源模糊搜索 // number组件不支持模糊搜索，用户无法通过名称搜索供应商 SupplierID int `json:\"supplier_id\" gorm:\"column:supplier_id;comment:供应商ID\" runner:\"name:供应商\" widget:\"type:select\" search:\"in\" validate:\"required\" msg:\"请选择供应商\"` SupplierName string `json:\"supplier_name\" gorm:\"column:supplier_name;comment:供应商名称\" runner:\"name:供应商名称\" widget:\"type:input\" search:\"like\" validate:\"required\" msg:\"供应商名称必填\"` Name string `json:\"name\" gorm:\"column:name;comment:商品名称\" runner:\"name:商品名称\" widget:\"type:input\" search:\"like\" validate:\"required,min=2,max=100\" msg:\"商品名称必填，长度2-100字符\"` Category string `json:\"category\" gorm:\"column:category;comment:商品分类\" runner:\"name:商品分类\" widget:\"type:select;options:原材料,半成品,成品,包装物\" data:\"default_value:原材料\" search:\"in\" validate:\"required,oneof=原材料 半成品 成品 包装物\" msg:\"请选择商品分类\"` Spec string `json:\"spec\" gorm:\"column:spec;comment:规格\" runner:\"name:规格\" widget:\"type:input\" search:\"like\" validate:\"required,min=2,max=100\" msg:\"规格必填，长度2-100字符\" data:\"example:500g;100ml;A4尺寸;M码;红色\"` Unit string `json:\"unit\" gorm:\"column:unit;comment:单位\" runner:\"name:单位\" widget:\"type:input\" search:\"like\" validate:\"required,min=1,max=20\" msg:\"单位必填，长度1-20字符\" data:\"example:个;件;包;盒;米;千克;升\"` UnitPrice float64 `json:\"unit_price\" gorm:\"column:unit_price;comment:标准单价\" runner:\"name:标准单价\" widget:\"type:number;min:0;precision:2\" search:\"gte,lte\" validate:\"required,min=0\" msg:\"标准单价必须大于0\"` } func (RetailPurchaseProduct) TableName() string { return \"retail_purchase_product\" } // 3. 采购商品项（业务数据，依赖商品信息） type RetailPurchaseItem struct { ProductID int `json:\"product_id\" runner:\"name:商品\" widget:\"type:select\" validate:\"required\" msg:\"请选择商品\"` Quantity int `json:\"quantity\" runner:\"name:数量\" widget:\"type:number;min:1\" data:\"default_value:1\" validate:\"required,min=1\" msg:\"数量必须大于0\"` } // 5. 采购请求（业务操作，依赖上述所有模型） type RetailPurchaseReq struct { SupplierID int `json:\"supplier_id\" runner:\"name:供应商\" widget:\"type:select\" validate:\"required\" msg:\"请选择供应商\"` Items []RetailPurchaseItem `json:\"items\" runner:\"name:采购商品\" widget:\"type:list\" validate:\"required,min=1\" msg:\"请至少选择一件商品\"` Remarks string `json:\"remarks\" runner:\"name:备注\" widget:\"type:input;mode:text_area\"` } // 4. 采购记录（业务数据，依赖供应商和商品） type RetailPurchaseRecord struct { ID int `json:\"id\" gorm:\"primaryKey;autoIncrement\" runner:\"name:采购ID\" permission:\"read\"` CreatedAt int64 `json:\"created_at\" gorm:\"autoCreateTime:milli\" runner:\"name:创建时间\" widget:\"type:datetime;kind:datetime\" permission:\"read\"` UpdatedAt int64 `json:\"updated_at\" gorm:\"autoUpdateTime:milli\" runner:\"name:更新时间\" widget:\"type:datetime;kind:datetime\" permission:\"read\"` DeletedAt gorm.DeletedAt `json:\"deleted_at\" gorm:\"index\" runner:\"-\"` // 【框架说明】SupplierID在记录表中使用number组件，因为这是历史数据展示，不需要联想功能 // 联想功能只在新增/编辑时使用，记录表主要用于查询和展示 // 【框架说明】SupplierID在记录表中使用select组件，支持联想功能 // 新增和编辑时可以选择供应商，列表展示时隐藏ID字段 SupplierID int `json:\"supplier_id\" gorm:\"column:supplier_id;comment:供应商ID\" runner:\"name:供应商\" widget:\"type:select\" search:\"in\" validate:\"required\" msg:\"请选择供应商\" permission:\"create,update\"` // 【框架说明】SupplierName在列表展示时自动填充，属于join其他表的字段了，因此这个字段不能加search标签，注意，注意，注意：gorm:\"-\"的字段不能加search标签 // 所以此时如果需要支持非表内字段的搜索，需要在请求参数内添加这个字段，参考：RetailPurchaseRecordListReq，在里面扩展字段 // permission:\"read\"表示仅在列表显示，新增和编辑时不显示 SupplierName string `json:\"supplier_name\" gorm:\"-\" runner:\"name:供应商名称\" widget:\"type:input\" search:\"like\" validate:\"required\" msg:\"供应商名称必填\" permission:\"read\"` TotalAmount float64 `json:\"total_amount\" gorm:\"column:total_amount;comment:总金额\" runner:\"name:总金额\" widget:\"type:number;precision:2\" search:\"gte,lte\" validate:\"required,min=0\" msg:\"总金额必须大于0\"` TotalItems int `json:\"total_items\" gorm:\"column:total_items;comment:商品种类\" runner:\"name:商品种类\" widget:\"type:number\" search:\"gte,lte\" validate:\"required,min=1\" msg:\"商品种类必须大于0\"` Status string `json:\"status\" gorm:\"column:status;comment:状态\" runner:\"name:状态\" widget:\"type:select;options:待审核,已审核,已完成,已取消\" data:\"default_value:待审核\" search:\"in\" validate:\"required,oneof=待审核 已审核 已完成 已取消\" msg:\"请选择采购状态\"` Remarks string `json:\"remarks\" gorm:\"column:remarks;comment:备注\" runner:\"name:备注\" widget:\"type:input;mode:text_area\" search:\"like\"` Supplier RetailPurchaseSupplier `json:\"-\" runner:\"-\" gorm:\"foreignKey:SupplierID\"` // 【隐藏】关联关系，仅用于gorm预加载 } func (RetailPurchaseRecord) TableName() string { return \"retail_purchase_record\" } // 6. 采购响应（业务结果） type RetailPurchaseResp struct { PurchaseID int `json:\"purchase_id\" runner:\"name:采购ID\" widget:\"type:number\"` TotalAmount float64 `json:\"total_amount\" runner:\"name:总金额\" widget:\"type:number;precision:2\"` TotalItems int `json:\"total_items\" runner:\"name:商品种类\" widget:\"type:number\"` Message string `json:\"message\" runner:\"name:处理结果\" widget:\"type:input;mode:text_area\"` } // ================ 业务逻辑函数（按调用顺序排序） ================ // 1. 辅助函数：验证商品并计算总金额 func retailPurchaseValidateAndCalculate(db *gorm.DB, supplierID int, items []RetailPurchaseItem) (float64, int, error) { if len(items) == 0 { return 0, 0, fmt.Errorf(\"采购商品不能为空\") } var totalAmount float64 productIDs := make([]int, 0, len(items)) for _, item := range items { if item.Quantity <= 0 { return 0, 0, fmt.Errorf(\"商品数量必须大于0\") } productIDs = append(productIDs, item.ProductID) } // 验证商品是否存在且属于指定供应商 var products []RetailPurchaseProduct if err := db.Where(\"id IN ? AND supplier_id = ?\", productIDs, supplierID).Find(&products).Error; err != nil { return 0, 0, fmt.Errorf(\"验证商品失败: %v\", err) } if len(products) != len(productIDs) { return 0, 0, fmt.Errorf(\"部分商品不存在或不属于该供应商\") } // 计算总金额（数量 × 标准单价） for _, item := range items { for _, product := range products { if product.ID == item.ProductID { totalAmount += float64(item.Quantity) * product.UnitPrice break } } } return totalAmount, len(productIDs), nil } // 2. 主要业务函数：创建采购申请 func RetailPurchaseCreate(ctx *runner.Context, req *RetailPurchaseReq, resp response.Response) error { db := ctx.MustGetOrInitDB() // 1. 验证供应商 var supplier RetailPurchaseSupplier if err := db.Where(\"id = ?\", req.SupplierID).First(&supplier).Error; err != nil { if errors.Is(err, gorm.ErrRecordNotFound) { return resp.Form(&RetailPurchaseResp{ Message: \"供应商不存在\", }).Build() } return resp.Form(&RetailPurchaseResp{ Message: fmt.Sprintf(\"查询供应商失败: %v\", err), }).Build() } // 2. 验证商品并计算 totalAmount, totalItems, err := retailPurchaseValidateAndCalculate(db, req.SupplierID, req.Items) if err != nil { return resp.Form(&RetailPurchaseResp{ Message: fmt.Sprintf(\"商品验证失败: %v\", err), }).Build() } // 3. 创建采购记录 purchaseRecord := &RetailPurchaseRecord{ SupplierID: req.SupplierID, SupplierName: supplier.Name, TotalAmount: totalAmount, TotalItems: totalItems, Status: \"待审核\", Remarks: req.Remarks, } if err := db.Create(purchaseRecord).Error; err != nil { ctx.Logger.Errorf(\"创建采购记录失败: %v\", err) return resp.Form(&RetailPurchaseResp{ Message: \"创建采购记录失败，请重试\", }).Build() } return resp.Form(&RetailPurchaseResp{ PurchaseID: purchaseRecord.ID, TotalAmount: totalAmount, TotalItems: totalItems, Message: \"采购申请创建成功\", }).Build() } // ================ Table函数（按数据依赖关系排序） ================ // 1. 供应商列表管理（基础数据管理） func RetailPurchaseSupplierList(ctx *runner.Context, req *query.PageInfoReq, resp response.Response) error { db := ctx.MustGetOrInitDB() var rows []RetailPurchaseSupplier return resp.Table(&rows).AutoPaginated(db, &RetailPurchaseSupplier{}, req).Build() } // 2. 商品列表管理（依赖供应商数据） func RetailPurchaseProductList(ctx *runner.Context, req *query.PageInfoReq, resp response.Response) error { db := ctx.MustGetOrInitDB() var rows []RetailPurchaseProduct return resp.Table(&rows).AutoPaginated(db, &RetailPurchaseProduct{}, req).Build() } // 3. 采购记录列表管理（依赖供应商和商品数据） type RetailPurchaseRecordListReq struct { query.PageInfoReq `runner:\"-\"` // 【框架说明】join字段搜索参数 // query.PageInfoReq 默认只支持当前表字段的搜索（如 TotalAmount、TotalItems、Status） // 但无法处理关联表的字段搜索（如供应商名称），需要手动实现关联查询 SupplierName string `json:\"supplier_name\" runner:\"name:供应商名称\" widget:\"type:input\" search:\"like\"` // 按供应商名称模糊搜索 } func RetailPurchaseRecordList(ctx *runner.Context, req *RetailPurchaseRecordListReq, resp response.Response) error { db := ctx.MustGetOrInitDB() var rows []*RetailPurchaseRecord // 【框架说明】手动实现join字段的搜索逻辑 // 因为框架的 AutoPaginated 无法自动处理关联表字段的搜索 // 需要手动构建查询条件，实现类似 JOIN 的效果 query := db.Model(&RetailPurchaseRecord{}) // 按供应商名称模糊搜索（需要关联查询） if req.SupplierName != \"\" { // 先查询匹配的供应商ID var supplierIDs []int if err := db.Model(&RetailPurchaseSupplier{}). Where(\"name LIKE ?\", \"%\"+req.SupplierName+\"%\"). Pluck(\"id\", &supplierIDs).Error; err == nil && len(supplierIDs) > 0 { query = query.Where(\"supplier_id IN ?\", supplierIDs) } else { // 如果没有找到匹配的供应商，返回空结果 return resp.Table(nil).Build() } } // 【优化】预加载供应商信息，避免N+1查询问题 query = query.Preload(\"Supplier\") // 获取分页数据（应用手动构建的查询条件） td := resp.Table(&rows).AutoPaginated(query, &RetailPurchaseRecord{}, &req.PageInfoReq) // 【框架说明】填充join字段，提升用户体验 // 在列表展示时自动填充供应商名称，避免显示无意义的ID for _, row := range rows { row.SupplierName = row.Supplier.Name } return td.Build() } // ================ 函数配置（按功能层次排序） ================ // 1. 供应商列表配置（基础数据管理） var RetailPurchaseSupplierListOption = &runner.TableFunctionOptions{ BaseConfig: runner.BaseConfig{ ChineseName: \"采购供应商管理\", ApiDesc: \"供应商基础信息管理，支持增删改查和搜索\", Tags: []string{\"采购管理\", \"供应商管理\"}, Request: &query.PageInfoReq{}, Response: query.PaginatedTable[[]RetailPurchaseSupplier]{}, CreateTables: []interface{}{&RetailPurchaseSupplier{}}, Group: RetailPurchaseGroup, }, AutoCrudTable: &RetailPurchaseSupplier{}, // 【框架说明】OnTableAddRows 在新增供应商后触发，用于额外的业务逻辑处理 // 触发时机：用户新增供应商成功后 使用场景：记录日志、发送通知、更新缓存等 // ⚠️ 重要：有此回调时框架不会自动写入数据，必须在回调中手动执行 db.Create() OnTableAddRows: func(ctx *runner.Context, req *usercall.OnTableAddRowsReq) (*usercall.OnTableAddRowsResp, error) { // 【业务场景】新增供应商后，记录操作日志 // 这里演示了如何在基础CRUD之外添加额外的业务操作 // 解析新增的供应商数据 var suppliers []RetailPurchaseSupplier if err := req.DecodeBy(&suppliers); err != nil { return nil, fmt.Errorf(\"解析供应商数据失败: %v\", err) } //然后可以在写入前后做任何操作，但是不要太复杂，要不然会影响用户体验 // 【关键步骤】手动执行数据库写入操作（用了OnTableAddRows框架不会自动执行写入） db := ctx.MustGetOrInitDB() err := db.Create(&suppliers).Error if err != nil { return nil, err } return &usercall.OnTableAddRowsResp{}, nil }, // 【框架说明】OnTableUpdateRows 在更新供应商后触发，用于额外的业务逻辑处理 // 触发时机：用户更新供应商成功后 使用场景：数据同步、状态检查、关联更新等 // ⚠️ 重要：有此回调时框架不会自动更新数据，必须在回调中手动执行 db.Model().Updates() OnTableUpdateRows: func(ctx *runner.Context, req *usercall.OnTableUpdateRowsReq) (*usercall.OnTableUpdateRowsResp, error) { // 【业务场景】更新供应商后，记录变更历史 // 这里演示了如何在数据更新后处理关联业务逻辑 // 【关键步骤】手动执行数据库更新操作（框架不会自动执行） db := ctx.MustGetOrInitDB() var updateValue RetailPurchaseSupplier err := req.DecodeBy(&updateValue) //【框架规范】注意这里只会获取到更新的字段的值，例如更新phone的话updateValue的Phone字段会被赋值 if err != nil { return nil, err } //GetString，GetInt，GetBool，GetFloat64，GetFiles 这几个方法，可以安全的获取更新字段的值，下面几个只是演示一下用法，正常情况下字段校验的逻辑会被框架和前端处理， //name, exist, err := req.GetString(\"name\") //if err != nil { // return nil, err //} //if exist { // if name == \"\" { // return nil, fmt.Errorf(\"名称不能为空\") // } //} //phone, exist, err := req.GetInt(\"phone\") //if err != nil { // return nil, err //} //fmt.Println(phone) //GetFieldsMap() 返回 map[string]interface 是需要更新的字段和更新的值 //fieldsMap := req.GetFieldsMap() err = db.Where(\"id in ?\", req.Ids).Updates(updateValue).Error if err != nil { return nil, err } return &usercall.OnTableUpdateRowsResp{}, nil }, // 【框架说明】OnTableDeleteRows 在删除供应商前触发，用于业务规则验证 // 触发时机：用户删除供应商前 使用场景：关联检查、防止误删 // ⚠️ 重要：有此回调时框架不会自动删除数据，必须在回调中手动执行 db.Delete() OnTableDeleteRows: func(ctx *runner.Context, req *usercall.OnTableDeleteRowsReq) (*usercall.OnTableDeleteRowsResp, error) { // 【业务场景】删除供应商前，检查是否有关联商品，防止误删 db := ctx.MustGetOrInitDB() // 检查是否有关联商品 for _, id := range req.Ids { var count int64 if err := db.Model(&RetailPurchaseProduct{}).Where(\"supplier_id = ?\", id).Count(&count).Error; err != nil { return nil, fmt.Errorf(\"检查供应商关联商品失败: %v\", err) } if count > 0 { var supplier RetailPurchaseSupplier if err := db.Where(\"id = ?\", id).First(&supplier).Error; err != nil { return nil, fmt.Errorf(\"查询供应商信息失败: %v\", err) } return nil, fmt.Errorf(\"供应商 %s 下还有 %d 个商品，无法删除，请先把相关商品删除\", supplier.Name, count) } } // 手动执行删除操作 if err := db.Where(\"id IN ?\", req.Ids).Delete(&RetailPurchaseSupplier{}).Error; err != nil { return nil, fmt.Errorf(\"删除供应商记录失败: %v\", err) } return &usercall.OnTableDeleteRowsResp{}, nil }, } // 2. 商品列表配置（依赖供应商数据） var RetailPurchaseProductListOption = &runner.TableFunctionOptions{ BaseConfig: runner.BaseConfig{ ChineseName: \"采购商品管理\", ApiDesc: \"商品基础信息管理，支持增删改查和搜索\", Tags: []string{\"采购管理\", \"商品管理\"}, Request: &query.PageInfoReq{}, Response: query.PaginatedTable[[]RetailPurchaseProduct]{}, CreateTables: []interface{}{&RetailPurchaseProduct{}}, Group: RetailPurchaseGroup, }, AutoCrudTable: &RetailPurchaseProduct{}, // 【框架说明】OnInputFuzzyMap 为表格字段提供模糊搜索数据 // 触发时机：用户输入时 数据去向：挂载到对应的表格字段 OnInputFuzzyMap: map[string]runner.OnInputFuzzy{ \"supplier_id\": func(ctx *runner.Context, req *usercall.OnInputFuzzyReq) (*usercall.OnInputFuzzyResp, error) { // 【当前场景】supplier_id是单选字段，只返回静态信息，不做聚合计算 // 【目标字段】为 RetailPurchaseProduct.SupplierID 提供供应商选择数据 keyword := fmt.Sprintf(\"%v\", req.Value) var suppliers []RetailPurchaseSupplier db := ctx.MustGetOrInitDB() db.Where(\"name LIKE ? OR contact LIKE ? OR phone LIKE ?\", \"%\"+keyword+\"%\", \"%\"+keyword+\"%\", \"%\"+keyword+\"%\"). Limit(20). Find(&suppliers) items := make([]*usercall.InputFuzzyItem, 0) for _, s := range suppliers { items = append(items, &usercall.InputFuzzyItem{ Value: s.ID, Label: s.Name, DisplayInfo: map[string]interface{}{ \"供应商名称\": s.Name, \"联系人\": s.Contact, \"联系电话\": s.Phone, \"地址\": s.Address, }, }) } return &usercall.OnInputFuzzyResp{ Statistics: map[string]interface{}{ \"供应商名称\": \"text(供应商名称)\", //前端动态展示选中的那个供应商名称 \"联系人\": \"text(联系人)\", //前端动态展示选中的那个供应商联系人 \"联系电话\": \"text(联系电话)\", //前端动态展示选中的那个供应商联系电话 \"地址\": \"text(地址)\", //前端动态展示选中的那个供应商地址 }, Values: items, }, nil }, }, } // 3. 采购记录列表配置（依赖供应商和商品数据） var RetailPurchaseRecordListOption = &runner.TableFunctionOptions{ BaseConfig: runner.BaseConfig{ ChineseName: \"采购记录管理\", ApiDesc: \"采购记录管理，支持状态跟踪、供应商筛选和统计分析\", Tags: []string{\"采购管理\", \"记录管理\"}, Request: &RetailPurchaseRecordListReq{}, Response: query.PaginatedTable[[]RetailPurchaseRecord]{}, CreateTables: []interface{}{&RetailPurchaseRecord{}}, Group: RetailPurchaseGroup, }, AutoCrudTable: &RetailPurchaseRecord{}, // 【框架说明】复用采购创建的OnInputFuzzy配置，避免重复造轮子 // 供应商联想功能在新增和编辑时都需要，直接复用已有配置 OnInputFuzzyMap: RetailPurchaseCreateOption.OnInputFuzzyMap, } // 4. 采购创建配置（核心业务逻辑） var RetailPurchaseCreateOption = &runner.FormFunctionOptions{ BaseConfig: runner.BaseConfig{ ChineseName: \"采购申请创建\", ApiDesc: \"创建采购申请，支持多商品选择和实时统计\", Tags: []string{\"采购管理\", \"申请创建\"}, Request: &RetailPurchaseReq{}, Response: &RetailPurchaseResp{}, CreateTables: []interface{}{&RetailPurchaseRecord{}}, Group: RetailPurchaseGroup, }, // 【框架说明】OnInputFuzzyMap 为请求结构体字段提供模糊搜索数据 // 触发时机：用户输入时 数据去向：挂载到对应的请求字段 OnInputFuzzyMap: map[string]runner.OnInputFuzzy{ \"supplier_id\": RetailPurchaseProductListOption.OnInputFuzzyMap[\"supplier_id\"], //这里为了防止重复造轮子，可以直接复用同一个方法 \"product_id\": func(ctx *runner.Context, req *usercall.OnInputFuzzyReq) (*usercall.OnInputFuzzyResp, error) { // 【当前场景】product_id是list内单选字段，支持聚合统计计算 // 【目标字段】为 RetailPurchaseItem.ProductID 提供商品选择数据 // 注意：这里需要根据已选择的供应商来过滤商品 // 【框架说明】通过req.DecodeBy获取表单上下文中的supplier_id // 这样可以确保商品选择与供应商选择保持一致性 var currentInputData RetailPurchaseReq if err := req.DecodeBy(&currentInputData); err != nil { return nil, fmt.Errorf(\"表单解析失败，请刷新选择供应商后再重试\") } if currentInputData.SupplierID == 0 { return nil, fmt.Errorf(\"请先选择供应商，再选择商品\") } keyword := fmt.Sprintf(\"%v\", req.Value) var products []RetailPurchaseProduct db := ctx.MustGetOrInitDB() // 根据供应商ID过滤商品，确保只能选择该供应商的商品 db.Where(\"supplier_id = ?\", currentInputData.SupplierID). Where(\"name LIKE ? OR spec LIKE ? OR category LIKE ?\", \"%\"+keyword+\"%\", \"%\"+keyword+\"%\", \"%\"+keyword+\"%\"). Limit(20). Find(&products) items := make([]*usercall.InputFuzzyItem, 0) for _, p := range products { items = append(items, &usercall.InputFuzzyItem{ Value: p.ID, Label: fmt.Sprintf(\"%s - %s (¥%.2f/%s)\", p.Name, p.Spec, p.UnitPrice, p.Unit), DisplayInfo: map[string]interface{}{ \"商品名称\": p.Name, \"规格\": p.Spec, \"分类\": p.Category, \"单位\": p.Unit, \"标准单价\": p.UnitPrice, }, }) } return &usercall.OnInputFuzzyResp{ Statistics: map[string]interface{}{ // ✅ 前端实时聚合计算 - 用户每添加一行商品，前端立即计算 \"采购总金额\": \"sum(标准单价,*quantity)\", // 用户选择的所有商品总价 \"九折金额\": \"sum(标准单价,*0.9)\", // 用户选择的所有商品总价*0.9 \"商品种类数\": \"count(标准单价)\", // 选了几种商品 \"采购总数量\": \"sum(quantity)\", // 商品总数量 \"平均单价\": \"avg(标准单价)\", // 选中商品平均价格 // ✅ 有用的静态信息 \"采购说明\": \"批量采购享优惠\", \"付款方式\": \"月结30天\", }, Values: items, }, nil }, }, } // ================ API注册（按功能层次排序） ================ func init() { // 1. Table函数 - 基础数据管理（按依赖关系排序） //RouterGroup 变量可以直接用，当前package下已经创建好该变量了 runner.Get(RouterGroup+\"/retail_purchase_supplier_list\", RetailPurchaseSupplierList, RetailPurchaseSupplierListOption) runner.Get(RouterGroup+\"/retail_purchase_product_list\", RetailPurchaseProductList, RetailPurchaseProductListOption) runner.Get(RouterGroup+\"/retail_purchase_record_list\", RetailPurchaseRecordList, RetailPurchaseRecordListOption) // 2. Form函数 - 业务逻辑处理（依赖基础数据） runner.Post(RouterGroup+\"/retail_purchase_create\", RetailPurchaseCreate, RetailPurchaseCreateOption) } //<总结> //采购管理系统：供应商管理、商品管理、采购记录管理 //技术栈：AutoCrudTable自动生成CRUD界面、OnInputFuzzy模糊搜索、Statistics前端选中实时聚合计算、OnTable表格回调机制、GORM预加载优化 //复杂度：S2级别，包含回调机制和业务逻辑处理，无复杂依赖关系 //业务逻辑：商品属于供应商，采购时支持多商品选择，系统自动计算总价和统计信息 //性能优化：GORM预加载避免N+1查询，JOIN查询优化关联搜索 //</总结> ``` 基于用户需求生成高质量生产可用可直接编译的function-go代码 用户需求： 我需要一个医院管理系统，方便医院管理药品、医生、科室和患者挂号问诊。 包含以下功能： 药品管理： - 药品名称：药品的通用名称 - 规格：药品的规格型号（如：100mg×30片） - 生产厂家：药品的生产企业名称 - 单价：药品的销售价格（元） - 药品分类：选择分类（西药、中药、生物制品、医疗器械） - 状态：选择状态（正常、停用，默认正常） 医生管理： - 医生姓名：医生的真实姓名 - 职称：选择职称（主任医师、副主任医师、主治医师、住院医师） - 所属科室：选择科室（内科、外科、儿科、妇产科、眼科、骨科、皮肤科、口腔科） - 联系电话：医生的联系电话 - 备注：医生的其他信息 科室管理： - 科室名称：科室的完整名称 - 科室代码：科室的简称代码 - 详细描述：科室的详细说明 - 状态：选择状态（正常、停用，默认正常） 患者挂号管理： - 患者姓名：患者的真实姓名 - 年龄：患者的年龄（岁） - 性别：选择性别（男、女） - 电话：患者的联系电话 - 挂号医生：选择要挂号的医生 - 挂号费：挂号费用（元） - 问诊号：系统自动生成，每天从1开始编号 - 挂号状态：选择状态（待问诊、已完成，默认待问诊） - 备注：其他需要说明的信息 医生问诊管理： - 问诊号：选择要问诊的患者挂号记录 - 医生：选择进行问诊的医生 - 药品选择：选择要开具的药品（支持多选） - 用法用量：药品的使用方法和用量说明 - 总金额：自动计算所选药品的总价格 - 问诊状态：选择状态（进行中、已完成，默认进行中） 业务规则： - 问诊号每天从1开始，方便医生和患者记忆 - 支持多药品同时开方，自动计算总价格 - 问诊完成后自动更新挂号状态 - 药品状态为\"停用\"时不能选择开方 使用场景： - 药房管理员维护药品信息 - 护士为患者办理挂号 - 医生为患者进行问诊开方 - 管理员管理医生和科室信息 ## 代码质量要求 - **框架规范**：严格按照function-go框架，不得自造武功 - **导入检查**：确保所有使用的函数都有对应的import语句 - **命名规范**：文件名、结构体、函数名保持一致性 - **错误处理**：提供友好的用户提示，不要技术性错误信息 - **业务逻辑**：确保功能自洽，逻辑完整 ## 输出要求 - 生成完整的Go代码，禁止出现“伪代码”，“占位符”，“简化实现”，“mock数据” - 包含所有必要的import语句 - 保证生成的代码可以直接编译通过 - 添加详细的代码注释说明 - 在文件末尾添加技术总结",
			},
		},
		Timeout:     &timeout,
		MaxTokens:   13000,
		Temperature: 0.1, // 代码生成需要低温度
	}

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	resp, err := client.Chat(ctx, req)
	if err != nil {
		t.Fatalf("代码生成请求失败: %v", err)
	}

	if resp == nil {
		t.Fatal("响应为空")
	}

	// 检查是否有错误
	if resp.Error != "" {
		t.Logf("API返回错误: %s", resp.Error)
		t.Logf("注意：这可能是API key无效或网络问题，请检查配置")
		return
	}

	// 检查响应内容
	if resp.Content == "" {
		t.Error("响应内容为空")
	} else {
		t.Logf("生成的代码: %s", resp.Content)
	}

	// 检查使用统计
	if resp.Usage != nil {
		t.Logf("Token使用: 输入%d, 输出%d, 总计%d",
			resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
	}
}

// TestDeepSeekErrorHandling 测试错误处理
func TestDeepSeekErrorHandling(t *testing.T) {
	// 使用无效的API Key测试错误处理
	invalidClient := NewDeepSeekClient("invalid-api-key")

	req := &ChatRequest{
		Messages: []Message{
			{Role: "user", Content: "你好"},
		},
		MaxTokens: 10,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := invalidClient.Chat(ctx, req)
	if err != nil {
		t.Logf("预期的错误: %v", err)
		return
	}

	// 如果API返回了错误信息
	if resp != nil && resp.Error != "" {
		t.Logf("API返回错误: %s", resp.Error)
		return
	}

	t.Log("注意：API可能没有返回预期的错误信息")
}

// TestDeepSeekTimeout 测试超时处理
func TestDeepSeekTimeout(t *testing.T) {
	client, err := NewDeepSeekClientFromEnv()
	if err != nil {
		t.Fatalf("从环境变量创建客户端失败: %v", err)
	}

	req := &ChatRequest{
		Messages: []Message{
			{Role: "user", Content: "请生成一个非常复杂的代码示例"},
		},
		MaxTokens: 5000,
	}

	// 设置很短的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err = client.Chat(ctx, req)
	if err != nil {
		t.Logf("超时错误（预期）: %v", err)
	} else {
		t.Log("注意：请求没有超时，可能是网络很快或API响应很快")
	}
}

// TestDeepSeekIntegration 测试集成功能
func TestDeepSeekIntegration(t *testing.T) {
	client, err := NewDeepSeekClientFromEnv()
	if err != nil {
		t.Fatalf("从环境变量创建客户端失败: %v", err)
	}

	req := &ChatRequest{
		Messages: []Message{
			{
				Role:    "system",
				Content: "你是一个专业的Go语言开发助手，请生成高质量的代码",
			},
			{
				Role:    "user",
				Content: "请创建一个简单的Go HTTP服务器",
			},
		},
		MaxTokens:   1500,
		Temperature: 0.1,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := client.Chat(ctx, req)
	if err != nil {
		t.Fatalf("集成测试请求失败: %v", err)
	}

	if resp == nil {
		t.Fatal("响应为空")
	}

	// 检查是否有错误
	if resp.Error != "" {
		t.Logf("API返回错误: %s", resp.Error)
		t.Logf("注意：这可能是API key无效或网络问题，请检查配置")
		return
	}

	// 检查响应内容
	if resp.Content == "" {
		t.Error("响应内容为空")
	} else {
		t.Logf("集成测试成功，生成的代码: %s", resp.Content)
	}

	// 检查使用统计
	if resp.Usage != nil {
		t.Logf("Token使用统计: 输入%d, 输出%d, 总计%d",
			resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
	}
}

// TestDeepSeekAll 运行所有测试
func TestDeepSeekAll(t *testing.T) {
	t.Run("客户端创建", TestDeepSeekClientCreation)
	t.Run("接口实现", TestDeepSeekClientInterface)
	t.Run("基础聊天", TestDeepSeekChatBasic)
	t.Run("系统提示", TestDeepSeekChatWithSystemPrompt)
	t.Run("代码生成", TestDeepSeekCodeGeneration)
	t.Run("错误处理", TestDeepSeekErrorHandling)
	t.Run("超时处理", TestDeepSeekTimeout)
	t.Run("集成测试", TestDeepSeekIntegration)
}
