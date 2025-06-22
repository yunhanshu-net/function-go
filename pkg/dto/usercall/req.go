package usercall

import (
	"encoding/json"
	"fmt"
	"github.com/yunhanshu-net/pkg/x/jsonx"
)

type OnPageLoadReq struct {
}

type OnPageLoadResp struct {
	Request interface{} `json:"request"`  //会初始化前端的表单参数
	AutoRun bool        `json:"auto_run"` //是否自动运行
}

type OnApiCreatedReq struct {
	//Method string `json:"method"`
	//Router string `json:"router"`
}

type OnApiUpdatedReq struct {
	Method string `json:"method"`
	Router string `json:"router"`
}

type BeforeApiDeleteReq struct {
	Method string `json:"method"`
	Router string `json:"router"`
}

type AfterApiDeletedReq struct {
	Method string `json:"method"`
	Router string `json:"router"`
}

type BeforeRunnerCloseReq struct {
}

type AfterRunnerCloseReq struct {
}

type Change struct {
	Method string `json:"method"`
	Router string `json:"router"`
	Type   string `json:"type"`
}

func (c *Change) String() string {
	return fmt.Sprintf(`{"method": "%s", "router": "%s","type","%s"}`, c.Method, c.Router, c.Type)
}

type OnVersionChangeReq struct {
	Change []Change `json:"change"`
}

type OnInputFuzzyReq struct {
	Code  string `json:"code"`
	Value string `json:"value"`
}

type OnInputValidateReq struct {
	Code  string `json:"code"`
	Value string `json:"value"`
}

type OnTableDeleteRowsReq struct {
	Ids []int `json:"ids"`
}

type OnTableUpdateRowsReq struct {
	Ids    []int                  `json:"ids"`
	Fields map[string]interface{} `json:"fields"` // 要更新的字段和值的映射
}
type OnTableAddRowsReq struct {
	Rows interface{} `json:"rows"`
}

type OnTableAddRowsResp struct {
}

func (r *OnTableAddRowsReq) DecodeBy(el interface{}) error {
	err := jsonx.Convert(r.Rows, el)
	if err != nil {
		return err
	}
	return nil
}

type OnTableSearchReq struct {
	Cond map[string]string `json:"cond"`
}
type Request struct {
	Method string      `json:"method"`
	Router string      `json:"router"`
	Type   string      `json:"type"`
	Body   interface{} `json:"body"`
}

func (c *Request) DecodeData(el interface{}) error {
	marshal, err := json.Marshal(c.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, &el)
	if err != nil {
		return err
	}
	return nil
}

type Response struct {
	Request  interface{} `json:"request"`
	Response interface{} `json:"response"`
}

type InputFuzzyItem struct {
	Value string `json:"value"`
}
type OnInputFuzzyResp struct {
	Values []*InputFuzzyItem `json:"values"`
}

type OnInputValidateResp struct {
	Msg string `json:"msg"`
}

type OnTableDeleteRowsResp struct {
}

type OnTableUpdateRowsResp struct {
}

type OnTableSearchResp struct {
}
