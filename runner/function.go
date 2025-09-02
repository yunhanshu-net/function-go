package runner

import (
	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/function-go/pkg/dto/usercall"
	"strings"
)

func Get[ReqPtr any](router string, handler func(ctx *Context, req ReqPtr, resp response.Response) error, options ...Option) {
	initRunner()
	r.get(router, handler, options...)
}

func Post[ReqPtr any](router string, handler func(ctx *Context, req ReqPtr, resp response.Response) error, options ...Option) {
	initRunner()
	r.post(router, handler, options...)
}

func Put[ReqPtr any](router string, handler func(ctx *Context, req ReqPtr, resp response.Response) error, options ...Option) {
	initRunner()
	r.put(router, handler, options...)
}

func Delete[ReqPtr any](router string, handler func(ctx *Context, req ReqPtr, resp response.Response) error, options ...Option) {
	initRunner()
	r.delete(router, handler, options...)
}

func Patch[ReqPtr any](router string, handler func(ctx *Context, req ReqPtr, resp response.Response) error, options ...Option) {
	initRunner()
	r.patch(router, handler, options...)
}

func (r *Runner) initOptions(method string, router string, handel interface{}, options ...Option) *routerInfo {

	key := fmtKey(router, method)
	worker := &routerInfo{
		key:    key,
		Handel: handel,
		Method: method,
		Router: router,
	}
	if len(options) > 0 && options[0] != nil {
		setEnName(router, options[0])
		worker.Option = options[0]
		// 处理 AutoUpdateConfig
		if options[0].GetBaseConfig().AutoUpdateConfig != nil {
			r.registerAutoUpdateConfig(router, method, options[0].GetBaseConfig().AutoUpdateConfig)
		}
	}
	return worker
}

func setEnName(router string, option Option) {
	split := strings.Split(router, "/")
	enname := split[len(split)-1]
	config := option.GetBaseConfig()
	if config != nil {
		config.EnglishName = enname
	}
}

func (r *Runner) get(router string, handel interface{}, options ...Option) {
	key := fmtKey(router, "GET")
	_, ok := r.routerMap[key]
	if !ok {
		r.routerMap[key] = r.initOptions("GET", router, handel, options...)
	} else {
		r.routerMap[key].Handel = handel
	}
}

func (r *Runner) post(router string, handel interface{}, options ...Option) {
	key := fmtKey(router, "POST")
	_, ok := r.routerMap[key]
	if !ok {
		r.routerMap[key] = r.initOptions("POST", router, handel, options...)
	} else {
		r.routerMap[key].Handel = handel
	}
}

func (r *Runner) put(router string, handel interface{}, options ...Option) {
	key := fmtKey(router, "PUT")
	_, ok := r.routerMap[key]
	if !ok {
		r.routerMap[key] = r.initOptions("PUT", router, handel, options...)
	} else {
		r.routerMap[key].Handel = handel
	}
}

func (r *Runner) delete(router string, handel interface{}, options ...Option) {
	key := fmtKey(router, "DELETE")
	_, ok := r.routerMap[key]
	if !ok {
		r.routerMap[key] = r.initOptions("DELETE", router, handel, options...)
	} else {
		r.routerMap[key].Handel = handel
	}
}

func (r *Runner) patch(router string, handel interface{}, options ...Option) {
	key := fmtKey(router, "PATCH")
	_, ok := r.routerMap[key]
	if !ok {
		r.routerMap[key] = r.initOptions("PATCH", router, handel, options...)

	} else {
		r.routerMap[key].Handel = handel
	}
}

// registerAutoUpdateConfig 注册自动更新配置
func (r *Runner) registerAutoUpdateConfig(router string, method string, autoConfig *AutoUpdateConfig) {
	configKey := usercall.GenerateConfigKey(router, method)
	configManager := GetConfigManager()
	if autoConfig.BeforeConfigChange != nil {
		configManager.RegisterCallback(configKey, autoConfig.BeforeConfigChange)
	}
	if autoConfig.ConfigStruct != nil {
		configManager.RegisterConfigStruct(configKey, autoConfig.ConfigStruct)
	}
}
