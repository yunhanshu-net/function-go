# Function-Go Options 系统设计（基于现有代码）

## 🎯 设计目标

1. **向后兼容** - 保持现有 `FunctionOptions` 不变
2. **类型安全** - 编译时就能发现配置错误
3. **扩展性** - 轻松支持新的函数类型
4. **简洁性** - 直接赋值，简单明了
5. **清晰性** - 回调分类清晰，易于理解

## 🏗️ 现有代码分析

### 1. 现有的 FunctionOptions 结构

```go
// 现有的 FunctionOptions（保持不变）
type FunctionOptions struct {
    AutoUpdateConfig *AutoUpdateConfig                  `json:"auto_update_config"`
    Router           string                             `json:"router"`
    Method           string                             `json:"method"`
    ApiDesc          string                             `json:"api_desc"`
    IsPublicApi      bool                               `json:"is_public_api"`
    ChineseName      string                             `json:"chinese_name"`
    EnglishName      string                             `json:"english_name"`
    Classify         string                             `json:"classify"`
    Tags             []string                           `json:"tags"`
    Async            bool                               `json:"async"`
    FunctionType     FunctionType                       `json:"function_type"`
    Timeout          int                                `json:"timeout"`
    RenderType       string                             `json:"widget"`
    CreateTables     []interface{}                      `json:"create_tables"`
    UseTables        []interface{}                      `json:"use_tables"`
    OperateTables    map[interface{}][]OperateTableType `json:"-"`
    AutoRun          bool                               `json:"-"`
    Request          interface{}                        `json:"-"`
    Response         interface{}                        `json:"-"`
    AutoCrudTable    interface{}                        `json:"-"`
    
    // 回调函数
    OnPageLoad        OnPageLoad                        `json:"-"`
    OnApiCreated      OnApiCreated                      `json:"-"`
    OnApiUpdated      OnApiUpdated                      `json:"-"`
    BeforeApiDelete   BeforeApiDelete                   `json:"-"`
    AfterApiDeleted   AfterApiDeleted                   `json:"-"`
    BeforeRunnerClose BeforeRunnerClose                 `json:"-"`
    AfterRunnerClose  AfterRunnerClose                  `json:"-"`
    OnVersionChange   OnVersionChange                   `json:"-"`
    OnTableDeleteRows OnTableDeleteRows                 `json:"-"`
    OnTableUpdateRows OnTableUpdateRows                 `json:"-"`
    OnTableAddRows    OnTableAddRows                    `json:"-"`
    OnTableSearch     OnTableSearch                     `json:"-"`
    OnInputFuzzyMap   map[string]OnInputFuzzy          `json:"-"`
    OnInputValidateMap map[string]OnInputValidate       `json:"-"`
    OnDryRun          OnDryRun                          `json:"-"`
}
```

## 🏗️ 新设计架构

### 1. 函数组结构体

```go
// FunctionGroup 函数组配置（简化版）
type FunctionGroup struct {
    Name string `json:"name"` // 组名称，如 "JSON转换"
    // 后续按需要添加字段
    // Description string `json:"description"` // 组描述
    // Version     string `json:"version"`     // 版本
    // Author      string `json:"author"`      // 作者
    // Tags        []string `json:"tags"`      // 标签
    // IsAtomic    bool `json:"is_atomic"`    // 是否原子组
}
```

### 2. 新的选项接口

```go
// Option 选项接口
type Option interface {
    GetFunctionType() FunctionType
    GetRenderType() string
    GetBaseConfig() *BaseConfig
    Validate() error
}

// BaseConfig 基础配置（所有函数通用）
type BaseConfig struct {
    // 路由配置
    Router      string `json:"router" validate:"required"`
    Method      string `json:"method" validate:"required"`
    
    // 名称配置
    EnglishName string   `json:"english_name" validate:"required"`
    ChineseName string   `json:"chinese_name" validate:"required"`
    ApiDesc     string   `json:"api_desc"`
    Tags        []string `json:"tags"`
    Classify    string   `json:"classify"`
    
    // 函数组配置
    Group       *FunctionGroup `json:"group"`        // 函数组配置
    
    // 执行配置
    Async        bool         `json:"async"`
    FunctionType FunctionType `json:"function_type"` // 函数类型：static/dynamic/pure
    Timeout      int          `json:"timeout"`
    
    // 权限配置
    IsPublicApi bool `json:"is_public_api"`
    
    // 请求响应
    Request  interface{} `json:"-"`
    Response interface{} `json:"-"`
    
    // 数据库配置
    CreateTables  []interface{}                      `json:"create_tables"`
    UseTables     []interface{}                      `json:"use_tables"`
    OperateTables map[interface{}][]OperateTableType `json:"-"`
    
    // 自动更新配置
    AutoUpdateConfig *AutoUpdateConfig `json:"auto_update_config"`
    
    // 自动运行
    AutoRun bool `json:"-"`
}
```

### 3. 回调分类设计

```go
// FunctionLifecycleCallback 函数生命周期回调（所有函数通用）
type FunctionLifecycleCallback struct {
    // API生命周期回调
    OnApiCreated    OnApiCreated    `json:"-"` // API创建时
    OnApiUpdated    OnApiUpdated    `json:"-"` // API更新时
    BeforeApiDelete BeforeApiDelete `json:"-"` // API删除前
    AfterApiDeleted AfterApiDeleted `json:"-"` // API删除后
    
    // 运行器生命周期回调
    BeforeRunnerClose BeforeRunnerClose `json:"-"` // 运行器关闭前
    AfterRunnerClose  AfterRunnerClose  `json:"-"` // 运行器关闭后
    OnVersionChange   OnVersionChange   `json:"-"` // 版本变更时
}

// CommonCallback 通用回调（所有函数通用）
type CommonCallback struct {
    OnPageLoad OnPageLoad `json:"-"` // 页面加载时，优先级最高
}

// ComponentCallback 组件级回调（表单函数专用）
type ComponentCallback struct {
    // 字段级模糊搜索回调
    OnInputFuzzyMap map[string]OnInputFuzzy `json:"-"` // key是字段code
    
    // 字段级验证回调
    OnInputValidateMap map[string]OnInputValidate `json:"-"` // key是字段code
}

// FormSpecificCallback 表单专用回调（表单函数专用）
type FormSpecificCallback struct {
    OnDryRun OnDryRun `json:"-"` // 预览模式，用于危险操作预览（仅在form函数中使用）
}

// TableSpecificCallback 表格专用回调（表格函数专用）
type TableSpecificCallback struct {
    BeforeTableDeleteRows OnTableDeleteRows `json:"-"` // 删除行前
    BeforeTableUpdateRows OnTableUpdateRows `json:"-"` // 更新行前
    BeforeTableAddRows    OnTableAddRows    `json:"-"` // 添加行前
    BeforeTableSearch     OnTableSearch     `json:"-"` // 搜索前
}
```

### 4. 专用配置

```go
// FormConfig 表单专用配置
type FormConfig struct {
    // 组件级回调
    ComponentCallback `json:",inline"`
    
    // 表单专用回调
    FormSpecificCallback `json:",inline"`
}

// TableConfig 表格专用配置
type TableConfig struct {
    // 组件级回调（table函数也可能需要字段级回调，如搜索字段的模糊搜索）
    ComponentCallback `json:",inline"`
    
    // 表格专用回调
    TableSpecificCallback `json:",inline"`
    
    // 表格特有配置
    AutoCrudTable interface{} `json:"-"` // 自动CRUD表格
}
```

### 5. 完整选项结构

```go
// FormFunctionOptions 表单函数选项
type FormFunctionOptions struct {
    BaseConfig                `json:",inline"`
    FunctionLifecycleCallback `json:",inline"`
    CommonCallback            `json:",inline"`
    FormConfig                `json:",inline"`
}

// TableFunctionOptions 表格函数选项
type TableFunctionOptions struct {
    BaseConfig                `json:",inline"`
    FunctionLifecycleCallback `json:",inline"`
    CommonCallback            `json:",inline"`
    TableConfig               `json:",inline"`
}
```

### 6. 接口实现

```go
// FormFunctionOptions 实现
func (opt *FormFunctionOptions) GetFunctionType() FunctionType {
    return opt.FunctionType
}

func (opt *FormFunctionOptions) GetRenderType() string {
    return "form"
}

func (opt *FormFunctionOptions) GetBaseConfig() *BaseConfig {
    return &opt.BaseConfig
}

func (opt *FormFunctionOptions) Validate() error {
    if opt.Router == "" {
        return errors.New("router is required")
    }
    if opt.EnglishName == "" {
        return errors.New("english_name is required")
    }
    return nil
}

// TableFunctionOptions 实现
func (opt *TableFunctionOptions) GetFunctionType() FunctionType {
    return opt.FunctionType
}

func (opt *TableFunctionOptions) GetRenderType() string {
    return "table"
}

func (opt *TableFunctionOptions) GetBaseConfig() *BaseConfig {
    return &opt.BaseConfig
}

func (opt *TableFunctionOptions) Validate() error {
    if opt.Router == "" {
        return errors.New("router is required")
    }
    if opt.EnglishName == "" {
        return errors.New("english_name is required")
    }
    if opt.AutoCrudTable == nil {
        return errors.New("auto_crud_table is required for table functions")
    }
    return nil
}
```

### 7. 便捷构造函数

```go
// NewFormOptions 创建表单选项
func NewFormOptions() *FormFunctionOptions {
    return &FormFunctionOptions{
        BaseConfig: BaseConfig{
            Method:       "POST",
            FunctionType: FunctionTypeDynamic, // 默认动态函数
            Timeout:      30000,
        },
    }
}

// NewTableOptions 创建表格选项
func NewTableOptions() *TableFunctionOptions {
    return &TableFunctionOptions{
        BaseConfig: BaseConfig{
            Method:       "GET",
            FunctionType: FunctionTypeDynamic, // 默认动态函数
            Timeout:      30000,
        },
    }
}
```

### 8. 使用示例

```go
// 预定义的函数组
var (
    JsonConverterGroup = &FunctionGroup{
        Name: "JSON转换",
    }
    
    ProductManagementGroup = &FunctionGroup{
        Name: "产品管理系统",
    }
)

// JSON转换工具组示例
func JsonConverterExamples() {
    // JSON转CSV（纯函数）
    var json2csvOption = &FormFunctionOptions{
        BaseConfig: BaseConfig{
            Router:        "/api/demo/form/json2csv",
            Method:        "POST",
            EnglishName:   "json2csv",
            ChineseName:   "JSON转CSV",
            ApiDesc:       "将JSON数据转换为CSV格式",
            Tags:          []string{"JSON转换", "数据转换"},
            Group:         JsonConverterGroup,
            Request:       &Json2CsvReq{},
            Response:      &Json2CsvResp{},
            CreateTables:  []interface{}{&ConversionRecord{}},
            Timeout:       60000,
            Async:         false,
            FunctionType:  FunctionTypePure, // 纯函数，输入输出可预测
        },
        CommonCallback: CommonCallback{
            OnPageLoad: func(ctx *Context, resp response.Response) (initData *usercall.OnPageLoadResp, err error) {
                // 返回默认的请求参数，预填充表单
                return &usercall.OnPageLoadResp{
                    Request: &Json2CsvReq{
                        InputJson: `{"name":"张三","age":25,"city":"北京"}`,
                        Delimiter: ",",
                    },
                }, nil
            },
        },
    }
    
    // JSON转YAML（纯函数）
    var json2yamlOption = &FormFunctionOptions{
        BaseConfig: BaseConfig{
            Router:        "/api/demo/form/json2yaml",
            Method:        "POST",
            EnglishName:   "json2yaml",
            ChineseName:   "JSON转YAML",
            ApiDesc:       "将JSON数据转换为YAML格式",
            Tags:          []string{"JSON转换", "数据转换"},
            Group:         JsonConverterGroup,
            Request:       &Json2YamlReq{},
            Response:      &Json2YamlResp{},
            CreateTables:  []interface{}{&ConversionRecord{}},
            Timeout:       60000,
            Async:         false,
            FunctionType:  FunctionTypePure, // 纯函数
        },
    }
    
    // JSON转XML（纯函数）
    var json2xmlOption = &FormFunctionOptions{
        BaseConfig: BaseConfig{
            Router:        "/api/demo/form/json2xml",
            Method:        "POST",
            EnglishName:   "json2xml",
            ChineseName:   "JSON转XML",
            ApiDesc:       "将JSON数据转换为XML格式",
            Tags:          []string{"JSON转换", "数据转换"},
            Group:         JsonConverterGroup,
            Request:       &Json2XmlReq{},
            Response:      &Json2XmlResp{},
            CreateTables:  []interface{}{&ConversionRecord{}},
            Timeout:       60000,
            Async:         false,
            FunctionType:  FunctionTypePure, // 纯函数
        },
    }
    
    // 手动注册路由
    runner.Post("/api/demo/form/json2csv", Json2CsvHandler, json2csvOption)
    runner.Post("/api/demo/form/json2yaml", Json2YamlHandler, json2yamlOption)
    runner.Post("/api/demo/form/json2xml", Json2XmlHandler, json2xmlOption)
}

// 产品管理系统组示例
func ProductManagementExamples() {
    // 产品登记（动态函数）
    var productRegisterOption = &FormFunctionOptions{
        BaseConfig: BaseConfig{
            Router:        "/api/demo/form/product_register",
            Method:        "POST",
            EnglishName:   "product_register",
            ChineseName:   "产品登记",
            ApiDesc:       "新产品登记表单",
            Tags:          []string{"产品管理", "登记"},
            Group:         ProductManagementGroup,
            Request:       &ProductRegisterReq{},
            Response:      &ProductRegisterResp{},
            CreateTables:  []interface{}{&Product{}},
            Timeout:       60000,
            Async:         false,
            FunctionType:  FunctionTypeDynamic, // 动态函数，结果不可预测
        },
        ComponentCallback: ComponentCallback{
            OnInputFuzzyMap: map[string]OnInputFuzzy{
                "company_name": func(ctx *Context, req *usercall.OnInputFuzzyReq) (*usercall.OnInputFuzzyResp, error) {
                    // 根据输入的公司名称模糊搜索
                    return &usercall.OnInputFuzzyResp{
                        Values: []*usercall.InputFuzzyItem{
                            {Value: "阿里巴巴"},
                            {Value: "腾讯科技"},
                            {Value: "字节跳动"},
                        },
                    }, nil
                },
            },
            OnInputValidateMap: map[string]OnInputValidate{
                "product_name": func(ctx *Context, req *usercall.OnInputValidateReq) (*usercall.OnInputValidateResp, error) {
                    // 验证产品名称是否已存在
                    productName := req.Value.(string)
                    if productName == "" {
                        return &usercall.OnInputValidateResp{
                            ErrorMsg: "产品名称不能为空",
                        }, nil
                    }
                    return &usercall.OnInputValidateResp{
                        ErrorMsg: "", // 空字符串表示验证通过
                    }, nil
                },
            },
        },
    }
    
    // 产品列表（动态函数）
    var productListOption = &TableFunctionOptions{
        BaseConfig: BaseConfig{
            Router:        "/api/demo/table/product_list",
            Method:        "GET",
            EnglishName:   "product_list",
            ChineseName:   "产品列表",
            ApiDesc:       "产品列表管理",
            Tags:          []string{"产品管理", "列表"},
            Group:         ProductManagementGroup,
            Request:       &ProductListReq{},
            Response:      &ProductListResp{},
            CreateTables:  []interface{}{&Product{}},
            Timeout:       30000,
            Async:         false,
            FunctionType:  FunctionTypeDynamic, // 动态函数，数据随时变化
        },
        TableConfig: TableConfig{
            AutoCrudTable: &Product{},
            ComponentCallback: ComponentCallback{
                OnInputFuzzyMap: map[string]OnInputFuzzy{
                    "category_search": func(ctx *Context, req *usercall.OnInputFuzzyReq) (*usercall.OnInputFuzzyResp, error) {
                        // 产品分类模糊搜索
                        return &usercall.OnInputFuzzyResp{
                            Values: []*usercall.InputFuzzyItem{
                                {Value: "电子产品"},
                                {Value: "服装鞋帽"},
                                {Value: "家居用品"},
                                {Value: "食品饮料"},
                            },
                        }, nil
                    },
                },
            },
        },
    }
    
    // 手动注册路由
    runner.Post("/api/demo/form/product_register", ProductRegisterHandler, productRegisterOption)
    runner.Get("/api/demo/table/product_list", ProductListHandler, productListOption)
}

// 静态函数示例
func StaticFunctionExamples() {
    // 获取系统时间（静态函数）
    var getSystemTimeOption = &FormFunctionOptions{
        BaseConfig: BaseConfig{
            Router:        "/api/demo/form/get_system_time",
            Method:        "GET",
            EnglishName:   "get_system_time",
            ChineseName:   "获取系统时间",
            ApiDesc:       "获取当前系统时间",
            Tags:          []string{"系统工具", "时间"},
            Request:       &GetSystemTimeReq{},
            Response:      &GetSystemTimeResp{},
            Timeout:       5000,
            Async:         false,
            FunctionType:  FunctionTypeStatic, // 静态函数，结果恒定
        },
    }
    
    runner.Get("/api/demo/form/get_system_time", GetSystemTimeHandler, getSystemTimeOption)
}
```

## 🎯 函数类型说明

### 1. **静态函数** (FunctionTypeStatic)
- **特点**: 无需参数，或者输入参数，但是结果永远恒定
- **示例**: 获取系统时间、获取版本号、获取配置信息
- **用途**: 系统信息查询、配置获取等

### 2. **动态函数** (FunctionTypeDynamic)
- **特点**: 请求参数不可预测，响应参数不可预测
- **示例**: 查询用户信息、产品列表、订单管理
- **用途**: 业务数据处理、数据库查询等

### 3. **纯函数** (FunctionTypePure)
- **特点**: 输入输出可预测，如数学函数
- **示例**: JSON转换、数学计算、格式转换
- **用途**: 数据处理、格式转换、计算等

## 🎯 函数组设计说明

### 1. **简化设计**
```go
type FunctionGroup struct {
    Name string `json:"name"` // 组名称，如 "JSON转换"
    // 后续按需要添加字段
}
```

### 2. **使用场景**
```go
// 预定义组，复用方便
var JsonConverterGroup = &FunctionGroup{
    Name: "JSON转换",
}

// 在函数中使用
Group: JsonConverterGroup,
```

### 3. **扩展性**
- 当前：只有name字段
- 后续按需要添加：description、version、author等

## 🎯 设计优势

### 1. **向后兼容**
- 保持现有 `FunctionOptions` 不变
- 新设计作为补充，不破坏现有代码

### 2. **分类清晰**
- 按作用范围和触发时机分类
- 每个分类职责明确
- 易于理解和维护

### 3. **扩展性好**
- 新增函数类型时只需添加对应的SpecificCallback
- 通用回调可以复用
- 组件级回调可以跨函数类型复用

### 4. **类型安全**
- 编译时就能发现配置错误
- 不同函数类型有不同的回调需求

### 5. **简洁直观**
- 直接赋值，一目了然
- 分类明确，易于查找

### 6. **函数组支持**
- 简化设计，只有name字段
- 预定义组，复用方便
- 为后续扩展预留空间

### 7. **OnDryRun 优化**
- 从通用回调中移除
- 只在表单函数中使用
- 避免在 table 函数中的无意义使用

## 🚀 迁移策略

1. **第一阶段**：实现新的选项系统（FormFunctionOptions、TableFunctionOptions）
2. **第二阶段**：保持现有 FunctionOptions 向后兼容
3. **第三阶段**：逐步迁移现有代码到新选项系统
4. **第四阶段**：最终废弃旧的 FunctionOptions（可选）

这个设计怎么样？现在 `OnDryRun` 只在表单函数中使用，避免了在 table 函数中的无意义使用，你觉得还有什么需要调整的吗？ 