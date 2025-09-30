# function-go - AI原生全栈开发框架

> **专为大模型设计的Go全栈开发框架**
> AI原生结构体标签系统 | 一个文件 = 一个完整Web应用 | 大模型生成代码直接编译部署 | 框架内部有运行入口，专注业务函数

## 🏷️ 标签系统 - 大模型代码生成核心

### 设计理念

function-go通过结构体标签实现**声明式配置**，让大模型能够：

- **零样板代码**：一个结构体 = 完整Web应用
- **AI友好**：标签即配置，大模型直接理解
- **业务导向**：专注业务逻辑，框架处理技术细节

### 💡 大模型使用指南

**记住：标签就是配置，配置就是功能！**

| 标签 | 作用         | 大模型使用场景   |
| ------ | -------------- | ------------------ |
| `runner:"name:字段名"`     | 前端显示名称 | 所有字段都需要   |
| `widget:"type:组件类型"`     | UI组件选择   | 根据数据类型选择 |
| `validate:"required"`     | 数据验证     | 必填字段必须加   |
| `search:"like"`     | 搜索功能     | 需要搜索的字段   |

### 快速上手模板

```go
// 大模型代码生成模板
type YourModel struct {
ID        int    `json:"id" gorm:"primaryKey" runner:"name:ID" permission:"read"` //这个字段只读
Name      string `json:"name" runner:"name:名称" widget:"type:input" search:"like" validate:"required"`
Status    string `json:"status" runner:"name:状态" widget:"type:select;options:启用,禁用" search:"in"`
CreatedAt int64  `json:"created_at" gorm:"autoCreateTime:milli" runner:"name:创建时间" permission:"read"`
}
```

#### 2. **UI渲染**

- **widget标签**：选择最合适的UI组件进行渲染
- **data标签**：自动设置默认值、示例值、动态值，格式化的值，
- **runner标签**：自动生成中文标签、单位、显示名称

#### 3. **数据验证引擎**

- **validate标签**：自动生成前端和后端验证规则
- **search标签**：自动生成搜索和过滤功能
- **permission标签**：自动控制字段在不同场景的显示权限

#### 4. **数据库操作自动化**

- **gorm标签**：自动生成数据库表结构
- **CreateTables**：服务启动时自动建表
- **AutoCrudTable**：自动生成增删改查操作

#### 5. **回调函数集成**

- **OnInputFuzzy**：自动集成模糊搜索和聚合计算
- **OnInputValidate**：自动集成实时字段验证
- **OnTableAddRows**：table函数新增记录回调
- **OnTableUpdateRows** table函数更新记录回调
- **OnTableDeleteRows** table函数删除记录回调

### 🔄 标签系统的工作流程

```
结构体定义 → 标签解析 → 代码生成 → 运行时执行
     ↓           ↓         ↓         ↓
  业务模型   配置信息   前端界面   完整应用
  数据库表   验证规则   API接口   业务逻辑
```

### 🌟 标签系统的优势

| 传统开发方式     | function-go标签方式 |
| ------------------ | --------------------- |
| 手动编写CRUD代码 | 自动生成CRUD代码    |
| 手动编写验证逻辑 | 标签声明验证规则    |
| 手动设计UI界面   | 自动渲染UI界面      |
| 手动管理数据库   | 自动管理数据库      |
| 代码量大、易出错 | 代码简洁、零错误    |

通过标签系统，开发者只需要关注**业务逻辑**，框架自动处理所有**技术细节**，真正实现了"一个文件 = 一个完整Web应用"的愿景。

### 标签顺序建议

```
json → gorm → runner → widget → data → search → permission → validate
```

### 核心标签说明

#### runner标签 - 业务逻辑配置

| 属性     | 格式 | 示例 | 说明                         |
| ---------- | ------ | ------ | ------------------------------ |
| 字段名称 | `name:显示名称`     | `runner:"name:用户名"`     | 设置字段在前端的显示名称     |
| 字段单位 | `desc:字段介绍`     | `runner:"name:年龄;desc:年龄0-100"`     | 设置字段的详细介绍（非必要） |

#### widget标签 - UI组件配置

| 属性     | 格式 | 示例 | 说明           |
| ---------- | ------ | ------ | ---------------- |
| 组件类型 | `type:组件类型`     | `widget:"type:input"`     | 设置UI组件类型 |

#### data标签 - 数据和值配置

| 功能       | 格式 | 示例 | 说明               |
| ------------ | ------ | ------ | -------------------- |
| 默认值     | `default_value:值`     | `data:"default_value:默认值"`     | 设置字段默认值     |
| 示例值     | `example:示例值`     | `data:"example:示例文本"`     | 设置示例值         |
| 动态默认值 | `default_value:$变量`     | `data:"default_value:$now"`     | 使用变量作为默认值 |
| 格式化     | `format:格式化类型`     | `format:markdown`     | `设置格式化类型，csv或者markdown`                   |

#### validate标签 - 验证规则

| 规则     | 格式 | 示例 | 说明         |
| ---------- | ------ | ------ | -------------- |
| 必填验证 | `required`     | `validate:"required"`     | 字段必填     |
| 长度验证 | `min=值,max=值`     | `validate:"min=2,max=50"`     | 长度范围验证 |
| 数值验证 | `min=值,max=值`     | `validate:"min=1,max=120"`     | 数值范围验证 |
| 枚举验证 | `oneof=值1 值2`     | `validate:"oneof=男 女"`     | 枚举值验证   |

#### search标签 - 搜索配置（仅table函数）

| 搜索类型 | 格式 | 示例 | 说明                       |
| ---------- | ------ | ------ | ---------------------------- |
| 模糊搜索 | `like`     | `search:"like"`     | 启用模糊搜索               |
| 精确搜索 | `eq`     | `search:"eq"`     | 启用精确搜索               |
| 区间搜索 | `gte,lte`     | `search:"gte,lte"`     | 启用大于等于、小于等于搜索 |
| 多选搜索 | `in`     | `search:"in"`     | 启用多选搜索               |

#### permission标签 - 权限控制（仅table函数）

| 权限类型 | 格式 | 示例   | 说明                   |
| ---------- | ------ | -------- | ------------------------ |
| 仅可读   | `read`     | `permission:"read"`       | 仅列表显示，不能编辑   |
| 仅可创建 | `create`     | `permission:"create"`       | 仅新增显示，列表不显示 |
| 仅可更新 | `update`     | `permission:"update"`       | 仅编辑显示，列表不显示 |
| 全权限   | 不写 | 无标签 | 列表、新增、编辑都显示 |

## 🧩 组件系统

### 基础输入组件

#### input组件 - 文本输入

| 类型     | 配置 | 示例 | 说明           |
| ---------- | ------ | ------ | ---------------- |
| 单行文本 | `type:input`     | `widget:"type:input"`     | 基础文本输入框 |
| 多行文本 | `mode:text_area`     | `widget:"type:input;mode:text_area"`     | 多行文本区域   |
| 密码输入 | `mode:password`     | `widget:"type:input;mode:password"`     | 密码输入框     |

#### number组件 - 数字输入

| 类型     | 配置 | 示例 | 说明         |
| ---------- | ------ | ------ | -------------- |
| 整数输入 | `type:number`     | `widget:"type:number;min:1;max:120;unit:岁"`     | 整数输入框   |
| 小数输入 | `precision:小数位`     | `widget:"type:number;min:0;precision:2;prefix:￥"`     | 小数输入框   |
| 百分比   | `suffix:%`     | `widget:"type:number;min:0;max:100;precision:1;suffix:%"`     | 百分比输入框 |

#### select组件 - 下拉选择

| 类型     | 配置 | 示例 | 说明       |
| ---------- | ------ | ------ | ------------ |
| 单选下拉 | `type:select`     | `widget:"type:select;options:男,女"`     | 单选下拉框 |
| 多选下拉 | `multiple:true`     | `widget:"type:select;options:技术,产品,设计;multiple:true"`     | 多选下拉框 |

#### datetime组件 - 日期时间

| 类型     | 配置 | 示例 | 说明           |
| ---------- | ------ | ------ | ---------------- |
| 日期选择 | `kind:date`     | `widget:"type:datetime;kind:date;format:yyyy-MM-dd"`     | 日期选择器     |
| 时间选择 | `kind:time`     | `widget:"type:datetime;kind:time;format:HH:mm"`     | 时间选择器     |
| 日期时间 | `kind:datetime`     | `widget:"type:datetime;kind:datetime"`     | 日期时间选择器 |
| 日期范围 | `kind:daterange`     | `widget:"type:datetime;kind:daterange;format:yyyy-MM-dd"`     | 日期范围选择器 |

### 高级组件

#### multiselect组件 - 多选组件

| 配置         | 示例 | 说明               |
| -------------- | ------ | -------------------- |
| 静态多选     | `widget:"type:multiselect;options:紧急,重要,API,UI"`     | 固定选项多选       |
| 可创建新选项 | `widget:"type:multiselect;options:Java,Python,Go;allow_create:true"`     | 支持自定义创建选项 |

#### color组件 - 颜色选择器

| 格式     | 配置 | 示例 | 说明         |
| ---------- | ------ | ------ | -------------- |
| Hex格式  | `format:hex`     | `widget:"type:color;format:hex;show_alpha:false"`     | 6位hex颜色   |
| RGBA格式 | `format:rgba`     | `widget:"type:color;format:rgba;show_alpha:true"`     | RGBA颜色格式 |
| HSL格式  | `format:hsl`     | `widget:"type:color;format:hsl;show_alpha:false"`     | HSL颜色格式  |

#### file_upload组件 - 文件上传

| 配置       | 示例 | 说明       |
| ------------ | ------ | ------------ |
| 单文件上传 | `widget:"type:file_upload;accept:.jpg,.png;max_size:5MB"`     | 单文件上传 |
| 多文件上传 | `widget:"type:file_upload;accept:.pdf,.doc;multiple:true;max_size:10MB"`     | 多文件上传 |

#### list组件 - 列表输入

| 类型     | 示例 | 说明             |
| ---------- | ------ | ------------------ |
| 简单列表 | `widget:"type:list"`     | 字符串或数字列表 |
| 复杂列表 | `widget:"type:list"`     | 结构体列表       |

#### form组件 - 嵌套表单

| 示例 | 说明                               |
| ------ | ------------------------------------ |
| `widget:"type:form"`     | 嵌套表单结构，对应数据结构是结构体 |

### 其他组件

#### switch组件 - 开关

| 配置       | 示例 | 说明           |
| ------------ | ------ | ---------------- |
| 基础开关   | `widget:"type:switch"`     | 布尔值开关     |
| 自定义标签 | `widget:"type:switch;true_label:启用;false_label:禁用"`     | 自定义开关标签 |

#### radio组件 - 单选框

| 配置       | 示例 | 说明               |
| ------------ | ------ | -------------------- |
| 基础单选框 | `widget:"type:radio;options:男,女"`     | 单选按钮组         |
| 水平排列   | `widget:"type:radio;options:男,女;direction:horizontal"`     | 水平排列的单选按钮 |

#### checkbox组件 - 复选框

| 配置       | 示例 | 说明         |
| ------------ | ------ | -------------- |
| 基础复选框 | `widget:"type:checkbox;options:阅读,音乐,运动"`     | 多选复选框组 |

#### slider组件 - 滑块

| 配置     | 示例 | 说明             |
| ---------- | ------ | ------------------ |
| 数值滑块 | `widget:"type:slider;min:0;max:100;unit:%"`     | 数值范围滑块     |
| 评分滑块 | `widget:"type:slider;min:1;max:5;step:0.5;unit:分"`     | 带步进的评分滑块 |

## 📝 使用示例

### 基础字段配置

## Form函数模型示例 - 大模型代码生成模板

#### 用户注册模型

```go
// 请求结构体 - 用户输入
type UserRegisterReq struct {
    // 基础信息
    Username string `json:"username" runner:"name:用户名" widget:"type:input" data:"example:john_doe" validate:"required,min=3,max=20"`
    Password string `json:"password" runner:"name:密码" widget:"type:input;mode:password" data:"example:123456" validate:"required,min=6,max=20"`
    Email    string `json:"email" runner:"name:邮箱" widget:"type:input" data:"example:john@example.com" validate:"required,email"`
    
    // 个人信息
    RealName string `json:"real_name" runner:"name:真实姓名" widget:"type:input" data:"example:张三" validate:"required,min=2,max=20"`
    Age      int    `json:"age" runner:"name:年龄" widget:"type:number;min:18;max:65;unit:岁" data:"example:25" validate:"required,min=18,max=65"`
    Gender   string `json:"gender" runner:"name:性别" widget:"type:radio;options:男,女;direction:horizontal" data:"example:男" validate:"required,oneof=男 女"`
    
    // 工作信息
    Department string `json:"department" runner:"name:部门" widget:"type:select;options:技术部,产品部,设计部,运营部" data:"default_value:技术部" validate:"required"`
    Position   string `json:"position" runner:"name:职位" widget:"type:input" data:"example:软件工程师" validate:"required"`
    Salary     int    `json:"salary" runner:"name:期望薪资" widget:"type:number;min:3000;max:50000;unit:元" data:"example:15000" validate:"required,min=3000,max=50000"`
    
    // 技能标签
    Skills []string `json:"skills" runner:"name:技能标签" widget:"type:multiselect;options:Java,Python,Go,JavaScript,React,Vue" data:"example:Java,Go" validate:"required,min=1"`
    
    // 附件上传
    Resume *files.Files `json:"resume" runner:"name:简历" widget:"type:file_upload;accept:.pdf,.doc,.docx;max_size:10MB" validate:"required"`
    Avatar *files.Files `json:"avatar" runner:"name:头像" widget:"type:file_upload;accept:.jpg,.png,.gif;max_size:5MB"`
    
    // 其他信息
    Bio       string `json:"bio" runner:"name:个人简介" widget:"type:input;mode:text_area" data:"example:热爱编程，有3年开发经验"`
    AgreeTerms bool  `json:"agree_terms" runner:"name:同意条款" widget:"type:switch;true_label:同意;false_label:不同意" data:"example:true" validate:"required"`
}

// 响应结构体 - 处理结果
type UserRegisterResp struct {
    // 处理结果
    Success   bool   `json:"success" runner:"name:是否成功" widget:"type:switch;true_label:成功;false_label:失败"`
    Message   string `json:"message" runner:"name:处理结果" widget:"type:input;mode:text_area"`
    
    // 用户信息
    UserID    int    `json:"user_id" runner:"name:用户ID" widget:"type:number"`
    Username  string `json:"username" runner:"name:用户名" widget:"type:input"`
    
    // 系统信息
    CreatedAt int64  `json:"created_at" runner:"name:注册时间" widget:"type:datetime;kind:datetime"`
    Token     string `json:"token" runner:"name:访问令牌" widget:"type:input;mode:password"`
}
```

#### 采购申请模型

```go
// 请求结构体 - 采购申请
type PurchaseReq struct {
    // 基础信息
    Title       string `json:"title" runner:"name:采购标题" widget:"type:input" data:"example:办公用品采购" validate:"required,min=5,max=100"`
    Department  string `json:"department" runner:"name:申请部门" widget:"type:select;options:技术部,产品部,设计部,运营部" validate:"required"`
    Priority    string `json:"priority" runner:"name:优先级" widget:"type:select;options:低,中,高,紧急" data:"default_value:中" validate:"required"`
    
    // 供应商信息
    SupplierID int `json:"supplier_id" runner:"name:供应商" widget:"type:select" validate:"required"`
    
    // 采购商品列表
    Items []PurchaseItem `json:"items" runner:"name:采购商品" widget:"type:list" validate:"required,min=1"`
    
    // 其他信息
    ExpectedDate int64  `json:"expected_date" runner:"name:期望到货日期" widget:"type:datetime;kind:date;format:yyyy-MM-dd" validate:"required"`
    Remarks      string `json:"remarks" runner:"name:备注说明" widget:"type:input;mode:text_area"`
}

// 采购商品项
type PurchaseItem struct {
    ProductID int     `json:"product_id" runner:"name:商品" widget:"type:select" validate:"required"`
    Quantity  int     `json:"quantity" runner:"name:数量" widget:"type:number;min:1" data:"default_value:1" validate:"required,min=1"`
    UnitPrice float64 `json:"unit_price" runner:"name:单价" widget:"type:number;min:0;precision:2;prefix:￥" validate:"required,min=0"`
    Remarks   string  `json:"remarks" runner:"name:备注" widget:"type:input"`
}

// 响应结构体 - 采购结果
type PurchaseResp struct {
    // 处理结果
    Success      bool   `json:"success" runner:"name:是否成功" widget:"type:switch;true_label:成功;false_label:失败"`
    Message      string `json:"message" runner:"name:处理结果" widget:"type:input;mode:text_area"`
    
    // 采购信息
    PurchaseID   int     `json:"purchase_id" runner:"name:采购单号" widget:"type:number"`
    TotalAmount  float64 `json:"total_amount" runner:"name:总金额" widget:"type:number;precision:2;prefix:￥"`
    TotalItems   int     `json:"total_items" runner:"name:商品种类" widget:"type:number"`
    
    // 状态信息
    Status       string `json:"status" runner:"name:采购状态" widget:"type:input"`
    CreatedAt    int64  `json:"created_at" runner:"name:创建时间" widget:"type:datetime;kind:datetime"`
}
```


### Form函数配置模板

```go
var YourFormOption = &runner.FormFunctionOptions{
    BaseConfig: runner.BaseConfig{
        ChineseName: "功能名称",
        ApiDesc:     "功能描述",
        Tags:        []string{"标签1", "标签2"},
        Request:     &YourReq{},
        Response:    &YourResp{},
        CreateTables: []interface{}{&YourModel{}}, // 如果需要建表
        Group:       YourGroup, // 如果使用函数组
    },
}
```


## Function-Go 命名规范

下面定义了 Function-Go 项目中的命名规范，确保代码的一致性和可维护性（注意：需要严格遵循命名规范）

## 1. 单文件单函数

### 文件命名

```go
// 文件：pdf_to_image.go（具体功能）
package pdf
```

### 路由命名

```go
// 路由和文件名称保持一致
RouterGroup+"/pdf_to_image"
```

### 结构体和函数命名

```go
// 结构体和函数用 PdfToImage 开头（具体功能）
type PdfToImageReq struct { ... }
type PdfToImageResp struct { ... }

func PdfToImage(ctx *runner.Context, req *PdfToImageReq, resp response.Response) error { ... }

var PdfToImageOption = &runner.FormFunctionOptions{ ... }
```


## 2. 单文件多函数（需要用函数组）

### 文件命名

```go
// 文件：pdf_tools.go
package pdf
```

### 函数组配置

```go
var PdfToolsGroup = &runner.FunctionGroup{
    CnName: "PDF工具集",
    EnName: "pdf_tools",  // 【框架规范】与文件名一致
}
```

### 路由命名

```go
// 路由用 pdf_tools_xxx 开头
RouterGroup+"/pdf_tools_convert"
RouterGroup+"/pdf_tools_merge"
RouterGroup+"/pdf_tools_split"
```


### 结构体和函数命名

```go
// 结构体和函数用 PdfTools 开头（文件名称的驼峰格式开头）
type PdfToolsConvertReq struct { ... }
type PdfToolsMergeReq struct { ... }
type PdfToolsSplitReq struct { ... }

func PdfToolsConvert(ctx *runner.Context, req *PdfToolsConvertReq, resp response.Response) error { ... }
func PdfToolsMerge(ctx *runner.Context, req *PdfToolsMergeReq, resp response.Response) error { ... }
func PdfToolsSplit(ctx *runner.Context, req *PdfToolsSplitReq, resp response.Response) error { ... }

var PdfToolsConvertOption = &runner.FormFunctionOptions{ ... }
var PdfToolsMergeOption = &runner.FormFunctionOptions{ ... }
var PdfToolsSplitOption = &runner.FormFunctionOptions{ ... }
```

## 总结

- **单文件单函数**：具体功能（如 `pdf_to_image.go`），命名用具体功能前缀（如 `PdfToImage`），路由用 `RouterGroup+"/pdf_to_image"`
- **单文件多函数**：抽象工具集（如 `pdf_tools.go`），命名用抽象前缀（如 `PdfTools`），路由用 `RouterGroup+"/pdf_tools_xxx"`，需要用函数组来归类这一系列相关函数
- **命名一致性**：确保文件名、包名、路由名、结构体名、函数名都遵循相同的命名模式
## 🎯 最佳实践

### 1. 标签配置原则

- **必填字段**：添加 `validate:"required"`
- **搜索字段**：根据类型选择合适的 `search` 标签
- **权限控制**：使用 `permission` 标签控制字段显示
- **默认值**：使用 `data:"default_value:值"` 设置默认值

### 2. 组件选择原则

- **文本输入**：使用 `input` 组件
- **数字输入**：使用 `number` 组件
- **选择输入**：使用 `select`、`radio`、`checkbox` 组件
- **日期时间**：使用 `datetime` 组件
- **文件处理**：使用 `file_upload` 组件
- **多选场景**：使用 `multiselect` 组件

### 3. 搜索配置原则

- **文本字段**：使用 `like` 模糊搜索
- **状态字段**：使用 `in` 多选搜索
- **数值字段**：使用 `gte,lte` 区间搜索
- **时间字段**：使用 `gte,lte` 时间范围搜索

### 4. 验证规则原则

- **必填验证**：必填字段必须添加 `required`
- **长度限制**：防止过长输入，使用 `min`、`max`
- **格式验证**：邮箱、URL等特殊格式使用相应验证规则
- **业务规则**：符合实际业务需求的验证规则