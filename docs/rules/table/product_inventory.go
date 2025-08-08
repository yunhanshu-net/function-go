package widgets

import (
	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/function-go/runner"
	"github.com/yunhanshu-net/pkg/query"
	"gorm.io/gorm"
)

// Product 产品数据模型
type Product struct {
	ID          int            `json:"id" gorm:"primaryKey;autoIncrement" runner:"code:id;name:产品ID" permission:"read"`
	CreatedAt   int64          `json:"created_at" gorm:"autoCreateTime:milli" runner:"code:created_at;name:创建时间" widget:"type:datetime;kind:datetime" data:"type:number;example:1705292200000" search:"gte,lte" permission:"read"`
	UpdatedAt   int64          `json:"updated_at" gorm:"autoUpdateTime:milli" runner:"code:updated_at;name:更新时间" widget:"type:datetime;kind:datetime" data:"type:number;example:1705292200000" search:"gte,lte" permission:"read"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index" runner:"-"`
	ProductCode string         `json:"product_code" gorm:"uniqueIndex;size:50" runner:"code:product_code;name:产品编码" widget:"type:input;placeholder:请输入产品编码" data:"type:string;example:PROD001" search:"like" validate:"required,min=3,max=50"`
	Name        string         `json:"name" gorm:"size:100" runner:"code:name;name:产品名称" widget:"type:input;placeholder:请输入产品名称" data:"type:string;example:iPhone 15 Pro 256GB" search:"like" validate:"required,min=2,max=100"`
	Category    string         `json:"category" gorm:"size:50" runner:"code:category;name:产品分类" widget:"type:select;options:手机,电脑,配件,服装,食品,家居;placeholder:请选择分类" data:"type:string;example:手机" search:"eq,in" validate:"required,oneof=手机 电脑 配件 服装 食品 家居"`
	Brand       string         `json:"brand" gorm:"size:50" runner:"code:brand;name:品牌" widget:"type:input;placeholder:请输入品牌" data:"type:string;example:Apple" search:"like" validate:"required,min=2,max=50"`
	Model       string         `json:"model" gorm:"size:100" runner:"code:model;name:型号" widget:"type:input;placeholder:请输入型号" data:"type:string;example:A3092" search:"like" validate:"required,min=2,max=100"`
	Specs       string         `json:"specs" gorm:"size:200" runner:"code:specs;name:规格参数" widget:"type:input;mode:text_area;placeholder:请输入规格参数" data:"type:string;example:256GB,深空黑色,钛金属" search:"like"`
	Stock       int            `json:"stock" runner:"code:stock;name:库存数量" widget:"type:number;min:0;unit:件" data:"type:number;example:150" search:"eq,gte,lte" validate:"required,min=0"`
	MinStock    int            `json:"min_stock" runner:"code:min_stock;name:最低库存" widget:"type:number;min:0;unit:件" data:"type:number;example:10" search:"eq,gte,lte" validate:"required,min=0"`
	CostPrice   float64        `json:"cost_price" runner:"code:cost_price;name:成本价" widget:"type:number;min:0;precision:2;prefix:￥" data:"type:float;example:7500.00" search:"eq,gte,lte" validate:"required,min=0"`
	SellPrice   float64        `json:"sell_price" runner:"code:sell_price;name:销售价" widget:"type:number;min:0;precision:2;prefix:￥" data:"type:float;example:8999.00" search:"eq,gte,lte" validate:"required,min=0"`
	Supplier    string         `json:"supplier" gorm:"size:100" runner:"code:supplier;name:供应商" widget:"type:input;placeholder:请输入供应商" data:"type:string;example:苹果官方授权经销商" search:"like" validate:"required,min=2,max=100"`
	Location    string         `json:"location" gorm:"size:100" runner:"code:location;name:存放位置" widget:"type:input;placeholder:请输入存放位置" data:"type:string;example:A区-01-01" search:"like" validate:"required,min=2,max=100"`
	Status      string         `json:"status" runner:"code:status;name:状态" widget:"type:select;options:正常,缺货,停售,下架;placeholder:请选择状态" data:"type:string;example:正常" search:"eq,in" validate:"required,oneof=正常 缺货 停售 下架"`
	Remark      string         `json:"remark" gorm:"size:200" runner:"code:remark;name:备注" widget:"type:input;mode:text_area;placeholder:请输入备注" data:"type:string;example:热销产品，注意及时补货" search:"like"`
}

func (p *Product) TableName() string {
	return "product"
}

// ProductListReq 请求结构体
type ProductListReq struct {
	query.PageInfoReq `runner:"-"`
}

// ProductList 产品列表处理逻辑
func ProductList(ctx *runner.Context, req *ProductListReq, resp response.Response) error {
	db := ctx.MustGetOrInitDB()
	var products []Product
	return resp.Table(&products).AutoPaginated(db, &Product{}, &req.PageInfoReq).Build()
}

// ProductListOption 注册API
var ProductListOption = &runner.FunctionOptions{
	Tags:          []string{"产品管理", "库存管理", "表格演示"},
	EnglishName:   "product_list",
	ChineseName:   "产品库存列表",
	ApiDesc:       "展示产品库存列表，支持分页、搜索和CRUD操作。",
	Request:       &ProductListReq{},
	Response:      query.PaginatedTable[[]Product]{},
	RenderType:    response.RenderTypeTable,
	CreateTables:  []interface{}{&Product{}},
	AutoCrudTable: &Product{},
}

func init() {
	runner.Get("/widgets/product_list", ProductList, ProductListOption)
} 