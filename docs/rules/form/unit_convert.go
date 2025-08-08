// 单位换算工具
// 输入数值和单位（米/厘米/毫米），输出换算为米的结果

package form

import (
	"fmt"
	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/function-go/runner"
)

// 请求结构体
// 数值：0.01-10000，单位：米/厘米/毫米
// example: 数值100，单位"厘米"
type UnitConvertReq struct {
    Value  float64 `json:"value" runner:"code:value;name:数值" widget:"type:number;min:0.01;max:10000;step:0.01;placeholder:请输入数值" data:"type:float;default_value:100;example:100" validate:"required,min=0.01,max=10000"`
    Unit   string  `json:"unit" runner:"code:unit;name:单位" widget:"type:select;options:米,厘米,毫米;placeholder:请选择单位" data:"type:string;default_value:厘米;example:厘米" validate:"required,oneof=米 厘米 毫米"`
}

// 响应结构体
// example: 1
// 结果单位为米
type UnitConvertResp struct {
    Result float64 `json:"result" runner:"code:result;name:换算结果(米)" widget:"type:number;precision:4;placeholder:自动计算" data:"type:float;example:1"`
}

// 业务处理函数
func UnitConvertHandler(ctx *runner.Context, req *UnitConvertReq, resp response.Response) error {
    var valueInMeter float64
    switch req.Unit {
    case "米":
        valueInMeter = req.Value
    case "厘米":
        valueInMeter = req.Value / 100
    case "毫米":
        valueInMeter = req.Value / 1000
    default:
        return fmt.Errorf("不支持的单位类型")
    }
    return resp.Form(&UnitConvertResp{Result: valueInMeter}).Build()
}

// API注册
var UnitConvertOption = &runner.FunctionOptions{
    ChineseName: "单位换算工具",
    EnglishName: "unit_convert",
    ApiDesc:     "输入数值和单位（米/厘米/毫米），自动换算为米。",
    Request:     &UnitConvertReq{},
    Response:    &UnitConvertResp{},
    RenderType:  response.RenderTypeForm,
    Tags:        []string{"单位", "换算", "工具"},
}

func init() {
    runner.Post("/form/unit_convert", UnitConvertHandler, UnitConvertOption)
} 