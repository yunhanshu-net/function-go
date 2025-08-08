// BMI（身体质量指数）计算器
// 输入身高（cm）和体重（kg），输出BMI值

package form

import (
	"fmt"
	"math"

	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/function-go/runner"
)

// 请求结构体
// 身高：100-250cm，体重：20-200kg
// example: 身高170，体重65
// BMI = 体重(kg) / (身高(m) * 身高(m))
type BMIReq struct {
    Height float64 `json:"height" runner:"code:height;name:身高(cm)" widget:"type:number;min:100;max:250;step:0.1;unit:cm;placeholder:请输入身高" data:"type:float;default_value:170;example:170" validate:"required,min=100,max=250"`
    Weight float64 `json:"weight" runner:"code:weight;name:体重(kg)" widget:"type:number;min:20;max:200;step:0.1;unit:kg;placeholder:请输入体重" data:"type:float;default_value:65;example:65" validate:"required,min=20,max=200"`
}

// 响应结构体
// example: 22.49
type BMIResp struct {
    BMI float64 `json:"bmi" runner:"code:bmi;name:BMI值" widget:"type:number;precision:2;placeholder:自动计算" data:"type:float;example:22.49"`
}

// 业务处理函数
func BMIHandler(ctx *runner.Context, req *BMIReq, resp response.Response) error {
    // BMI = 体重(kg) / (身高(m) * 身高(m))
    if req.Height <= 0 {
        return fmt.Errorf("身高必须大于0")
    }
    h := req.Height / 100.0
    bmi := req.Weight / (h * h)
    bmi = math.Round(bmi*100) / 100 // 保留两位小数
    return resp.Form(&BMIResp{BMI: bmi}).Build()
}

// API注册
var BMIOption = &runner.FunctionOptions{
    ChineseName: "BMI计算器",
    EnglishName: "bmi",
    ApiDesc:     "输入身高（cm）和体重（kg），自动计算BMI值。",
    Request:     &BMIReq{},
    Response:    &BMIResp{},
    RenderType:  response.RenderTypeForm,
    Tags:        []string{"健康", "工具", "BMI"},
}

func init() {
    runner.Post("/form/bmi", BMIHandler, BMIOption)
} 