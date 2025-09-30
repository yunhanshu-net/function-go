# Function-Go 新 Options 系统使用指南（扁平化设计）

## 🎯 概述

新的 Options 系统采用**扁平化设计**，提供更清晰、更类型安全的函数配置方式，同时**大幅提升大模型代码生成效率**。

## 🏗️ 核心特性

### 1. **扁平化设计** ⭐
- **减少嵌套层级**：从4层嵌套减少到2层
- **提升大模型理解**：更简单的结构，更容易生成准确代码
- **降低复杂度**：减少上下文窗口占用

### 2. **类型安全**
- 编译时就能发现配置错误
- 不同函数类型有不同的配置需求

### 3. **扩展性好**
- 轻松支持新的函数类型
- 回调可以跨函数类型复用

### 4. **函数组支持**
- 简化设计，只有name字段
- 预定义组，复用方便

### 5. **路由配置优化** ⭐
- **路由信息只在注册时指定**，避免 BaseConfig 和注册时的重复配置
- **消除冗余**：防止路由信息不一致的错误
- **简化配置**：开发者只需要在一个地方指定路由信息

## 📝 使用方法

### 1. **表单函数配置**

```go
// 创建表单选项（扁平化设计）
var formOption = &FormFunctionOptions{
    BaseConfig: BaseConfig{
        EnglishName:   "example_form",
        ChineseName:   "示例表单",
        ApiDesc:       "这是一个示例表单",
        Tags:          []string{"示例", "表单"},
        Group:         JsonConverterGroup,
        Request:       &ExampleReq{},
        Response:      &ExampleResp{},
        CreateTables:  []interface{}{&ExampleTable{}},
        Timeout:       30000,
        Async:         false,
        FunctionType:  FunctionTypeDynamic,
    },
    // 直接设置回调，无需嵌套
    OnPageLoad: func(ctx *Context, resp response.Response) (initData *usercall.OnPageLoadResp, err error) {
        return &usercall.OnPageLoadResp{
            Request: &ExampleReq{
                DefaultField: "默认值",
            },
        }, nil
    },
    OnInputFuzzyMap: map[string]OnInputFuzzy{
        "field_name": func(ctx *Context, req *usercall.OnInputFuzzyReq) (*usercall.OnInputFuzzyResp, error) {
            return &usercall.OnInputFuzzyResp{
                Values: []*usercall.InputFuzzyItem{
                    {Value: "选项1"},
                    {Value: "选项2"},
                },
            }, nil
        },
    },
    OnInputValidateMap: map[string]OnInputValidate{
        "field_name": func(ctx *Context, req *usercall.OnInputValidateReq) (*usercall.OnInputValidateResp, error) {
            return &usercall.OnInputValidateResp{
                ErrorMsg: "", // 空字符串表示验证通过
            }, nil
        },
    },
    OnDryRun: func(ctx *Context, req *usercall.OnDryRunReq) (*usercall.OnDryRunResp, error) {
        return &usercall.OnDryRunResp{
            Valid:   true,
            Message: "预览操作",
        }, nil
    },
}

// 注册路由
runner.Post("/api/demo/form/example", ExampleHandler, formOption)
```

### 2. **表格函数配置（扁平化）**

```go
// 创建表格选项（扁平化设计）
var tableOption = &TableFunctionOptions{
    BaseConfig: BaseConfig{
        EnglishName:   "example_table",
        ChineseName:   "示例表格",
        ApiDesc:       "这是一个示例表格",
        Tags:          []string{"示例", "表格"},
        Group:         ProductManagementGroup,
        Request:       &ExampleListReq{},
        Response:      &ExampleListResp{},
        CreateTables:  []interface{}{&ExampleTable{}},
        Timeout:       30000,
        Async:         false,
        FunctionType:  FunctionTypeDynamic,
    },
    // 直接设置回调，无需嵌套
    OnPageLoad: func(ctx *Context, resp response.Response) (initData *usercall.OnPageLoadResp, err error) {
        return &usercall.OnPageLoadResp{
            Request: &ExampleListReq{
                PageInfoReq: query.SearchFilterPageReq{PageSize: 10},
            },
        }, nil
    },
    OnInputFuzzyMap: map[string]OnInputFuzzy{
        "search_field": func(ctx *Context, req *usercall.OnInputFuzzyReq) (*usercall.OnInputFuzzyResp, error) {
            return &usercall.OnInputFuzzyResp{
                Values: []*usercall.InputFuzzyItem{
                    {Value: "搜索结果1"},
                    {Value: "搜索结果2"},
                },
            }, nil
        },
    },
    AutoCrudTable: &ExampleTable{},
    BeforeTableDeleteRows: func(ctx *Context, req *usercall.OnTableDeleteRowsReq) (*usercall.OnTableDeleteRowsResp, error) {
        return &usercall.OnTableDeleteRowsResp{}, nil
    },
    BeforeTableUpdateRows: func(ctx *Context, req *usercall.OnTableUpdateRowsReq) (*usercall.OnTableUpdateRowsResp, error) {
        return &usercall.OnTableUpdateRowsResp{}, nil
    },
}

// 注册路由
runner.Get("/api/demo/table/example", ExampleListHandler, tableOption)
```


## 🎯 函数类型说明

### 1. **静态函数** (FunctionTypeStatic)
- 无需参数，或者输入参数，但是结果永远恒定
- 示例：获取系统时间、获取版本号、获取配置信息

### 2. **动态函数** (FunctionTypeDynamic)
- 请求参数不可预测，响应参数不可预测
- 示例：查询用户信息、产品列表、订单管理

### 3. **纯函数** (FunctionTypePure)
- 输入输出可预测，如数学函数
- 示例：JSON转换、数学计算、格式转换

## 🏷️ 函数组使用

### 1. **预定义函数组**

```go
// 使用预定义函数组
var (
    JsonConverterGroup = &FunctionGroup{
        Name: "JSON转换",
    }
    
    ProductManagementGroup = &FunctionGroup{
        Name: "产品管理系统",
    }
)

// 在函数中使用
Group: JsonConverterGroup,
```

### 2. **自定义函数组**

```go
// 创建自定义函数组
var MyCustomGroup = &FunctionGroup{
    Name: "我的自定义组",
}

// 在函数中使用
Group: MyCustomGroup,
```

## 🔄 迁移策略

### 1. **第一阶段**：使用新系统
- 新函数使用 `FormFunctionOptions` 或 `TableFunctionOptions`
- 现有函数继续使用 `FunctionOptions`

### 2. **第二阶段**：逐步迁移
- 逐步将现有函数迁移到新系统
- 保持向后兼容

### 3. **第三阶段**：完全迁移
- 所有函数使用新系统
- 可选：废弃旧的 `FunctionOptions`

## ✅ 优势对比

| 特性 | 旧系统 | 新系统（扁平化） |
|------|--------|------------------|
| 嵌套层级 | 4层 | 2层 |
| 大模型理解 | ❌ 复杂 | ✅ 简单 |
| 代码生成 | ❌ 容易出错 | ✅ 准确率高 |
| 类型安全 | ❌ 运行时错误 | ✅ 编译时检查 |
| 回调分类 | ❌ 混在一起 | ✅ 分类清晰 |
| 函数组 | ❌ 不支持 | ✅ 支持 |
| 扩展性 | ❌ 难以扩展 | ✅ 易于扩展 |
| 代码复用 | ❌ 重复代码 | ✅ 高度复用 |
| 开发体验 | ❌ 复杂 | ✅ 直观 |
| 路由配置 | ❌ 重复配置 | ✅ 单一配置 |

## 🚀 大模型代码生成示例

### 1. **表单函数生成**
```go
// 大模型可以更容易地生成这样的代码
var calculatorOption = &FormFunctionOptions{
    BaseConfig: BaseConfig{
        EnglishName:   "calculator",
        ChineseName:   "计算器",
        ApiDesc:       "基础数学计算",
        Tags:          []string{"工具", "计算"},
        Group:         SystemToolsGroup,
        Request:       &CalculatorReq{},
        Response:      &CalculatorResp{},
        Timeout:       30000,
        FunctionType:  FunctionTypePure,
    },
    OnPageLoad: func(ctx *Context, resp response.Response) (initData *usercall.OnPageLoadResp, err error) {
        return &usercall.OnPageLoadResp{
            Request: &CalculatorReq{
                Expression: "1+1",
            },
        }, nil
    },
}
```

### 2. **表格函数生成**
```go
// 大模型可以更容易地生成这样的代码
var userListOption = &TableFunctionOptions{
    BaseConfig: BaseConfig{
        EnglishName:   "user_list",
        ChineseName:   "用户列表",
        ApiDesc:       "用户管理列表",
        Tags:          []string{"用户管理", "列表"},
        Group:         UserManagementGroup,
        Request:       &UserListReq{},
        Response:      &UserListResp{},
        CreateTables:  []interface{}{&User{}},
        Timeout:       30000,
        FunctionType:  FunctionTypeDynamic,
    },
    AutoCrudTable: &User{},
    OnInputFuzzyMap: map[string]OnInputFuzzy{
        "name_search": func(ctx *Context, req *usercall.OnInputFuzzyReq) (*usercall.OnInputFuzzyResp, error) {
            return &usercall.OnInputFuzzyResp{
                Values: []*usercall.InputFuzzyItem{
                    {Value: "张三"},
                    {Value: "李四"},
                },
            }, nil
        },
    },
}
```

## 🚀 下一步

1. **试用新系统**：在新函数中使用扁平化的 `FormFunctionOptions` 和 `TableFunctionOptions`
2. **反馈优化**：根据大模型生成效果反馈优化设计
3. **逐步迁移**：将现有函数逐步迁移到新系统
4. **完善文档**：根据实际使用情况完善文档

这个扁平化的 Options 系统专门为大模型代码生成优化，提供了更好的生成效率和准确性！ 