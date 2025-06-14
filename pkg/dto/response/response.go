package response

type RunFunctionResp struct {
	MetaData   map[string]interface{} `json:"meta_data"`
	Headers    map[string]string      `json:"headers"`
	Code       int                    `json:"code"`
	Msg        string                 `json:"msg"`
	TraceID    string                 `json:"trace_id"`
	RenderType string                 `json:"render_type"`
	Data       interface{}            `json:"data"`
	DataList   []interface{}          `json:"data_list"`
	Multiple   bool                   `json:"multiple"`
}

func (r *RunFunctionResp) GetData() interface{} {
	if r.Multiple {
		return r.DataList
	}
	return r.Data
}

type RunFunctionRespWithData[T any] struct {
	MetaData   map[string]interface{} `json:"meta_data"`
	Headers    map[string]string      `json:"headers"`
	Code       int                    `json:"code"`
	Msg        string                 `json:"msg"`
	TraceID    string                 `json:"trace_id"`
	RenderType string                 `json:"render_type"`
	Data       T                      `json:"data"`
	Multiple   bool                   `json:"multiple"`
}

type Builder interface {
	Build() error
}

type Response interface {
	Form(data interface{}) Form
	Table(resultList interface{}, title ...string) Table
}

func (r *RunFunctionResp) Form(data interface{}) Form {
	return newForm(data, r)
}
