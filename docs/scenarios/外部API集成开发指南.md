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


### 示例4: Form函数-工作空间管理系统（调用外部API，后台管理类，替代postman类只能程序员用的情况，生成的是非技术人可以用的页面功能）







示例1: 权限管理系统
用户需求
我需要一个权限管理系统，用于批量删除权限规则。下面是接口请求示例：


1 权限删除
支持批量删除操作
API调用示例：
```shell

curl 'https://techlink.corp.kuaishou.com/api/authority/authorization/revoked' \
-H 'Content-Type: application/json' \
-H 'X-Halo-Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...' \
--data-raw '{
  "notify_type": 1,
  "rules": [
    {
      "ruleId": 12345,
      "productPid": "product_001",
      "subject": "user_456",
      "object": "resource_123", 
      "rolePid": "role_789",
      "roleProductPid": "product_001"
    }
  ]
}' \
--insecure
```

```json
{
  "code": 0,
  "data": {}
}
```

生成代码如下：

```go
// 文件名：auth_permission_manage.go
package email

import (
   "fmt"
   "time"

   "github.com/yunhanshu-net/function-go/pkg/dto/response"
   "github.com/yunhanshu-net/function-go/pkg/dto/usercall"
   "github.com/yunhanshu-net/function-go/runner"
   "github.com/yunhanshu-net/pkg/x/httpx"
)

// ==================== 配置管理 ====================
// 权限管理API配置：管理远程API的连接参数，包括域名和Token认证
type AuthPermissionManageAPIConfig struct {
   // API连接配置
   RevokedURL string `json:"revoked_url" runner:"name:权限删除API地址" widget:"type:input" data:"default_value:https://techlink.corp.kuaishou.com/api/authority/authorization/revoked"`
   Token      string `json:"token" runner:"name:认证Token" widget:"type:input;mode:password" data:"default_value:请替换成真实token"`

   // 连接配置
   TimeoutSeconds int `json:"timeout_seconds" runner:"name:超时时间(秒)" widget:"type:number;min:5;max:300" data:"default_value:30"`
}

// ==================== 数据结构 ====================
// 权限规则结构体
type AuthPermissionRule struct {
   // 规则ID
   RuleID int64 `json:"rule_id" runner:"name:规则ID" widget:"type:number" validate:"required"`
   // 产品唯一标识
   ProductPid string `json:"product_pid" runner:"name:产品标识" widget:"type:input" validate:"required"`
   // 主体唯一标识
   Subject string `json:"subject" runner:"name:主体标识" widget:"type:input" validate:"required"`
   // 资源对象唯一标识
   Object string `json:"object" runner:"name:对象标识" widget:"type:input" validate:"required"`
   // 角色唯一标识
   RoleRid string `json:"role_rid" runner:"name:角色标识" widget:"type:input" validate:"required"`
   // 角色所属产品唯一标识
   RoleProductPid string `json:"role_product_pid" runner:"name:角色产品标识" widget:"type:input"`

   // 处理状态（内部使用）
   Status string `json:"status" runner:"name:处理状态" widget:"type:input" permission:"read"`
}

// API响应结构体
type AuthPermissionAPIResponse[T any] struct {
   Code int `json:"code"`
   Data T   `json:"data"`
}

// ==================== 请求响应结构体 ====================
// 权限删除请求
type AuthPermissionRevokeReq struct {
   // 权限规则列表
   Rules []*AuthPermissionRule `json:"rules" runner:"name:权限规则列表" widget:"type:list" validate:"required,min=1"`
   // 通知类型
   NotifyType int `json:"notify_type" runner:"name:通知类型" widget:"type:select;options:1,2" data:"default_value:1" validate:"required,oneof=1 2"`
}

// 权限删除响应
type AuthPermissionRevokeResp struct {
   // 删除结果
   Message string `json:"message" runner:"name:删除结果" widget:"type:input;mode:text_area"`
   // 删除成功的规则
   SuccessRules []*AuthPermissionRule `json:"success_rules" runner:"name:删除成功规则" widget:"type:list"`
   // 删除失败的规则
   FailedRules []*AuthPermissionRule `json:"failed_rules" runner:"name:删除失败规则" widget:"type:list"`
   // 删除统计
   TotalCount   int `json:"total_count" runner:"name:总数量" widget:"type:number"`
   SuccessCount int `json:"success_count" runner:"name:成功数量" widget:"type:number"`
   FailedCount  int `json:"failed_count" runner:"name:失败数量" widget:"type:number"`
   // API状态
   APIStatus string `json:"api_status" runner:"name:API状态" widget:"type:input"`
}

// ==================== 工具函数 ====================
// Token验证函数：检查用户是否设置了真实的Token
func validateToken(config AuthPermissionManageAPIConfig) error {
   // 检查Token是否为空或默认值
   if config.Token == "" {
      return fmt.Errorf("Token不能为空，请在配置中设置真实的API Token")
   }

   // 检查是否为默认提示文本
   if config.Token == "请替换成真实token" {
      return fmt.Errorf("请将Token替换为真实的API Token，当前使用的是默认提示文本")
   }

   // 检查Token长度是否合理（至少8位）
   if len(config.Token) < 8 {
      return fmt.Errorf("Token长度过短，请检查是否设置了正确的API Token")
   }

   return nil
}

// ==================== 核心业务逻辑：权限删除 ====================
// 权限删除函数
func AuthPermissionRevoke(ctx *runner.Context, req *AuthPermissionRevokeReq, resp response.Response) error {
   // 【框架规范】配置获取：从上下文获取配置信息
   config := ctx.GetConfig().(AuthPermissionManageAPIConfig)

   // 【业务逻辑】Token验证：检查配置是否有效
   if err := validateToken(config); err != nil {
      return resp.Form(&AuthPermissionRevokeResp{
         Message:   err.Error(),
         APIStatus: "配置错误",
      }).Build()
   }

   startTime := time.Now()
   successRules := make([]*AuthPermissionRule, 0)
   failedRules := make([]*AuthPermissionRule, 0)

   // 删除每个规则
   for i, rule := range req.Rules {
      // 构建删除请求
      deleteReq := map[string]interface{}{
         "notify_type": req.NotifyType,
         "rules": []map[string]interface{}{
            {
               "ruleId":         rule.RuleID,
               "productPid":     rule.ProductPid,
               "subject":        rule.Subject,
               "object":         rule.Object,
               "rolePid":        rule.RoleRid,
               "roleProductPid": rule.RoleProductPid,
            },
         },
      }

      // 【框架规范】httpx库使用：框架提供的HTTP客户端库
      // 【Why】为什么用httpx：支持链式调用、直接绑定响应结构体、统一错误处理
      // 【What】httpx做什么：提供优雅的HTTP请求API，支持GET/POST/PUT/DELETE等方法
      // 【How】如何使用：链式调用Post().Header().Timeout().Body().Do(响应结构体)
      // 【业务逻辑】POST请求删除权限，包含完整的请求头设置和错误处理
      var apiResp AuthPermissionAPIResponse[struct{}]
      httpResult, err := httpx.Post(config.RevokedURL).
         Header("Content-Type", "application/json").
         Header("X-Halo-Token", config.Token).
         Timeout(time.Duration(config.TimeoutSeconds) * time.Second).
         Body(deleteReq).
         Do(&apiResp)

      if err != nil {
         ctx.Logger.Errorf("删除权限失败 RuleID=%v: %v", rule.RuleID, err)
         rule.Status = fmt.Sprintf("删除失败: %v", err)
         failedRules = append(failedRules, rule)
         continue
      }

      // 检查HTTP状态码
      if !httpResult.OK() {
         ctx.Logger.Errorf("删除权限HTTP错误 RuleID=%v: %d", rule.RuleID, httpResult.Code)
         rule.Status = fmt.Sprintf("HTTP错误: %d", httpResult.Code)
         failedRules = append(failedRules, rule)
         continue
      }

      // 检查API返回状态
      if apiResp.Code != 0 {
         ctx.Logger.Errorf("删除权限API错误 RuleID=%v: %d", rule.RuleID, apiResp.Code)
         rule.Status = fmt.Sprintf("API错误: %d", apiResp.Code)
         failedRules = append(failedRules, rule)
         continue
      }

      // 删除成功
      rule.Status = "删除成功"
      successRules = append(successRules, rule)
      ctx.Logger.Infof("删除成功第%d个: RuleID=%v", i+1, rule.RuleID)
   }

   // 构建响应
   result := &AuthPermissionRevokeResp{
      Message:      fmt.Sprintf("删除完成，耗时: %v", time.Now().Sub(startTime)),
      SuccessRules: successRules,
      FailedRules:  failedRules,
      TotalCount:   len(req.Rules),
      SuccessCount: len(successRules),
      FailedCount:  len(failedRules),
      APIStatus:    "删除完成",
   }

   return resp.Form(result).Build()
}

// ==================== 函数配置 ====================
// 权限删除函数配置
var AuthPermissionRevokeOption = &runner.FormFunctionOptions{
   BaseConfig: runner.BaseConfig{
      ChineseName: "权限删除",
      ApiDesc:     "批量删除权限规则，支持权限规则列表输入和删除结果统计。",
      Tags:        []string{"权限管理", "删除", "API调用"},
      Request:     &AuthPermissionRevokeReq{},
      Response:    &AuthPermissionRevokeResp{},
      AutoUpdateConfig: &runner.AutoUpdateConfig{
         ConfigStruct: AuthPermissionManageAPIConfig{
            RevokedURL:     "https://techlink.corp.kuaishou.com/api/authority/authorization/revoked",
            Token:          "请替换成真实token",
            TimeoutSeconds: 30,
         },
      },
   },

   // 【框架规范】DryRun回调：框架提供的API测试机制
   // 【Why】为什么需要DryRun：POST等写操作有风险，需要先测试连接和参数，避免误操作
   // 【What】DryRun做什么：模拟API调用，测试连接状态，验证参数格式，不执行实际业务
   // 【How】如何使用DryRun：前端自动提供DryRun按钮，点击后触发OnDryRun回调
   // 【触发时机】用户点击DryRun按钮时自动触发，无需用户输入DryRun参数
   // 【返回要求】必须返回Valid状态和测试案例，框架自动展示测试结果
   OnDryRun: func(ctx *runner.Context, req *usercall.OnDryRunReq) (*usercall.OnDryRunResp, error) {
      // 【框架规范】配置获取：从上下文获取配置信息
      config := ctx.GetConfig().(AuthPermissionManageAPIConfig)

      // 【业务逻辑】Token验证：检查配置是否有效
      if err := validateToken(config); err != nil {
         return &usercall.OnDryRunResp{
            Valid:   false,
            Message: err.Error(),
         }, nil
      }

      // 【框架规范】参数解码：从请求中解码用户输入参数
      var revokeReq AuthPermissionRevokeReq
      if err := req.DecodeBody(&revokeReq); err != nil {
         return &usercall.OnDryRunResp{
            Valid:   false,
            Message: fmt.Sprintf("参数解码失败: %v", err),
         }, nil
      }

      // 【业务逻辑】参数验证：检查业务参数是否有效
      if len(revokeReq.Rules) == 0 {
         return &usercall.OnDryRunResp{
            Valid:   false,
            Message: "请至少提供一个权限规则",
         }, nil
      }

      // 【业务逻辑】构建API请求：根据参数构建完整的请求体
      testRule := revokeReq.Rules[0]
      deleteReq := map[string]interface{}{
         "notify_type": revokeReq.NotifyType,
         "rules": []map[string]interface{}{
            {
               "ruleId":         testRule.RuleID,
               "productPid":     testRule.ProductPid,
               "subject":        testRule.Subject,
               "object":         testRule.Object,
               "rolePid":        testRule.RoleRid,
               "roleProductPid": testRule.RoleProductPid,
            },
         },
      }

      // 【框架规范】httpx DryRun：使用httpx库构建测试案例
      // 【Why】为什么用httpx：httpx提供ConnectivityCheck()和DryRun()方法，自动测试连接
      // 【What】httpx DryRun做什么：模拟HTTP请求，测试网络连接，验证请求格式
      // 【How】如何使用：链式调用Post().Header().Body().ConnectivityCheck().DryRun()
      // 【ConnectivityCheck底层实现】通过HEAD方法测试接口可用性和网络连通性
      // 【环境痛点解决】即使代码正确，环境问题（网络、防火墙、DNS等）也会导致API调用失败
      // 【用户价值】让用户提前发现环境问题，避免实际执行时的失败，提供保险机制
      dryRunCase := httpx.Post(config.RevokedURL).
         Header("Content-Type", "application/json").
         Header("X-Halo-Token", config.Token).
         Timeout(time.Duration(config.TimeoutSeconds) * time.Second).
         Body(deleteReq).
         ConnectivityCheck().
         DryRun()

      // 【框架规范】DryRun响应：返回测试结果和案例
      return &usercall.OnDryRunResp{
         Valid:   true,
         Message: fmt.Sprintf("预览权限删除，共 %d 条规则", len(revokeReq.Rules)),
         Cases:   []usercall.DryRunCase{dryRunCase},
      }, nil
   },
}

// ==================== 路由注册 ====================
func init() {
   // 权限删除
   runner.Post(RouterGroup+"/auth_permission_revoke", AuthPermissionRevoke, AuthPermissionRevokeOption)
}

//<总结>
//权限管理系统：专注于权限规则批量删除功能
//技术栈：AutoUpdateConfig配置管理、DryRun回调测试、httpx外部API调用、批量处理
//复杂度：S2级别，包含完整的业务逻辑处理，支持配置热更新和API测试
//设计模式：使用函数组管理相关功能，配置与业务逻辑分离，支持动态配置更新
//重要功能：权限删除支持批量删除，包含详细的错误处理和统计信息
//安全特性：使用Token认证，支持超时配置，提供DryRun测试功能避免误操作
//用户体验：提供详细的删除结果统计，支持配置管理界面，操作结果清晰展示
//</总结>

```





示例2: 工作空间管理
用户需求：

我有两个api，一个创建工作空间，一个获取工作空间列表，需要你帮我搞成对应的功能，方便用户使用，下面是接口请求示例

1. 创建工作空间

curl 'http://func-ai.geeleo.com/api/v1/runner'
-H 'Content-Type: application/json'
-H 'Token: 这里可以从配置里管理token参数'
--data-raw '{"title":"清华大学科研工作空间","name":"qinghuadaxue_keyan","description":"主要是为了科研"}'
--insecure

返回值：

{
"code": 0,
"msg": "成功",
"data": {
"id": 5
}
}

2. 获取工作空间列表

curl 'http://func-ai.geeleo.com/api/v1/runner?page_size=100'
-H 'Token: 这里可以从配置里管理token参数'
--insecure

返回值：
{
"code": 0,
"msg": "成功",
"data": {
"items": [
{
"id": 1,
"created_at": "2025-09-03 01:06:25",
"updated_at": "2025-09-03 21:18:41",
"created_by": "beiluo",
"title": "测试空间",
"name": "demo6",
"description": "测试",
"version": "v10",
"status": 1,   /1是已经启用，2停用
"runcher_id": null,
"is_public": false

},
{
"id": 2,
"created_at": "2025-09-03 02:03:21",
"updated_at": "2025-09-03 22:53:49",
"created_by": "beiluo",
"title": "test1",
"name": "test1",
"description": "测试",
"version": "v9",
"status": 1,
"runcher_id": null,
"is_public": false

}
],
"current_page": 0,
"total_count": 2,
"total_pages": 1,
"page_size": 100
}
}

```go
// 文件名：workspace_admin_manage.go


package workspace_admin

import (
	"fmt"
	"time"

	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/function-go/pkg/dto/usercall"
	"github.com/yunhanshu-net/function-go/runner"
	"github.com/yunhanshu-net/pkg/x/httpx" //调用http的需要用封装的库
)

// ==================== 函数组：工作空间管理 ====================
var WorkspaceAdminManageGroup = &runner.FunctionGroup{CnName: "工作空间管理", EnName: "workspace_admin_manage"}

// ==================== 配置管理：API连接配置 ====================

// <rag-api>
// 工作空间API配置：管理远程API的连接参数，包括域名和Token认证
// 【框架规范】AutoUpdateConfig配置管理：框架提供的配置热更新机制
// 【Why】为什么需要配置管理：外部API的域名、Token等参数经常变化，需要支持动态配置
// 【What】配置管理做什么：提供配置界面，支持配置热更新，自动持久化到本地文件
// 【How】如何使用配置：通过ctx.GetConfig()获取配置，框架自动管理配置生命周期
// 【业务逻辑】定义API连接参数：基础URL、认证Token、超时设置等
// 【数据来源】管理员通过配置界面设置，框架自动持久化到本地文件
// 【使用场景】外部API调用、Token认证、连接测试等远程服务管理场景
type WorkspaceAdminManageAPIConfig struct {
	// API连接配置
	BaseURL string `json:"base_url" runner:"name:API域名" widget:"type:input" data:"default_value:http://func-ai.geeleo.com/api/v1/runner"`
	Token   string `json:"token" runner:"name:认证Token" widget:"type:input;mode:password" data:"default_value:请替换成真实token"`

	// 连接配置
	TimeoutSeconds int `json:"timeout_seconds" runner:"name:超时时间(秒)" widget:"type:number;min:5;max:300" data:"default_value:30"`
}

// </rag-api>

// ==================== 外部API数据结构 ====================

// 外部API返回的工作空间信息结构体（用于解析外部API响应）
type ExternalWorkspaceInfo struct {
	ID          int    `json:"id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	CreatedBy   string `json:"created_by"`
	Title       string `json:"title"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Status      int    `json:"status"`
	RuncherID   *int   `json:"runcher_id"`
	IsPublic    bool   `json:"is_public"`
}

// 外部API响应结构体
type ExternalAPIResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Items       []ExternalWorkspaceInfo `json:"items"`
		CurrentPage int                     `json:"current_page"`
		TotalCount  int                     `json:"total_count"`
		TotalPages  int                     `json:"total_pages"`
		PageSize    int                     `json:"page_size"`
	} `json:"data"`
}

// 创建工作空间API响应
type ExternalCreateResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		ID int `json:"id"`
	} `json:"data"`
}

// ==================== 用户响应数据结构 ====================

// 返回给用户的工作空间信息结构体（用于form响应）
type WorkspaceInfo struct {
	ID          int    `json:"id" runner:"name:工作空间ID"`
	CreatedAt   string `json:"created_at" runner:"name:创建时间" widget:"type:datetime;kind:datetime"`
	UpdatedAt   string `json:"updated_at" runner:"name:更新时间" widget:"type:datetime;kind:datetime"`
	CreatedBy   string `json:"created_by" runner:"name:创建者" widget:"type:input"`
	Title       string `json:"title" runner:"name:标题" widget:"type:input"`
	Name        string `json:"name" runner:"name:名称" widget:"type:input"`
	Description string `json:"description" runner:"name:描述" widget:"type:input;mode:text_area"`
	Version     string `json:"version" runner:"name:版本" widget:"type:input"`
	Status      string `json:"status" runner:"name:状态" widget:"type:input"` // 转换为字符串显示
	RuncherID   *int   `json:"runcher_id" runner:"name:运行器ID" widget:"type:number"`
	IsPublic    bool   `json:"is_public" runner:"name:是否公开" widget:"type:switch;true_label:公开;false_label:私有"`
}

// ==================== 请求响应结构体 ====================

// <rag-api>
// 创建工作空间请求参数：包含工作空间的基本信息
// 【业务逻辑】系统自动调用外部API创建工作空间，支持Token认证
// 【使用建议】标题和名称建议使用有意义的标识，描述可用于详细说明工作空间用途
type WorkspaceAdminManageCreateReq struct {
	Title       string `json:"title" runner:"name:工作空间标题" widget:"type:input" data:"example:清华大学科研工作空间" validate:"required,min=2,max=100"`
	Name        string `json:"name" runner:"name:工作空间名称" widget:"type:input" data:"example:qinghuadaxue_keyan" validate:"required,min=2,max=50"`
	Description string `json:"description" runner:"name:工作空间描述" widget:"type:input;mode:text_area" data:"example:主要是为了科研"`
}

// </rag-api>

// <rag-api>
// 获取工作空间列表请求参数：支持分页查询
// 【业务逻辑】系统自动调用外部API获取工作空间列表，支持分页和筛选
// 【使用建议】page_size建议根据实际需要设置，避免一次性获取过多数据
type WorkspaceAdminManageListReq struct {
	PageSize int `json:"page_size" runner:"name:每页数量" widget:"type:number;min:1;max:1000;unit:个" data:"default_value:20;example:100" validate:"required,min=1,max=1000"`
}

// </rag-api>

// <rag-api>
// 创建工作空间响应结果：包含操作结果、工作空间ID、API状态等
// 【业务逻辑】根据API调用结果返回创建成功的工作空间ID和状态信息
// 【使用建议】通过响应信息了解创建结果，工作空间ID用于后续操作
type WorkspaceAdminManageCreateResp struct {
	Message     string `json:"message" runner:"name:操作结果" widget:"type:input;mode:text_area"`
	APIStatus   string `json:"api_status" runner:"name:API状态" widget:"type:input"`
	ConfigInfo  string `json:"config_info" runner:"name:配置信息" widget:"type:input"`
	WorkspaceID int    `json:"workspace_id" runner:"name:工作空间ID" widget:"type:number"`
}

// 获取工作空间列表响应结果：包含操作结果、工作空间列表、分页信息、API状态等
// 【业务逻辑】根据API调用结果返回工作空间列表和分页信息
// 【使用建议】通过响应信息了解查询结果，分页信息用于前端展示
type WorkspaceAdminManageListResp struct {
	Message     string          `json:"message" runner:"name:操作结果" widget:"type:input;mode:text_area"`
	APIStatus   string          `json:"api_status" runner:"name:API状态" widget:"type:input"`
	ConfigInfo  string          `json:"config_info" runner:"name:配置信息" widget:"type:input"`
	Workspaces  []WorkspaceInfo `json:"workspaces" runner:"name:工作空间列表" widget:"type:list"`
	TotalCount  int             `json:"total_count" runner:"name:总数量" widget:"type:number"`
	CurrentPage int             `json:"current_page" runner:"name:当前页" widget:"type:number"`
	TotalPages  int             `json:"total_pages" runner:"name:总页数" widget:"type:number"`
	PageSize    int             `json:"page_size" runner:"name:每页数量" widget:"type:number"`
}

// </rag-api>

// ==================== 工具函数：数据转换和API调用 ====================

// Token验证函数：检查用户是否设置了真实的Token
func validateToken(config WorkspaceAdminManageAPIConfig) error {
	// 检查Token是否为空或默认值
	if config.Token == "" {
		return fmt.Errorf("Token不能为空，请在配置中设置真实的API Token")
	}

	// 检查是否为默认提示文本
	if config.Token == "请替换成真实token" {
		return fmt.Errorf("请将Token替换为真实的API Token，当前使用的是默认提示文本")
	}

	// 检查Token长度是否合理（至少8位）
	if len(config.Token) < 8 {
		return fmt.Errorf("Token长度过短，请检查是否设置了正确的API Token")
	}

	return nil
}

// 数据转换函数：将外部API数据转换为用户友好的格式
func convertExternalToWorkspaceInfo(external ExternalWorkspaceInfo) WorkspaceInfo {
	// 状态转换：1=启用，2=停用
	statusText := "停用"
	if external.Status == 1 {
		statusText = "启用"
	}

	return WorkspaceInfo{
		ID:          external.ID,
		CreatedAt:   external.CreatedAt,
		UpdatedAt:   external.UpdatedAt,
		CreatedBy:   external.CreatedBy,
		Title:       external.Title,
		Name:        external.Name,
		Description: external.Description,
		Version:     external.Version,
		Status:      statusText,
		RuncherID:   external.RuncherID,
		IsPublic:    external.IsPublic,
	}
}

// ==================== 核心业务逻辑：工作空间管理 ====================

// ==================== 框架适配层：Form函数 ====================

// 创建工作空间
// 【框架说明】Form函数：处理工作空间创建请求，调用外部API创建工作空间
// 【业务逻辑】内联业务逻辑，方便大模型学习和理解完整的工作空间创建流程
func WorkspaceAdminManageCreate(ctx *runner.Context, req *WorkspaceAdminManageCreateReq, resp response.Response) error {
	// 【框架规范】配置获取：从上下文获取配置信息
	// 【Why】为什么需要配置：外部API的域名、Token等参数需要动态配置，不能硬编码
	// 【What】配置获取做什么：从框架配置管理中获取API连接参数
	// 【How】如何使用：ctx.GetConfig().(配置结构体类型)，框架自动管理配置生命周期
	config := ctx.GetConfig().(WorkspaceAdminManageAPIConfig)

	// 【业务逻辑】Token验证：检查配置是否有效
	if err := validateToken(config); err != nil {
		return resp.Form(&WorkspaceAdminManageCreateResp{
			Message:    err.Error(),
			APIStatus:  "配置错误",
			ConfigInfo: fmt.Sprintf("API: %s, 超时: %ds", config.BaseURL, config.TimeoutSeconds),
		}).Build()
	}

	// 构建API URL
	apiURL := config.BaseURL

	// 构建请求体，也可以直接用req
	requestBody := map[string]interface{}{
		"title":       req.Title,
		"name":        req.Name,
		"description": req.Description,
	}

	// 【框架规范】httpx库使用：框架提供的HTTP客户端库
	// 【Why】为什么用httpx：支持链式调用、直接绑定响应结构体、统一错误处理
	// 【What】httpx做什么：提供优雅的HTTP请求API，支持GET/POST/PUT/DELETE等方法
	// 【How】如何使用：链式调用Post().Header().Timeout().Body().Do(响应结构体)
	// 【业务逻辑】POST请求创建工作空间，包含完整的请求头设置和错误处理
	var apiResp ExternalCreateResponse
	httpResult, err := httpx.Post(apiURL).
		Header("Content-Type", "application/json").
		Header("Token", config.Token).
		Timeout(time.Duration(config.TimeoutSeconds) * time.Second).
		Body(requestBody).
		Do(&apiResp)

	if err != nil {
		return resp.Form(&WorkspaceAdminManageCreateResp{
			Message:    fmt.Sprintf("创建工作空间失败: %v", err),
			APIStatus:  "连接失败",
			ConfigInfo: fmt.Sprintf("API: %s, 超时: %ds", config.BaseURL, config.TimeoutSeconds),
		}).Build()
	}

	// 检查HTTP状态码
	if !httpResult.OK() {
		return resp.Form(&WorkspaceAdminManageCreateResp{
			Message:    fmt.Sprintf("API返回错误状态码: %d, 响应: %s", httpResult.Code, httpResult.ResBodyString),
			APIStatus:  "HTTP错误",
			ConfigInfo: fmt.Sprintf("API: %s, 超时: %ds", config.BaseURL, config.TimeoutSeconds),
		}).Build()
	}

	// 检查API返回状态
	if apiResp.Code != 0 {
		return resp.Form(&WorkspaceAdminManageCreateResp{
			Message:    fmt.Sprintf("API返回错误: %s", apiResp.Msg),
			APIStatus:  "API业务错误",
			ConfigInfo: fmt.Sprintf("API: %s, 超时: %ds", config.BaseURL, config.TimeoutSeconds),
		}).Build()
	}

	// 构建成功响应
	result := &WorkspaceAdminManageCreateResp{
		Message:    fmt.Sprintf("工作空间创建成功！ID: %d", apiResp.Data.ID),
		APIStatus:  "调用成功",
		ConfigInfo: fmt.Sprintf("API: %s, 超时: %ds", config.BaseURL, config.TimeoutSeconds),
	}

	return resp.Form(result).Build()
}

// 获取工作空间列表
// 【框架说明】Form函数：处理工作空间列表查询请求，调用外部API获取列表数据
// 【业务逻辑】内联业务逻辑，方便大模型学习和理解完整的工作空间查询流程
func WorkspaceAdminManageList(ctx *runner.Context, req *WorkspaceAdminManageListReq, resp response.Response) error {
	// 【框架规范】配置获取：从上下文获取配置信息
	// 【Why】为什么需要配置：外部API的域名、Token等参数需要动态配置，不能硬编码
	// 【What】配置获取做什么：从框架配置管理中获取API连接参数
	// 【How】如何使用：ctx.GetConfig().(配置结构体类型)，框架自动管理配置生命周期
	config := ctx.GetConfig().(WorkspaceAdminManageAPIConfig)

	// 【业务逻辑】Token验证：检查配置是否有效
	if err := validateToken(config); err != nil {
		return resp.Form(&WorkspaceAdminManageListResp{
			Message:    err.Error(),
			APIStatus:  "配置错误",
			ConfigInfo: fmt.Sprintf("API: %s, 超时: %ds", config.BaseURL, config.TimeoutSeconds),
		}).Build()
	}

	// 构建API URL
	apiURL := fmt.Sprintf("%s?page_size=%d", config.BaseURL, req.PageSize)

	// 【框架规范】httpx库使用：框架提供的HTTP客户端库
	// 【Why】为什么用httpx：支持链式调用、直接绑定响应结构体、统一错误处理
	// 【What】httpx做什么：提供优雅的HTTP请求API，支持GET/POST/PUT/DELETE等方法
	// 【How】如何使用：链式调用Get().Header().Timeout().Do(响应结构体)
	// 【业务逻辑】GET请求获取工作空间列表，包含完整的请求头设置和错误处理
	var apiResp ExternalAPIResponse
	httpResult, err := httpx.Get(apiURL).
		Header("Content-Type", "application/json").
		Header("Token", config.Token).
		Timeout(time.Duration(config.TimeoutSeconds) * time.Second).
		Do(&apiResp)

	if err != nil {
		return resp.Form(&WorkspaceAdminManageListResp{
			Message:    fmt.Sprintf("获取工作空间列表失败: %v", err),
			APIStatus:  "连接失败",
			ConfigInfo: fmt.Sprintf("API: %s, 超时: %ds", config.BaseURL, config.TimeoutSeconds),
		}).Build()
	}

	// 检查HTTP状态码
	if !httpResult.OK() {
		return resp.Form(&WorkspaceAdminManageListResp{
			Message:    fmt.Sprintf("API返回错误状态码: %d, 响应: %s", httpResult.Code, httpResult.ResBodyString),
			APIStatus:  "HTTP错误",
			ConfigInfo: fmt.Sprintf("API: %s, 超时: %ds", config.BaseURL, config.TimeoutSeconds),
		}).Build()
	}

	// 检查API返回状态
	if apiResp.Code != 0 {
		return resp.Form(&WorkspaceAdminManageListResp{
			Message:    fmt.Sprintf("API返回错误: %s", apiResp.Msg),
			APIStatus:  "API业务错误",
			ConfigInfo: fmt.Sprintf("API: %s, 超时: %ds", config.BaseURL, config.TimeoutSeconds),
		}).Build()
	}

	// 转换数据格式
	workspaces := make([]WorkspaceInfo, 0, len(apiResp.Data.Items))
	for _, external := range apiResp.Data.Items {
		workspaces = append(workspaces, convertExternalToWorkspaceInfo(external))
	}

	// 构建成功响应
	result := &WorkspaceAdminManageListResp{
		Message:     fmt.Sprintf("成功获取 %d 个工作空间", len(workspaces)),
		Workspaces:  workspaces,
		TotalCount:  apiResp.Data.TotalCount,
		CurrentPage: apiResp.Data.CurrentPage,
		TotalPages:  apiResp.Data.TotalPages,
		PageSize:    apiResp.Data.PageSize,
		APIStatus:   "调用成功",
		ConfigInfo:  fmt.Sprintf("API: %s, 超时: %ds", config.BaseURL, config.TimeoutSeconds),
	}

	return resp.Form(result).Build()
}

// ==================== 配置和注册 ====================

// 创建工作空间配置
var WorkspaceAdminManageCreateOption = &runner.FormFunctionOptions{
	BaseConfig: runner.BaseConfig{
		ChineseName: "工作空间管理-创建",
		ApiDesc:     "创建工作空间，支持设置标题、名称、描述等基本信息。",
		Tags:        []string{"工作空间管理", "创建", "API调用"},
		Request:     &WorkspaceAdminManageCreateReq{},
		Response:    &WorkspaceAdminManageCreateResp{},
		Group:       WorkspaceAdminManageGroup,
		AutoUpdateConfig: &runner.AutoUpdateConfig{ //【框架规范】如果需要配置管理功能需要用这个，ConfigStruct是对应的配置，框架会自动热更新配置
			ConfigStruct: WorkspaceAdminManageAPIConfig{
				BaseURL:        "http://func-ai.geeleo.com/api/v1/runner",
				Token:          "请替换成真实token",
				TimeoutSeconds: 30,
			},
		},
	},
	// 【框架规范】DryRun回调：框架提供的API测试机制
	// 【Why】为什么需要DryRun：POST等写操作有风险，需要先测试连接和参数，避免误操作
	// 【What】DryRun做什么：模拟API调用，测试连接状态，验证参数格式，不执行实际业务
	// 【How】如何使用DryRun：前端自动提供DryRun按钮，点击后触发OnDryRun回调
	// 【触发时机】用户点击DryRun按钮时自动触发，无需用户输入DryRun参数
	// 【返回要求】必须返回Valid状态和测试案例，框架自动展示测试结果
	OnDryRun: func(ctx *runner.Context, req *usercall.OnDryRunReq) (*usercall.OnDryRunResp, error) {
		// 【框架规范】配置获取：从上下文获取配置信息
		config := ctx.GetConfig().(WorkspaceAdminManageAPIConfig)

		// 【业务逻辑】Token验证：检查配置是否有效
		if err := validateToken(config); err != nil {
			return &usercall.OnDryRunResp{
				Valid:   false,
				Message: err.Error(),
			}, nil
		}

		// 【框架规范】参数解码：从请求中解码用户输入参数
		var createReq WorkspaceAdminManageCreateReq
		if err := req.DecodeBody(&createReq); err != nil {
			return &usercall.OnDryRunResp{
				Valid:   false,
				Message: fmt.Sprintf("参数解码失败: %v", err),
			}, nil
		}

		// 【框架规范】httpx DryRun：使用httpx库构建测试案例
		// 【Why】为什么用httpx：httpx提供ConnectivityCheck()和DryRun()方法，自动测试连接
		// 【What】httpx DryRun做什么：模拟HTTP请求，测试网络连接，验证请求格式
		// 【How】如何使用：链式调用Post().Header().Body().ConnectivityCheck().DryRun()
		// 【ConnectivityCheck底层实现】通过HEAD方法测试接口可用性和网络连通性
		// 【环境痛点解决】即使代码正确，环境问题（网络、防火墙、DNS等）也会导致API调用失败
		// 【用户价值】让用户提前发现环境问题，避免实际执行时的失败，提供保险机制
		dryRunCase := httpx.Post(config.BaseURL).
			Header("Content-Type", "application/json").
			Header("Token", config.Token).
			Timeout(time.Duration(config.TimeoutSeconds) * time.Second).
			Body(createReq).
			ConnectivityCheck().
			DryRun()

		// 【框架规范】DryRun响应：返回测试结果和案例
		return &usercall.OnDryRunResp{
			Valid:   true,
			Message: fmt.Sprintf("预览创建工作空间：%s (%s)", createReq.Title, createReq.Name),
			Cases:   []usercall.DryRunCase{dryRunCase},
		}, nil
	},
}

// 获取工作空间列表配置
var WorkspaceAdminManageListOption = &runner.FormFunctionOptions{
	BaseConfig: runner.BaseConfig{
		ChineseName: "工作空间管理-列表",
		ApiDesc:     "获取工作空间列表，支持分页查询和详细信息展示。",
		Tags:        []string{"工作空间管理", "列表", "API调用"},
		Request:     &WorkspaceAdminManageListReq{},
		Response:    &WorkspaceAdminManageListResp{},
		Group:       WorkspaceAdminManageGroup,
		AutoUpdateConfig: &runner.AutoUpdateConfig{
			ConfigStruct: WorkspaceAdminManageAPIConfig{ //【框架规范】如果需要配置管理功能需要用这个，ConfigStruct是对应的配置，框架会自动热更新配置
				BaseURL:        "http://func-ai.geeleo.com/api/v1/runner",
				Token:          "请替换成真实token",
				TimeoutSeconds: 30,
			},
		},
	},
	// 【框架规范】DryRun回调：框架提供的API测试机制
	// 【Why】为什么需要DryRun：GET请求也需要测试连接，验证参数有效性，确保API可用
	// 【What】DryRun做什么：模拟API调用，测试连接状态，验证参数格式，不执行实际业务
	// 【How】如何使用DryRun：前端自动提供DryRun按钮，点击后触发OnDryRun回调
	// 【触发时机】用户点击DryRun按钮时自动触发，无需用户输入DryRun参数
	// 【返回要求】必须返回Valid状态和测试案例，框架自动展示测试结果
	OnDryRun: func(ctx *runner.Context, req *usercall.OnDryRunReq) (*usercall.OnDryRunResp, error) {
		// 【框架规范】配置获取：从上下文获取配置信息
		config := ctx.GetConfig().(WorkspaceAdminManageAPIConfig)

		// 【业务逻辑】Token验证：检查配置是否有效
		if err := validateToken(config); err != nil {
			return &usercall.OnDryRunResp{
				Valid:   false,
				Message: err.Error(),
			}, nil
		}

		// 【框架规范】参数解码：从请求中解码用户输入参数
		var listReq WorkspaceAdminManageListReq
		if err := req.DecodeBody(&listReq); err != nil {
			return &usercall.OnDryRunResp{
				Valid:   false,
				Message: fmt.Sprintf("参数解码失败: %v", err),
			}, nil
		}

		// 【业务逻辑】参数验证：检查业务参数是否有效
		if listReq.PageSize < 1 || listReq.PageSize > 1000 {
			return &usercall.OnDryRunResp{
				Valid:   false,
				Message: "每页数量必须在1-1000之间",
			}, nil
		}

		// 【业务逻辑】构建API URL：根据参数构建完整的请求URL
		apiURL := fmt.Sprintf("%s?page_size=%d", config.BaseURL, listReq.PageSize)

		// 【框架规范】httpx DryRun：使用httpx库构建测试案例
		// 【Why】为什么用httpx：httpx提供ConnectivityCheck()和DryRun()方法，自动测试连接
		// 【What】httpx DryRun做什么：模拟HTTP请求，测试网络连接，验证请求格式
		// 【How】如何使用：链式调用Get().Header().ConnectivityCheck().DryRun()
		// 【ConnectivityCheck底层实现】通过HEAD方法测试接口可用性和网络连通性
		// 【环境痛点解决】即使代码正确，环境问题（网络、防火墙、DNS等）也会导致API调用失败
		// 【用户价值】让用户提前发现环境问题，避免实际执行时的失败，提供保险机制
		dryRunCase := httpx.Get(apiURL).
			Header("Content-Type", "application/json").
			Header("Token", config.Token).
			Timeout(time.Duration(config.TimeoutSeconds) * time.Second).
			ConnectivityCheck().
			DryRun()

		// 【框架规范】DryRun响应：返回测试结果和案例
		return &usercall.OnDryRunResp{
			Valid:   true,
			Message: fmt.Sprintf("预览获取工作空间列表：每页 %d 个", listReq.PageSize),
			Cases:   []usercall.DryRunCase{dryRunCase},
		}, nil
	},
}

// ==================== 路由注册 ====================

func init() {
	// 工作空间管理操作
	runner.Post(RouterGroup+"/workspace_admin_manage_create", WorkspaceAdminManageCreate, WorkspaceAdminManageCreateOption)
	runner.Post(RouterGroup+"/workspace_admin_manage_list", WorkspaceAdminManageList, WorkspaceAdminManageListOption)
}

//<总结>
//这里是个S2级别的工作空间管理系统，包含工作空间创建、列表查询、配置管理等功能
//技术栈：AutoUpdateConfig配置管理、DryRun回调测试、httpx外部API调用、数据转换
//复杂度：S2级别，包含基础回调机制，简单的业务逻辑处理，完全独立无依赖
//设计模式：使用AutoUpdateConfig管理API配置，支持实时配置更新，符合框架最佳实践
//重要提醒：DryRun回调独立于正常处理逻辑，用于测试验证，前端自动提供DryRun按钮
//外部API：使用httpx库调用远程API，支持Token认证、超时设置和连通性检查
//功能说明：支持工作空间创建、列表查询、API配置管理、连接测试等完整工作空间管理功能
//</总结>


```
