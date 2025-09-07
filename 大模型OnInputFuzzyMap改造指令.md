# 大模型OnInputFuzzyMap改造指令

## 🎯 改造任务

请帮我将OnInputFuzzyMap回调从旧用法升级为新用法。

## 📋 改造要求

### 核心目标
1. **将单一模糊查询模式升级为三层查询逻辑**
2. **修复GORM链式调用问题**
3. **优化查询性能**
4. **保持原有业务逻辑不变**

### 必须遵循的规则
- 严格按照新用法模式进行改造
- 不能改变原有的业务逻辑
- 必须添加标准注释
- 确保代码语法正确

## 🔍 识别旧用法

### 旧用法特征
查找以下代码模式：
```go
// ❌ 旧用法特征
keyword := fmt.Sprintf("%v", req.Value)
db.Where("name LIKE ? OR category LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
    Where("status = ?", "上架").
    Limit(20).
    Find(&products)
```

### 关键识别点
1. 使用 `keyword := fmt.Sprintf("%v", req.Value)`
2. 直接链式调用 `db.Where()`
3. 只支持模糊查询
4. 没有查询类型区分

## ✅ 应用新用法

### 标准新用法模式
```go
// ✅ 新用法标准模式
if req.IsByFiledValues() {
    // 多值查询：使用 IN 查询，这里必须用 db = db.Where
    db = db.Where("id in ?", req.GetFiledValues())
} else if req.IsByFiledValue() {
    // 单值查询：使用等值查询，性能最优
    db = db.Where("id = ?", req.GetFiledValue()).Limit(1)
} else {
    // 模糊查询：关键字搜索
    db = db.Where("name LIKE ? OR category LIKE ?", "%"+req.Keyword()+"%", "%"+req.Keyword()+"%").
        Where("status = ?", "上架").
        Limit(20)
}
db.Find(&products)
```

## 🔧 具体改造步骤

### 步骤1：删除旧代码
删除以下代码：
```go
keyword := fmt.Sprintf("%v", req.Value)
```

### 步骤2：替换查询逻辑
将原有的单一查询替换为三层查询逻辑：
```go
if req.IsByFiledValues() {
    // 多值查询：使用 IN 查询，这里必须用 db = db.Where
    db = db.Where("id in ?", req.GetFiledValues())
} else if req.IsByFiledValue() {
    // 单值查询：使用等值查询，性能最优
    db = db.Where("id = ?", req.GetFiledValue()).Limit(1)
} else {
    // 模糊查询：关键字搜索
    db = db.Where("name LIKE ? OR category LIKE ?", "%"+req.Keyword()+"%", "%"+req.Keyword()+"%").
        Where("status = ?", "上架").
        Limit(20)
}
```

### 步骤3：修复GORM链式调用
确保所有 `db.Where()` 都改为 `db = db.Where()`：
```go
// ❌ 错误
db.Where("id = ?", 1)
db.Where("status = ?", "正常")

// ✅ 正确
db = db.Where("id = ?", 1)
db = db.Where("status = ?", "正常")
```

### 步骤4：添加标准注释
在函数开头添加：
```go
// 【当前场景】field_name是单选字段，只返回静态信息，不做聚合计算
// 【目标字段】为 YourRequest.FieldName 提供数据选择
```

## 📝 完整改造示例

### 改造前（旧用法）
```go
"product_id": func(ctx *runner.Context, req *usercall.OnInputFuzzyReq) (*usercall.OnInputFuzzyResp, error) {
    keyword := fmt.Sprintf("%v", req.Value)
    var products []Product
    
    db := ctx.MustGetOrInitDB()
    db.Where("name LIKE ? OR category LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
        Where("status = ?", "上架").
        Limit(20).
        Find(&products)
    
    // 构建返回数据...
}
```

### 改造后（新用法）
```go
"product_id": func(ctx *runner.Context, req *usercall.OnInputFuzzyReq) (*usercall.OnInputFuzzyResp, error) {
    // 【当前场景】product_id是单选字段，只返回静态信息，不做聚合计算
    // 【目标字段】为 YourRequest.ProductID 提供商品选择数据
    var products []Product
    db := ctx.MustGetOrInitDB()

    if req.IsByFiledValues() {
        // 多值查询：使用 IN 查询，这里必须用 db = db.Where
        db = db.Where("id in ?", req.GetFiledValues())
    } else if req.IsByFiledValue() {
        // 单值查询：使用等值查询，性能最优
        db = db.Where("id = ?", req.GetFiledValue()).Limit(1)
    } else {
        // 模糊查询：关键字搜索
        db = db.Where("name LIKE ? OR category LIKE ?", "%"+req.Keyword()+"%", "%"+req.Keyword()+"%").
            Where("status = ?", "上架").
            Limit(20)
    }
    db.Find(&products)
    
    // 构建返回数据...
}
```

## ✅ 验证清单

改造完成后，请确保：
- [ ] 删除了 `keyword := fmt.Sprintf("%v", req.Value)`
- [ ] 使用了三层查询逻辑
- [ ] 所有 `db.Where()` 都改为 `db = db.Where()`
- [ ] 单值查询使用了等值查询 + `Limit(1)`
- [ ] 保持了原有的业务逻辑
- [ ] 添加了【当前场景】和【目标字段】注释
- [ ] 代码语法正确，没有编译错误

## 🚨 常见错误

### 错误1：忘记重新赋值
```go
// ❌ 错误
db.Where("id = ?", 1)
db.Where("status = ?", "正常")

// ✅ 正确
db = db.Where("id = ?", 1)
db = db.Where("status = ?", "正常")
```

### 错误2：使用错误的字段
```go
// ❌ 错误
req.Value

// ✅ 正确
req.Keyword()
req.GetFiledValue()
req.GetFiledValues()
```

### 错误3：缺少查询类型判断
```go
// ❌ 错误
db.Where("name LIKE ?", "%"+req.Keyword()+"%")

// ✅ 正确
if req.IsByFiledValues() {
    db = db.Where("id in ?", req.GetFiledValues())
} else if req.IsByFiledValue() {
    db = db.Where("id = ?", req.GetFiledValue()).Limit(1)
} else {
    db = db.Where("name LIKE ?", "%"+req.Keyword()+"%")
}
```

## 📚 参考文档

- [OnInputFuzzyMap改造计划](./OnInputFuzzyMap改造计划.md)
- [OnInputFuzzy回调最佳实践](./OnInputFuzzy回调最佳实践.md)
- [GORM常见陷阱](./GORM常见陷阱.md)

---

**重要提醒**：请严格按照以上步骤进行改造，确保代码质量和性能优化。
