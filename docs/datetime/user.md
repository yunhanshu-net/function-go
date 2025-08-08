# 日期时间组件 (datetime)

## 概述

日期时间组件用于处理各种时间相关的输入，支持日期选择、时间选择、日期时间选择等多种模式。

## 主要特性

- **多种格式**：支持日期、时间、日期时间、日期范围等多种格式
- **时间戳存储**：使用 int64 毫秒时间戳存储，便于数据库操作
- **灵活显示**：通过 format 标签控制显示格式，节省表格空间
- **智能默认值**：支持 $today、$now 等特殊默认值
- **范围限制**：支持设置最小和最大可选日期
- **快捷选项**：支持预设的快捷选择选项

## 使用场景

### 1. 业务日期字段
适用于开始日期、截止日期、生日等业务日期，在表格中只显示日期部分：

```go
type TaskReq struct {
    StartDate int64 `json:"start_date" runner:"code:start_date;name:开始日期" widget:"type:datetime;kind:date;format:yyyy-MM-dd;placeholder:请选择开始日期" data:"type:number;default_value:$today" validate:"required"`
    DueDate int64 `json:"due_date" runner:"code:due_date;name:截止日期" widget:"type:datetime;kind:date;format:yyyy-MM-dd;placeholder:请选择截止日期" data:"type:number;default_value:$7_days_later" validate:"required"`
    Birthday int64 `json:"birthday" runner:"code:birthday;name:生日" widget:"type:datetime;kind:date;format:yyyy-MM-dd;placeholder:请选择生日" data:"type:number"`
}
```

### 2. 时间字段
适用于下课时间、上班时间等只需要显示时间的场景：

```go
type ScheduleReq struct {
    ClassEndTime int64 `json:"class_end_time" runner:"code:class_end_time;name:下课时间" widget:"type:datetime;kind:time;format:HH:mm;placeholder:请选择结束时间" data:"type:number"`
    WorkStartTime int64 `json:"work_start_time" runner:"code:work_start_time;name:上班时间" widget:"type:datetime;kind:time;format:HH:mm;placeholder:请选择上班时间" data:"type:number"`
}
```

### 3. 会议时间字段
适用于会议时间、预约时间等需要显示日期和时间的场景：

```go
type MeetingReq struct {
    MeetingTime int64 `json:"meeting_time" runner:"code:meeting_time;name:会议时间" widget:"type:datetime;kind:datetime;placeholder:请选择会议时间" data:"type:number"`
}
```

### 4. 系统时间字段
适用于创建时间、更新时间等系统自动生成的时间字段：

```go
type UserReq struct {
    CreatedAt int64 `json:"created_at" runner:"code:created_at;name:创建时间" widget:"type:datetime;kind:datetime" data:"type:number" permission:"read"`
    UpdatedAt int64 `json:"updated_at" runner:"code:updated_at;name:更新时间" widget:"type:datetime;kind:datetime" data:"type:number" permission:"read"`
}
```

## 参数说明

### 基础参数
- **type:datetime** (必需)：标识为日期时间组件
- **kind** (可选)：具体格式类型
  - `date`：日期选择
  - `time`：时间选择
  - `datetime`：日期时间选择
  - `daterange`：日期范围选择
  - `datetimerange`：日期时间范围选择
  - `month`：月份选择
  - `year`：年份选择
  - `week`：周选择

### 显示格式参数
- **format** (可选)：显示格式
  - `yyyy-MM-dd`：仅显示日期，如 2025-01-15
  - `HH:mm`：仅显示时间，如 18:30
  - `HH:mm:ss`：显示时间（含秒），如 18:30:25
  - `yyyy-MM-dd HH:mm`：日期时间（不含秒），如 2025-01-15 18:30
  - `yyyy-MM-dd HH:mm:ss`：完整日期时间，如 2025-01-15 18:30:25

### 交互参数
- **placeholder** (可选)：占位符文本
- **start_placeholder** (可选)：范围选择时开始日期占位符
- **end_placeholder** (可选)：范围选择时结束日期占位符
- **disabled** (可选)：是否禁用

### 默认值参数
- **default_value** (可选)：默认值
  - `$today`：今天
  - `$now`：当前时间
  - `$7_days_later`：7天后
  - 具体时间戳：如 1705292200000

### 限制参数
- **min_date** (可选)：最小可选日期
- **max_date** (可选)：最大可选日期

### 其他参数
- **separator** (可选)：日期范围分隔符，默认"至"
- **shortcuts** (可选)：快捷选项配置
- **default_time** (可选)：选中日期后的默认具体时刻

## 最佳实践

### 1. 时间戳存储
所有时间字段都使用 int64 毫秒时间戳存储：

```go
// ✅ 正确：使用时间戳
type TaskReq struct {
    StartDate int64 `json:"start_date" widget:"type:datetime;kind:date" data:"type:number"`
}

// ❌ 错误：使用字符串
type TaskReq struct {
    StartDate string `json:"start_date" widget:"type:datetime;kind:date" data:"type:string"`
}
```

### 2. 格式选择
根据业务需求选择合适的显示格式：

```go
// 业务日期：只显示日期，节省表格空间
DueDate int64 `widget:"type:datetime;kind:date;format:yyyy-MM-dd"`

// 时间：只显示时间
EndTime int64 `widget:"type:datetime;kind:time;format:HH:mm"`

// 系统时间：显示完整格式
CreatedAt int64 `widget:"type:datetime;kind:datetime"`
```

### 3. 默认值设置
使用智能默认值提高用户体验：

```go
// 开始日期默认为今天
StartDate int64 `widget:"type:datetime;kind:date;default_value:$today"`

// 截止日期默认为7天后
DueDate int64 `widget:"type:datetime;kind:date;default_value:$7_days_later"`
```

### 4. 权限控制
系统时间字段通常设置为只读：

```go
// 创建时间只读
CreatedAt int64 `widget:"type:datetime;kind:datetime" permission:"read"`

// 更新时间只读
UpdatedAt int64 `widget:"type:datetime;kind:datetime" permission:"read"`
```

## 注意事项

1. **时间戳格式**：所有时间字段使用 int64 毫秒时间戳存储
2. **GORM标签**：系统时间字段使用 `autoCreateTime:milli` 和 `autoUpdateTime:milli`
3. **format标签**：仅在业务日期字段中使用，系统时间字段保持默认格式
4. **权限控制**：系统时间字段通常设置为只读权限
5. **默认值**：使用 $today、$now 等特殊值提高用户体验 