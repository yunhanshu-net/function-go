## 时间字段（仅时间场景）统一存储策略

目标：统一“仅时间”场景的前后端协议与存储，确保排序/筛选一致并便于扩展。

### 约定
- 所有时间相关字段（含仅时间的场景）统一使用毫秒时间戳（int64）在 API 与数据库中传递/存储。
- 前端在仅时间场景（如上课时间 08:30）也提交一个时间戳（固定基准日期，如 1970-01-01 的本地时区），后端照常存储为毫秒时间戳。
- 前端渲染时按照组件的 `format` 与 `kind` 进行格式化展示（如 `kind:time; format:HH:mm`），与日期无关。

### 原因
- 统一类型：后端/数据库统一 int64 毫秒，避免 string/时间戳混用导致的歧义。
- 可排序：即便是“仅时间”，也能直接按毫秒排序/筛选（固定基准日保证排序稳定）。
- 易扩展：未来需要跨天区间或与日期拼接时无需做类型转换。

### 标签与示例
- 仅时间的配置字段（推荐写法）：
  ```go
  // 仅时间（统一使用毫秒时间戳），前端会提交固定日期的时间戳，例如 1970-01-01 08:30
  ClassStartAt int64 `json:"class_start_at" runner:"code:class_start_at;name:上课时间" widget:"type:datetime;kind:time;format:HH:mm"`
  ```

- 仅时间的请求字段默认值：
  ```go
  ArriveAt int64 `json:"arrive_at" runner:"code:arrive_at;name:到校时间" widget:"type:datetime;kind:datetime;disabled:true" data:"default_value:$now"`
  ```

- GET 绑定（如使用 GET）：字段需补 `form:"field_name"`，依现有规范执行。

### 后端处理建议
- 仅时间配置与日期字段组合：
  1) 配置存储为毫秒时间戳 `class_start_at`（基准日任意，如 1970-01-01 08:30）。
  2) 业务需要与“当天日期”组合时，从 `class_start_at` 中解析出 时/分，再与当天日期拼接生成当天的时间点用于计算。
  ```go
  base := time.UnixMilli(classStartAt)
  hour, minute := base.Hour(), base.Minute()
  classStartToday := time.Date(day.Year(), day.Month(), day.Day(), hour, minute, 0, 0, day.Location())
  ```

### 与现有代码的关系
- 现有示例中如使用了 string（HH:mm）作为配置，不影响功能；后续迭代可平滑切换为毫秒时间戳以与统一策略保持一致。
- 前端已支持将毫秒时间戳按 `format` 渲染为仅时间文本（HH:mm）。


