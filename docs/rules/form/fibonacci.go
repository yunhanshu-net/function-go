// 斐波那契数列生成器
// 输入起始位、截止位、分隔符，输出对应区间的斐波那契数列字符串

package form

import (
	"fmt"
	"strings"

	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/function-go/runner"
)

// 请求结构体
// 所有数值范围校验全部用validate标签，widget只做UI
// 起始位、截止位：1-1000
// 分隔符：必填
// 其他业务校验在业务逻辑中处理
type FibonacciReq struct {
    Start     int    `json:"start" runner:"code:start;name:起始位" widget:"type:number;placeholder:请输入起始位" data:"type:number;default_value:1;example:1" validate:"required,min=1,max=1000"`
    End       int    `json:"end" runner:"code:end;name:截止位" widget:"type:number;placeholder:请输入截止位" data:"type:number;default_value:3;example:3" validate:"required,min=1,max=1000"`
    Separator string `json:"separator" runner:"code:separator;name:分隔符" widget:"type:input;placeholder:请输入分隔符" data:"type:string;default_value:,;example:," validate:"required"`
}

// 响应结构体
type FibonacciResp struct {
    Result string `json:"result" runner:"code:result;name:结果" widget:"type:input;mode:text_area" data:"type:string;example:1,1,2"`
}

// 业务处理函数
func FibonacciHandler(ctx *runner.Context, req *FibonacciReq, resp response.Response) error {
	// 参数校验
	if req.Start > req.End {
		return fmt.Errorf("起始位不能大于截止位，请检查输入")
	}
	if req.End-req.Start > 1000 {
		return fmt.Errorf("区间长度不能超过1000")
	}
	// 生成斐波那契数列
	fibs := fibonacci(req.End)
	if req.Start < 1 || req.End > len(fibs) {
		return fmt.Errorf("起始位或截止位超出范围")
	}
	sub := fibs[req.Start-1 : req.End]
	// 拼接为字符串
	strArr := make([]string, len(sub))
	for i, v := range sub {
		strArr[i] = fmt.Sprintf("%d", v)
	}
	result := strings.Join(strArr, req.Separator)
	return resp.Form(&FibonacciResp{Result: result}).Build()
}

// 斐波那契数列生成工具
func fibonacci(n int) []int {
	if n <= 0 {
		return []int{}
	}
	fibs := make([]int, n)
	fibs[0] = 1
	if n > 1 {
		fibs[1] = 1
		for i := 2; i < n; i++ {
			fibs[i] = fibs[i-1] + fibs[i-2]
		}
	}
	return fibs
}

// API注册
var FibonacciOption = &runner.FunctionOptions{
	ChineseName: "斐波那契数列生成器",
	EnglishName: "fibonacci",
	ApiDesc:     "输入起始位、截止位、分隔符，输出对应区间的斐波那契数列字符串",
	Request:     &FibonacciReq{},
	Response:    &FibonacciResp{},
	RenderType:  response.RenderTypeForm,
	Tags:        []string{"数学", "工具", "斐波那契"},
}

func init() {
	runner.Post("/form/fibonacci", FibonacciHandler, FibonacciOption)
}
