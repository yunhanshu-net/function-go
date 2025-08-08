package widgets

import (
	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/function-go/runner"
	"github.com/yunhanshu-net/pkg/query"
	"gorm.io/gorm"
)

// Customer 客户数据模型
type Customer struct {
	ID           int            `json:"id" gorm:"primaryKey;autoIncrement" runner:"code:id;name:客户ID" permission:"read"`
	CreatedAt    int64          `json:"created_at" gorm:"autoCreateTime:milli" runner:"code:created_at;name:创建时间" widget:"type:datetime;kind:datetime" data:"type:number;example:1705292200000" search:"gte,lte" permission:"read"`
	UpdatedAt    int64          `json:"updated_at" gorm:"autoUpdateTime:milli" runner:"code:updated_at;name:更新时间" widget:"type:datetime;kind:datetime" data:"type:number;example:1705292200000" search:"gte,lte" permission:"read"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index" runner:"-"`
	CustomerNo   string         `json:"customer_no" gorm:"uniqueIndex;size:20" runner:"code:customer_no;name:客户编号" widget:"type:input;placeholder:请输入客户编号" data:"type:string;example:CUST001" search:"like" validate:"required,min=3,max=20"`
	Name         string         `json:"name" gorm:"size:50" runner:"code:name;name:客户姓名" widget:"type:input;placeholder:请输入客户姓名" data:"type:string;example:李四" search:"like" validate:"required,min=2,max=20"`
	Gender       string         `json:"gender" gorm:"size:10" runner:"code:gender;name:性别" widget:"type:select;options:男,女;placeholder:请选择性别" data:"type:string;example:男" search:"eq" validate:"required,oneof=男 女"`
	Age          int            `json:"age" runner:"code:age;name:年龄" widget:"type:number;min:18;max:80;unit:岁" data:"type:number;example:35" search:"eq,gte,lte" validate:"required,min=18,max=80"`
	Phone        string         `json:"phone" gorm:"size:20" runner:"code:phone;name:手机号" widget:"type:input;placeholder:请输入手机号" data:"type:string;example:13900139000" search:"like" validate:"required,len=11"`
	Email        string         `json:"email" gorm:"size:100" runner:"code:email;name:邮箱" widget:"type:input;placeholder:请输入邮箱" data:"type:string;example:lisi@email.com" search:"like" validate:"required,email"`
	Wechat       string         `json:"wechat" gorm:"size:50" runner:"code:wechat;name:微信号" widget:"type:input;placeholder:请输入微信号" data:"type:string;example:lisi_wechat" search:"like"`
	Company      string         `json:"company" gorm:"size:100" runner:"code:company;name:公司名称" widget:"type:input;placeholder:请输入公司名称" data:"type:string;example:某某科技有限公司" search:"like"`
	Position     string         `json:"position" gorm:"size:50" runner:"code:position;name:职位" widget:"type:input;placeholder:请输入职位" data:"type:string;example:技术总监" search:"like"`
	Level        string         `json:"level" runner:"code:level;name:客户等级" widget:"type:select;options:普通客户,银卡客户,金卡客户,钻石客户;placeholder:请选择等级" data:"type:string;example:金卡客户" search:"eq,in" validate:"required,oneof=普通客户 银卡客户 金卡客户 钻石客户"`
	Source       string         `json:"source" runner:"code:source;name:来源渠道" widget:"type:select;options:官网,搜索引擎,朋友推荐,广告投放,展会,其他;placeholder:请选择来源" data:"type:string;example:朋友推荐" search:"eq,in" validate:"required,oneof=官网 搜索引擎 朋友推荐 广告投放 展会 其他"`
	Address      string         `json:"address" gorm:"size:200" runner:"code:address;name:地址" widget:"type:input;mode:text_area;placeholder:请输入地址" data:"type:string;example:北京市海淀区中关村大街1号" search:"like"`
	Status       string         `json:"status" runner:"code:status;name:客户状态" widget:"type:select;options:潜在客户,意向客户,成交客户,流失客户;placeholder:请选择状态" data:"type:string;example:意向客户" search:"eq,in" validate:"required,oneof=潜在客户 意向客户 成交客户 流失客户"`
	Remark       string         `json:"remark" gorm:"size:200" runner:"code:remark;name:备注" widget:"type:input;mode:text_area;placeholder:请输入备注" data:"type:string;example:对产品功能很感兴趣，需要跟进" search:"like"`
}

func (c *Customer) TableName() string {
	return "customer"
}

// CustomerListReq 请求结构体
type CustomerListReq struct {
	query.PageInfoReq `runner:"-"`
}

// CustomerList 客户列表处理逻辑
func CustomerList(ctx *runner.Context, req *CustomerListReq, resp response.Response) error {
	db := ctx.MustGetOrInitDB()
	var customers []Customer
	return resp.Table(&customers).AutoPaginated(db, &Customer{}, &req.PageInfoReq).Build()
}

// CustomerListOption 注册API
var CustomerListOption = &runner.FunctionOptions{
	Tags:          []string{"客户管理", "CRM", "表格演示"},
	EnglishName:   "customer_list",
	ChineseName:   "客户列表",
	ApiDesc:       "展示客户列表，支持分页、搜索和CRUD操作。",
	Request:       &CustomerListReq{},
	Response:      query.PaginatedTable[[]Customer]{},
	RenderType:    response.RenderTypeTable,
	CreateTables:  []interface{}{&Customer{}},
	AutoCrudTable: &Customer{},
}

func init() {
	runner.Get("/widgets/customer_list", CustomerList, CustomerListOption)
} 