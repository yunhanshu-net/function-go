# input 组件标签参数规则（面向大模型/自动化）

---

## 1. runner 标签
- 必须包含 code、name 两个参数。
- code：字段唯一英文标识，英文小写+下划线，必填。
- name：字段中文名，必填。
- 只能出现一次，不能有其他参数。

【正例】
```go
// 完整结构体示例
type RunnerReq struct {
    Username string `runner:"code:username;name:用户名" widget:"type:input" data:"type:string" validate:"required"`
}
```
【反例】
```go
// 缺少 code，其他标签都规范
type RunnerReq struct {
    Username string `runner:"name:用户名" widget:"type:input" data:"type:string" validate:"required"` // 缺少code
}
// 缺少 name，其他标签都规范
type RunnerReq struct {
    Username string `runner:"code:username" widget:"type:input" data:"type:string" validate:"required"` // 缺少name
}
```

---

## 2. widget 标签
- 必须包含 type 参数，且为 input。
- mode 参数可选，默认值为 line_text，可选值：line_text、text_area、password。
- placeholder、prefix、suffix 可选。
- 不能包含 default_value、example、type 等数据相关参数。

【正例】
```go
// 完整结构体示例
type WidgetReq struct {
    Username string `runner:"code:username;name:用户名" widget:"type:input;placeholder:请输入用户名;mode:line_text" data:"type:string" validate:"required"`
    Password string `runner:"code:password;name:密码" widget:"type:input;mode:password" data:"type:string" validate:"required,min=6"`
    Description string `runner:"code:description;name:描述" widget:"type:input;mode:text_area;placeholder:请输入描述" data:"type:string"`
}
```
【反例】
```go
// default_value 错误位置，其他标签都规范
type WidgetReq struct {
    Username string `runner:"code:username;name:用户名" widget:"type:input;default_value:张三" data:"type:string" validate:"required"` // default_value 错误位置
}
// mode 非法枚举值，其他标签都规范
type WidgetReq struct {
    Username string `runner:"code:username;name:用户名" widget:"type:input;mode:email" data:"type:string" validate:"required"` // mode 非法枚举值
}
```

---

## 3. data 标签
- 必须包含 type 参数，且与Go字段类型一致。
- default_value、example 可选，default_value 必须与 type/options 保持一致。
- source、format 仅在有动态数据源/日期字段时使用。
- 不能包含 type:input、placeholder、mode 等UI参数。

【正例】
```go
// 完整结构体示例
type DataReq struct {
    Username string `runner:"code:username;name:用户名" widget:"type:input" data:"type:string;default_value:张三;example:李四" validate:"required"`
}
```
【反例】
```go
// type 错误，其他标签都规范
type DataReq struct {
    Username string `runner:"code:username;name:用户名" widget:"type:input" data:"type:input" validate:"required"` // type 错误
}
// mode 错误位置，其他标签都规范
type DataReq struct {
    Username string `runner:"code:username;name:用户名" widget:"type:input" data:"type:string;mode:text_area" validate:"required"` // mode 错误位置
}
```

---

## 4. validate 标签
- 只用于字段级校验。
- 支持 required、min、max、email 等 go-playground/validator 语法。
- 不能包含 UI 或数据相关参数。

【正例】
```go
// 完整结构体示例
type ValidateReq struct {
    Username string `runner:"code:username;name:用户名" widget:"type:input" data:"type:string" validate:"required,min=2,max=20"`
    Email string `runner:"code:email;name:邮箱" widget:"type:input" data:"type:string" validate:"required,email"`
}
```
【反例】
```go
// type 错误位置，其他标签都规范
type ValidateReq struct {
    Username string `runner:"code:username;name:用户名" widget:"type:input" data:"type:string" validate:"type:input"` // type 错误位置
}
```

---

## 5. 组合正例
```go
// 完整结构体示例
type InputAllTagReq struct {
    // 用户名输入框（单行文本，mode 可省略或为 line_text）
    Username string `json:"username" runner:"code:username;name:用户名" widget:"type:input;placeholder:请输入用户名;mode:line_text" data:"type:string;default_value:张三;example:李四" validate:"required,min=2,max=20"`
    // 密码输入框
    Password string `json:"password" runner:"code:password;name:密码" widget:"type:input;placeholder:请输入密码;mode:password" data:"type:string" validate:"required,min=6"`
    // 多行文本输入框
    Description string `json:"description" runner:"code:description;name:描述" widget:"type:input;mode:text_area;placeholder:请输入描述" data:"type:string"`
}
```

---

## 6. 其他规则
- 每个字段必须有 runner、widget、data 标签。
- 所有参数顺序建议：runner、widget、data、validate。
- 任何违反上述规则的写法都视为不合规。
- 生成时如有疑问，优先参考正例。 