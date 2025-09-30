package response

import (
	"fmt"
	"reflect"

	"github.com/yunhanshu-net/pkg/query"
	"github.com/yunhanshu-net/pkg/x/tagx"
	"gorm.io/gorm"
)

type Table interface {
	Builder
	AutoPaginated(dbAndWhere *gorm.DB, model interface{}, pageInfo *query.SearchFilterPageReq) Table
}

type column struct {
	Idx  int    `json:"idx"`
	Name string `json:"name"`
	Code string `json:"code"`
}
type paginated struct {
	CurrentPage int `json:"current_page"` // 当前页码
	TotalCount  int `json:"total_count"`  // 总数据量
	TotalPages  int `json:"total_pages"`  // 总页数
	PageSize    int `json:"page_size"`    // 每页数量
}

type tableData struct {
	err  error
	val  interface{}
	resp *RunFunctionResp
	Data table
}
type table struct {
	Title      string                   `json:"title"`
	Column     []column                 `json:"column"`
	Values     map[string][]interface{} `json:"values"`
	Pagination paginated                `json:"pagination"`
}

func newTable(resp *RunFunctionResp, resultList interface{}, title ...string) *tableData {
	titleStr := ""
	if len(title) > 0 {
		titleStr = title[0]
	}
	return &tableData{resp: resp, val: resultList, Data: table{Title: titleStr}}
}
func (r *RunFunctionResp) Table(resultList interface{}, title ...string) Table {
	return newTable(r, resultList, title...)
}

func (t *tableData) Build() error {
	if t.err != nil {
		return t.err
	}
	if t.val == nil {
		return build(t.resp, t.Data, RenderTypeTable)
	}
	sliceVal := reflect.ValueOf(t.val)
	if sliceVal.Kind() == reflect.Pointer {
		sliceVal = sliceVal.Elem()
	}
	if sliceVal.Kind() != reflect.Slice {
		return fmt.Errorf("类型错误")
	}
	if sliceVal.Len() == 0 {
		return build(t.resp, t.Data, RenderTypeTable)
	}
	row := sliceVal.Index(0)
	columns := parserTableInfo(row.Interface())

	values := make(map[string][]interface{}, len(columns))
	for i := 0; i < sliceVal.Len(); i++ {
		row := sliceVal.Index(i)
		// 确保row不是指针值
		if row.Kind() == reflect.Pointer {
			row = row.Elem()
		}
		for _, col := range columns {
			field := row.Field(col.Idx) //根据idx取field
			if _, ok := values[col.Code]; !ok {
				values[col.Code] = make([]interface{}, 0, sliceVal.Len())
			}
			values[col.Code] = append(values[col.Code], field.Interface())
		}
	}

	t.Data.Column = columns
	t.Data.Values = values
	return build(t.resp, t.Data, RenderTypeTable)
}

func (t *tableData) AutoPaginated(db *gorm.DB, model interface{}, pageInfo *query.SearchFilterPageReq) Table {

	if pageInfo == nil {
		pageInfo = new(query.SearchFilterPageReq)
	}

	// 修复：在应用搜索条件之前先克隆数据库连接，避免污染原始连接
	dbClone := db.Session(&gorm.Session{})

	// 使用query库的公开方法应用搜索条件
	dbWithConditions, err := query.ApplySearchConditions(dbClone, pageInfo)
	if err != nil {
		t.err = fmt.Errorf("AutoPaginated.ApplySearchConditions failed: %v", err)
		return t
	}

	// 获取分页大小
	pageSize := pageInfo.GetLimit()
	offset := pageInfo.GetOffset()

	// 查询总数
	var totalCount int64
	if err := dbWithConditions.Model(model).Count(&totalCount).Error; err != nil {
		t.err = fmt.Errorf("AutoPaginated.Count :%+v failed to count records: %v", t.val, err)
		return t
	}

	// 应用排序
	if pageInfo.GetSorts() != "" {
		dbWithConditions = dbWithConditions.Order(pageInfo.GetSorts())
	}

	// 查询当前页数据
	queryDB := dbWithConditions.Offset(offset).Limit(pageSize)

	if err := queryDB.Find(t.val).Error; err != nil {
		t.err = fmt.Errorf("AutoPaginated.Find :%+v failed to find records: %v", t.val, err)
		return t
	}

	// 计算总页数
	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize != 0 {
		totalPages++
	}

	// 构造分页结果
	t.Data.Pagination = paginated{
		CurrentPage: pageInfo.Page,
		TotalCount:  int(totalCount),
		TotalPages:  totalPages,
		PageSize:    pageSize,
	}

	return t
}

func parserTableInfo(row interface{}) []column {
	of := reflect.TypeOf(row)

	// 添加指针类型处理
	if of.Kind() == reflect.Pointer {
		of = of.Elem() // 获取指针指向的类型
	}

	var columns []column
	for i := 0; i < of.NumField(); i++ {
		field := of.Field(i)

		// 检查runner标签是否为"-"，如果是则忽略该字段
		runnerTag := field.Tag.Get("runner")
		if runnerTag == "-" {
			continue
		}

		kv := tagx.ParserKv(runnerTag)
		name := kv["name"]
		code := kv["code"]
		if name == "" {
			name = field.Tag.Get("json")
			if name == "" {
				name = field.Name
			}
		}
		if code == "" {
			code = field.Tag.Get("json")
			if code == "" {
				code = field.Name
			}
		}

		columns = append(columns, column{Name: name, Code: code, Idx: i})
	}
	return columns
}
