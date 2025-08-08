// 字符串反转工具
// 输入字符串，输出其反转结果

package form

import (
	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/function-go/runner"
)

// 请求结构体
// 只需一个输入框，必填，长度2-20
// example: "hello"
type ReverseReq struct {
    Text string `json:"text" runner:"code:text;name:待反转内容" widget:"type:input;placeholder:请输入内容" data:"type:string;default_value:hello;example:hello" validate:"required,min=2,max=20"`
}

// 响应结构体
// example: "olleh"
type ReverseResp struct {
    Result string `json:"result" runner:"code:result;name:反转结果" widget:"type:input;mode:text_area" data:"type:string;example:olleh"`
}

// 业务处理函数
func ReverseHandler(ctx *runner.Context, req *ReverseReq, resp response.Response) error {
    // 只做业务处理，不做参数校验
    runes := []rune(req.Text)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    reversed := string(runes)
    return resp.Form(&ReverseResp{Result: reversed}).Build()
}

// API注册
var ReverseOption = &runner.FunctionOptions{
    ChineseName: "字符串反转工具",
    EnglishName: "reverse",
    ApiDesc:     "输入字符串，输出其反转结果。内容不能为空，长度2-20。",
    Request:     &ReverseReq{},
    Response:    &ReverseResp{},
    RenderType:  response.RenderTypeForm,
    Tags:        []string{"字符串", "工具", "反转"},
}

func init() {
    runner.Post("/form/reverse", ReverseHandler, ReverseOption)
} 