# 聚合计算功能使用指南

## 📋 目录

- [概述](#概述)
- [基础语法](#基础语法)
- [聚合函数](#聚合函数)
- [使用场景](#使用场景)
- [完整示例](#完整示例)
- [最佳实践](#最佳实践)
- [常见问题](#常见问题)

---

## 概述

聚合计算功能允许在用户选择数据时实时显示统计信息，通过 `OnInputFuzzyResp.Statistics` 字段提供动态计算和静态展示信息。

### 核心特性

- **实时计算**：用户选择时即时显示统计结果
- **多种聚合**：支持求和、平均值、最大值、最小值、计数等
- **乘法运算**：支持字段间的乘法计算（如价格×数量）
- **静态信息**：支持显示固定的业务信息
- **人性化展示**：自动格式化数字，支持单位和前缀

---

## 基础语法

### Statistics 字段结构

```go
Statistics: map[string]interface{}{
    "显示名称": "聚合公式或静态值",
}
```

### 聚合公式语法

```go
// 基础聚合
"sum(字段名)"          // 求和
"avg(字段名)"          // 平均值  
"min(字段名)"          // 最小值
"max(字段名)"          // 最大值
"count(字段名)"        // 计数

// 乘法运算
"sum(字段名,*数量字段)"     // 求和后乘以数量
"sum(字段名,*数量字段,*0.9)" // 求和后乘以数量再乘以0.9(九折)

// 静态信息
"配送说明": "全国包邮"      // 直接显示静态文本
```

---

## 聚合函数

### 1. 基础聚合函数

| 函数 | 语法 | 说明 | 示例 |
|------|------|------|------|
| **sum** | `sum(字段名)` | 求和 | `sum(价格)` → 所有价格的总和 |
| **avg** | `avg(字段名)` | 平均值 | `avg(价格)` → 价格的平均值 |
| **min** | `min(字段名)` | 最小值 | `min(价格)` → 最低价格 |
| **max** | `max(字段名)` | 最大值 | `max(价格)` → 最高价格 |
| **count** | `count(字段名)` | 计数 | `count(价格)` → 商品种类数 |

### 2. 乘法运算

#### 乘法变量来源说明

**关键概念**：乘法运算中的变量（如`*quantity`、`*count`等）来自于**当前请求结构体的同级字段**，不是DisplayInfo中的字段。

```go
// 示例1：quantity变量
type ProductSelectReq struct {
    ProductIDs []int `json:"product_ids"`     // 选择的商品ID列表
    Quantity   int   `json:"quantity"`        // 这个字段名就是聚合公式中的变量名
}
// 聚合公式：sum(价格,*quantity)  ← quantity变量来自同级的Quantity字段

// 示例2：count变量
type OrderReq struct {
    ProductIDs []int `json:"product_ids"`
    Count      int   `json:"count"`           // 字段名是count
}
// 聚合公式：sum(价格,*count)     ← count变量来自同级的Count字段

// 示例3：List组件中的变量
type ProductQuantity struct {
    ProductID int `json:"product_id"`
    Quantity  int `json:"quantity"`          // List项中的字段
    Weight    int `json:"weight"`            // 可以作为另一个变量
}
// 聚合公式：sum(价格,*quantity) 或 sum(价格,*weight)
```

#### 乘法运算语法

```go
// 基础乘法：DisplayInfo字段 × 同级结构体字段
"sum(价格,*quantity)"           // 价格来自DisplayInfo，quantity来自同级字段

// 使用不同的变量名
"sum(价格,*count)"              // count来自同级的Count字段
"sum(价格,*num)"                // num来自同级的Num字段

// 多重乘法：字段值 × 变量 × 常数系数  
"sum(价格,*quantity,*0.9)"      // 九折价格
"sum(重量,*count,*2.5)"         // 运费计算

// 使用多个不同变量
"sum(价格,*quantity,*discount)" // discount也来自同级字段
```

#### 变量的不同使用场景

**场景1：非List组件 - 同级字段作为变量**
```go
// 请求结构体
type OrderReq struct {
    ProductIDs []int   `json:"product_ids"`  // 多选商品
    Quantity   int     `json:"quantity"`     // 变量quantity
    Discount   float64 `json:"discount"`     // 变量discount
}

// 聚合计算：变量来自同级字段
"sum(价格,*quantity)"           // quantity变量 ← OrderReq.Quantity字段
"sum(价格,*quantity,*discount)" // discount变量 ← OrderReq.Discount字段
```

**场景2：List组件 - List项字段作为变量**
```go
// 请求结构体
type ListReq struct {
    Items []ProductItem `json:"items"`
}

type ProductItem struct {
    ProductID int     `json:"product_id"`
    Quantity  int     `json:"quantity"`     // 变量quantity
    Weight    float64 `json:"weight"`       // 变量weight  
}

// 聚合计算：变量来自List项的字段
"sum(价格,*quantity)"  // quantity变量 ← ProductItem.Quantity字段
"sum(价格,*weight)"    // weight变量 ← ProductItem.Weight字段
```

### 3. 静态信息

```go
// 业务规则展示
"会员折扣": "9折优惠",
"配送说明": "全国包邮",
"服务承诺": "30天包换",

// 动态提示
"支付方式": "余额扣款",
"到账时间": "实时到账",
```

---

## 使用场景

### 场景1：电商购物车

```go
Statistics: map[string]interface{}{
    // 核心价格信息
    "商品原价": "sum(价格,*quantity)",
    "会员价格": "sum(价格,*quantity,*0.9)", 
    "节省金额": "sum(价格,*quantity,*0.1)",
    
    // 重要提示
    "会员折扣": "9折优惠",
}
```

**用户看到的效果**：
```
商品原价: ¥168.50
会员价格: ¥151.65
节省金额: ¥16.85
会员折扣: 9折优惠
```

### 场景2：会员卡管理

```go
Statistics: map[string]interface{}{
    // 核心信息
    "会员总数": "count(余额)",
    "总余额": "sum(余额)",
    
    // 重要提示  
    "支付方式": "余额扣款",
    "会员折扣": "9折优惠",
}
```

### 场景3：物流配送

```go
Statistics: map[string]interface{}{
    // 重量计算
    "总重量(kg)": "sum(重量,*quantity)",
    "配送费(元)": "sum(重量,*quantity,*2.5)",
    
    // 配送信息
    "配送范围": "全国配送", 
    "配送时效": "1-3个工作日",
}
```

---

## 乘法变量详解

### 核心概念

**关键理解**：聚合计算中的乘法变量（如`*quantity`、`*count`等）来自**当前请求结构体的同级字段**，不是回调函数中的DisplayInfo。

**变量命名规则**：
- 变量名必须与结构体字段名完全一致（区分大小写）
- 变量名前加`*`符号表示乘法运算
- 字段名`Quantity` → 变量名`*quantity`（JSON tag决定变量名）
- 字段名`Count` → 变量名`*count`
- 字段名`Weight` → 变量名`*weight`

### 完整的变量使用示例

```go
// 1. 用户请求结构体定义（变量来源）
type ProductSelectReq struct {
    ProductIDs []int   `json:"product_ids" widget:"type:multiselect"`  // 多选商品
    Quantity   int     `json:"quantity" widget:"type:number"`          // 变量quantity ← 字段名决定变量名
    Discount   float64 `json:"discount" widget:"type:number"`          // 变量discount ← 字段名决定变量名
}

// 2. 回调函数中的DisplayInfo（聚合目标字段）
"product_ids": func(ctx *runner.Context, req *usercall.OnInputFuzzyReq) (*usercall.OnInputFuzzyResp, error) {
    items := []*usercall.InputFuzzyItem{
        {
            Value: 1,
            Label: "苹果 - ¥8.5",
            DisplayInfo: map[string]interface{}{
                "商品名称": "苹果",
                "价格":   8.5,        // ← 聚合目标字段
                "重量":   0.2,        // ← 聚合目标字段
                // 注意：变量不在DisplayInfo中！
            },
        },
    }
    
    return &usercall.OnInputFuzzyResp{
        Statistics: map[string]interface{}{
            // 3. 聚合计算：DisplayInfo字段 × 同级结构体字段
            "总价": "sum(价格,*quantity)",              // quantity变量 ← ProductSelectReq.Quantity字段
            "总重量": "sum(重量,*quantity)",            // quantity变量 ← ProductSelectReq.Quantity字段
            "折扣价": "sum(价格,*quantity,*discount)",  // discount变量 ← ProductSelectReq.Discount字段
        },
        Values: items,
    }, nil
}

// 4. 用户操作示例
// 用户选择：苹果(¥8.5)、香蕉(¥6.8)
// 用户输入：数量=3, 折扣=0.9
// 
// 聚合计算结果：
// 总价 = sum(8.5, 6.8) * 3 = 15.3 * 3 = 45.9
// 折扣价 = sum(8.5, 6.8) * 3 * 0.9 = 45.9 * 0.9 = 41.31
```

### List组件中的变量使用

```go
// 1. List组件的请求结构体（变量来源）
type ListProductReq struct {
    Items []ProductQuantity `json:"items" widget:"type:list"`
}

type ProductQuantity struct {
    ProductID int     `json:"product_id" widget:"type:select"`     // 商品选择
    Quantity  int     `json:"quantity" widget:"type:number"`       // 变量quantity ← 字段名决定变量名
    Weight    float64 `json:"weight" widget:"type:number"`         // 变量weight ← 字段名决定变量名
}

// 2. 回调函数
"product_id": func(ctx *runner.Context, req *usercall.OnInputFuzzyReq) (*usercall.OnInputFuzzyResp, error) {
    // ... DisplayInfo设置 ...
    
    return &usercall.OnInputFuzzyResp{
        Statistics: map[string]interface{}{
            // 3. List聚合：每个item使用各自的同级字段作为变量
            "总价": "sum(价格,*quantity)",      // quantity变量 ← ProductQuantity.Quantity字段
            "总运费": "sum(运费,*weight)",      // weight变量 ← ProductQuantity.Weight字段
            "总重量": "sum(单重,*quantity)",    // quantity变量 ← ProductQuantity.Quantity字段
        },
        Values: items,
    }, nil
}

// 4. 用户操作示例
// 用户添加：
// Item1: 苹果(¥8.5) × 2个, 重量1.5kg
// Item2: 香蕉(¥6.8) × 3个, 重量2.0kg
//
// 聚合计算结果：
// 总价 = (8.5 * 2) + (6.8 * 3) = 17 + 20.4 = 37.4
// 总运费 = (运费单价 * 1.5) + (运费单价 * 2.0)
```

---

## 完整示例

### 示例1：商品选择聚合计算

```go
"product_ids": func(ctx *runner.Context, req *usercall.OnInputFuzzyReq) (*usercall.OnInputFuzzyResp, error) {
    // 查询商品数据
    var products []Product
    db.Where("status = ?", "上架").Where("stock > ?", 0).Find(&products)
    
    // 构建选项
    items := make([]*usercall.InputFuzzyItem, 0)
    for _, p := range products {
        items = append(items, &usercall.InputFuzzyItem{
            Value: p.ID,
            Label: fmt.Sprintf("%s - ¥%.2f", p.Name, p.Price),
            DisplayInfo: map[string]interface{}{
                "商品名称": p.Name,
                "价格":   p.Price,
                "重量":   p.Weight,
                "库存":   p.Stock,
            },
        })
    }
    
    return &usercall.OnInputFuzzyResp{
        Statistics: map[string]interface{}{
            // 核心价格信息 - 用户最关心
            "商品原价": "sum(价格,*quantity)",
            "会员价格": "sum(价格,*quantity,*0.9)", 
            "节省金额": "sum(价格,*quantity,*0.1)",
            
            // 重量信息 - 影响配送
            "总重量": "sum(重量,*quantity)",
            
            // 重要提示 - 业务规则
            "会员折扣": "9折优惠",
            "配送说明": "满99元包邮",
        },
        Values: items,
    }, nil
}
```

### 示例2：会员卡选择聚合计算

```go
"member_id": func(ctx *runner.Context, req *usercall.OnInputFuzzyReq) (*usercall.OnInputFuzzyResp, error) {
    // 查询会员数据
    var members []MemberCard
    db.Where("status = ?", "正常").Find(&members)
    
    // 构建选项
    items := make([]*usercall.InputFuzzyItem, 0)
    for _, m := range members {
        items = append(items, &usercall.InputFuzzyItem{
            Value: m.ID,
            Label: fmt.Sprintf("%s - 余额:¥%.2f", m.CustomerName, m.Balance),
            DisplayInfo: map[string]interface{}{
                "客户姓名": m.CustomerName,
                "余额":   m.Balance,
                "卡号":   m.CardNumber,
            },
        })
    }
    
    return &usercall.OnInputFuzzyResp{
        Statistics: map[string]interface{}{
            // 核心信息
            "会员总数": "count(余额)",
            "总余额": "sum(余额)",
            
            // 重要提示
            "支付方式": "余额扣款", 
            "会员折扣": "9折优惠",
        },
        Values: items,
    }, nil
}
```

---

## 最佳实践

### 1. 信息精简原则

❌ **错误示例**：信息过多，用户困惑
```go
Statistics: map[string]interface{}{
    "会员总数": "count(余额)",
    "总余额(元)": "sum(余额)", 
    "平均余额(元)": "avg(余额)",
    "最高余额(元)": "max(余额)",
    "最低余额(元)": "min(余额)",
    "会员权益": "专享9折优惠",
    "支付方式": "余额扣款",
    "查询方式": "支持卡号/手机号/姓名搜索", 
    "会员状态": "仅显示正常状态会员",
    "服务热线": "400-888-8888",
    "营业时间": "9:00-22:00",
    "积分政策": "消费1元得1积分",
    "充值优惠": "充值满100送10",
}
```

✅ **正确示例**：信息精简，重点突出
```go
Statistics: map[string]interface{}{
    // 核心信息
    "会员总数": "count(余额)",
    "总余额": "sum(余额)",
    
    // 重要提示
    "支付方式": "余额扣款",
    "会员折扣": "9折优惠", 
}
```

### 2. 层次分明原则

```go
Statistics: map[string]interface{}{
    // 第一层：核心数据 - 用户最关心的
    "商品原价": "sum(价格,*quantity)",
    "会员价格": "sum(价格,*quantity,*0.9)",
    "节省金额": "sum(价格,*quantity,*0.1)",
    
    // 第二层：重要提示 - 关键业务规则
    "会员折扣": "9折优惠",
}
```

### 3. 字段名称规范

```go
// ✅ 好的命名：简洁明了
"商品原价": "sum(价格,*quantity)",
"会员价格": "sum(价格,*quantity,*0.9)", 
"节省金额": "sum(价格,*quantity,*0.1)",

// ❌ 不好的命名：过于详细
"商品原价(元)": "sum(价格,*quantity)",
"会员价格(元)": "sum(价格,*quantity,*0.9)",
"优惠金额(元)": "sum(价格,*quantity,*0.1)",
```

### 4. 业务场景适配

不同场景关注点不同，聚合信息应该匹配：

**收银场景**：关注价格和优惠
```go
"商品原价": "sum(价格,*quantity)",
"会员价格": "sum(价格,*quantity,*0.9)",
"节省金额": "sum(价格,*quantity,*0.1)",
```

**充值场景**：关注余额和优惠政策
```go
"会员总数": "count(余额)", 
"总余额": "sum(余额)",
"充值优惠": "满100送10",
```

**库存场景**：关注数量和库存
```go
"商品种类": "count(库存)",
"总库存": "sum(库存)",
"库存预警": "低于10件需补货",
```

---

## 常见问题

### Q1: 聚合计算不生效？

**原因**：DisplayInfo中的字段名与聚合公式不匹配

```go
// ❌ 错误：字段名不匹配
DisplayInfo: map[string]interface{}{
    "price": p.Price,  // 字段名是 price
}
Statistics: map[string]interface{}{
    "总价": "sum(价格)",  // 但聚合用的是 价格
}

// ✅ 正确：字段名匹配
DisplayInfo: map[string]interface{}{
    "价格": p.Price,
}
Statistics: map[string]interface{}{
    "总价": "sum(价格)",
}
```

### Q2: 乘法运算结果不对？

**原因**：变量来源不正确或变量名与字段名不匹配

**重要**：乘法变量来自**当前请求结构体的同级字段**，不是DisplayInfo！变量名必须与字段名完全一致。

```go
// ❌ 错误理解：以为变量在DisplayInfo中
DisplayInfo: map[string]interface{}{
    "价格": p.Price,
    "quantity": 1,  // 错误！变量不在DisplayInfo中
}

// ✅ 正确理解：变量来自请求结构体字段
type ProductSelectReq struct {
    ProductIDs []int `json:"product_ids"`
    Quantity   int   `json:"quantity"`     // 变量quantity ← 字段名决定变量名
    Count      int   `json:"count"`        // 变量count ← 字段名决定变量名
}

// 或者List组件中的变量
type ProductQuantity struct {
    ProductID int `json:"product_id"`
    Quantity  int `json:"quantity"`       // 变量quantity ← 字段名决定变量名
    Weight    int `json:"weight"`         // 变量weight ← 字段名决定变量名
}

// 聚合公式：变量名必须与字段名一致
"总价": "sum(价格,*quantity)",  // quantity变量 ← 对应Quantity字段
"总数": "sum(价格,*count)",     // count变量 ← 对应Count字段
"总重": "sum(价格,*weight)",    // weight变量 ← 对应Weight字段
```

**调试步骤**：
1. 检查请求结构体中是否有对应的字段（如quantity对应Quantity字段）
2. 确认变量名与字段名完全一致（区分大小写）
3. List组件确保每个item都有对应字段
4. 确认字段类型为int或float64
5. 变量名前必须加`*`符号

### Q3: 静态信息不显示？

**原因**：静态信息直接使用字符串值，不需要聚合函数

```go
// ✅ 正确：静态信息
"会员折扣": "9折优惠",
"配送说明": "全国包邮",

// ❌ 错误：对静态信息使用聚合函数  
"会员折扣": "sum(9折优惠)",  // 错误！
```

### Q4: 如何调试聚合计算？

1. **检查字段名**：确保DisplayInfo中的字段名与聚合公式中的一致
2. **检查数据类型**：数值字段确保是number类型
3. **检查语法**：聚合公式语法是否正确
4. **逐步测试**：先测试简单的count，再测试复杂的sum

```go
// 调试步骤示例
Statistics: map[string]interface{}{
    // 第1步：测试基础计数
    "数据总数": "count(价格)",
    
    // 第2步：测试简单求和  
    "价格总和": "sum(价格)",
    
    // 第3步：测试乘法运算
    "总价格": "sum(价格,*quantity)",
    
    // 第4步：测试复杂计算
    "会员价格": "sum(价格,*quantity,*0.9)",
}
```

---

## 版本更新记录

- **v1.0** (2024-01): 初始版本，支持基础聚合函数
- **v1.1** (2024-02): 新增乘法运算支持  
- **v1.2** (2024-03): 新增静态信息展示
- **v1.3** (2024-04): 优化用户体验，精简信息展示

---

## 相关文档

- [Form函数开发指南](./form-functions.md)
- [OnInputFuzzy回调机制](./callback-mechanism.md)  
- [Widget标签使用手册](./widget-tags.md)
- [聚合计算演示示例](../soft/beiluo/demo4/code/api/formexample/aggregation_demo.go)

---

*最后更新：2024年12月*
