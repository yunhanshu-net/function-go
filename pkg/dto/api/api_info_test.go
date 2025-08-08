package api

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/pkg/query"
	"github.com/yunhanshu-net/pkg/typex/files"
)

type AddReq struct {
	A        int    `json:"a" form:"a" runner:"code:a;name:值a;type:number;example:100;placeholder:请输入值a" validate:"required,min=-1000,max=10000"`
	B        int    `json:"b" form:"b" runner:"code:b;name:值b;type:number;example:200;placeholder:请输入值b" validate:"required,min=-1000,max=10000"`
	Receiver string `json:"receiver" form:"receiver" runner:"code:receiver;name:接收人;widget:select;default_value:beiluo;options:admin,beiluo,user;type:string;placeholder:请输入接收人"`
	Desc     string `json:"desc" form:"desc" runner:"code:desc;name:描述;type:string;placeholder:请描述此次计算;callback:OnInputFuzzy"`
}

type AddResp struct {
	Result   int    `json:"result" runner:"code:result;name:计算结果;example:30000"`
	Receiver string `json:"receiver" form:"receiver" runner:"code:receiver;name:接收人"`
	Desc     string `json:"desc" form:"desc" runner:"code:desc;name:描述"`
}

// FileUploadReq 文件上传请求
type FileUploadReq struct {
	Files *files.Files `json:"files" runner:"code:files;name:文件列表;widget:file_upload;type:files;multiple:true;max_size:10240;placeholder:请选择文件"`
}

// FileListResp 文件列表响应
type FileListResp struct {
	Files *files.Files `json:"files" runner:"code:files;name:文件列表;widget:file_display;type:files;display_mode:card;preview:true;download:true"`
}

func TestReq(t *testing.T) {
	params, err := NewRequestParams(&query.PageInfoReq{}, response.RenderTypeTable)
	if err != nil {
		panic(err)
	}
	fmt.Println(params)
	marshal, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(marshal))
}

func TestResp(t *testing.T) {
	params, err := NewResponseParams(&query.PaginatedTable[[]AddResp]{}, response.RenderTypeTable)
	if err != nil {
		panic(err)
	}
	fmt.Println(params)
	marshal, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(marshal))
}

func TestTableReq(t *testing.T) {
	params, err := NewTableRequestParams(&query.PaginatedTable[[]AddResp]{})
	if err != nil {
		panic(err)
	}
	fmt.Println(params)
	marshal, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(marshal))
}

func TestName(t *testing.T) {
	params, err := NewRequestParams(&AddReq{}, "")
	if err != nil {
		panic(err)
	}
	fmt.Println(params)
	marshal, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(marshal))
}

func TestResponse(t *testing.T) {
	params, err := NewResponseParams(&AddResp{}, "")
	if err != nil {
		panic(err)
	}
	fmt.Println(params)
	marshal, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(marshal))
}

// TestFileUpload 测试文件上传请求参数解析
func TestFileUpload(t *testing.T) {
	params, err := NewRequestParams(&FileUploadReq{}, response.RenderTypeForm)
	if err != nil {
		panic(err)
	}
	fmt.Println("文件上传请求参数：")
	marshal, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(marshal))
}

// TestFileList 测试文件列表响应参数解析
func TestFileList(t *testing.T) {
	params, err := NewResponseParams(&FileListResp{}, response.RenderTypeForm)
	if err != nil {
		panic(err)
	}
	fmt.Println("文件列表响应参数：")
	marshal, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(marshal))
}
