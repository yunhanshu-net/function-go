# 项目函数/示例/训练数据注意事项

## 1. 错误用法反例

### 反例：用 resp.Form 返回错误信息（禁止）
```go
// 错误示范：不要用 resp.Form 返回错误
if req.Text == "" {
    return resp.Form(&ReverseResp{
        Result: "",
        Message: "内容不能为空",
    }).Build() // 错误！应直接 return error
}
```

### 正确写法
```go
if req.Text == "" {
    // 错误信息必须用户友好、易懂，便于前端用户理解，也便于大模型自动排障
    return fmt.Errorf("内容不能为空，请输入内容")
}
```

## 2. resp.Form 正确用法标准代码块

- 仅用于返回正常业务结果或提示，不用于错误场景。
- 推荐完整示例：包含 struct、handler、注册，便于标准参考。

```go
package main

import (
    "fmt"
    "github.com/function-go/runner"
    "github.com/function-go/pkg/dto/response"
)

func init() {
    runner.Post("/demo/reverse", ReverseHandler, ReverseOption)
}

// 请求结构体
type ReverseReq struct {
    Text string `json:"text" runner:"code:text;name:待反转内容" widget:"type:input;placeholder:请输入内容" data:"type:string;default_value:hello" validate:"required"`
}

// 响应结构体
type ReverseResp struct {
    Result string `json:"result" runner:"code:result;name:反转结果" widget:"type:input;mode:text_area" data:"type:string"`
    Message string `json:"message" runner:"code:message;name:提示信息" widget:"type:input" data:"type:string"`
}

// 字符串反转
func reverseString(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

// 业务处理函数
func ReverseHandler(ctx *runner.Context, req *ReverseReq, resp response.Response) error {
    if req.Text == "" {
        return fmt.Errorf("内容不能为空，请输入内容")
    }
    reversed := reverseString(req.Text)
    return resp.Form(&ReverseResp{
        Result: reversed,
        Message: "反转成功！",
    }).Build()
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
```

## 3. 训练示例禁止项
- 严禁伪代码、逻辑不通、无法编译的代码。
- 示例必须真实可运行，结构体、handler、注册、注释齐全。
- 代码必须加必要注释，便于理解和后续维护。

## 4. 错误信息规范
- 错误输出要详细、准确，**必须让前端用户能看懂**，同时便于大模型和人工排障。
- 例如：
```go
return fmt.Errorf("起始位不能大于截止位，请检查输入")
```
- 不要只返回"失败"或"出错"，要说明具体原因，且用语通俗。

## 5. 其它注意事项
- 所有函数示例必须加注释，说明业务意图、参数含义、边界情况。
- 代码风格统一，标签用法规范。
- 后续如有新要求，持续补充本文档。 