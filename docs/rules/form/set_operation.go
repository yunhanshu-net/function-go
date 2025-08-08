// 集合运算器
// 输入两个集合（字符串，用分隔符分割）和运算方式，输出运算结果
// 分隔符支持逗号、分号、空格等，适配不同用户的数据格式

package form

import (
	"fmt"
	"sort"
	"strings"
	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/function-go/runner"
)

// 请求结构体
// 集合A、集合B：字符串输入，用分隔符分割
// 分隔符：适配不同用户数据格式，如逗号、分号、空格等
// 运算方式：交集、并集、差集、补集
// example: 集合A"苹果,橙子"，集合B"苹果,葡萄"，分隔符","，运算方式"交集"
type SetOperationReq struct {
    SetA      string `json:"set_a" runner:"code:set_a;name:集合A" widget:"type:input;placeholder:请输入集合A元素，用分隔符分割" data:"type:string;default_value:苹果,橙子;example:苹果,橙子" validate:"required"`
    SetB      string `json:"set_b" runner:"code:set_b;name:集合B" widget:"type:input;placeholder:请输入集合B元素，用分隔符分割" data:"type:string;default_value:苹果,葡萄;example:苹果,葡萄" validate:"required"`
    Separator string `json:"separator" runner:"code:separator;name:分隔符" widget:"type:input;placeholder:请输入分隔符，如逗号、分号等" data:"type:string;default_value:,;example:," validate:"required"`
    Operation string `json:"operation" runner:"code:operation;name:运算方式" widget:"type:select;options:交集,并集,差集,补集;placeholder:请选择运算方式" data:"type:string;default_value:交集;example:交集" validate:"required,oneof=交集 并集 差集 补集"`
}

// 响应结构体
// example: 交集运算结果：苹果
type SetOperationResp struct {
    Result string `json:"result" runner:"code:result;name:运算结果" widget:"type:input;mode:text_area;placeholder:自动计算" data:"type:string;example:交集运算结果：苹果"`
}

// 业务处理函数
func SetOperationHandler(ctx *runner.Context, req *SetOperationReq, resp response.Response) error {
    // 分割字符串为数组
    setA := strings.Split(strings.TrimSpace(req.SetA), req.Separator)
    setB := strings.Split(strings.TrimSpace(req.SetB), req.Separator)
    
    // 去除空元素和空格
    setA = removeEmptyAndTrim(setA)
    setB = removeEmptyAndTrim(setB)
    
    if len(setA) == 0 || len(setB) == 0 {
        return fmt.Errorf("请至少为每个集合输入一个元素")
    }
    
    var result []string
    switch req.Operation {
    case "交集":
        result = intersection(setA, setB)
    case "并集":
        result = union(setA, setB)
    case "差集":
        result = difference(setA, setB)
    case "补集":
        result = complement(setA, setB)
    default:
        return fmt.Errorf("不支持的运算方式")
    }
    
    // 排序并格式化输出
    sort.Strings(result)
    if len(result) == 0 {
        return resp.Form(&SetOperationResp{Result: fmt.Sprintf("%s运算结果：空集", req.Operation)}).Build()
    }
    return resp.Form(&SetOperationResp{Result: fmt.Sprintf("%s运算结果：%s", req.Operation, strings.Join(result, "、"))}).Build()
}

// 去除空元素和空格
func removeEmptyAndTrim(items []string) []string {
    var result []string
    for _, item := range items {
        trimmed := strings.TrimSpace(item)
        if trimmed != "" {
            result = append(result, trimmed)
        }
    }
    return result
}

// 交集运算
func intersection(setA, setB []string) []string {
    setAMap := make(map[string]bool)
    for _, item := range setA {
        setAMap[item] = true
    }
    var result []string
    for _, item := range setB {
        if setAMap[item] {
            result = append(result, item)
        }
    }
    return result
}

// 并集运算
func union(setA, setB []string) []string {
    unionMap := make(map[string]bool)
    for _, item := range setA {
        unionMap[item] = true
    }
    for _, item := range setB {
        unionMap[item] = true
    }
    var result []string
    for item := range unionMap {
        result = append(result, item)
    }
    return result
}

// 差集运算 A-B
func difference(setA, setB []string) []string {
    setBMap := make(map[string]bool)
    for _, item := range setB {
        setBMap[item] = true
    }
    var result []string
    for _, item := range setA {
        if !setBMap[item] {
            result = append(result, item)
        }
    }
    return result
}

// 补集运算（A的补集，相对于A∪B）
func complement(setA, setB []string) []string {
    unionSet := union(setA, setB)
    return difference(unionSet, setA)
}

// API注册
var SetOperationOption = &runner.FunctionOptions{
    ChineseName: "集合运算器",
    EnglishName: "set_operation",
    ApiDesc:     "输入两个集合（字符串，用分隔符分割）和运算方式，输出集合运算结果。",
    Request:     &SetOperationReq{},
    Response:    &SetOperationResp{},
    RenderType:  response.RenderTypeForm,
    Tags:        []string{"集合", "运算", "数学"},
}

func init() {
    runner.Post("/form/set_operation", SetOperationHandler, SetOperationOption)
} 