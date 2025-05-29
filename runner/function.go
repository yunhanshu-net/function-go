package runner

import (
	"github.com/yunhanshu-net/function-go/pkg/dto/response"
)

func Get[ReqPtr any](router string, handler func(ctx *Context, req ReqPtr, resp response.Response) error, config ...*FunctionInfo) {
	initRunner()
	r.get(router, handler, config...)
}

func Post[ReqPtr any](router string, handler func(ctx *Context, req ReqPtr, resp response.Response) error, config ...*FunctionInfo) {
	initRunner()
	r.post(router, handler, config...)
}

func (r *Runner) get(router string, handel interface{}, config ...*FunctionInfo) {
	key := fmtKey(router, "GET")
	_, ok := r.routerMap[key]
	if !ok {
		worker := &routerInfo{
			key:     key,
			Handel:  handel,
			Method:  "GET",
			Router:  router,
			ApiInfo: &FunctionInfo{},
		}
		if len(config) > 0 && config[0] != nil {
			worker.ApiInfo = config[0]
		}

		r.routerMap[key] = worker
	} else {
		r.routerMap[key].Handel = handel
	}
}

func (r *Runner) post(router string, handel interface{}, config ...*FunctionInfo) {
	key := fmtKey(router, "POST")
	_, ok := r.routerMap[key]
	if !ok {
		worker := &routerInfo{
			key:     key,
			Handel:  handel,
			Method:  "POST",
			Router:  router,
			ApiInfo: &FunctionInfo{},
		}
		if len(config) > 0 && config[0] != nil {
			worker.ApiInfo = config[0]
		}

		r.routerMap[key] = worker
	} else {
		r.routerMap[key].Handel = handel
	}
}
