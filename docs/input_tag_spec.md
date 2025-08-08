# input 组件标签参数说明与标准示例

---

## 一、人类查阅版（表格+注释）

### runner 标签
| 参数名 | 说明         | 是否必填 | 可选值/格式 | 示例           |
|--------|--------------|----------|-------------|----------------|
| code   | 字段唯一英文标识 | 是       | 英文小写+下划线 | code:name      |
| name   | 字段中文名   | 是       | 任意中文     | name:用户名    |

### widget 标签
| 参数名      | 说明                   | 是否必填 | 可选值/格式           | 示例                        |
|-------------|------------------------|----------|-----------------------|-----------------------------|
| type        | 组件类型，固定为 input | 是       | input                 | type:input                  |
| placeholder | 占位提示               | 否       | 任意字符串            | placeholder:请输入用户名     |
| prefix      | 前缀符号               | 否       | 任意字符串            | prefix:￥                   |
| suffix      | 后缀符号               | 否       | 任意字符串            | suffix:%                    |
| mode        | 特殊模式               | 否       | text_area等           | mode:text_area              |

### data 标签
| 参数名       | 说明                   | 是否必填 | 可选值/格式           | 示例                        |
|--------------|------------------------|----------|-----------------------|-----------------------------|
| type         | 数据类型               | 是       | string/number/float/boolean | type:string      |
| default_value| 默认值                 | 否       | 与type/options一致    | default_value:张三          |
| example      | 示例值                 | 否       | 任意                  | example:李四                |
| source       | 数据源（动态选项）     | 否       | api://xxx             | source:api://users          |
| format       | 格式（仅业务日期字段） | 否       | yyyy-MM-dd等          | format:yyyy-MM-dd           |

### validate 标签
| 参数名   | 说明                   | 是否必填 | 可选值/格式           | 示例                        |
|----------|------------------------|----------|-----------------------|-----------------------------|
| required | 是否必填               | 否       | required              | required                    |
| min      | 最小值/最小长度        | 否       | min=数字              | min=1                       |
| max      | 最大值/最大长度        | 否       | max=数字              | max=20                      |
| email    | 邮箱格式校验           | 否       | email                 | email                       |
| ...      | 其他go-playground/validator支持的参数 | 否 | 见官方文档           |                             |

---

## 二、适合大模型Prompt/训练的规则说明

【规则说明】
- default_value 只能在 data 标签，不能在 widget 标签。
- options 只能在 widget 标签，不能在 data 标签。
- 每个字段必须有 runner、widget、data 标签。
- type 必须在 data 标签，且与字段Go类型一致。
- validate 只用于字段级校验。
- placeholder、prefix、suffix、mode 只在 widget 标签。
- example/source/format 只在 data 标签。

【错误示例】
Name string `json:"name" widget:"type:input;default_value:张三"`

【正确示例】
Name string `json:"name" runner:"code:name;name:姓名" widget:"type:input;placeholder:请输入姓名" data:"type:string;default_value:张三" validate:"required,min=2,max=10"`

---

## 三、标准代码块示例

```go
// 用户名输入框
Username string `json:"username" runner:"code:username;name:用户名" widget:"type:input;placeholder:请输入用户名" data:"type:string;default_value:张三;example:李四" validate:"required,min=2,max=20"`

// 邮箱输入框（带邮箱校验）
Email string `json:"email" runner:"code:email;name:邮箱" widget:"type:input;placeholder:请输入邮箱" data:"type:string;example:test@example.com" validate:"required,email"`

// 多行文本输入框
Description string `json:"description" runner:"code:description;name:描述" widget:"type:input;mode:text_area;placeholder:请输入描述" data:"type:string"`
``` 