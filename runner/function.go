package runner

import (
	"github.com/yunhanshu-net/function-go/pkg/dto/response"
)

func Get[ReqPtr any](router string, handler func(ctx *Context, req ReqPtr, resp response.Response) error, options ...*FunctionOptions) {
	initRunner()
	r.get(router, handler, options...)
}

func Post[ReqPtr any](router string, handler func(ctx *Context, req ReqPtr, resp response.Response) error, options ...*FunctionOptions) {
	initRunner()
	r.post(router, handler, options...)
}

func (r *Runner) get(router string, handel interface{}, options ...*FunctionOptions) {
	key := fmtKey(router, "GET")
	_, ok := r.routerMap[key]
	if !ok {
		worker := &routerInfo{
			key:          key,
			Handel:       handel,
			Method:       "GET",
			Router:       router,
			FunctionInfo: &FunctionOptions{},
		}
		if len(options) > 0 && options[0] != nil {
			worker.FunctionInfo = options[0]
		}

		r.routerMap[key] = worker
	} else {
		r.routerMap[key].Handel = handel
	}
}

func (r *Runner) post(router string, handel interface{}, options ...*FunctionOptions) {
	key := fmtKey(router, "POST")
	_, ok := r.routerMap[key]
	if !ok {
		worker := &routerInfo{
			key:          key,
			Handel:       handel,
			Method:       "POST",
			Router:       router,
			FunctionInfo: &FunctionOptions{},
		}
		if len(options) > 0 && options[0] != nil {
			worker.FunctionInfo = options[0]
		}

		r.routerMap[key] = worker
	} else {
		r.routerMap[key].Handel = handel
	}
}
