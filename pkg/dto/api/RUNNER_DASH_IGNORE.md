# runner:"-" 字段忽略功能

## 功能说明

当结构体字段标记为 `runner:"-"` 时，该字段将在表单和表格配置生成过程中被完全忽略，不会出现在最终的API响应中。

## 使用场景

### 1. 内部字段
```go
type User struct {
    ID       int    `json:"id" runner:"code:id;name:用户ID" data:"type:number"`
    Name     string `json:"name" runner:"code:name;name:用户名" widget:"type:input" data:"type:string"`
    Password string `json:"-" runner:"-"` // 密码字段，不在API中暴露
}
```

### 2. 临时字段
```go
type Product struct {
    ID          int    `json:"id" runner:"code:id;name:产品ID" data:"type:number"`
    Name        string `json:"name" runner:"code:name;name:产品名称" widget:"type:input" data:"type:string"`
    TempData    string `runner:"-"` // 临时数据，不需要在表单中显示
}
```

### 3. 系统字段
```go
type Order struct {
    ID          int       `json:"id" runner:"code:id;name:订单ID" data:"type:number"`
    Amount      float64   `json:"amount" runner:"code:amount;name:金额" widget:"type:input" data:"type:float"`
    InternalRef string    `json:"internal_ref" runner:"-"` // 内部引用，不对外显示
    SystemFlag  bool      `runner:"-"` // 系统标志，不在表单中显示
}
```

## 支持的忽略场景

### ✅ 支持的标签组合

1. **仅runner忽略**
   ```go
   Field string `runner:"-"`
   ```

2. **json + runner忽略**
   ```go
   Field string `json:"field" runner:"-"`
   ```

3. **json和runner都忽略**
   ```go
   Field string `json:"-" runner:"-"`
   ```

4. **有其他标签但runner忽略**
   ```go
   Field string `json:"field" form:"field" runner:"-" validate:"required"`
   ```

5. **复杂标签组合**
   ```go
   Field string `json:"field" form:"field" runner:"-" widget:"type:input" data:"type:string"`
   ```

### 🔍 忽略优先级

`runner:"-"` 具有最高优先级，无论字段有多少其他标签，只要包含 `runner:"-"`，该字段就会被完全忽略。

## 实现原理

### 解析阶段忽略

在 `MultiTagParser.ParseStruct()` 方法中，解析器会在处理每个字段时首先检查 `runner` 标签：

```go
// 检查runner标签，如果是"-"则跳过该字段
if runnerTag := field.Tag.Get("runner"); runnerTag == "-" {
    continue
}
```

### 优势

1. **早期过滤**：在解析阶段就过滤掉不需要的字段，避免后续处理
2. **性能优化**：减少不必要的字段处理和内存占用
3. **安全性**：确保敏感字段不会意外暴露在API中
4. **简洁性**：API响应更加简洁，只包含必要的字段

## 测试验证

### 测试用例覆盖

```go
type TestComplexIgnoreFields struct {
    // 正常字段
    ID   int    `json:"id" runner:"code:id;name:ID" data:"type:number"`
    Name string `json:"name" runner:"code:name;name:名称" widget:"type:input" data:"type:string"`
    
    // 各种忽略场景
    InternalField1 string `runner:"-"`                                 // 仅runner忽略
    InternalField2 string `json:"internal2" runner:"-"`                // json+runner忽略
    InternalField3 string `json:"-" runner:"-"`                        // json和runner都忽略
    PasswordField  string `json:"password" form:"password" runner:"-"` // 有json和form但runner忽略
    TempData       []byte `runner:"-" validate:"required"`             // 有其他标签但runner忽略
    
    // 正常字段继续
    Status    string     `json:"status" runner:"code:status;name:状态" widget:"type:switch" data:"type:string"`
    CreatedAt typex.Time `json:"created_at" runner:"code:created_at;name:创建时间" widget:"type:datetime" data:"type:string" permission:"read"`
}
```

### 测试结果

```bash
=== RUN   TestComplexIgnoreFieldsScenarios
    ✅ Form响应正确忽略了 5 个字段，保留了 4 个字段: [id name status created_at]
    ✅ Table响应正确忽略了 5 个字段，保留了 4 个列: [id name status created_at]
    ✅ 所有 runner:"-" 字段忽略测试通过！
--- PASS: TestComplexIgnoreFieldsScenarios (0.00s)
```

## 使用示例

### 示例1：用户管理
```go
type User struct {
    ID           int        `json:"id" runner:"code:id;name:用户ID" data:"type:number"`
    Username     string     `json:"username" runner:"code:username;name:用户名" widget:"type:input" data:"type:string" validate:"required"`
    Email        string     `json:"email" runner:"code:email;name:邮箱" widget:"type:input" data:"type:string" validate:"required,email"`
    Password     string     `json:"-" runner:"-"` // 密码不在API中暴露
    Salt         string     `runner:"-"`          // 盐值不在API中暴露
    LastLoginIP  string     `runner:"-"`          // 内部追踪信息
    CreatedAt    typex.Time `json:"created_at" runner:"code:created_at;name:创建时间" widget:"type:datetime" data:"type:string" permission:"read"`
}
```

### 示例2：订单管理
```go
type Order struct {
    ID          int        `json:"id" runner:"code:id;name:订单ID" data:"type:number"`
    OrderNo     string     `json:"order_no" runner:"code:order_no;name:订单号" widget:"type:input" data:"type:string"`
    Amount      float64    `json:"amount" runner:"code:amount;name:订单金额" widget:"type:input;prefix:¥" data:"type:float" validate:"required,min=0"`
    Status      string     `json:"status" runner:"code:status;name:订单状态" widget:"type:select;options:待付款,已付款,已发货,已完成,已取消" data:"type:string"`
    
    // 内部字段，不对外暴露
    InternalRef    string `json:"internal_ref" runner:"-"`    // 内部引用号
    PaymentSecret  string `runner:"-"`                        // 支付密钥
    SystemNotes    string `runner:"-"`                        // 系统备注
    
    CreatedAt   typex.Time `json:"created_at" runner:"code:created_at;name:创建时间" widget:"type:datetime" data:"type:string" permission:"read"`
}
```

## 最佳实践

### 1. 安全性优先
- 所有密码、密钥、敏感信息字段都应该使用 `runner:"-"`
- 内部系统字段不应该暴露给前端

### 2. 性能考虑
- 大型对象或不必要的字段使用 `runner:"-"` 减少传输量
- 临时计算字段不需要在API中暴露

### 3. 维护性
- 使用清晰的注释说明为什么某个字段被忽略
- 定期审查哪些字段需要忽略

### 4. 一致性
- 在同一个项目中保持一致的忽略策略
- 团队成员都应该了解这个功能的使用规范

## 总结

`runner:"-"` 字段忽略功能提供了一种简单而强大的方式来控制哪些字段应该在API响应中暴露。它具有以下优势：

- ✅ **安全性**：防止敏感信息泄露
- ✅ **性能**：减少不必要的数据传输
- ✅ **简洁性**：保持API响应的简洁性
- ✅ **灵活性**：支持多种标签组合场景
- ✅ **可靠性**：在解析阶段就过滤，确保字段不会意外暴露

通过合理使用这个功能，可以构建更安全、更高效的API接口。 