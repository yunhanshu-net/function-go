# OnInputFuzzyMap 改造计划

## 📋 问题背景

### 当前问题
在 `function-go` 框架中，`OnInputFuzzyMap` 回调函数存在以下问题：

1. **查询逻辑不统一**：不同开发者使用不同的查询模式，缺乏标准化
2. **性能问题**：没有针对不同查询场景进行优化
3. **GORM链式调用错误**：经常忘记重新赋值 `db = db.Where()`，导致条件不生效
4. **缺乏查询类型区分**：没有区分多值查询、单值查询和模糊查询

### 影响范围
- 所有使用 `OnInputFuzzyMap` 的表单函数
- 前端用户体验（查询响应慢、结果不准确）
- 开发效率（需要反复调试查询逻辑）

## 🎯 改造目标

### 核心目标
1. **统一查询模式**：建立标准化的三层查询逻辑
2. **性能优化**：针对不同场景使用最优查询方式
3. **错误预防**：通过标准模式避免常见GORM错误
4. **开发效率**：提供清晰的代码模板和最佳实践

### 具体目标
- 所有 `OnInputFuzzyMap` 回调使用统一的三层查询逻辑
- 查询性能提升 30% 以上
- 减少 90% 的GORM链式调用错误
- 提供完整的代码模板和文档

## 🔧 解决方案

### 新的三层查询逻辑

#### 1. 查询类型分类
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
        Where("stock > ?", 0).
        Limit(20)
}
db.Find(&products)
```

#### 2. 性能优化策略

| 查询类型 | SQL语句 | 性能特点 | 适用场景 |
|---------|---------|----------|----------|
| 单值查询 | `id = ?` + `Limit(1)` | 性能最优 | 单选字段，精确匹配 |
| 多值查询 | `id IN ?` | 批量查询 | 多选字段，减少查询次数 |
| 模糊查询 | `LIKE` + `Limit(20)` | 控制结果数量 | 关键字搜索，提升响应速度 |

#### 3. 关键改进点

**GORM链式调用修复**：
```go
// ❌ 错误：没有重新赋值，条件不会生效
db.Where("id = ?", 1)
db.Where("status = ?", "正常")
db.Find(&items) // 这里会查询所有数据，忽略上面的条件！

// ✅ 正确：重新赋值，条件生效
db = db.Where("id = ?", 1)
db = db.Where("status = ?", "正常")
db.Find(&items) // 这里会应用所有条件
```

**记忆口诀**：`db = db.Where()` 而不是 `db.Where()`

### 旧用法 vs 新用法对比

#### 旧用法（retail_simple_cashier.go）
```go
"product_id": func(ctx *runner.Context, req *usercall.OnInputFuzzyReq) (*usercall.OnInputFuzzyResp, error) {
    keyword := fmt.Sprintf("%v", req.Value)
    var products []RetailSimpleCashierProduct
    
    db := ctx.MustGetOrInitDB()
    // ❌ 旧用法：直接链式调用，没有区分查询类型
    db.Where("name LIKE ? OR category LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
        Where("status = ?", "上架").
        Where("stock > ?", 0).
        Limit(20).
        Find(&products)
    
    // ... 构建返回数据
}
```

#### 新用法（healthcare_medicine.go）
```go
"medicine_id": func(ctx *runner.Context, req *usercall.OnInputFuzzyReq) (*usercall.OnInputFuzzyResp, error) {
    var medicines []HealthcareMedicine
    db := ctx.MustGetOrInitDB()

    if req.IsByFiledValues() {
        // 多值查询：使用 IN 查询，这里必须用 db = db.Where
        db = db.Where("id in ?", req.GetFiledValues())
    } else if req.IsByFiledValue() {
        // 单值查询：使用等值查询，性能最优
        db = db.Where("id = ?", req.GetFiledValue()).Limit(1)
    } else {
        // 模糊查询：关键字搜索
        db = db.Where("name LIKE ? OR code LIKE ?", "%"+req.Keyword()+"%", "%"+req.Keyword()+"%").
            Where("status = ?", "正常").
            Limit(20)
    }
    db.Find(&medicines)
    
    // ... 构建返回数据
}
```

#### 关键差异说明

1. **查询类型区分**：
   - 旧用法：只支持模糊查询，使用 `req.Value` 和 `fmt.Sprintf`
   - 新用法：支持三种查询类型，使用 `req.IsByFiledValues()`、`req.IsByFiledValue()`、`req.Keyword()`

2. **GORM链式调用**：
   - 旧用法：直接链式调用，可能忘记重新赋值
   - 新用法：明确使用 `db = db.Where()` 确保条件生效

3. **性能优化**：
   - 旧用法：所有查询都使用 `LIKE`，性能较差
   - 新用法：单值查询使用等值查询，性能最优

4. **字段使用**：
   - 旧用法：使用 `req.Value` 和 `fmt.Sprintf("%v", req.Value)`
   - 新用法：使用 `req.GetFiledValues()`、`req.GetFiledValue()`、`req.Keyword()`

## 🤖 大模型改造指南

### 改造任务描述
请帮我将以下OnInputFuzzyMap回调从旧用法升级为新用法：

**改造目标**：
- 将单一模糊查询模式升级为三层查询逻辑
- 修复GORM链式调用问题
- 优化查询性能
- 保持原有业务逻辑不变

### 改造步骤

#### 步骤1：识别旧用法模式
查找以下特征：
```go
// ❌ 旧用法特征
keyword := fmt.Sprintf("%v", req.Value)
db.Where("name LIKE ? OR category LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
    Where("status = ?", "上架").
    Limit(20).
    Find(&products)
```

#### 步骤2：应用新用法模式
替换为以下标准模式：
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

#### 步骤3：关键修改点
1. **删除**：`keyword := fmt.Sprintf("%v", req.Value)`
2. **替换**：`req.Value` → `req.Keyword()`
3. **添加**：三层查询逻辑判断
4. **修复**：`db.Where()` → `db = db.Where()`
5. **优化**：单值查询使用等值查询 + `Limit(1)`

#### 步骤4：验证改造结果
确保以下要点：
- [ ] 使用了三层查询逻辑
- [ ] 所有 `db.Where()` 都改为 `db = db.Where()`
- [ ] 单值查询使用了等值查询
- [ ] 保持了原有的业务逻辑
- [ ] 注释包含了【当前场景】和【目标字段】

### 改造示例

#### 改造前（旧用法）
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

#### 改造后（新用法）
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

## 📝 实施计划

### 阶段一：文档和模板更新（已完成）
- [x] 更新 `function-go文档.md` 中的示例代码
- [x] 创建标准化的代码模板
- [x] 编写最佳实践文档

### 阶段二：示例代码改造（已完成）
- [x] 更新 `formexample_input_fuzzy.go` 中的所有回调
- [x] 更新 `healthcare_medicine.go` 中的回调
- [x] 添加联想产品等测试数据

### 阶段三：框架规则更新（进行中）
- [x] 更新 `complete_framework_guide.mdc` 中的规则
- [x] 更新 `OnInputFuzzy回调最佳实践` 规则
- [x] 更新 `OnInputFuzzy聚合计算使用规范` 规则

### 阶段四：全面推广（待进行）
- [ ] 扫描所有现有代码，识别需要改造的回调
- [ ] 批量更新现有代码
- [ ] 建立代码审查标准
- [ ] 培训开发团队

## 🏆 最佳实践

### 1. 注释规范
每个回调函数都应该包含：
```go
// 【当前场景】field_name是单选字段，只返回静态信息，不做聚合计算
// 【目标字段】为 YourRequest.FieldName 提供数据选择
```

### 2. 聚合计算规则
- **单选字段**：只返回静态信息，不做聚合计算
- **多选字段**：支持聚合计算（count、sum、avg等）
- **List内字段**：支持聚合计算，基于用户选择的数据

### 3. 错误处理
```go
// 数据库连接检查
db := ctx.MustGetOrInitDB()

// 查询结果验证
if len(items) == 0 {
    return &usercall.OnInputFuzzyResp{
        Statistics: map[string]interface{}{
            "提示": "未找到匹配的数据",
        },
        Values: []*usercall.InputFuzzyItem{},
    }, nil
}
```

## 📊 预期效果

### 性能提升
- 单值查询响应时间减少 50%
- 多值查询减少数据库连接次数
- 模糊查询结果数量控制，提升前端渲染速度

### 开发效率
- 减少 90% 的GORM链式调用错误
- 统一的代码模式，降低学习成本
- 清晰的注释规范，提升代码可维护性

### 用户体验
- 查询响应更快
- 结果更准确
- 聚合计算实时更新

## 🔍 验证方法

### 1. 功能测试
- 测试三种查询类型的正确性
- 验证聚合计算的准确性
- 检查错误处理逻辑

### 2. 性能测试
- 对比改造前后的查询性能
- 测试大数据量下的响应时间
- 验证内存使用情况

### 3. 代码质量
- 检查GORM链式调用的正确性
- 验证注释的完整性
- 确保代码风格一致性

## 📚 相关文档

- [OnInputFuzzy回调最佳实践](./OnInputFuzzy回调最佳实践.md)
- [OnInputFuzzy聚合计算使用规范](./OnInputFuzzy聚合计算使用规范.md)
- [GORM常见陷阱](./GORM常见陷阱.md)
- [function-go框架完整指南](./complete_framework_guide.mdc)

## 🎯 总结

通过实施OnInputFuzzyMap改造计划，我们将：

1. **统一查询模式**：建立标准化的三层查询逻辑
2. **提升性能**：针对不同场景优化查询策略
3. **减少错误**：通过标准模式避免常见问题
4. **改善体验**：为用户提供更快更准确的查询结果

这个改造计划将显著提升 `function-go` 框架的稳定性和开发效率，为后续功能开发奠定坚实基础。
