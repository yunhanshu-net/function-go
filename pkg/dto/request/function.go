package request

import (
	"github.com/yunhanshu-net/pkg/dto/runnerproject"
	"strings"
)

type RunFunctionReq struct {
	RunnerID   string                `json:"runner_id"`
	Runner     *runnerproject.Runner `json:"runner"`
	TraceID    string                `json:"trace_id"`
	Router     string                `json:"router"`
	Method     string                `json:"method"`
	Headers    map[string]string     `json:"headers"`
	BodyType   string                `json:"body_type"`
	Body       interface{}           `json:"body"`
	UrlQuery   string                `json:"url_query"`
	NatsHeader map[string][]string   `json:"nats_header"`
}

func (r *RunFunctionReq) IsMethodGet() bool {
	return strings.ToLower(r.Method) == "get"
}

type NoData struct {
}

type ApiInfoRequest struct {
	Router string `json:"router" form:"router"` // API路由路径
	Method string `json:"method" form:"method"` // HTTP方法（GET/POST）
}
