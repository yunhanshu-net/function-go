// 订单价格计算器
// 输入商品名称、单价、数量、折扣类型，输出总价

package form

import (
	"fmt"
	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/function-go/runner"
)

// 请求结构体
// 商品名称：必填，2-20字
// 单价：1-10000元
// 数量：1-1000
// 折扣类型：无折扣/9折/8折
// example: 商品名称"苹果手机"，单价6999，数量2，折扣类型"9折"
type OrderPriceReq struct {
    ProductName string  `json:"product_name" runner:"code:product_name;name:商品名称" widget:"type:input;placeholder:请输入商品名称" data:"type:string;default_value:苹果手机;example:苹果手机" validate:"required,min=2,max=20"`
    UnitPrice   float64 `json:"unit_price" runner:"code:unit_price;name:单价(元)" widget:"type:number;min:1;max:10000;step:0.01;unit:元;placeholder:请输入单价" data:"type:float;default_value:6999;example:6999" validate:"required,min=1,max=10000"`
    Quantity    int     `json:"quantity" runner:"code:quantity;name:数量" widget:"type:number;min:1;max:1000;unit:件;placeholder:请输入数量" data:"type:number;default_value:2;example:2" validate:"required,min=1,max=1000"`
    Discount    string  `json:"discount" runner:"code:discount;name:折扣类型" widget:"type:select;options:无折扣,9折,8折;placeholder:请选择折扣类型" data:"type:string;default_value:无折扣;example:9折" validate:"required,oneof=无折扣 9折 8折"`
}

// 响应结构体
// example: 12598.2
type OrderPriceResp struct {
    TotalPrice float64 `json:"total_price" runner:"code:total_price;name:总价(元)" widget:"type:number;precision:2;placeholder:自动计算" data:"type:float;example:12598.2"`
}

// 业务处理函数
func OrderPriceHandler(ctx *runner.Context, req *OrderPriceReq, resp response.Response) error {
    var discountRate float64 = 1.0
    switch req.Discount {
    case "9折":
        discountRate = 0.9
    case "8折":
        discountRate = 0.8
    }
    total := req.UnitPrice * float64(req.Quantity) * discountRate
    // 保留两位小数
    total = float64(int(total*100+0.5)) / 100
    return resp.Form(&OrderPriceResp{TotalPrice: total}).Build()
}

// API注册
var OrderPriceOption = &runner.FunctionOptions{
    ChineseName: "订单价格计算器",
    EnglishName: "order_price",
    ApiDesc:     "输入商品名称、单价、数量、折扣类型，自动计算订单总价。",
    Request:     &OrderPriceReq{},
    Response:    &OrderPriceResp{},
    RenderType:  response.RenderTypeForm,
    Tags:        []string{"订单", "价格", "计算器"},
}

func init() {
    runner.Post("/form/order_price", OrderPriceHandler, OrderPriceOption)
} 