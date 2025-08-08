package widgets

import (
	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/function-go/runner"
	"github.com/yunhanshu-net/pkg/query"
	"gorm.io/gorm"
)

// Employee 员工数据模型
type Employee struct {
	ID           int            `json:"id" gorm:"primaryKey;autoIncrement" runner:"code:id;name:员工ID" permission:"read"`
	CreatedAt    int64          `json:"created_at" gorm:"autoCreateTime:milli" runner:"code:created_at;name:创建时间" widget:"type:datetime;kind:datetime" data:"type:number;example:1705292200000" search:"gte,lte" permission:"read"`
	UpdatedAt    int64          `json:"updated_at" gorm:"autoUpdateTime:milli" runner:"code:updated_at;name:更新时间" widget:"type:datetime;kind:datetime" data:"type:number;example:1705292200000" search:"gte,lte" permission:"read"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index" runner:"-"`
	EmployeeNo   string         `json:"employee_no" gorm:"uniqueIndex;size:20" runner:"code:employee_no;name:工号" widget:"type:input;placeholder:请输入工号" data:"type:string;example:EMP001" search:"like" validate:"required,min=3,max=20"`
	Name         string         `json:"name" gorm:"size:50" runner:"code:name;name:姓名" widget:"type:input;placeholder:请输入姓名" data:"type:string;example:张三" search:"like" validate:"required,min=2,max=20"`
	Gender       string         `json:"gender" gorm:"size:10" runner:"code:gender;name:性别" widget:"type:select;options:男,女;placeholder:请选择性别" data:"type:string;example:男" search:"eq" validate:"required,oneof=男 女"`
	Age          int            `json:"age" runner:"code:age;name:年龄" widget:"type:number;min:18;max:65;unit:岁" data:"type:number;example:28" search:"eq,gte,lte" validate:"required,min=18,max=65"`
	Phone        string         `json:"phone" gorm:"size:20" runner:"code:phone;name:手机号" widget:"type:input;placeholder:请输入手机号" data:"type:string;example:13800138000" search:"like" validate:"required,len=11"`
	Email        string         `json:"email" gorm:"size:100" runner:"code:email;name:邮箱" widget:"type:input;placeholder:请输入邮箱" data:"type:string;example:zhangsan@company.com" search:"like" validate:"required,email"`
	Department   string         `json:"department" gorm:"size:50" runner:"code:department;name:部门" widget:"type:select;options:技术部,产品部,运营部,人事部,财务部;placeholder:请选择部门" data:"type:string;example:技术部" search:"eq,in" validate:"required,oneof=技术部 产品部 运营部 人事部 财务部"`
	Position     string         `json:"position" gorm:"size:50" runner:"code:position;name:职位" widget:"type:select;options:员工,主管,经理,总监,VP;placeholder:请选择职位" data:"type:string;example:员工" search:"eq,in" validate:"required,oneof=员工 主管 经理 总监 VP"`
	Salary       float64        `json:"salary" runner:"code:salary;name:月薪" widget:"type:number;min:3000;precision:2;prefix:￥" data:"type:float;example:15000.00" search:"eq,gte,lte" validate:"required,min=3000"`
	HireDate     int64          `json:"hire_date" runner:"code:hire_date;name:入职时间" widget:"type:datetime;kind:date" data:"type:number;example:1704067200000" search:"gte,lte" validate:"required"`
	Status       string         `json:"status" runner:"code:status;name:在职状态" widget:"type:select;options:在职,离职,试用期;placeholder:请选择状态" data:"type:string;example:在职" search:"eq,in" validate:"required,oneof=在职 离职 试用期"`
	Address      string         `json:"address" gorm:"size:200" runner:"code:address;name:住址" widget:"type:input;mode:text_area;placeholder:请输入住址" data:"type:string;example:北京市朝阳区xxx街道xxx号" search:"like"`
	Remark       string         `json:"remark" gorm:"size:200" runner:"code:remark;name:备注" widget:"type:input;mode:text_area;placeholder:请输入备注" data:"type:string;example:技术能力强，沟通良好" search:"like"`
}

func (e *Employee) TableName() string {
	return "employee"
}

// EmployeeListReq 请求结构体
type EmployeeListReq struct {
	query.PageInfoReq `runner:"-"`
}

// EmployeeList 员工列表处理逻辑
func EmployeeList(ctx *runner.Context, req *EmployeeListReq, resp response.Response) error {
	db := ctx.MustGetOrInitDB()
	var employees []Employee
	return resp.Table(&employees).AutoPaginated(db, &Employee{}, &req.PageInfoReq).Build()
}

// EmployeeListOption 注册API
var EmployeeListOption = &runner.FunctionOptions{
	Tags:          []string{"员工管理", "表格演示"},
	EnglishName:   "employee_list",
	ChineseName:   "员工列表",
	ApiDesc:       "展示员工列表，支持分页、搜索和CRUD操作。",
	Request:       &EmployeeListReq{},
	Response:      query.PaginatedTable[[]Employee]{},
	RenderType:    response.RenderTypeTable,
	CreateTables:  []interface{}{&Employee{}},
	AutoCrudTable: &Employee{},
}

func init() {
	runner.Get("/widgets/employee_list", EmployeeList, EmployeeListOption)
} 