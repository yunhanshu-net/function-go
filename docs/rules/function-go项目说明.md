# function-go 框架介绍

## 项目背景

function-go是一个**Go语言云函数开发框架**，主要为了方便大模型依照这个框架生成可以直接构建的应用。

## 核心能力

### 1. 文件即是函数（应用）
在function-go中，一个Go文件就是一个完整的应用，包含：
- 数据模型定义
- API接口注册
- 业务逻辑处理
- 自动生成Web界面

### 2. 结构体参数即是组件
通过结构体标签定义UI组件，自动生成表单和表格：

```go
type UserReq struct {
    // 自动生成输入框组件
    Name string `json:"name" runner:"code:name;name:用户名" widget:"type:input;placeholder:请输入用户名" data:"type:string" validate:"required"`
    
    // 自动生成数字输入框组件
    Age int `json:"age" runner:"code:age;name:年龄" widget:"type:number;min:1;max:120;unit:岁" data:"type:number" validate:"required"`
    
    // 自动生成下拉选择框组件
    Role string `json:"role" runner:"code:role;name:用户角色" widget:"type:select;options:管理员,普通用户,访客;placeholder:请选择角色" data:"type:string" validate:"required"`
}
```

### 3. 自动建表能力（CreateTables）
只需将结构体放入 `CreateTables`，框架会自动在数据库中创建对应表结构，无需手写SQL。例如：

```go
var CallbackDemoOption = &runner.FunctionOptions{
    // ...
    CreateTables: []interface{}{&UserProfile{}, &Company{}},
    // ...
}
```
- 结构体字段支持GORM标签，主键、索引、注释等自动生成。
- 适合大模型和开发者一键生成业务表。

### 4. 数据库集成与操作
- ctx（*runner.Context）提供 `MustGetOrInitDB()` 方法，返回的就是**gorm的*gorm.DB对象**，可直接用GORM所有能力进行增删改查、事务、复杂查询等。
- 推荐配合CreateTables自动建表，开发体验极佳。

```go
// 获取gorm.DB对象
db := ctx.MustGetOrInitDB()

// 直接用GORM操作
var users []User
err := db.Where("age > ?", 18).Find(&users).Error
```

### 5. 提供丰富的组件类型
- **input**: 文本输入框（支持普通文本、密码、多行文本）
- **number**: 数字输入框（支持范围限制、单位、精度）
- **select**: 下拉选择框（支持单选、多选）
- **checkbox**: 复选框（支持多选）
- **radio**: 单选框
- **switch**: 开关组件
- **slider**: 滑块组件
- **color**: 颜色选择器
- **datetime**: 日期时间选择器
- **multiselect**: 多选组件（支持固定选项多选）
- **tag**: 标签组件
- **file_upload**: 文件上传组件
- **file_display**: 文件展示组件
- **list_input**: 列表输入组件
- **form**: 表单组件

### 6. 渲染类型
- **form**: 渲染出表单界面，用于数据输入
- **table**: 渲染出标准的Element Plus表格，支持搜索、排序、分页等

### 7. 函数级别回调和字段级别回调
- **函数级别回调**: 在API创建时、页面加载时执行
- **字段级别回调**: 在字段值变化时执行自定义逻辑

## 快速开始

### 1. form函数实战示例：字符串反转

#### 用户需求
我需要一个字符串反转的功能，可以输入字符串，然后给出反转的结果。

#### 示例代码

```go
package demo

import (
    "github.com/yunhanshu-net/function-go/pkg/dto/response"
    "github.com/yunhanshu-net/function-go/runner"
)

// 1. 请求结构体（决定表单字段和校验规则）
type ReverseReq struct {
    Text string `json:"text" runner:"code:text;name:待反转内容" widget:"type:input;placeholder:请输入内容" data:"type:string;default_value:hello" validate:"required"`
}

// 2. 响应结构体（决定结果展示）
type ReverseResp struct {
    Result string `json:"result" runner:"code:result;name:反转结果" widget:"type:input;mode:text_area" data:"type:string"`
}

// 3. 业务处理函数（只做核心业务，校验交给标签）
func ReverseHandler(ctx *runner.Context, req *ReverseReq, resp response.Response) error {
    // 只做业务处理，不做参数校验
    reversed := reverseString(req.Text)
    return resp.Form(&ReverseResp{
        Result: reversed,
    }).Build()
}

// 字符串反转工具函数
func reverseString(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

// 4. 注册API
var ReverseOption = &runner.FunctionOptions{
    Tags:        []string{"input", "字符串", "表单演示"},
    EnglishName: "reverse_demo",
    ChineseName: "字符串反转演示",
    ApiDesc:     "输入内容，返回其反转结果。内容不能为空。",
    Request:     &ReverseReq{},
    Response:    &ReverseResp{},
    RenderType:  response.RenderTypeForm,
}

func init() {
	runner.Post("/demo/reverse", ReverseHandler, ReverseOption)
}
```

#### 前端表现说明
- 该函数会在前端自动渲染为一个表单，包含“待反转内容”输入框。
- 用户输入内容后点击“运行”或“提交”，后端处理后，结果会以表单形式展示（如“反转结果”多行文本框）。
- 校验规则自动生效，输入为空时前端会自动提示。

#### 规范说明
- 参数校验全部靠 validate 标签，业务函数只做核心逻辑。
- 结构体标签写全，前端自动生成表单和校验。
- resp.Form(...) 返回结果，自动渲染为表单。
- 注册API用 FunctionOptions，推荐写清楚 Request/Response/RenderType。

---

### 2. table函数实战示例：图书管理系统

#### 用户需求
我需要一个图书管理系统，字段如下：  
书名：填写图书名称  
作者：填写作者名字  
价格：填写图书价格  
库存：填写库存数量  
状态：选择“热销”、“预售”或“已售罄”  
上架：开关控制图书是否上架销售

#### 示例代码

```go
package widgets

import (
	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/function-go/runner"
	"github.com/yunhanshu-net/pkg/query"
	"gorm.io/gorm"
)

// Book 1. 图书数据模型
type Book struct {
	ID        int            `json:"id" gorm:"primaryKey;autoIncrement" runner:"code:id;name:ID" permission:"read"`
	CreatedAt int64          `json:"created_at" gorm:"autoCreateTime:milli" runner:"code:created_at;name:创建时间" widget:"type:datetime;kind:datetime" data:"type:number;example:1705292200000" search:"gte,lte" permission:"read"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index" runner:"-"` // runner:"-" 表示不渲染该字段
	Title     string         `json:"title" runner:"code:title;name:书名" widget:"type:input;placeholder:请输入书名" data:"type:string" search:"like" validate:"required,min=2,max=50"`                          // 书名，必填，2-50字
	Author    string         `json:"author" runner:"code:author;name:作者" widget:"type:input;placeholder:请输入作者" data:"type:string" search:"like" validate:"required,min=2,max=20"`                        // 作者，必填，2-20字
	Price     float64        `json:"price" runner:"code:price;name:价格" widget:"type:number;min:0;precision:2;prefix:￥" data:"type:float" search:"eq,gte,lte,gt,lt" validate:"required,min=0"`            // 价格，必填，最小0元，支持区间/等值搜索
	Stock     int            `json:"stock" runner:"code:stock;name:库存" widget:"type:number;min:0;unit:本" data:"type:number" search:"eq,gte,lte,gt,lt" validate:"required,min=0"`                         // 库存，必填，最小0，支持区间/等值搜索
	Status    string         `json:"status" runner:"code:status;name:状态" widget:"type:select;options:热销,预售,已售罄;placeholder:请选择状态" data:"type:string" search:"eq,in" validate:"required,oneof=热销 预售 已售罄"` // 状态，必填，枚举，支持等值/多选筛选
	OnSale    bool           `json:"on_sale" runner:"code:on_sale;name:上架" widget:"type:switch;true_label:已上架;false_label:未上架" data:"type:boolean;default_value:true" search:"eq" validate:"required"`   // 上架，开关组件，支持等值筛选
}

func (b *Book) TableName() string {
	return "book"
}

// BookListReq 2. 请求结构体（自动包含分页和搜索）
type BookListReq struct {
	query.PageInfoReq `runner:"-"`
}

// BookList 3. table函数处理逻辑
func BookList(ctx *runner.Context, req *BookListReq, resp response.Response) error {
	db := ctx.MustGetOrInitDB() // gorm db
	var books []Book
	// 自动分页和搜索
	return resp.Table(&books).AutoPaginated(db, &Book{}, &req.PageInfoReq).Build()
}

// BookListOption 4. 注册API
var BookListOption = &runner.FunctionOptions{
	Tags:          []string{"图书管理", "表格演示"},
	EnglishName:   "book_list",
	ChineseName:   "图书列表",
	ApiDesc:       "展示图书列表，支持分页和搜索。",
	Request:       &BookListReq{},
	Response:      query.PaginatedTable[[]Book]{},
	RenderType:    response.RenderTypeTable, // table函数
	CreateTables:  []interface{}{&Book{}},   //这里注册的表会在程序构建时候自动创建
	AutoCrudTable: &Book{},                  //围绕着Book自动生成Book表的的新增，修改，删除操作接口，前端可以直接操作
}

func init() {
	runner.Get("/widgets/book_list", BookList, BookListOption)
}
```

#### 前端表现说明
- 该API会在前端自动渲染为一个表格，展示所有图书的ID、书名、作者、价格、库存、状态、上架状态、创建时间等字段。
- 表格上方会自动生成搜索栏：书名/作者为模糊搜索，价格/库存为区间搜索，状态为下拉多选，上架为开关筛选，创建时间为区间日期选择器。
- 支持分页、排序、翻页等交互，表格样式为标准Element Plus表格。
- 注册了 AutoCrudTable 后，前端自动支持新增、修改、删除等操作，无需手写接口。
- 软删除字段（DeletedAt）前端默认不展示，但后端自动支持“假删除”。
- 用户可直接在表格上方输入关键词、选择区间、切换开关等方式进行多条件筛选，体验与企业级管理后台一致。

#### 规范说明
- BaseModel 建议所有业务表都嵌入，统一ID、创建时间、软删除等通用字段。
- search 标签如 search:"like"、search:"eq" 等，前端会自动生成对应的搜索输入框或下拉筛选。
- query.PageInfoReq runner:"-"：table函数的固定写法，自动支持分页、排序、搜索等功能。
- TableName 必须实现，否则自动建表和查询会失败。
- resp.Table(...) 返回结果，自动渲染为表格。
- 注册API用 FunctionOptions，推荐写清楚 Request/Response/RenderType/CreateTables/AutoCrudTable。

## 一句话总结

**function-go = Go云函数框架 + 自动UI生成器 + 自动建表 + 数据库ORM**

让大模型能够快速生成完整的Web应用，一个Go文件就是一个完整的应用。

---

## 回调机制简介

> 所有回调相关配置均在 `FunctionOptions` 结构体中设置，通常在注册 API 时一并传入。例如：
>
> ```go
> var DemoOption = &runner.FunctionOptions{
>     // ...
>     OnPageLoad: ...,
>     OnInputFuzzyMap: ...,
> }
> runner.Post("/demo/path", DemoHandler, DemoOption)
> ```
>
> 这样可以将所有与该API相关的回调逻辑集中管理，便于维护和理解。

function-go 支持灵活的回调机制，极大提升了动态业务场景下的自动化和智能化能力。

### 1. OnPageLoad（函数级别回调）
- 用于页面加载时初始化参数，常见于需要动态获取数据作为初始值的场景。
- 回调函数签名：
  ```go
  OnPageLoad: func(ctx *runner.Context, resp response.Response) (initData *usercall.OnPageLoadResp, err error)
  ```
  - ctx：请求上下文，包含用户信息、请求参数等
  - resp：用于返回前端的初始化数据（如表单、提示等）
  - 返回值：
    - initData：结构体，通常包含 Request 字段（类型为你的请求结构体），用于给前端表单字段动态赋初值
    - err：错误信息，返回非nil时前端会提示
- 典型业务场景：
  - 用户进入页面时，自动查询数据库获取最近一次操作记录，作为表单初始值
  - 根据当前登录用户信息，动态填充部门、姓名等字段
  - 结合外部API结果，动态生成推荐参数
- 推荐用法：通过 `ctx.GetUserInfo()` 获取用户信息（返回 UserInfo 结构体），无需判断未登录报错，未登录时直接返回空初始化参数即可。
- 示例代码（参考 @callback_demo.go）：

```go
OnPageLoad: func(ctx *runner.Context, resp response.Response) (initData *usercall.OnPageLoadResp, err error) {
    // 1. 获取当前用户信息
    userInfo := ctx.GetUserInfo()
    // 2. 查询数据库，获取最近一条记录（如已登录）
    var lastProfile UserProfile
    db := ctx.MustGetOrInitDB()
    if userInfo.IsLoggedIn {
        db.Where("username = ?", userInfo.Username).Order("created_at desc").First(&lastProfile)
    }
    // 3. 构造初始化参数
    initReq := CallbackDemoReq{}
    if userInfo.IsLoggedIn {
        initReq.Username = userInfo.Username
        initReq.Department = lastProfile.Department
        // ...其他字段
    }
    // 4. 返回给前端，自动填充表单（未登录时返回空参数即可）
    return &usercall.OnPageLoadResp{Request: initReq}, nil
}
```
- 注意事项：
  - OnPageLoad 返回的 Request 字段会自动作为表单初始值，优先级高于 default_value
  - 可以结合 ctx.GetUserInfo() 获取用户、组织、权限等上下文信息，实现个性化初始化

### 2. OnInputFuzzyMap（字段级别回调）
- 用于输入框等组件的"动态联想/模糊搜索"场景。
- 只有在用户输入时才会触发，实时从数据库或其他数据源返回候选项。
- 回调函数签名：
  ```go
  OnInputFuzzyMap: map[string]runner.OnInputFuzzy{
      "字段名": func(ctx *runner.Context, req *usercall.OnInputFuzzyReq) (*usercall.OnInputFuzzyResp, error) {
          // ...
      },
  }
  ```
  - ctx：请求上下文，包含用户信息、请求参数等
  - req：包含当前输入内容（req.Value）、其他上下文参数
  - 返回值：
    - *usercall.OnInputFuzzyResp：结构体，Values 字段为候选项列表（每项为 InputFuzzyItem，仅包含 Value 字段，Value 既作为下拉展示文本也作为实际取值）
    - error：错误信息，返回非nil时前端会提示
- 典型业务场景：
  - 公司名称输入框，用户输入2个字后，自动联想并下拉展示相关公司
  - 用户搜索输入框，输入关键字后动态展示匹配的用户
  - 地点、标签等字段的智能补全
- 示例代码（参考 @callback_demo.go）：

```go
OnInputFuzzyMap: map[string]runner.OnInputFuzzy{
    "company_search": func(ctx *runner.Context, req *usercall.OnInputFuzzyReq) (*usercall.OnInputFuzzyResp, error) {
        // 公司名称模糊搜索
        db := ctx.MustGetOrInitDB()
        var companies []Company
        if req.Value != "" {
            db.Where("name LIKE ? OR industry LIKE ?", "%"+req.Value+"%", "%"+req.Value+"%").Limit(10).Find(&companies)
        } else {
            db.Limit(10).Find(&companies)
        }
        resp := &usercall.OnInputFuzzyResp{}
        for _, company := range companies {
            resp.Values = append(resp.Values, &usercall.InputFuzzyItem{
                Value: company.Name, // 只需赋值 Value 字段
            })
        }
        return resp, nil
    },
    "user_search": func(ctx *runner.Context, req *usercall.OnInputFuzzyReq) (*usercall.OnInputFuzzyResp, error) {
        // 用户名称模糊搜索
        db := ctx.MustGetOrInitDB()
        var users []UserProfile
        if req.Value != "" {
            db.Where("username LIKE ? OR email LIKE ? OR department LIKE ?", "%"+req.Value+"%", "%"+req.Value+"%", "%"+req.Value+"%").Limit(10).Find(&users)
        } else {
            db.Limit(10).Find(&users)
        }
        resp := &usercall.OnInputFuzzyResp{}
        for _, user := range users {
            resp.Values = append(resp.Values, &usercall.InputFuzzyItem{
                Value: user.Username, // 只需赋值 Value 字段
            })
        }
        return resp, nil
    },
}
```
- 注意事项：
  - 字段名需与请求结构体中的字段一致，如 `company_search`、`user_search` 等
  - 返回的 Values 会自动作为下拉候选项展示，Value 既是展示文本也是实际取值
  - 可结合 ctx.GetUserInfo() 实现个性化联想，如只展示当前用户有权限的数据

---

## 标签使用规范

### 1. 核心标签说明

#### json 标签
- **用途**：JSON序列化/反序列化
- **格式**：`json:"字段名"`

#### runner 标签
- **用途**：业务逻辑配置
- **格式**：`runner:"code:字段代码;name:显示名称;desc:描述"`

#### widget 标签
- **用途**：UI组件配置
- **格式**：`widget:"type:组件类型;参数1:值1;参数2:值2"`

#### data 标签
- **用途**：数据类型和值配置
- **格式**：`data:"type:数据类型;default_value:默认值;example:示例值"`

#### validate 标签
- **用途**：验证规则配置
- **格式**：`validate:"规则1,规则2,规则3"`

#### search 标签
- **用途**：搜索配置（仅table函数）
- **格式**：`search:"搜索类型1,搜索类型2"`

#### permission 标签
- **用途**：权限控制配置（仅table函数）
- **格式**：`permission:"权限类型"`

### 2. 权限控制原则（permission标签，仅table函数）

#### 权限类型说明
- **permission:"read"**：仅可读权限（仅列表显示）
- **permission:"create"**：仅可创建权限（仅新增显示）
- **permission:"update"**：仅可更新权限（仅编辑显示）
- **permission:"create,update"**：可创建可更新权限（新增编辑显示，列表不显示）
- **不写permission标签**：全权限（列表、新增、编辑都显示）

#### 使用场景
```go
type User struct {
    // 系统字段：只读权限
    ID        int    `permission:"read"`
    CreatedAt int64  `permission:"read"`
    
    // 敏感字段：仅创建权限（如密码）
    Password  string `permission:"create"`
    
    // 业务字段：全权限（不写permission标签）
    Name      string
    Email     string
    
    // 特殊字段：仅编辑权限
    LastLogin int64  `permission:"update"`
}
```

### 3. 搜索配置原则（search标签，仅table函数）

#### 搜索类型说明
- **文本字段**：使用 `like` 模糊搜索，`eq` 精确搜索，`not_like` 否定模糊搜索，`not_eq` 否定精确搜索
- **数值字段**：使用 `eq` 等值搜索，`gt,gte,lt,lte` 区间搜索，`not_eq` 否定搜索
- **枚举字段**：使用 `eq` 等值搜索，`in` 多选搜索，`not_eq` 否定等值搜索，`not_in` 否定多选搜索
- **时间字段**：使用 `eq` 等值搜索，`gt,gte,lt,lte` 时间范围搜索，`not_eq` 否定搜索

#### 使用示例
```go
type Product struct {
    // 文本搜索
    Name        string  `search:"like,eq"`
    Description string  `search:"like,not_like"`
    
    // 数值搜索
    Price       float64 `search:"eq,gt,gte,lt,lte"`
    Stock       int     `search:"eq,gt,gte,lt,lte"`
    
    // 枚举搜索
    Status      string  `search:"eq,in,not_eq"`
    Category    string  `search:"eq,in,not_in"`
    
    // 时间搜索
    CreatedAt   int64   `search:"eq,gt,gte,lt,lte"`
}
```

### 4. 验证规则原则（validate标签）

#### 常用验证规则
- **required**：必填验证
- **min/max**：数值范围验证，格式：`min=值, max=值`
- **len**：长度验证，格式：`len=值`
- **oneof**：枚举值验证，格式：`oneof=值1 值2 值3`
- **email**：邮箱格式验证
- **url**：URL格式验证

#### 使用示例
```go
type UserReq struct {
    // 必填验证
    Username string `validate:"required"`
    
    // 长度验证
    Password string `validate:"required,min=6,max=20"`
    
    // 数值范围验证
    Age      int    `validate:"required,min=1,max=120"`
    
    // 枚举验证
    Gender   string `validate:"required,oneof=男 女"`
    
    // 格式验证
    Email    string `validate:"required,email"`
}
```

---

> 本文仅介绍了最常用的两类回调，后续文档会对回调机制进行更详细的说明和高级用法示例。
