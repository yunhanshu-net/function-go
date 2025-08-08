# 日期时间组件规则说明

## 标签参数详解（按标签类型分组）

### widget 标签参数
- **type:datetime** (必需)：标识为日期时间组件
- **kind** (可选)：具体格式类型，枚举值：date, datetime, time, daterange, datetimerange, month, year, week
- **format** (可选)：显示格式，如 yyyy-MM-dd, HH:mm, yyyy-MM-dd HH:mm:ss
- **placeholder** (可选)：占位符文本
- **start_placeholder** (可选)：范围选择时开始日期占位符
- **end_placeholder** (可选)：范围选择时结束日期占位符
- **default_value** (可选)：默认值，支持 $today, $now, $7_days_later 等特殊值
- **default_time** (可选)：选中日期后的默认具体时刻
- **min_date** (可选)：最小可选日期
- **max_date** (可选)：最大可选日期
- **separator** (可选)：日期范围分隔符，默认"至"
- **shortcuts** (可选)：快捷选项配置
- **disabled** (可选)：是否禁用

### data 标签参数
- **type:number** (必需)：标识字段为数字类型（时间戳）
- **default_value** (可选)：设置默认值，如 $today, $now, $7_days_later
- **example** (可选)：提供示例值，如 1705292200000

### validate 标签参数
- **required** (可选)：必填验证
- **min** (可选)：最小值验证，格式：min=时间戳值
- **max** (可选)：最大值验证，格式：max=时间戳值

## 标签规则
1. **widget 标签**：必须包含 `type:datetime`，可选参数如上
2. **data 标签**：必须包含 `type:number`（时间戳），可包含 `default_value` 和 `example`
3. **validate 标签**：可包含数值验证规则

## 正例（完整结构体定义）
```go
// 工作流任务配置
type WorkflowTaskConfig struct {
    StartDate int64 `json:"start_date" runner:"code:start_date;name:开始日期" widget:"type:datetime;kind:date;format:yyyy-MM-dd;placeholder:请选择开始日期" data:"type:number;default_value:$today;example:1705292200000" validate:"required"`
    DueDate int64 `json:"due_date" runner:"code:due_date;name:截止日期" widget:"type:datetime;kind:date;format:yyyy-MM-dd;placeholder:请选择截止日期" data:"type:number;default_value:$7_days_later;example:1705897600000" validate:"required"`
    Birthday int64 `json:"birthday" runner:"code:birthday;name:生日" widget:"type:datetime;kind:date;format:yyyy-MM-dd;placeholder:请选择生日" data:"type:number;example:694224000000"`
    ClassEndTime int64 `json:"class_end_time" runner:"code:class_end_time;name:下课时间" widget:"type:datetime;kind:time;format:HH:mm;placeholder:请选择结束时间" data:"type:number;example:1705309800000"`
    MeetingTime int64 `json:"meeting_time" runner:"code:meeting_time;name:会议时间" widget:"type:datetime;kind:datetime;placeholder:请选择会议时间" data:"type:number;example:1705467600000"`
    CreatedAt int64 `json:"created_at" runner:"code:created_at;name:创建时间" widget:"type:datetime;kind:datetime" data:"type:number;example:1705292200000" permission:"read"`
}
```

## 反例（完整结构体定义）
```go
// 错误：使用字符串类型存储时间
type DateTimeErrorConfig struct {
    StartDate string `json:"start_date" runner:"code:start_date;name:开始日期" widget:"type:datetime;kind:date" data:"type:string;example:2025-01-15" validate:"required"`
}

// 错误：format标签使用错误
type FormatErrorConfig struct {
    DueDate int64 `json:"due_date" runner:"code:due_date;name:截止日期" widget:"type:datetime;kind:date;format:invalid-format" data:"type:number;example:1705897600000" validate:"required"`
}

// 错误：default_value与kind不匹配
type DefaultValueErrorConfig struct {
    EndTime int64 `json:"end_time" runner:"code:end_time;name:结束时间" widget:"type:datetime;kind:time;default_value:$today" data:"type:number;example:1705309800000"`
}
```

## 注意事项
1. **时间戳存储**：所有时间字段必须使用 int64 类型存储毫秒时间戳
2. **GORM标签**：系统时间字段使用 `autoCreateTime:milli` 和 `autoUpdateTime:milli`
3. **format标签**：仅在业务日期字段中使用，系统时间字段保持默认格式
4. **权限控制**：系统时间字段通常设置为 `permission:"read"`
5. **默认值**：使用 $today、$now、$7_days_later 等特殊值
6. **显示格式**：
   - `format:yyyy-MM-dd`：仅显示日期，适用于业务日期字段
   - `format:HH:mm`：仅显示时间，适用于下课时间、上班时间等
   - `format:HH:mm:ss`：显示时间（含秒），适用于精确时间记录
   - `format:yyyy-MM-dd HH:mm`：日期时间（不含秒），适用于会议时间等
   - 无 format 标签：显示完整日期时间，适用于系统时间字段

## 使用场景分类

### 业务日期字段
- **开始日期**：`widget:"type:datetime;kind:date;format:yyyy-MM-dd"`
- **截止日期**：`widget:"type:datetime;kind:date;format:yyyy-MM-dd"`
- **生日**：`widget:"type:datetime;kind:date;format:yyyy-MM-dd"`

### 时间字段
- **下课时间**：`widget:"type:datetime;kind:time;format:HH:mm"`
- **上班时间**：`widget:"type:datetime;kind:time;format:HH:mm"`
- **精确时间**：`widget:"type:datetime;kind:time;format:HH:mm:ss"`

### 会议时间字段
- **会议时间**：`widget:"type:datetime;kind:datetime"`
- **预约时间**：`widget:"type:datetime;kind:datetime"`

### 系统时间字段
- **创建时间**：`widget:"type:datetime;kind:datetime" permission:"read"`
- **更新时间**：`widget:"type:datetime;kind:datetime" permission:"read"`
- **完成时间**：`widget:"type:datetime;kind:datetime" permission:"read"`

## 最佳实践

### 1. 时间戳格式
```go
// ✅ 正确：使用 int64 时间戳
type TaskReq struct {
    StartDate int64 `json:"start_date" widget:"type:datetime;kind:date" data:"type:number"`
}

// ❌ 错误：使用字符串
type TaskReq struct {
    StartDate string `json:"start_date" widget:"type:datetime;kind:date" data:"type:string"`
}
```

### 2. 格式选择
```go
// 业务日期：只显示日期，节省表格空间
DueDate int64 `widget:"type:datetime;kind:date;format:yyyy-MM-dd"`

// 时间：只显示时间
EndTime int64 `widget:"type:datetime;kind:time;format:HH:mm"`

// 系统时间：显示完整格式
CreatedAt int64 `widget:"type:datetime;kind:datetime"`
```

### 3. 默认值设置
```go
// 开始日期默认为今天
StartDate int64 `widget:"type:datetime;kind:date;default_value:$today"`

// 截止日期默认为7天后
DueDate int64 `widget:"type:datetime;kind:date;default_value:$7_days_later"`

// 当前时间
MeetingTime int64 `widget:"type:datetime;kind:datetime;default_value:$now"`
```

### 4. 权限控制
```go
// 创建时间只读
CreatedAt int64 `widget:"type:datetime;kind:datetime" permission:"read"`

// 更新时间只读
UpdatedAt int64 `widget:"type:datetime;kind:datetime" permission:"read"`
```

## 常见错误

### 1. 数据类型错误
```go
// ❌ 错误：使用字符串存储时间
StartDate string `widget:"type:datetime;kind:date" data:"type:string"`
```

### 2. 格式标签错误
```go
// ❌ 错误：format标签格式不正确
DueDate int64 `widget:"type:datetime;kind:date;format:invalid-format"`
```

### 3. 默认值错误
```go
// ❌ 错误：默认值与kind不匹配
EndTime int64 `widget:"type:datetime;kind:time;default_value:$today"`
```

### 4. 权限设置错误
```go
// ❌ 错误：系统时间字段未设置只读权限
CreatedAt int64 `widget:"type:datetime;kind:datetime"` // 应该设置 permission:"read"
``` 