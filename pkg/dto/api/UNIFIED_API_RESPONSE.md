# 统一API响应结构设计方案

## 问题背景

当前系统存在多套不同的响应结构：
1. **Table Response** 使用 `AddFormConfig` 字段
2. **Form Request** 使用 `FormRequestParamInfo` 结构  
3. **新的Form Builder** 使用 `FieldInfo` 结构

这导致了结构不一致，维护困难的问题。

## 解决方案

### 1. 统一数据结构

使用 `FieldInfo` 作为统一的字段描述结构：

```go
type FieldInfo struct {
    // 基础信息
    Code string `json:"code"` // 字段代码
    Name string `json:"name"` // 字段显示名称
    Desc string `json:"desc"` // 字段描述

    // Widget配置
    Widget WidgetConfig `json:"widget"`

    // 数据配置
    Data DataConfig `json:"data"`

    // 权限配置
    Permission *PermissionConfig `json:"permission"`

    // 回调配置
    Callbacks []CallbackConfig `json:"callbacks"`

    // 验证配置
    Validation string `json:"validation"`
}
```

### 2. 统一API响应

```go
type UnifiedAPIResponse struct {
    RenderType       string       `json:"render_type"`       // form, table
    Fields           []*FieldInfo `json:"fields"`            // 表单字段
    Columns          []*FieldInfo `json:"columns"`           // 表格列
    SearchConditions []string     `json:"search_conditions"` // 搜索条件
}
```

### 3. 响应结构对比

#### 旧的Table响应（问题）
```json
{
    "code": "updated_at",
    "name": "更新时间",
    "value_type": "object",
    "widget_type": "input",
    "widget_config": {
        "mode": "line_text",
        "placeholder": "",
        "default_value": ""
    },
    "add_form_config": {
        "code": "updated_at",
        "desc": "",
        "name": "更新时间",
        "show": "",
        "hidden": "",
        "example": "",
        "required": false,
        "callbacks": "",
        "validates": "",
        "value_type": "object",
        "widget_type": "input",
        "default_value": "",
        "widget_config": {
            "mode": "line_text",
            "placeholder": "",
            "default_value": ""
        }
    }
}
```

#### 新的统一响应（解决方案）
```json
{
    "render_type": "table",
    "columns": [
        {
            "code": "updated_at",
            "name": "更新时间",
            "desc": "",
            "widget": {
                "type": "datetime",
                "config": {
                    "format": "datetime",
                    "placeholder": "请选择更新时间"
                }
            },
            "data": {
                "type": "string",
                "example": "2025-01-15 10:30:00",
                "default_value": "$now"
            },
            "permission": {
                "read": true,
                "update": false,
                "create": false
            },
            "callbacks": [],
            "validation": ""
        }
    ]
}
```

## 技术优势

### 1. 结构统一
- Form和Table使用相同的 `FieldInfo` 结构
- 消除了 `add_form_config` 的冗余嵌套
- 统一的字段描述格式

### 2. 配置清晰
- **Widget配置**：组件类型和个性化配置分离
- **数据配置**：类型、默认值、示例值统一管理
- **权限配置**：简化为读、写、创建三权限
- **回调配置**：支持字段级别的事件回调

### 3. 扩展性强
- 新增组件类型只需在 `widget.config` 中添加配置
- 新增数据类型只需在 `data.type` 中定义
- 新增权限控制只需在 `permission` 中扩展

### 4. 维护简单
- 单一数据结构，减少维护成本
- 统一的解析逻辑，减少重复代码
- 清晰的职责分离，便于调试

## 使用方法

### Form响应
```go
response, err := NewUnifiedFormResponse(req, "form")
if err != nil {
    return err
}
return c.JSON(200, response)
```

### Table响应
```go
response, err := NewUnifiedTableResponse(resp)
if err != nil {
    return err
}
return c.JSON(200, response)
```

## 迁移建议

1. **渐进式迁移**：保持旧接口兼容，新功能使用统一结构
2. **测试验证**：确保新结构与前端组件兼容
3. **文档更新**：更新API文档和使用示例
4. **性能优化**：统一结构后可以进行缓存优化

## 总结

通过统一API响应结构，我们实现了：
- ✅ 消除结构冗余和不一致
- ✅ 简化维护和扩展
- ✅ 提高代码复用性
- ✅ 增强系统的可维护性

这个方案既解决了当前的问题，又为未来的功能扩展奠定了良好的基础。 