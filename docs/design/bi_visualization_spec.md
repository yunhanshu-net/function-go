## 通用 BI 可视化与供数规范（提案）

目标：提供与 form/table 并列的一类“BI 函数”，用于通用的数据可视化供数与渲染，不局限于 NPS。强调协议通用、逐步可演进（缓存/预聚合）、与现有框架自然融合。

### 1. 背景与诉求
- 现有 form/table 已覆盖表单与列表场景，但缺少“图表/看板”类型的标准输出。
- 希望一个协议即可覆盖多业务（NPS、满意度、访问分析、库存看板等），并能渐进增强（缓存、预聚合、联动）。

### 2. 设计原则
- 通用抽象：维度/指标/筛选/时间粒度 → 图表/表格/KPI。
- 开闭扩展：新增图表类型/编码字段不破坏既有协议。
- 低耦合：与业务结构体解耦，统一渲染数据结构 BIResult。
- 易落地：先以 form 返回 BIResult，后续引入 BIFunctionOptions 与 render_type: bi。

### 3. 核心抽象

#### 3.1 BIFunctionOptions（选项，提案）
- 作用：声明 BI 渲染函数，配置默认图表、缓存、预聚合策略。
- 字段（建议）：
  - BaseConfig（与现一致）
  - RenderType: "bi"
  - AllowedChartTypes: []string（kpi/line/bar/stacked_bar/pie/hist/table）
  - Defaults: {time_grain, dimensions, measures}
  - CacheTTL: int（秒）
  - PreAggregate: {enabled, table, grain, keys, measures}
  - Security: {max_rows, allow_export}

注：在未引入新 Options 之前，可用 FormFunctionOptions 返回 BIResult 先跑通。

#### 3.2 通用请求（建议）
```json
{
  "dimensions": ["date", "channel"],
  "measures": ["nps", "total", "promoters", "detractors"],
  "filters": {"project_id": 1, "start_at": 1722787200000, "end_at": 1754928000000, "channel": "全部"},
  "time_grain": "day"  // none/day/week/month
}
```

#### 3.3 BIResult（统一响应）
```json
{
  "kpis": [
    {"name": "NPS", "value": 25, "unit": "", "trend": {"wow": 5.2}}
  ],
  "charts": [
    {
      "type": "line",
      "title": "NPS 趋势",
      "encoding": {"x": "date", "y": "nps"},
      "data": [{"date": "2025-08-01", "nps": 18}, {"date": "2025-08-02", "nps": 25}],
      "options": {"smooth": true, "unit": ""}
    },
    {
      "type": "stacked_bar",
      "title": "人群构成",
      "encoding": {"x": "date", "y": "count", "series": "segment", "stack": true, "color": "segment"},
      "data": [
        {"date": "2025-08-01", "segment": "Promoters", "count": 10},
        {"date": "2025-08-01", "segment": "Passives", "count": 3},
        {"date": "2025-08-01", "segment": "Detractors", "count": 2}
      ]
    }
  ],
  "tables": [
    {"title": "Top 关键词", "columns": ["keyword", "count"], "rows": [{"keyword": "物流", "count": 26}]}
  ],
  "meta": {"query_time_ms": 12, "sample_size": 100, "note": "样本不足30仅供参考"}
}
```

Chart.type 建议集：kpi, line, bar, stacked_bar, pie, hist, table。encoding 字段用于声明 x/y/series/stack/color/y2 等映射。

### 4. 与 function-go 的集成
- 短期：使用 FormFunctionOptions 返回 BIResult，前端用轻量渲染器展示 KPI/折线/堆叠柱。
- 中期：新增 BIFunctionOptions，`BaseConfig.RenderType = "bi"`；前端内置 BI 渲染组件（ECharts/G2Plot）。
- 仍保留 CreateTables、OnInputFuzzy、OnInputValidate 等既有能力。

### 5. 性能与演进
- 缓存：对相同 filters+time_grain 的结果做秒级/分钟级缓存。
- 预聚合：
  - 日级聚合表 `survey_nps_daily(project_id, date, channel, total, promoters, passives, detractors, score_counts_json, nps)`。
  - 更新策略：明细表实时写入；定时任务/异步队列合并增量到 daily 表；查询优先日表 + 增量补差。
- 索引建议（以 NPS 为例）：
  - 保留 `(project_id, created_at)` 复合索引，覆盖“项目+时间”的主路径。
  - score/channel 索引按压测再定，避免写放大。

### 6. 安全与上限
- 限制最大返回行数与数据量，防止前端渲染阻塞。
- 控制导出（CSV/PNG）权限与开关。
- 样本过小提示（< 30）与口径注释。

### 7. 迁移路线
1) 用现有 form 返回 BIResult（快速上线）。
2) 加入渲染组件与 RenderType: bi，实现配置化可视化。
3) 引入 CacheTTL，再按项目量级逐步落地日级预聚合。
4) 扩展到更多业务看板，验证通用性。

### 8. 示例映射（以 NPS 为例）
- summary → kpis（nps / total / promoters / detractors）
- timeseries(day) → line(nps) + stacked_bar(三类人群)
- distribution(0~10) → hist（分数桶）
- breakdown(channel/region) → bar/stacked_bar（维度拆解）
- comments → table/词云（可选）

### 9. 未来扩展
- 图表联动、高亮与刷选；
- 阈值告警（连续 N 天跌落）；
- 主题聚类与诊断；
- BI 布局编排（多个 charts 组合成仪表盘）。

---

工程反馈（现状简评）：
- 现有工程规范清晰：单文件系统、命名/表名/路径一致，回调/建表机制完善；日志携带 trace_id，排查链路直观。
- 统计口径与参数绑定（POST/GET+form）已固化，默认值（如 channel）贴合常用场景。
- 建议引入本规范后，前端补一个通用图表渲染器，即可在多业务快速复用，向“全自动 AI 应用构建”迈进一步。


