package widgets

import (
	"fmt"
	"time"
	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/function-go/runner"
	"github.com/yunhanshu-net/pkg/query"
	"gorm.io/gorm"
)

// Order 订单数据模型
type Order struct {
	ID          int            `json:"id" gorm:"primaryKey;autoIncrement" runner:"code:id;name:订单ID" permission:"read"`
	CreatedAt   int64          `json:"created_at" gorm:"autoCreateTime:milli" runner:"code:created_at;name:下单时间" widget:"type:datetime;kind:datetime" data:"type:number;example:1705292200000" search:"gte,lte" permission:"read"`
	UpdatedAt   int64          `json:"updated_at" gorm:"autoUpdateTime:milli" runner:"code:updated_at;name:更新时间" widget:"type:datetime;kind:datetime" data:"type:number;example:1705292200000" search:"gte,lte" permission:"read"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index" runner:"-"`
	OrderNo     string         `json:"order_no" gorm:"uniqueIndex;size:32" runner:"code:order_no;name:订单号" data:"type:string;example:ORD202401150001" search:"like" permission:"read"` // 自动生成，只读
	CustomerName string        `json:"customer_name" gorm:"size:50" runner:"code:customer_name;name:客户姓名" widget:"type:input;placeholder:请输入客户姓名" data:"type:string;example:张三" search:"like" validate:"required,min=2,max=20"`
	CustomerPhone string       `json:"customer_phone" gorm:"size:20" runner:"code:customer_phone;name:联系电话" widget:"type:input;placeholder:请输入联系电话" data:"type:string;example:13800138000" search:"like" validate:"required,len=11"`
	ProductName  string        `json:"product_name" gorm:"size:100" runner:"code:product_name;name:商品名称" widget:"type:input;placeholder:请输入商品名称" data:"type:string;example:iPhone 15 Pro" search:"like" validate:"required,min=2,max=50"`
	Quantity     int           `json:"quantity" runner:"code:quantity;name:数量" widget:"type:number;min:1;unit:件" data:"type:number;example:2" search:"eq,gte,lte" validate:"required,min=1"`
	UnitPrice    float64       `json:"unit_price" runner:"code:unit_price;name:单价" widget:"type:number;min:0;precision:2;prefix:￥" data:"type:float;example:8999.00" search:"eq,gte,lte" validate:"required,min=0"`
	TotalAmount  float64       `json:"total_amount" runner:"code:total_amount;name:总金额" widget:"type:number;precision:2;prefix:￥" data:"type:float;example:17998.00" search:"eq,gte,lte" validate:"required,min=0"`
	Status       string        `json:"status" runner:"code:status;name:订单状态" widget:"type:select;options:待付款,已付款,已发货,已完成,已取消;placeholder:请选择状态" data:"type:string;example:已付款" search:"eq,in" validate:"required,oneof=待付款 已付款 已发货 已完成 已取消"`
	PaymentMethod string       `json:"payment_method" runner:"code:payment_method;name:支付方式" widget:"type:select;options:微信支付,支付宝,银行卡,现金;placeholder:请选择支付方式" data:"type:string;example:微信支付" search:"eq,in" validate:"required,oneof=微信支付 支付宝 银行卡 现金"`
	Remark       string        `json:"remark" gorm:"size:200" runner:"code:remark;name:备注" widget:"type:input;mode:text_area;placeholder:请输入备注信息" data:"type:string;example:客户要求尽快发货" search:"like"`
}

// BeforeCreate GORM钩子函数：创建订单前自动生成订单号
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	// 生成订单号：ORD + 年月日 + 4位序号
	now := time.Now()
	dateStr := now.Format("20060102")
	
	// 查询当天最大序号
	var maxOrder Order
	tx.Where("order_no LIKE ?", "ORD"+dateStr+"%").Order("order_no DESC").First(&maxOrder)
	
	// 生成新序号
	var sequence int
	if maxOrder.OrderNo != "" {
		// 从现有订单号中提取序号
		fmt.Sscanf(maxOrder.OrderNo, "ORD%s%04d", &dateStr, &sequence)
		sequence++
	} else {
		sequence = 1
	}
	
	// 生成新订单号
	o.OrderNo = fmt.Sprintf("ORD%s%04d", dateStr, sequence)
	return nil
}

// BeforeSave GORM钩子函数：保存前自动计算总金额
func (o *Order) BeforeSave(tx *gorm.DB) error {
	// 自动计算总金额
	o.TotalAmount = float64(o.Quantity) * o.UnitPrice
	return nil
}

func (o *Order) TableName() string {
	return "order"
}

// OrderListReq 请求结构体
type OrderListReq struct {
	query.PageInfoReq `runner:"-"`
}

// OrderList 订单列表处理逻辑
func OrderList(ctx *runner.Context, req *OrderListReq, resp response.Response) error {
	db := ctx.MustGetOrInitDB()
	var orders []Order
	return resp.Table(&orders).AutoPaginated(db, &Order{}, &req.PageInfoReq).Build()
}

// OrderListOption 注册API
var OrderListOption = &runner.FunctionOptions{
	Tags:          []string{"订单管理", "表格演示"},
	EnglishName:   "order_list",
	ChineseName:   "订单列表",
	ApiDesc:       "展示订单列表，支持分页、搜索和CRUD操作。订单号自动生成，总金额自动计算。",
	Request:       &OrderListReq{},
	Response:      query.PaginatedTable[[]Order]{},
	RenderType:    response.RenderTypeTable,
	CreateTables:  []interface{}{&Order{}},
	AutoCrudTable: &Order{},
}

func init() {
	runner.Get("/widgets/order_list", OrderList, OrderListOption)
} 