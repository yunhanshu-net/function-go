# 简单CRUD开发指南

> **基于function-go框架的简单Table函数开发指南**  
> 单表数据管理，自动生成完整的增删改查界面，适合工单管理、用户管理等简单管理系统

## 📋 场景定位

### 什么是简单CRUD？
**简单CRUD = 单表 + Table函数 + 自动界面生成**

- **单表操作**：一个数据模型对应一张数据库表
- **Table函数**：框架自动生成列表、新增、编辑、删除界面
- **零业务逻辑**：纯数据管理，无复杂计算和关联

### 典型应用场景
- **工单管理系统**：客户工单列表，状态跟踪
- **用户管理系统**：用户信息维护，权限管理
- **产品管理系统**：商品基础信息管理
- **分类管理系统**：标签分类、内容分类
- **配置管理系统**：系统参数、字典数据

### 技术特点
- **L1级别**：最简单，学习门槛最低
- **自动建表**：框架自动创建数据库表
- **自动界面**：自动生成完整的管理界面
- **即开即用**：写完代码立即可用，无需额外配置

## 🎯 核心功能

### 1. 数据模型设计

#### 必须包含的系统字段
每个简单CRUD模型都必须包含以下4个系统字段：

```go
type YourModel struct {
    // 【必须字段】系统自动管理的4个基础字段
    ID        int            `json:"id" gorm:"primaryKey;autoIncrement;column:id" runner:"name:ID" permission:"read"`
    CreatedAt int64          `json:"created_at" gorm:"autoCreateTime:milli;column:created_at" runner:"name:创建时间" widget:"type:datetime;kind:datetime" permission:"read"`
    UpdatedAt int64          `json:"updated_at" gorm:"autoUpdateTime:milli;column:updated_at" runner:"name:更新时间" widget:"type:datetime;kind:datetime" permission:"read"`
    DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;column:deleted_at" runner:"-"`
    
    // 【业务字段】在这里添加你的业务字段...
}
```

#### 常用业务字段类型

```go
// 文本字段
Name string `json:"name" gorm:"column:name" runner:"name:名称" widget:"type:input" search:"like" validate:"required"`

// 长文本字段  
Description string `json:"description" gorm:"column:description" runner:"name:描述" widget:"type:input;mode:text_area"`

// 下拉选择字段
Status string `json:"status" gorm:"column:status" runner:"name:状态" widget:"type:select;options:启用,禁用" data:"default_value:启用" validate:"required"`

// 用户选择字段
CreateBy string `json:"create_by" gorm:"column:create_by" runner:"name:创建人" widget:"type:user"`

// 文件上传字段
Attachments files.Files `json:"attachments" gorm:"type:json;column:attachments" runner:"name:附件" widget:"type:file"`
```

### 2. 字段标签详解

#### 必须的标签
- **json**: JSON序列化字段名
- **gorm**: 数据库字段配置
- **runner**: 显示名称和功能配置

#### 常用标签组合

| 功能 | 标签配置 | 说明 |
|------|----------|------|
| 基础文本 | `widget:"type:input"` | 单行文本输入框 |
| 长文本 | `widget:"type:input;mode:text_area"` | 多行文本输入框 |
| 下拉选择 | `widget:"type:select;options:选项1,选项2"` | 下拉选择框 |
| 用户选择 | `widget:"type:user"` | 用户选择组件 |
| 文件上传 | `widget:"type:file"` | 文件上传组件 |
| 搜索支持 | `search:"like"` | 模糊搜索 |
| 搜索支持 | `search:"in"` | 精确搜索 |
| 字段验证 | `validate:"required"` | 必填验证 |
| 默认值 | `data:"default_value:默认值"` | 设置默认值 |
| 只读字段 | `permission:"read"` | 只能查看不能编辑 |

### 3. Table函数实现

#### 请求结构体
```go
// 标准Table函数请求结构体
type YourModelListReq struct {
    // 框架自动处理分页、搜索、排序参数
    query.SearchFilterPageReq `runner:"-"`
}
```

#### 处理函数
```go
// Table函数处理逻辑
func YourModelList(ctx *runner.Context, req *YourModelListReq, resp response.Response) error {
    var list []YourModel
    
    // 获取数据库连接
    db := ctx.MustGetOrInitDB()
    
    // 框架自动处理分页、搜索、排序
    paginate, err := query.AutoPaginate(ctx, db, &YourModel{}, &list, &req.SearchFilterPageReq)
    if err != nil {
        return err
    }
    
    // 返回分页结果
    return resp.Table(paginate).Build()
}
```

### 4. 配置选项

#### TableFunctionOptions配置
```go
var YourModelListOption = &runner.TableFunctionOptions{
    BaseConfig: runner.BaseConfig{
        EnglishName:  "your_model_list",
        ChineseName:  "数据管理",
        ApiDesc:      "数据列表管理，支持增删改查",
        Tags:         []string{"数据管理", "CRUD"},
        Request:      &YourModelListReq{},
        Response:     query.PaginatedTable[[]YourModel]{},
        CreateTables: []interface{}{&YourModel{}}, // 自动建表
    },
    // 自动CRUD配置
    AutoCrudTable: &YourModel{},
}
```

### 5. 路由注册

#### 路由注册代码
```go
// 在包的init函数中注册路由
func init() {
    RouterGroup.Post("/your_model_list", YourModelList, YourModelListOption)
}
```

### 6. 回调函数（可选）

#### 新增回调
```go
OnTableAddRows: func(ctx *runner.Context, req *usercall.OnTableAddRowsReq) (*usercall.OnTableAddRowsResp, error) {
    // 自动填充创建用户
    if user := ctx.GetString("user"); user != "" {
        req.SetString("create_by", user)
    }
    
    return &usercall.OnTableAddRowsResp{}, nil
},
```

#### 更新回调
```go
OnTableUpdateRows: func(ctx *runner.Context, req *usercall.OnTableUpdateRowsReq) (*usercall.OnTableUpdateRowsResp, error) {
    // 记录更新日志
    ctx.Logger.Infof("用户 %s 更新了记录 ID: %v", ctx.GetString("user"), req.GetInt("id"))
    
    return &usercall.OnTableUpdateRowsResp{}, nil
},
```

#### 删除回调
```go
OnTableDeleteRows: func(ctx *runner.Context, req *usercall.OnTableDeleteRowsReq) (*usercall.OnTableDeleteRowsResp, error) {
    // 删除前检查
    for _, id := range req.GetIDs() {
        ctx.Logger.Infof("删除记录 ID: %d", id)
    }
    
    return &usercall.OnTableDeleteRowsResp{}, nil
},
```

## 🛠️ 完整示例

### 基于CrmTicket的完整实现

```go
package crm

import (
    "github.com/yunhanshu-net/function-go/pkg/dto/response"
    "github.com/yunhanshu-net/function-go/pkg/dto/usercall"
    "github.com/yunhanshu-net/function-go/runner"
    "github.com/yunhanshu-net/pkg/query"
    "github.com/yunhanshu-net/pkg/typex/files"
    "gorm.io/gorm"
)

// 工单数据模型
type CrmTicket struct {
    // 系统字段
    ID        int            `json:"id" gorm:"primaryKey;autoIncrement;column:id" runner:"name:工单ID" permission:"read"`
    CreatedAt int64          `json:"created_at" gorm:"autoCreateTime:milli;column:created_at" runner:"name:创建时间" widget:"type:datetime;kind:datetime" permission:"read"`
    UpdatedAt int64          `json:"updated_at" gorm:"autoUpdateTime:milli;column:updated_at" runner:"name:更新时间" widget:"type:datetime;kind:datetime" permission:"read"`
    DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;column:deleted_at" runner:"-"`
    
    // 业务字段
    Title       string      `json:"title" gorm:"column:title" runner:"name:工单标题" widget:"type:input" search:"like" validate:"required"`
    Description string      `json:"description" gorm:"column:description" runner:"name:问题描述" widget:"type:input;mode:text_area" validate:"required"`
    Priority    string      `json:"priority" gorm:"column:priority" runner:"name:优先级" widget:"type:select;options:低,中,高" data:"default_value:中" validate:"required"`
    Status      string      `json:"status" gorm:"column:status" runner:"name:工单状态" widget:"type:select;options:待处理,处理中,已完成,已关闭" data:"default_value:待处理" validate:"required"`
    Phone       string      `json:"phone" gorm:"column:phone" runner:"name:联系电话" widget:"type:input" validate:"required"`
    CreateBy    string      `json:"create_by" gorm:"column:create_by" runner:"name:创建用户" widget:"type:user"`
    Attachments files.Files `json:"attachments" gorm:"type:json;column:attachments" runner:"name:附件" widget:"type:file"`
}

// 请求结构体
type CrmTicketListReq struct {
    query.SearchFilterPageReq `runner:"-"`
}

// 处理函数
func CrmTicketList(ctx *runner.Context, req *CrmTicketListReq, resp response.Response) error {
    var list []CrmTicket
    db := ctx.MustGetOrInitDB()
    
    paginate, err := query.AutoPaginate(ctx, db, &CrmTicket{}, &list, &req.SearchFilterPageReq)
    if err != nil {
        return err
    }
    
    return resp.Table(paginate).Build()
}

// 配置选项
var CrmTicketListOption = &runner.TableFunctionOptions{
    BaseConfig: runner.BaseConfig{
        EnglishName:  "crm_ticket_list",
        ChineseName:  "工单管理",
        ApiDesc:      "客户工单管理系统，支持工单的增删改查操作",
        Tags:         []string{"CRM", "工单管理", "客户服务"},
        Request:      &CrmTicketListReq{},
        Response:     query.PaginatedTable[[]CrmTicket]{},
        CreateTables: []interface{}{&CrmTicket{}},
    },
    AutoCrudTable: &CrmTicket{},
    
    // 新增回调
    OnTableAddRows: func(ctx *runner.Context, req *usercall.OnTableAddRowsReq) (*usercall.OnTableAddRowsResp, error) {
        if user := ctx.GetString("user"); user != "" {
            req.SetString("create_by", user)
        }
        return &usercall.OnTableAddRowsResp{}, nil
    },
}

// 路由注册
func init() {
    RouterGroup.Post("/crm_ticket_list", CrmTicketList, CrmTicketListOption)
}
```

## 📚 最佳实践

### 1. 数据模型设计原则
- **必须包含4个系统字段**：ID、CreatedAt、UpdatedAt、DeletedAt
- **字段命名规范**：数据库字段使用下划线，JSON字段保持一致
- **合理的字段类型**：根据实际需求选择合适的数据类型

### 2. 标签配置原则
- **标签顺序**：json → gorm → runner → widget → search → validate → data → permission
- **必填标签**：json、gorm、runner是每个字段的必填标签
- **功能标签**：根据需求添加widget、search、validate等功能标签

### 3. 验证规则设计
- **基础验证**：required（必填）、min/max（长度限制）
- **格式验证**：email（邮箱）、phone（手机号）
- **枚举验证**：oneof（枚举值验证）

### 4. 权限控制设计
- **只读字段**：系统字段通常设置为 `permission:"read"`
- **用户字段**：创建人、更新人等用户相关字段

### 5. 搜索功能设计
- **模糊搜索**：文本字段使用 `search:"like"`
- **精确搜索**：状态、分类等字段使用 `search:"in"`

## ⚠️ 常见问题

### 1. 数据库连接问题
**问题**：`ctx.MustGetOrInitDB()` 返回错误  
**解决**：检查数据库配置，确保数据库服务正常运行

### 2. 字段验证失败
**问题**：前端显示验证错误  
**解决**：检查 `validate` 标签配置，确保验证规则正确

### 3. 搜索功能不工作
**问题**：搜索条件不生效  
**解决**：检查 `search` 标签配置，确保字段支持搜索

### 4. 自动建表失败
**问题**：程序启动时建表失败  
**解决**：检查 `CreateTables` 配置，确保模型定义正确

## 🚀 快速开始

### 步骤1：定义数据模型
```go
type YourModel struct {
    ID        int            `json:"id" gorm:"primaryKey;autoIncrement;column:id" runner:"name:ID" permission:"read"`
    CreatedAt int64          `json:"created_at" gorm:"autoCreateTime:milli;column:created_at" runner:"name:创建时间" widget:"type:datetime;kind:datetime" permission:"read"`
    UpdatedAt int64          `json:"updated_at" gorm:"autoUpdateTime:milli;column:updated_at" runner:"name:更新时间" widget:"type:datetime;kind:datetime" permission:"read"`
    DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;column:deleted_at" runner:"-"`
    
    Name   string `json:"name" gorm:"column:name" runner:"name:名称" widget:"type:input" validate:"required"`
    Status string `json:"status" gorm:"column:status" runner:"name:状态" widget:"type:select;options:启用,禁用" data:"default_value:启用"`
}
```

### 步骤2：定义请求结构体
```go
type YourModelListReq struct {
    query.SearchFilterPageReq `runner:"-"`
}
```

### 步骤3：实现处理函数
```go
func YourModelList(ctx *runner.Context, req *YourModelListReq, resp response.Response) error {
    var list []YourModel
    db := ctx.MustGetOrInitDB()
    
    paginate, err := query.AutoPaginate(ctx, db, &YourModel{}, &list, &req.SearchFilterPageReq)
    if err != nil {
        return err
    }
    
    return resp.Table(paginate).Build()
}
```

### 步骤4：配置选项
```go
var YourModelListOption = &runner.TableFunctionOptions{
    BaseConfig: runner.BaseConfig{
        EnglishName:  "your_model_list",
        ChineseName:  "数据管理",
        ApiDesc:      "数据管理系统",
        Tags:         []string{"数据管理"},
        Request:      &YourModelListReq{},
        Response:     query.PaginatedTable[[]YourModel]{},
        CreateTables: []interface{}{&YourModel{}},
    },
    AutoCrudTable: &YourModel{},
}
```

### 步骤5：注册路由
```go
func init() {
    RouterGroup.Post("/your_model_list", YourModelList, YourModelListOption)
}
```

### 步骤6：（可选）添加回调函数
```go
OnTableAddRows: func(ctx *runner.Context, req *usercall.OnTableAddRowsReq) (*usercall.OnTableAddRowsResp, error) {
    // 添加业务逻辑
    return &usercall.OnTableAddRowsResp{}, nil
},
```

完成以上步骤后，框架会自动生成完整的CRUD界面，包括列表查看、新增、编辑、删除功能。