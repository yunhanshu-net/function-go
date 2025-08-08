// 兴趣爱好问卷
// 选择你的兴趣爱好（多选），选项固定

package form

import (
	"fmt"
	"strings"
	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/function-go/runner"
)

// 请求结构体
// 兴趣爱好：多选，选项固定，至少选1项
// example: 阅读,音乐
type HobbySurveyReq struct {
    	Hobbies []string `json:"hobbies" runner:"code:hobbies;name:兴趣爱好" widget:"type:checkbox;options:阅读,音乐,运动,旅行,美食" data:"type:array;default_value:阅读,音乐;example:阅读,音乐" validate:"required"`
}

// 响应结构体
// example: 你选择的兴趣爱好有：阅读、音乐
type HobbySurveyResp struct {
    Summary string `json:"summary" runner:"code:summary;name:结果总结" widget:"type:input;mode:text_area;placeholder:自动生成" data:"type:string;example:你选择的兴趣爱好有：阅读、音乐"`
}

// 业务处理函数
func HobbySurveyHandler(ctx *runner.Context, req *HobbySurveyReq, resp response.Response) error {
    if len(req.Hobbies) == 0 {
        return fmt.Errorf("请至少选择一个兴趣爱好")
    }
    summary := "你选择的兴趣爱好有：" + strings.Join(req.Hobbies, "、")
    return resp.Form(&HobbySurveyResp{Summary: summary}).Build()
}

// API注册
var HobbySurveyOption = &runner.FunctionOptions{
    ChineseName: "兴趣爱好问卷",
    EnglishName: "hobby_survey",
    ApiDesc:     "多选你的兴趣爱好，至少选择一项。",
    Request:     &HobbySurveyReq{},
    Response:    &HobbySurveyResp{},
    RenderType:  response.RenderTypeForm,
    Tags:        []string{"兴趣爱好", "问卷", "多选"},
}

func init() {
    runner.Post("/form/hobby_survey", HobbySurveyHandler, HobbySurveyOption)
} 