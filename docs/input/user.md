# input 组件标签参数说明（面向人类/业务）

---

## runner 标签
| 参数名 | 说明         | 是否必填 | 默认值 | 示例           |
|--------|--------------|----------|--------|----------------|
| code   | 字段唯一英文标识 | 是       | 无     | code:name      |
| name   | 字段中文名   | 是       | 无     | name:用户名    |

## widget 标签
| 参数名      | 说明                   | 是否必填 | 默认值    | 示例                        |
|-------------|------------------------|----------|-----------|-----------------------------|
| type        | 组件类型，固定为 input | 是       | input     | type:input                  |
| placeholder | 占位提示               | 否       | 无        | placeholder:请输入用户名     |
| prefix      | 前缀符号               | 否       | 无        | prefix:￥                   |
| suffix      | 后缀符号               | 否       | 无        | suffix:%                    |
| mode        | 特殊模式（可选：line_text=单行文本，text_area=多行文本，password=密码输入） | 否 | line_text | mode:text_area              |

## data 标签
| 参数名       | 说明                   | 是否必填 | 示例                        |
|--------------|------------------------|----------|-----------------------------|
| type         | 数据类型               | 是       | type:string                 |
| default_value| 默认值                 | 否       | default_value:张三          |
| example      | 示例值                 | 否       | example:李四                |

## validate 标签
| 参数名   | 说明                   | 是否必填 | 示例                        |
|----------|------------------------|----------|-----------------------------|
| required | 是否必填               | 否       | required                    |
| min      | 最小长度               | 否       | min=2                       |
| max      | 最大长度               | 否       | max=20                      |
| email    | 邮箱格式校验           | 否       | email                       |

---
# input 组件表单函数示例

## 业务场景
用户输入一个字符串，返回其反转结果。若输入为空，直接返回错误。

## 示例代码

```go
package demo

import (
    "fmt"
    "github.com/yunhanshu-net/function-go/pkg/dto/response"
    "github.com/yunhanshu-net/function-go/runner"
)

// 请求结构体
type ReverseReq struct {
    Text string `json:"text" runner:"code:text;name:待反转内容" widget:"type:input;placeholder:请输入内容" data:"type:string;default_value:hello" validate:"required"`
}

// 响应结构体
type ReverseResp struct {
    Result string `json:"result" runner:"code:result;name:反转结果" widget:"type:input;mode:text_area" data:"type:string"`
}

// 业务处理函数
func ReverseHandler(ctx *runner.Context, req *ReverseReq, resp response.Response) error {
    if req.Text == "" {
        // 校验失败，直接返回错误，前端自动弹窗
        return fmt.Errorf("内容不能为空")
    }
    reversed := reverseString(req.Text)
    return resp.Form(&ReverseResp{
        Result: reversed,
    }).Build()
}

// 字符串反转
func reverseString(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

// 注册API
var ReverseOption = &runner.FunctionOptions{
    Tags:        []string{"input", "字符串", "表单演示"},
    EnglishName: "reverse_demo",
    ChineseName: "字符串反转演示",
    ApiDesc:     "输入内容，返回其反转结果。内容不能为空。",
    Request:     &ReverseReq{},
    Response:    &ReverseResp{},
    RenderType:  response.RenderTypeForm,
}

func init() {
    runner.Post("/demo/reverse", ReverseHandler, ReverseOption)
}

## 规范说明
- 参数校验失败时直接 `return error`，前端自动弹窗。
- 业务正常时用 `resp.Form(...)` 返回结果。
- 标签用法与 input 组件规范一致。

## 标准代码块示例

```go
// 用户信息输入结构体标准示例
// Req 结尾标识输入对象，Resp 结尾标识输出对象

type UserInputReq struct {
    // 用户名输入框（单行文本，mode 可省略或为 line_text）
    Username string `json:"username" runner:"code:username;name:用户名" widget:"type:input;placeholder:请输入用户名;mode:line_text" data:"type:string;default_value:张三;example:李四" validate:"required,min=2,max=20"`
    // 密码输入框
    Password string `json:"password" runner:"code:password;name:密码" widget:"type:input;placeholder:请输入密码;mode:password" data:"type:string" validate:"required,min=6"`
    // 多行文本输入框
    Description string `json:"description" runner:"code:description;name:描述" widget:"type:input;mode:text_area;placeholder:请输入描述" data:"type:string"`
}
``` 