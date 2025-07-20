package runner

import (
	"context"
	"fmt"
	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/pkg/logger"
	"strings"
)

func Get[ReqPtr any](router string, handler func(ctx *Context, req ReqPtr, resp response.Response) error, options ...*FunctionOptions) {
	initRunner()
	r.get(router, handler, options...)
}

func Post[ReqPtr any](router string, handler func(ctx *Context, req ReqPtr, resp response.Response) error, options ...*FunctionOptions) {
	initRunner()
	r.post(router, handler, options...)
}

func Put[ReqPtr any](router string, handler func(ctx *Context, req ReqPtr, resp response.Response) error, options ...*FunctionOptions) {
	initRunner()
	r.put(router, handler, options...)
}

func Delete[ReqPtr any](router string, handler func(ctx *Context, req ReqPtr, resp response.Response) error, options ...*FunctionOptions) {
	initRunner()
	r.delete(router, handler, options...)
}

func Patch[ReqPtr any](router string, handler func(ctx *Context, req ReqPtr, resp response.Response) error, options ...*FunctionOptions) {
	initRunner()
	r.patch(router, handler, options...)
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

			// 处理 AutoUpdateConfig
			if options[0].AutoUpdateConfig != nil {
				r.registerAutoUpdateConfig(router, "GET", options[0].AutoUpdateConfig)
			}
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

			// 处理 AutoUpdateConfig
			if options[0].AutoUpdateConfig != nil {
				r.registerAutoUpdateConfig(router, "POST", options[0].AutoUpdateConfig)
			}
		}

		r.routerMap[key] = worker
	} else {
		r.routerMap[key].Handel = handel
	}
}

func (r *Runner) put(router string, handel interface{}, options ...*FunctionOptions) {
	key := fmtKey(router, "PUT")
	_, ok := r.routerMap[key]
	if !ok {
		worker := &routerInfo{
			key:          key,
			Handel:       handel,
			Method:       "PUT",
			Router:       router,
			FunctionInfo: &FunctionOptions{},
		}
		if len(options) > 0 && options[0] != nil {
			worker.FunctionInfo = options[0]

			// 处理 AutoUpdateConfig
			if options[0].AutoUpdateConfig != nil {
				r.registerAutoUpdateConfig(router, "PUT", options[0].AutoUpdateConfig)
			}
		}

		r.routerMap[key] = worker
	} else {
		r.routerMap[key].Handel = handel
	}
}

func (r *Runner) delete(router string, handel interface{}, options ...*FunctionOptions) {
	key := fmtKey(router, "DELETE")
	_, ok := r.routerMap[key]
	if !ok {
		worker := &routerInfo{
			key:          key,
			Handel:       handel,
			Method:       "DELETE",
			Router:       router,
			FunctionInfo: &FunctionOptions{},
		}
		if len(options) > 0 && options[0] != nil {
			worker.FunctionInfo = options[0]

			// 处理 AutoUpdateConfig
			if options[0].AutoUpdateConfig != nil {
				r.registerAutoUpdateConfig(router, "DELETE", options[0].AutoUpdateConfig)
			}
		}

		r.routerMap[key] = worker
	} else {
		r.routerMap[key].Handel = handel
	}
}

func (r *Runner) patch(router string, handel interface{}, options ...*FunctionOptions) {
	key := fmtKey(router, "PATCH")
	_, ok := r.routerMap[key]
	if !ok {
		worker := &routerInfo{
			key:          key,
			Handel:       handel,
			Method:       "PATCH",
			Router:       router,
			FunctionInfo: &FunctionOptions{},
		}
		if len(options) > 0 && options[0] != nil {
			worker.FunctionInfo = options[0]

			// 处理 AutoUpdateConfig
			if options[0].AutoUpdateConfig != nil {
				r.registerAutoUpdateConfig(router, "PATCH", options[0].AutoUpdateConfig)
			}
		}

		r.routerMap[key] = worker
	} else {
		r.routerMap[key].Handel = handel
	}
}

// registerAutoUpdateConfig 注册自动更新配置
func (r *Runner) registerAutoUpdateConfig(router string, method string, autoConfig *AutoUpdateConfig) {
	// 处理路由路径，将 / 替换为 . 以安全地用作配置键
	safeRouter := strings.ReplaceAll(router, "/", ".")
	// 移除前后的点
	safeRouter = strings.Trim(safeRouter, ".")

	// 生成配置键，包含method避免不同HTTP方法的冲突
	configKey := fmt.Sprintf("function.%s.%s", safeRouter, method)

	// 获取配置管理器
	configManager := GetConfigManager()

	// 注册配置变更回调
	if autoConfig.BeforeConfigChange != nil {
		configManager.RegisterCallback(configKey, autoConfig.BeforeConfigChange)
	}

	// 注册配置结构体
	if autoConfig.ConfigStruct != nil {
		configManager.RegisterConfigStruct(configKey, autoConfig.ConfigStruct)
	}
}
