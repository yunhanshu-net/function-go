# 云函数平台组件渲染架构文档

## 🏗️ 整体架构概览

本项目是一个**无代码云函数平台**，通过 **结构体标签驱动** 的方式自动生成前端UI组件，实现了从Go结构体到前端表单/表格的自动化渲染。

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Go 结构体      │───▶│   Runner 标签    │───▶│   Widget 组件    │
│  (业务模型)      │    │   (配置驱动)     │    │   (UI渲染)      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   数据验证       │    │   回调处理       │    │   前端界面       │
│  (Validation)   │    │  (Callbacks)    │    │  (Frontend UI)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🎯 核心设计思想

### 1. **零代码配置** - 通过标签驱动一切
```go
type UserReq struct {
    Name     string `runner:"code:name;name:用户名;widget:input;placeholder:请输入用户名"`
    Category string `runner:"code:category;name:分类;widget:select;options:A,B,C"`
    Active   bool   `runner:"code:active;name:启用状态;widget:switch"`
}
```

### 2. **组件化渲染** - 每个字段对应一个Widget
- 自动根据标签配置生成对应的UI组件
- 支持表单模式(Form)和表格模式(Table)两种渲染方式
- 统一的组件接口，易于扩展新组件类型

### 3. **动态交互** - 回调机制支持复杂业务逻辑
- 支持页面加载回调、输入模糊搜索、字段联动等
- 前后端解耦，回调逻辑在后端定义，前端自动处理

## 📁 目录结构说明

```
function-go/
├── view/widget/              # UI组件定义层
│   ├── widget.go            # 组件工厂 (Widget Factory)
│   ├── type.go              # 组件类型常量定义
│   ├── input.go             # 输入框组件
│   ├── select.go            # 下拉框组件
│   ├── checkbox.go          # 多选框组件
│   └── ...                  # 其他组件实现
├── pkg/dto/api/             # API定义层
│   ├── form_request.go      # 表单请求参数解析
│   ├── form_response.go     # 表单响应参数构建
│   └── ...                  # 其他API相关定义
└── soft/beiluo/lookup/code/api/widgets/  # 最佳实践示例层
    ├── callback_demo.go     # 回调函数演示
    ├── workflow_demo.go     # 工作流管理演示
    ├── crud_demo.go         # CRUD操作演示
    ├── data_analysis.go     # 数据分析演示
    └── file_manager.go      # 文件管理演示
```

## 🔧 核心组件详解

### 1. Widget组件系统 (`/view/widget/`)

#### 组件工厂模式 (`widget.go`)
```go
// 根据渲染类型和字段信息创建对应的Widget组件
func NewWidget(info *tagx.RunnerFieldInfo, renderType string) (Widget, error) {
    widgetType := info.Tags["widget"] // 从标签获取组件类型
    
    switch renderType {
    case response.RenderTypeTable:  // 表格模式
        // 根据widgetType创建对应组件
    case response.RenderTypeForm:   // 表单模式
        // 根据widgetType创建对应组件
    }
}
```

**设计特点：**
- **工厂模式**：统一的组件创建入口，便于管理和扩展
- **渲染模式分离**：同一个组件在表格和表单中可能有不同的表现
- **标签驱动**：完全基于结构体标签配置，零硬编码

#### 组件类型定义 (`type.go`)
```go
const (
    WidgetInput       = "input"        // 文本输入框
    WidgetSelect      = "select"       // 下拉框
    WidgetCheckbox    = "checkbox"     // 多选框
    WidgetRadio       = "radio"        // 单选框
    WidgetSwitch      = "switch"       // 开关
    WidgetSlider      = "slider"       // 滑块
    WidgetDateTime    = "datetime"     // 日期时间
    WidgetMultiSelect = "multiselect"  // 多选下拉
    WidgetFileUpload  = "file_upload"  // 文件上传
    WidgetFileDisplay = "file_display" // 文件展示
)
```

#### 具体组件实现示例

**输入框组件** (`input.go`)
```go
type InputWidget struct {
    Mode         string `json:"mode"`          // line_text/text_area/password
    Placeholder  string `json:"placeholder"`   // 占位符
    DefaultValue string `json:"default_value"` // 默认值
}

func NewInputWidget(info *tagx.RunnerFieldInfo) (*InputWidget, error) {
    return &InputWidget{
        Mode:         info.Tags["mode"],        // 从标签读取配置
        Placeholder:  info.Tags["placeholder"],
        DefaultValue: info.Tags["default_value"],
    }, nil
}
```

**下拉框组件** (`select.go`)
```go
type SelectWidget struct {
    Options      []string `json:"options"`       // 选项列表
    Multiple     bool     `json:"multiple"`      // 是否多选
    DefaultValue string   `json:"default_value"` // 默认值
}

func NewSelectWidget(info *tagx.RunnerFieldInfo) (Widget, error) {
    return &SelectWidget{
        Options:      strings.Split(info.Tags["options"], ","), // 解析逗号分隔的选项
        DefaultValue: info.Tags["default_value"],
        Multiple:     info.Tags["multiple"] != "",              // 检查是否设置多选
    }, nil
}
```

### 2. API定义系统 (`/pkg/dto/api/`)

#### 表单请求参数解析 (`form_request.go`)
```go
type FormRequestParamInfo struct {
    Code         string      `json:"code"`          // 字段标识
    Name         string      `json:"name"`          // 字段名称
    Required     bool        `json:"required"`      // 是否必填
    DefaultValue interface{} `json:"default_value"` // 默认值
    Callbacks    string      `json:"callbacks"`     // 回调配置
    Validates    string      `json:"validates"`     // 验证规则
    WidgetConfig interface{} `json:"widget_config"` // Widget配置
    WidgetType   string      `json:"widget_type"`   // 组件类型
    ValueType    string      `json:"value_type"`    // 数据类型
}

// 核心函数：将Go结构体转换为前端可用的表单配置
func NewFormRequestParams(el interface{}, renderType string) (*FormRequestParams, error) {
    // 1. 反射解析结构体
    typeOf := reflect.TypeOf(el)
    reqFields, err := tagx.ParseStructFieldsTypeOf(typeOf, "runner")
    
    // 2. 为每个字段创建Widget配置
    for _, field := range reqFields {
        widgetIns, err := widget.NewWidget(field, renderType)
        info := &FormRequestParamInfo{
            Code:         field.GetCode(),
            WidgetConfig: widgetIns,
            WidgetType:   widgetIns.GetWidgetType(),
            // ... 其他配置
        }
    }
}
```

**数据流说明：**
```
Go结构体 → 标签解析 → Widget创建 → JSON配置 → 前端渲染
```

### 3. 回调机制系统

#### 回调类型定义
```go
type FunctionInfo struct {
    // 页面加载回调 - 页面初始化时触发
    OnPageLoad func(ctx *Context, resp response.Response) (*OnPageLoadResp, error)
    
    // 输入模糊搜索回调 - 用户输入时触发
    OnInputFuzzyMap map[string]OnInputFuzzy
    
    // API创建后回调 - 用于初始化数据
    OnApiCreated func(ctx *Context, req *OnApiCreatedReq) error
}
```

#### 回调实现示例
```go
// 在字段标签中配置回调
Location string `runner:"..." callback:"OnInputFuzzy(field:location_suggest,delay:300,min:1)"`

// 在FunctionInfo中实现回调逻辑
OnInputFuzzyMap: map[string]runner.OnInputFuzzy{
    "location_suggest": func(ctx *Context, req *OnInputFuzzyReq) (*OnInputFuzzyResp, error) {
        // 模糊搜索逻辑
        locations := []string{"北京", "上海", "深圳", "广州"}
        resp := &OnInputFuzzyResp{}
        for _, location := range locations {
            if strings.Contains(location, req.Value) {
                resp.Values = append(resp.Values, &InputFuzzyItem{
                    Value: location,
                })
            }
        }
        return resp, nil
    },
}
```

## 🏷️ Runner标签系统详解

### 标签语法规范
```go
runner:"参数1:值1;参数2:值2;参数3:值3"
```

### 核心标签参数

| 参数 | 说明 | 示例 |
|------|------|------|
| `code` | 字段标识，用于前后端映射 | `code:username` |
| `name` | 字段显示名称 | `name:用户名` |
| `widget` | UI组件类型 | `widget:input` |
| `type` | 数据类型 | `type:string` |
| `placeholder` | 占位符提示 | `placeholder:请输入用户名` |
| `default_value` | 默认值 | `default_value:admin` |
| `options` | 选项列表(逗号分隔) | `options:选项1,选项2,选项3` |
| `show` | 显示场景控制 | `show:list` (仅在列表显示) |
| `hidden` | 隐藏控制 | `hidden:all` (全部隐藏) |
| `callback` | 回调配置 | `callback:OnInputFuzzy(...)` |

### 组件特定参数

#### Input组件
```go
`runner:"widget:input;mode:text_area;max:500;placeholder:请输入描述"`
```
- `mode`: `line_text`(单行) / `text_area`(多行) / `password`(密码)
- `max`: 最大字符数
- `min`: 最小字符数

#### Select组件
```go
`runner:"widget:select;options:选项1,选项2,选项3;default_value:选项1;multiple:true"`
```
- `options`: 选项列表(逗号分隔)
- `multiple`: 是否多选
- `placeholder`: 未选择时的提示

#### Switch组件
```go
`runner:"widget:switch;true_label:开启;false_label:关闭;default_value:true"`
```
- `true_label`: 开启状态文本
- `false_label`: 关闭状态文本

#### DateTime组件
```go
`runner:"widget:datetime;format:date;min_date:today;max_date:2025-12-31"`
```
- `format`: `date`(日期) / `time`(时间) / `datetime`(日期时间)
- `min_date`: 最小日期
- `max_date`: 最大日期

#### Slider组件
```go
`runner:"widget:slider;min:0;max:100;step:5;unit:%;show_tooltip:true"`
```
- `min`: 最小值
- `max`: 最大值
- `step`: 步长
- `unit`: 单位
- `show_tooltip`: 显示提示

## 📋 最佳实践示例分析

### 1. 回调函数演示 (`callback_demo.go`)

**核心特性：**
- 展示输入模糊搜索回调的实现
- 页面加载时的数据初始化
- 多种回调类型的组合使用

**关键代码片段：**
```go
// 地区字段配置模糊搜索回调
Location string `runner:"code:location;name:工作地点;widget:input" callback:"OnInputFuzzy(field:location_suggest,delay:300,min:1)"`

// 回调实现
OnInputFuzzyMap: map[string]runner.OnInputFuzzy{
    "location_suggest": func(ctx *Context, req *OnInputFuzzyReq) (*OnInputFuzzyResp, error) {
        // 模糊搜索逻辑实现
        locations := []string{"北京", "上海", "深圳"}
        // 返回匹配结果
    },
}
```

### 2. 工作流管理演示 (`workflow_demo.go`)

**核心特性：**
- 复杂表单的组件组合
- 业务逻辑驱动的智能建议
- 数据库操作与UI组件的结合

**关键代码片段：**
```go
// 复杂的表单字段配置
Priority string `runner:"code:priority;name:优先级;widget:select;options:低,中,高,紧急"`
Status   string `runner:"code:status;name:任务状态;widget:radio;options:待办,进行中,已完成"`
Progress int    `runner:"code:progress;name:进度;widget:slider;min:0;max:100;unit:%"`

// 业务逻辑处理
func WorkflowDemo(ctx *Context, req *WorkflowDemoReq, resp response.Response) error {
    // 根据请求参数生成智能建议
    suggestions := generateWorkflowSuggestions(req)
    return resp.Form(&WorkflowDemoResp{
        Suggestions: suggestions,
    }).Build()
}
```

### 3. CRUD操作演示 (`crud_demo.go`)

**核心特性：**
- 自动化的CRUD操作支持
- 表格模式的数据展示
- 分页、筛选、排序功能

**关键代码片段：**
```go
// 启用自动CRUD功能
AutoCrudTable: &Product{},

// 列表查询配置
type ProductListReq struct {
    query.PageInfoReq                    // 自动分页支持
    Name     string `runner:"code:name;name:产品名称;placeholder:按产品名称搜索"`
    Category string `runner:"code:category;name:分类;widget:select;options:手机,笔记本"`
    Status   string `runner:"code:status;name:状态;widget:select;options:启用,禁用"`
}
```

### 4. 数据分析演示 (`data_analysis.go`)

**核心特性：**
- 复杂的数据筛选配置
- 图表组件的使用
- 统计分析结果的展示

**关键代码片段：**
```go
// 数据分析配置
AnalysisType string `runner:"code:analysis_type;name:分析类型;widget:select;options:汇总分析,趋势分析"`
ChartType    string `runner:"code:chart_type;name:图表类型;widget:radio;options:柱状图,折线图,饼图"`
ShowTrend    bool   `runner:"code:show_trend;name:显示趋势;widget:switch"`

// 分析结果返回
return resp.Form(&DataAnalysisResp{
    ChartData:       analysis.ChartData,       // 图表JSON数据
    TrendAnalysis:   analysis.TrendAnalysis,   // 趋势分析结果
    Recommendations: analysis.Recommendations, // 智能建议
}).Build()
```

### 5. 文件管理演示 (`file_manager.go`)

**核心特性：**
- 文件上传组件的使用
- 多种处理模式的配置
- 文件处理结果的展示

**关键代码片段：**
```go
// 文件上传配置
Files *files.Files `runner:"code:files;name:文件列表;widget:file_upload;multiple:true;max_size:10240"`

// 处理模式配置
ProcessMode string `runner:"code:process_mode;name:处理模式;widget:radio;options:仅上传,分析文件,提取内容"`

// 文件处理逻辑
switch req.ProcessMode {
case "仅上传":
    err = processUploadOnly(ctx, writer, db, file, path, req)
case "分析文件":
    err = processWithAnalysis(ctx, writer, db, file, path, req)
}
```

## 🔄 数据流转过程

### 1. 前端界面生成流程
```
1. Go结构体定义 → 2. Runner标签解析 → 3. Widget组件创建 → 4. API配置生成 → 5. 前端界面渲染
```

### 2. 用户交互处理流程
```
1. 用户操作 → 2. 前端事件触发 → 3. 后端回调处理 → 4. 响应数据返回 → 5. 界面状态更新
```

### 3. 表单提交处理流程
```
1. 表单数据收集 → 2. 数据验证 → 3. 后端业务处理 → 4. 数据库操作 → 5. 响应结果返回
```

## 🚀 扩展指南

### 添加新的Widget组件

1. **定义组件类型常量** (`type.go`)
```go
const WidgetNewComponent = "new_component"
```

2. **实现组件结构体** (`new_component.go`)
```go
type NewComponentWidget struct {
    // 组件特有的配置属性
    CustomProperty string `json:"custom_property"`
}

func NewNewComponentWidget(info *tagx.RunnerFieldInfo) (*NewComponentWidget, error) {
    return &NewComponentWidget{
        CustomProperty: info.Tags["custom_property"],
    }, nil
}

func (w *NewComponentWidget) GetWidgetType() string {
    return WidgetNewComponent
}
```

3. **注册到工厂** (`widget.go`)
```go
case WidgetNewComponent:
    return NewNewComponentWidget(info)
```

### 添加新的回调类型

1. **定义回调接口**
```go
type OnNewCallback func(ctx *Context, req *OnNewCallbackReq) (*OnNewCallbackResp, error)
```

2. **在FunctionInfo中添加回调字段**
```go
type FunctionInfo struct {
    OnNewCallbackMap map[string]OnNewCallback
    // ... 其他字段
}
```

3. **在具体实现中使用**
```go
OnNewCallbackMap: map[string]OnNewCallback{
    "my_callback": func(ctx *Context, req *OnNewCallbackReq) (*OnNewCallbackResp, error) {
        // 回调逻辑实现
    },
}
```

## 🎯 设计原则总结

1. **配置驱动** - 通过标签配置而非硬编码
2. **组件化** - 每个UI元素都是独立的Widget组件
3. **类型安全** - 基于Go的强类型系统
4. **可扩展** - 易于添加新组件和回调类型
5. **分离关注** - UI逻辑与业务逻辑分离
6. **自动化** - 最大程度减少手动配置

这套架构实现了真正的"代码即配置"，让开发者可以专注于业务逻辑，而UI界面由系统自动生成，大大提高了开发效率和一致性。 