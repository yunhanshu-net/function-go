## 通用能力扩展方案（不改框架，平台侧先落地）

目标：在不修改框架核心的前提下，于平台侧先提供“定时任务、导入导出、通知、配置管理”等通用能力，兼容现有 form/table 模式，后续再逐步沉淀为 Options 标准。

---

### 1) 定时/异步任务（平台平替 JobFunctionOptions）

现状与思路：
- 平台在运行时提供“定时运行”按钮，选择一次性时间或 Cron 表达式；保存一个“计划 + 固定请求体”的实例。
- 到时由平台调度直接调用既有 form API（POST），请求体即用户在 UI 填好的表单值。无需后端新增 Job API。

建议的计划实体（平台侧实现）：
- 字段：id、name、api_path、http_method（POST）、fixed_payload（JSON）、cron/once_time、enabled、retry_policy、timeout_sec、notify_channel、created_at/updated_at。
- 运行记录：计划id、trace_id、start/end、status、error、cost_ms、stdout/brief。

调度保障：
- 幂等：平台对同一触发窗口加“互斥锁/去重键”。
- 失败策略：重试N次、进入DLQ（死信列表）人工重放。
- 并发：同一计划串行（默认），不同计划并行。

后续沉淀（可选）：定义 `JobFunctionOptions` 仅作为“声明”，实际调度仍由平台统一执行。

---

### 2) 数据导入/导出（平台平替 ImportExportFunctionOptions）

导出（平台前端即可完成）：
- table 函数返回的分页/列表 JSON，前端一键导出为 CSV/Excel（XLSX），支持选择列、筛选条件同步写入导出注记。

导入（平台前端即可完成）：
- 前端解析 Excel/CSV，提供“表头 → 字段”映射 UI，完成行级校验（必填/枚举/格式），合格行批量调用 form API/AutoCrud 新增；失败行回显并支持下载错误报告。

何时需要后端解析：
- 超大文件（> 10MB/几十万行）、复杂模板、需要服务端校验/权限/审计；再考虑后端文件上传 + 解析。

文件型业务（转换/生成类）：
- 继续使用 `*files.Files` 与 `gorm:"type:json"` 存储元数据（按既定规范）。

---

### 3) 通知/消息发送（后续可集成到 ctx）

平台先行：
- 计划执行结果、导入导出完成、阈值告警等，由平台统一发送通知（站内/邮件/IM），无需后端修改。

后续提案（不急改）：
- 在 runner.Context 上提供统一门面（示例）
  - `ctx.Notify().Email/SMS/IM(...)
  - `ctx.Notify().WithTemplate("tpl_code").Send(params)`
- 配置在平台管理：通道秘钥、模板、频控；后端仅调用统一门面。

---

### 4) 配置/特性开关（已具备）

现状：
- 已通过 `AutoUpdateConfig` 提供配置管理（渲染为表单，前端更新后，后端值实时更新）。
- 示例参考：`soft/beiluo/v115/code/api/widgets/config_demo.go`（`ConfigDemoOption.AutoUpdateConfig`）。

最佳实践：
- 区分环境/租户；敏感项（密钥）加密存储；变更审计与回滚；必要时在处理函数内加运行期校验（如生产环境禁用 debug）。

可进一步沉淀为 `ConfigFeatureFlagFunctionOptions`（声明式），但当前能力已满足大部分需求。

---

### 5) 与 BI 的协同
- 定时任务可用于构建“日级预聚合”，服务 BI 查询；导入导出作为数据通道；通知用于结果推送；配置管理承载阈值/开关。
- 具体 BI 协议见：`docs/design/bi_visualization_spec.md`。

---

### 6) 总结与演进
- 不改后端框架的前提下，平台即可提供：定时调度、导入导出、通知、配置管理，快速覆盖 80% 通用需求。
- 待业务沉淀后，再将调度/导入导出/通知抽象为 Options（如 `JobFunctionOptions`、`ImportExportFunctionOptions`、`NotificationFunctionOptions`），由后端提供声明，平台读取声明增强 UI/能力。


