# 云函数平台组件渲染系统架构文档

## 🏗️ 整体架构

本项目实现了一个**基于Go结构体标签驱动的UI组件自动生成系统**，通过在Go结构体字段上添加`runner`标签，自动生成对应的前端表单组件。

```
Go结构体 + Runner标签 → Widget组件 → 前端UI界面
```

## 📁 核心目录结构

```
function-go/
├── view/widget/              # UI组件定义层
│   ├── widget.go            # 组件工厂
│   ├── type.go              # 组件类型定义
│   ├── input.go             # 输入框组件
│   ├── select.go            # 下拉框组件
│   └── ...                  # 其他UI组件
├── pkg/dto/api/             # API配置层
│   ├── form_request.go      # 表单配置生成
│   └── form_response.go     # 响应配置生成
└── soft/.../widgets/        # 最佳实践示例
    ├── callback_demo.go     # 回调演示
    ├── workflow_demo.go     # 工作流演示
    ├── crud_demo.go         # CRUD演示
    ├── data_analysis.go     # 数据分析演示
    └── file_manager.go      # 文件管理演示
```

## 🔧 核心组件系统

### 1. Widget组件系统 (`/view/widget/`)

#### 工厂模式创建组件
```go
// widget.go - 组件工厂核心逻辑
func NewWidget(info *tagx.RunnerFieldInfo, renderType string) (Widget, error) {
    widgetType := info.Tags["widget"] // 从标签读取组件类型
    
    switch renderType {
    case response.RenderTypeForm:  // 表单模式
        switch widgetType {
        case WidgetInput:      return NewInputWidget(info)
        case WidgetSelect:     return NewSelectWidget(info)
        case WidgetCheckbox:   return newCheckboxWidget(info)
        case WidgetSwitch:     return newSwitchWidget(info)
        // ... 更多组件类型
        }
    case response.RenderTypeTable: // 表格模式
        // 表格模式下的组件配置
    }
}
```

#### 组件类型常量定义
```go
// type.go - 所有支持的组件类型
const (
    WidgetInput       = "input"        // 文本输入框
    WidgetSelect      = "select"       // 下拉框
    WidgetCheckbox    = "checkbox"     // 多选框
    WidgetRadio       = "radio"        // 单选框
    WidgetSwitch      = "switch"       // 开关按钮
    WidgetSlider      = "slider"       // 滑块
    WidgetDateTime    = "datetime"     // 日期时间选择
    WidgetMultiSelect = "multiselect"  // 多选下拉框
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
        Mode:         stringsx.DefaultString(info.Tags["mode"], "line_text"),
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
        DefaultValue: info.Tags["default_value"],
        Options:      strings.Split(info.Tags["options"], ","), // 解析逗号分隔的选项
        Multiple:     info.Tags["multiple"] != "",              // 是否设置了multiple标签
    }, nil
}
```

### 2. API配置生成系统 (`/pkg/dto/api/`)

#### 表单配置生成核心逻辑
```go
// form_request.go - 将Go结构体转换为前端表单配置
type FormRequestParamInfo struct {
    Code         string      `json:"code"`          // 字段标识
    Name         string      `json:"name"`          // 字段显示名称
    Required     bool        `json:"required"`      // 是否必填
    DefaultValue interface{} `json:"default_value"` // 默认值
    Callbacks    string      `json:"callbacks"`     // 回调配置
    Validates    string      `json:"validates"`     // 验证规则
    WidgetConfig interface{} `json:"widget_config"` // Widget配置对象
    WidgetType   string      `json:"widget_type"`   // 组件类型
    ValueType    string      `json:"value_type"`    // 数据类型
}

// 核心转换函数：Go结构体 → 前端表单配置JSON
func NewFormRequestParams(el interface{}, renderType string) (*FormRequestParams, error) {
    // 1. 反射解析结构体和runner标签
    typeOf := reflect.TypeOf(el)
    reqFields, err := tagx.ParseStructFieldsTypeOf(typeOf, "runner")
    
    // 2. 为每个字段创建对应的Widget配置
    children := make([]*FormRequestParamInfo, 0, len(reqFields))
    for _, field := range reqFields {
        // 创建Widget实例
        widgetIns, err := widget.NewWidget(field, renderType)
        
        // 构建字段配置
        info := &FormRequestParamInfo{
            Code:         field.GetCode(),
            Name:         field.GetName(),
            WidgetConfig: widgetIns,              // Widget配置对象
            WidgetType:   widgetIns.GetWidgetType(), // 组件类型
            // ... 其他配置
        }
        children = append(children, info)
    }
    
    return &FormRequestParams{
        RenderType: renderType,
        Children:   children,
    }, nil
}
```

## 🔐 字段权限控制系统

### 💡 设计理念
分离**显示**和**编辑**权限，使用语义化的标签让大模型和开发者都能直观理解字段的权限控制。

### 🏷️ 基础权限标签（推荐使用）

#### **只读字段** (`readonly:true`)
```go
// 系统字段，只显示不可编辑
ID        uint      `runner:"code:id;name:ID;readonly:true"`
CreatedAt time.Time `runner:"code:created_at;name:创建时间;widget:datetime;readonly:true"`
CreatedBy string    `runner:"code:created_by;name:创建人;readonly:true"`
Rating    float64   `runner:"code:rating;name:评分;widget:slider;readonly:true"`
```

#### **创建时可编辑** (`create_only:true`)
```go
// 创建时可设置，之后只读
UserType    string `runner:"code:user_type;name:用户类型;widget:select;options:管理员,普通用户;create_only:true"`
PublishDate string `runner:"code:publish_date;name:发布日期;widget:datetime;format:date;create_only:true"`
```

#### **更新时可编辑** (`update_only:true`)
```go
// 创建时不显示，只在更新时可编辑
LastLoginIP string `runner:"code:last_login_ip;name:最后登录IP;widget:input;update_only:true"`
Permissions string `runner:"code:permissions;name:权限设置;widget:multiselect;update_only:true"`
Config      string `runner:"code:config;name:配置;widget:input;mode:text_area;update_only:true"`
```

#### **完全可编辑** (`editable:true`)
```go
// 所有场景都可编辑（默认行为，可省略）
Name     string  `runner:"code:name;name:姓名;widget:input;editable:true;required:true"`
Price    float64 `runner:"code:price;name:价格;widget:slider;editable:true"`
Status   string  `runner:"code:status;name:状态;widget:switch;editable:true"`
```

#### **完全隐藏** (`hidden:true`)
```go
// 所有场景都不显示
InternalID string `runner:"code:internal_id;name:内部ID;hidden:true"`
Secret     string `runner:"code:secret;name:密钥;hidden:true"`
```

### 🎯 场景权限标签（高级用法）

#### **场景组合控制**
```go
// 列表中只读，表单中可编辑
Description string `runner:"code:description;name:描述;widget:input;mode:text_area;list:hidden;form:editable"`

// 创建时必填，更新时可选，列表中隐藏密码
Password string `runner:"code:password;name:密码;widget:input;mode:password;create:required;update:optional;list:hidden"`

// 详情页显示，其他场景隐藏
Detail string `runner:"code:detail;name:详细信息;widget:input;detail:readonly;list:hidden;form:hidden"`
```

#### **支持的场景标识**
- `list` - 列表页面
- `create` - 创建表单
- `update` - 更新表单  
- `detail` - 详情页面
- `form` - 表单页面（包含create和update）

#### **支持的权限级别**
- `readonly` - 只读显示
- `editable` - 可编辑
- `required` - 必填
- `optional` - 可选
- `hidden` - 隐藏

### 📊 权限标签对比表

| 权限需求 | 旧标签写法 | 新标签写法 | 语义清晰度 |
|---------|-----------|-----------|-----------|
| 只读字段 | `show:list` | `readonly:true` | ✅ 非常清晰 |
| 创建可编辑 | `show:create,list` | `create_only:true` | ✅ 一目了然 |
| 更新可编辑 | `show:update,list` | `update_only:true` | ✅ 一目了然 |
| 完全可编辑 | `show:list,create,update` | `editable:true` | ✅ 简洁明了 |
| 完全隐藏 | `hidden:list,create,update` | `hidden:true` | ✅ 简洁明了 |

### 🚀 实际应用示例

#### **用户管理系统**
```go
type User struct {
    // 系统字段 - 只读
    ID        uint      `json:"id" runner:"code:id;name:用户ID;readonly:true"`
    CreatedAt time.Time `json:"created_at" runner:"code:created_at;name:创建时间;widget:datetime;readonly:true"`
    
    // 基础信息 - 完全可编辑
    Username string `json:"username" runner:"code:username;name:用户名;widget:input;editable:true;required:true"`
    Email    string `json:"email" runner:"code:email;name:邮箱;widget:input;editable:true;required:true"`
    
    // 权限信息 - 创建时设置
    Role string `json:"role" runner:"code:role;name:角色;widget:select;options:管理员,普通用户;create_only:true"`
    
    // 状态信息 - 更新时可编辑
    LastLoginIP string `json:"last_login_ip" runner:"code:last_login_ip;name:最后登录IP;widget:input;update_only:true"`
    
    // 敏感信息 - 完全隐藏
    PasswordHash string `json:"password_hash" runner:"code:password_hash;name:密码哈希;hidden:true"`
}
```

#### **订单管理系统**
```go
type Order struct {
    // 订单号 - 只读
    OrderNo string `json:"order_no" runner:"code:order_no;name:订单号;readonly:true"`
    
    // 基础信息 - 创建时设置
    CustomerID uint   `json:"customer_id" runner:"code:customer_id;name:客户ID;widget:select;create_only:true"`
    TotalAmount float64 `json:"total_amount" runner:"code:total_amount;name:订单金额;widget:slider;create_only:true"`
    
    // 状态管理 - 可更新
    Status     string `json:"status" runner:"code:status;name:订单状态;widget:select;options:待付款,已付款,已发货,已完成;update_only:true"`
    Remark     string `json:"remark" runner:"code:remark;name:备注;widget:input;mode:text_area;editable:true"`
    
    // 物流信息 - 发货后可编辑
    TrackingNo string `json:"tracking_no" runner:"code:tracking_no;name:快递单号;widget:input;update_only:true"`
}
```

### 🎨 大模型提示词建议

**推荐提示词模板：**
```
请为用户管理系统生成代码，字段权限要求：
- ID、创建时间等系统字段：使用 readonly:true
- 用户名、邮箱等基础信息：使用 editable:true  
- 用户角色等关键信息：使用 create_only:true
- 最后登录信息：使用 update_only:true
- 密码哈希等敏感信息：使用 hidden:true
```

这样的标签设计让权限控制的语义更加清晰，大模型可以直观理解每个字段在不同场景下的行为！

## 🏷️ Runner标签系统详解

### 标签语法规范
```go
runner:"参数1:值1;参数2:值2;参数3:值3"
```

### 核心标签参数说明

| 参数 | 说明 | 示例 | 必需 |
|------|------|------|------|
| `code` | 字段唯一标识 | `code:username` | ✅ |
| `name` | 字段显示名称 | `name:用户名` | ✅ |
| `widget` | UI组件类型 | `widget:input` | ✅ |
| `type` | 数据类型 | `type:string` | ❌ |
| `placeholder` | 占位符提示 | `placeholder:请输入用户名` | ❌ |
| `default_value` | 默认值 | `default_value:admin` | ❌ |
| `options` | 选项列表 | `options:选项1,选项2,选项3` | ❌ |
| `callback` | 回调配置 | `callback:OnInputFuzzy(...)` | ❌ |

### 组件特定参数详解

#### Input输入框
```go
// 单行文本输入
Name string `runner:"code:name;name:姓名;widget:input;mode:line_text;placeholder:请输入姓名;max:50;min:2"`

// 多行文本域
Desc string `runner:"code:desc;name:描述;widget:input;mode:text_area;placeholder:请输入描述;max:500"`

// 密码输入
Password string `runner:"code:password;name:密码;widget:input;mode:password;placeholder:请输入密码"`
```
**支持参数：**
- `mode`: `line_text`(单行) / `text_area`(多行) / `password`(密码)
- `max`: 最大字符数限制
- `min`: 最小字符数限制
- `placeholder`: 占位符提示文本

#### Select下拉框
```go
// 单选下拉框
Category string `runner:"code:category;name:分类;widget:select;options:分类A,分类B,分类C;default_value:分类A;placeholder:请选择分类"`

// 多选下拉框
Tags []string `runner:"code:tags;name:标签;widget:select;options:标签1,标签2,标签3;multiple:true"`
```
**支持参数：**
- `options`: 逗号分隔的选项列表
- `multiple`: 是否多选模式
- `default_value`: 默认选中值
- `placeholder`: 未选择时的提示文本

#### MultiSelect多选下拉框（支持动态数据源）
```go
// 基础多选
Skills []string `runner:"code:skills;name:技能;widget:multiselect;options:Java,Python,Go,JavaScript;multiple_limit:5"`

// 支持创建新选项
CustomTags []string `runner:"code:custom_tags;name:自定义标签;widget:multiselect;allow_create:true;collapse_tags:true;placeholder:输入或选择标签"`
```
**支持参数：**
- `options`: 预定义选项列表
- `multiple_limit`: 最大选择数量限制(0为不限制)
- `allow_create`: 是否允许创建新选项
- `collapse_tags`: 是否折叠显示已选标签
- `default_value`: 默认选中值(逗号分隔)
- `placeholder`: 占位符文本

#### Checkbox多选框
```go
// 水平排列
Categories []string `runner:"code:categories;name:分类;widget:checkbox;options:文档,图片,视频,音频;direction:horizontal;show_select_all:true"`

// 垂直排列多列显示
Regions []string `runner:"code:regions;name:区域;widget:checkbox;options:华北,华东,华南,西南,东北,西北;direction:vertical;columns:3"`
```
**支持参数：**
- `options`: 选项列表
- `direction`: `horizontal`(水平) / `vertical`(垂直)
- `columns`: 垂直排列时的列数
- `show_select_all`: 是否显示"全选"选项

#### Radio单选框
```go
// 水平排列
WorkMode string `runner:"code:work_mode;name:工作模式;widget:radio;options:远程,办公室,混合;direction:horizontal;default_value:混合"`

// 垂直排列
Priority string `runner:"code:priority;name:优先级;widget:radio;options:低,中,高,紧急;direction:vertical"`
```
**支持参数：**
- `options`: 选项列表
- `direction`: `horizontal`(水平) / `vertical`(垂直)
- `default_value`: 默认选中值

#### Switch开关
```go
// 基础开关
Active bool `runner:"code:active;name:启用状态;widget:switch;true_label:启用;false_label:禁用;default_value:true"`

// 简化开关
Notification bool `runner:"code:notification;name:消息通知;widget:switch;true_label:开启;false_label:关闭"`
```
**支持参数：**
- `true_label`: 开启状态的文本
- `false_label`: 关闭状态的文本
- `default_value`: 默认值(true/false)

#### DateTime日期时间组件
```go
// 日期选择
StartDate string `runner:"code:start_date;name:开始日期;widget:datetime;format:date;min_date:today;max_date:2025-12-31"`

// 日期时间选择
CreateTime string `runner:"code:create_time;name:创建时间;widget:datetime;format:datetime;default_value:now"`

// 时间选择
WorkTime string `runner:"code:work_time;name:工作时间;widget:datetime;format:time;default_time:09:00:00"`

// 日期范围选择
DateRange string `runner:"code:date_range;name:日期范围;widget:datetime;format:daterange;separator:至;shortcuts:今天,昨天,最近7天,最近30天"`

// 年月选择
Year string `runner:"code:year;name:年份;widget:datetime;format:year;default_value:2024"`
Month string `runner:"code:month;name:月份;widget:datetime;format:month;default_value:2024-01"`

// 周选择
Week string `runner:"code:week;name:周;widget:datetime;format:week"`
```
**支持格式类型：**
- `date`: 日期选择 (YYYY-MM-DD)
- `datetime`: 日期时间选择 (YYYY-MM-DD HH:mm:ss)
- `time`: 时间选择 (HH:mm:ss)
- `daterange`: 日期范围选择
- `datetimerange`: 日期时间范围选择
- `month`: 年月选择 (YYYY-MM)
- `year`: 年份选择 (YYYY)
- `week`: 周选择

**支持参数：**
- `format`: 日期时间格式(必需)
- `default_value`: 默认值，支持特殊值`today`、`now`
- `default_time`: 选中日期后的默认时间
- `min_date`: 最小可选日期
- `max_date`: 最大可选日期
- `separator`: 范围选择的分隔符(默认"至")
- `shortcuts`: 快捷选项(逗号分隔)
- `placeholder`: 占位符文本
- `start_placeholder`: 范围选择开始日期占位符
- `end_placeholder`: 范围选择结束日期占位符

#### Slider滑块组件
```go
// 基础滑块
Progress int `runner:"code:progress;name:进度;widget:slider;min:0;max:100;step:5;unit:%;show_tooltip:true;show_marks:true"`

// 范围滑块
PriceRange string `runner:"code:price_range;name:价格区间;widget:slider;min:0;max:10000;step:100;range:true;default_value:1000,5000;unit:元"`

// 垂直滑块
Volume int `runner:"code:volume;name:音量;widget:slider;min:0;max:100;vertical:true;show_tooltip:true"`
```
**支持参数：**
- `min`: 最小值(默认0)
- `max`: 最大值(默认100)
- `step`: 步长(默认1)
- `default_value`: 默认值，范围选择用逗号分隔
- `unit`: 显示单位
- `show_tooltip`: 是否显示数值提示(默认true)
- `show_marks`: 是否显示刻度标记(默认false)
- `range`: 是否为范围选择(默认false)
- `vertical`: 是否垂直显示(默认false)

#### Color颜色选择器
```go
// 基础颜色选择
BgColor string `runner:"code:bg_color;name:背景色;widget:color;format:hex;default_value:#ffffff;show_alpha:false"`

// 支持透明度的颜色选择
TextColor string `runner:"code:text_color;name:文字颜色;widget:color;format:rgba;show_alpha:true;predefine:#ff0000,#00ff00,#0000ff"`

// 预设颜色选择
ThemeColor string `runner:"code:theme_color;name:主题色;widget:color;predefine:#409EFF,#67C23A,#E6A23C,#F56C6C;show_swatches:true"`
```
**支持格式：**
- `hex`: 十六进制格式 (#ffffff)
- `rgb`: RGB格式 (rgb(255,255,255))
- `rgba`: RGBA格式 (rgba(255,255,255,1))
- `hsl`: HSL格式 (hsl(0,0%,100%))
- `hsla`: HSLA格式 (hsla(0,0%,100%,1))

**支持参数：**
- `format`: 颜色格式(默认hex)
- `default_value`: 默认颜色值
- `show_alpha`: 是否显示透明度控制(默认false)
- `predefine`: 预定义颜色列表(逗号分隔)
- `show_swatches`: 是否显示色板(默认false)
- `allow_empty`: 是否允许清空(默认false)

#### FileUpload文件上传组件
```go
// 基础文件上传
Avatar *files.Files `runner:"code:avatar;name:头像;widget:file_upload;accept:.jpg,.png,.gif;max_size:2048;limit:1;placeholder:请选择头像"`

// 多文件上传
Documents *files.Files `runner:"code:documents;name:文档;widget:file_upload;multiple:true;accept:.pdf,.doc,.docx;max_size:10240;limit:5;drag:true"`

// 图片上传卡片模式
Images *files.Files `runner:"code:images;name:图片;widget:file_upload;multiple:true;accept:.jpg,.png,.gif;list_type:picture-card;auto_upload:false"`
```
**支持参数：**
- `accept`: 接受的文件类型(如:.jpg,.png,.pdf)
- `multiple`: 是否支持多文件上传(默认false)
- `max_size`: 单个文件大小限制(KB)
- `limit`: 文件数量限制
- `placeholder`: 占位符文本
- `auto_upload`: 是否自动上传(默认false)
- `action`: 上传接口地址
- `list_type`: 展示方式(`text`/`picture`/`picture-card`)
- `drag`: 是否支持拖拽上传(默认false)
- `button_text`: 上传按钮文字
- `tip`: 提示文字

#### FileDisplay文件展示组件
```go
// 列表模式文件展示
ProcessedFiles files.Writer `runner:"code:processed_files;name:处理结果;widget:file_display;display_mode:list;preview:true;download:true"`

// 卡片模式文件展示
AttachFiles files.Writer `runner:"code:attach_files;name:附件;widget:file_display;display_mode:card;show_size:true;show_time:true;max_preview:10"`
```
**支持参数：**
- `display_mode`: 展示模式(`list`列表/`card`卡片)
- `preview`: 是否支持预览(默认false)
- `download`: 是否支持下载(默认false)
- `show_size`: 是否显示文件大小(默认false)
- `show_time`: 是否显示上传时间(默认false)
- `show_type`: 是否显示文件类型(默认false)
- `max_preview`: 最大预览数量

### 表格组件（Table）
表格用于展示系统生成的数据列表，支持自动分页、排序、搜索等功能。**现在支持字段级别的组件渲染**！

**正确的表格响应格式：**
```go
// 表格配置中的Response必须使用 query.PaginatedTable[[]数据类型]
var ProductListConfig = &runner.FunctionInfo{
    Request:    &ProductListReq{},
    Response:   query.PaginatedTable[[]Product]{},  // 关键：必须用PaginatedTable包装切片类型
    RenderType: response.RenderTypeTable,
    // ...
}

// 数据模型定义 - 支持字段级别组件配置
type Product struct {
    ID          int        `json:"id" runner:"code:id;name:产品ID;show:list"`
    Name        string     `json:"name" runner:"code:name;name:产品名称;widget:input;placeholder:请输入产品名称;show:list"`
    Category    string     `json:"category" runner:"code:category;name:产品分类;widget:select;options:手机,笔记本,平板,耳机;show:list"`
    Price       float64    `json:"price" runner:"code:price;name:产品价格;widget:slider;min:0;max:50000;step:100;unit:元;show:list"`
    Stock       int        `json:"stock" runner:"code:stock;name:库存数量;widget:slider;min:0;max:1000;step:1;unit:件;show:list"`
    Status      string     `json:"status" runner:"code:status;name:产品状态;widget:switch;true_label:启用;false_label:禁用;show:list"`
    Tags        string     `json:"tags" runner:"code:tags;name:产品标签;widget:multiselect;options:热销,新品,折扣,限量;show:list"`
    CreatedAt   typex.Time `json:"created_at" runner:"code:created_at;name:创建时间;widget:datetime;format:datetime;show:list"`
    UpdatedAt   typex.Time `json:"updated_at" runner:"code:updated_at;name:更新时间;widget:datetime;format:datetime;show:list"`
}

// 请求参数结构体（继承分页功能）
type ProductListReq struct {
    query.PageInfoReq                            // 继承分页参数
    Name     string `form:"name" runner:"code:name;name:产品名称;widget:input;placeholder:按产品名称搜索"`
    Category string `form:"category" runner:"code:category;name:产品分类;widget:select;options:手机,笔记本,平板,耳机"`
    Status   string `form:"status" runner:"code:status;name:状态;widget:select;options:启用,禁用;default_value:启用"`
}
```

### 🎯 表格字段组件渲染规则

**所有表格字段都支持以下组件：**

#### 1. **Input输入框组件** (`widget:input`)
```go
Name string `runner:"code:name;name:产品名称;widget:input;mode:line_text;show:list"`
Desc string `runner:"code:desc;name:描述;widget:input;mode:text_area;show:list"`
```
- 表格中显示为可编辑的文本框
- `mode:line_text` - 单行文本显示
- `mode:text_area` - 多行文本显示（表格中会截断）

#### 2. **Select下拉选择器** (`widget:select`)
```go
Category string `runner:"code:category;name:产品分类;widget:select;options:手机,笔记本,平板,耳机;show:list"`
Status   string `runner:"code:status;name:状态;widget:select;options:启用,禁用;default_value:启用;show:list"`
```
- 表格中显示为下拉选择器
- 支持快速切换状态值
- 支持选项标签显示

#### 3. **Switch开关组件** (`widget:switch`)
```go
Status   string `runner:"code:status;name:产品状态;widget:switch;true_label:启用;false_label:禁用;show:list"`
IsActive bool   `runner:"code:is_active;name:是否激活;widget:switch;true_label:是;false_label:否;show:list"`
```
- 表格中显示为开关按钮
- 支持直接点击切换状态
- 布尔值和字符串都支持

#### 4. **DateTime日期时间组件** (`widget:datetime`)
```go
CreatedAt typex.Time `runner:"code:created_at;name:创建时间;widget:datetime;format:datetime;show:list"`
StartDate string     `runner:"code:start_date;name:开始日期;widget:datetime;format:date;show:list"`
```
- 表格中显示为格式化的日期时间
- 支持日期时间选择器
- 多种格式支持：`date`, `datetime`, `time`

#### 5. **Slider滑块组件** (`widget:slider`)
```go
Price float64 `runner:"code:price;name:价格;widget:slider;min:0;max:50000;step:100;unit:元;show:list"`
Stock int     `runner:"code:stock;name:库存;widget:slider;min:0;max:1000;step:1;unit:件;show:list"`
```
- 表格中显示为可拖拽的滑块
- 支持数值范围限制
- 直观的数值调整

#### 6. **MultiSelect多选组件** (`widget:multiselect`)
```go
Tags string `runner:"code:tags;name:标签;widget:multiselect;options:热销,新品,折扣,限量;allow_create:true;show:list"`
```
- 表格中显示为多选下拉框
- 支持标签形式展示已选项
- 支持动态添加选项

#### 7. **Color颜色选择器** (`widget:color`)
```go
ThemeColor string `runner:"code:theme_color;name:主题色;widget:color;format:hex;show_alpha:true;show:list"`
```
- 表格中显示为颜色块
- 点击弹出颜色选择器
- 支持多种颜色格式

#### 8. **FileDisplay文件展示** (`widget:file_display`)
```go
AttachFiles files.Writer `runner:"code:attach_files;name:附件;widget:file_display;display_mode:list;preview:true;show:list"`
```
- 表格中显示为文件列表
- 支持预览和下载操作
- 显示文件信息摘要

### 📊 表格字段渲染效果预览

| 字段类型 | 组件类型 | 表格中显示效果 | 交互能力 |
|---------|---------|--------------|---------|
| 文本字段 | `input` | 可编辑文本框 | ✅ 直接编辑 |
| 状态字段 | `select` | 下拉选择器 | ✅ 快速切换 |
| 布尔字段 | `switch` | 开关按钮 | ✅ 点击切换 |
| 日期字段 | `datetime` | 格式化日期 | ✅ 日期选择 |
| 数值字段 | `slider` | 滑块控制器 | ✅ 拖拽调整 |
| 标签字段 | `multiselect` | 标签列表 | ✅ 多选编辑 |
| 颜色字段 | `color` | 颜色块 | ✅ 颜色选择 |
| 文件字段 | `file_display` | 文件摘要 | ✅ 预览下载 |

## 🔄 回调机制系统

### 回调类型定义
```go
type FunctionInfo struct {
    // 页面加载回调 - 页面初始化时触发
    OnPageLoad func(ctx *Context, resp response.Response) (*OnPageLoadResp, error)
    
    // 输入模糊搜索回调 - 用户输入时实时触发
    OnInputFuzzyMap map[string]OnInputFuzzy
    
    // API创建后回调 - 用于初始化数据
    OnApiCreated func(ctx *Context, req *OnApiCreatedReq) error
}
```

### 输入模糊搜索回调示例
```go
// 1. 在字段标签中配置回调
Location string `runner:"..." callback:"OnInputFuzzy(field:location_suggest,delay:300,min:1)"`

// 2. 在FunctionInfo中实现回调逻辑
OnInputFuzzyMap: map[string]runner.OnInputFuzzy{
    "location_suggest": func(ctx *Context, req *OnInputFuzzyReq) (*OnInputFuzzyResp, error) {
        // 获取用户输入值
        inputValue := req.Value
        
        // 执行模糊搜索逻辑（这里是简化示例）
        locations := []string{"北京", "上海", "深圳", "广州", "杭州", "成都"}
        
        // 构建响应结果
        resp := &OnInputFuzzyResp{}
        for _, location := range locations {
            if strings.Contains(location, inputValue) {
                resp.Values = append(resp.Values, &InputFuzzyItem{
                    Value: location,
                })
            }
        }
        return resp, nil
    },
}
```

### 页面加载回调示例
```go
OnPageLoad: func(ctx *Context, resp response.Response) (*OnPageLoadResp, error) {
    // 设置初始响应数据
    resp.Form(&DemoResp{
        Message: "页面加载成功",
    }).Build()
    
    // 返回表单初始值
    return &OnPageLoadResp{
        Request: DemoReq{
            Name:     "默认用户",
            Category: "默认分类",
            Active:   true,
        },
    }, nil
}
```

## 📋 最佳实践示例解析

### 1. 回调演示 (`callback_demo.go`)
**展示功能：** 输入模糊搜索、页面加载回调、数据联动

```go
type CallbackDemoReq struct {
    // 支持模糊搜索的地点字段
    Location string `runner:"code:location;name:工作地点;widget:input" callback:"OnInputFuzzy(field:location_suggest,delay:300,min:1)"`
    
    // 支持模糊搜索的公司字段  
    CompanySearch string `runner:"code:company_search;name:公司搜索;widget:input" callback:"OnInputFuzzy(field:company_search,delay:500,min:2)"`
}
```

### 2. 工作流演示 (`workflow_demo.go`)
**展示功能：** 复杂表单、业务逻辑处理、智能建议生成

```go
type WorkflowDemoReq struct {
    // 优先级选择
    Priority string `runner:"code:priority;name:优先级;widget:select;options:低,中,高,紧急;default_value:中"`
    
    // 任务状态
    Status string `runner:"code:status;name:状态;widget:radio;options:待办,进行中,已完成;direction:horizontal"`
    
    // 进度滑块
    Progress int `runner:"code:progress;name:进度;widget:slider;min:0;max:100;unit:%;show_tooltip:true"`
}
```

### 3. CRUD演示 (`crud_demo.go`)
**展示功能：** 自动化CRUD、表格展示、分页筛选

```go
var ProductListConfig = &runner.FunctionInfo{
    Request:       &ProductListReq{},
    Response:      query.PaginatedTable[[]Product]{}, // 关键：正确的表格响应格式
    RenderType:    response.RenderTypeTable,
    AutoCrudTable: &Product{}, // 启用自动CRUD功能
    CreateTables:  []interface{}{&Product{}},
    // ...
}

type ProductListReq struct {
    query.PageInfoReq                            // 继承分页功能
    Name     string `form:"name" runner:"code:name;name:产品名称;placeholder:按产品名称搜索"`
    Category string `form:"category" runner:"code:category;name:产品分类;widget:select;options:手机,笔记本,平板,耳机"`
    Status   string `form:"status" runner:"code:status;name:状态;widget:select;options:启用,禁用"`
    MinPrice float64 `form:"min_price" runner:"code:min_price;name:最低价格;widget:number;min:0"`
    MaxPrice float64 `form:"max_price" runner:"code:max_price;name:最高价格;widget:number;min:0"`
}
```

### 4. 数据分析演示 (`data_analysis.go`)
**展示功能：** 复杂数据筛选、图表配置、分析结果展示

```go
type DataAnalysisReq struct {
    // 时间范围选择
    StartDate string `runner:"code:start_date;name:开始日期;widget:datetime;format:date"`
    EndDate   string `runner:"code:end_date;name:结束日期;widget:datetime;format:date"`
    
    // 分析类型
    AnalysisType string `runner:"code:analysis_type;name:分析类型;widget:select;options:汇总分析,趋势分析,对比分析"`
    
    // 图表类型
    ChartType string `runner:"code:chart_type;name:图表类型;widget:radio;options:柱状图,折线图,饼图"`
}
```

### 5. 文件管理演示 (`file_manager.go`)
**展示功能：** 文件上传、多种处理模式、结果展示

```go
type FileManagerReq struct {
    // 文件上传
    Files *files.Files `runner:"code:files;name:文件列表;widget:file_upload;multiple:true;max_size:10240"`
    
    // 处理模式
    ProcessMode string `runner:"code:process_mode;name:处理模式;widget:radio;options:仅上传,分析文件,提取内容,压缩处理"`
    
    // 高级选项
    AutoExtract bool `runner:"code:auto_extract;name:自动提取;widget:switch;true_label:开启;false_label:关闭"`
}
```

## 🔄 完整数据流转过程

### 1. 组件生成流程
```
Go结构体定义 → Runner标签解析 → Widget组件创建 → API配置生成 → 前端界面渲染
```

### 2. 用户交互流程
```
用户操作 → 前端事件 → 回调函数执行 → 响应数据返回 → 界面更新
```

### 3. 表单提交流程  
```
表单数据 → 参数验证 → 业务逻辑处理 → 数据库操作 → 结果响应
```

## 🚀 系统扩展指南

### 添加新Widget组件

1. **定义组件类型** (`type.go`)
```go
const WidgetNewComponent = "new_component"
```

2. **实现组件结构** (`new_component.go`)
```go
type NewComponentWidget struct {
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

### 添加新回调类型

1. **定义回调函数类型**
```go
type OnNewCallback func(ctx *Context, req *OnNewCallbackReq) (*OnNewCallbackResp, error)
```

2. **在FunctionInfo中添加**
```go
type FunctionInfo struct {
    OnNewCallbackMap map[string]OnNewCallback
}
```

3. **具体实现使用**
```go
OnNewCallbackMap: map[string]OnNewCallback{
    "my_callback": func(ctx *Context, req *OnNewCallbackReq) (*OnNewCallbackResp, error) {
        // 回调逻辑实现
        return &OnNewCallbackResp{}, nil
    },
}
```

## 🎯 核心设计原则

1. **配置驱动** - 通过标签配置而非硬编码实现UI生成
2. **组件化架构** - 每个UI元素都是独立可复用的Widget组件  
3. **类型安全** - 基于Go强类型系统，编译期检查错误
4. **易于扩展** - 工厂模式设计，新增组件类型非常简单
5. **关注分离** - UI逻辑与业务逻辑完全分离
6. **自动化优先** - 最大程度减少手动配置，提高开发效率

通过这套架构，开发者只需定义Go结构体和业务逻辑，UI界面由系统自动生成，真正实现了"代码即配置"的无代码开发体验。 