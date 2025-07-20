package runner

import (
	"fmt"
	"github.com/yunhanshu-net/function-go/pkg/dto/response"
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
	
	// 添加调试日志
	fmt.Printf("=== POST 路由注册 ===\n")
	fmt.Printf("路由: %s\n", router)
	fmt.Printf("路由键: %s\n", key)
	fmt.Printf("路由已存在: %v\n", ok)
	if len(options) > 0 && options[0] != nil {
		fmt.Printf("AutoUpdateConfig 存在: %v\n", options[0].AutoUpdateConfig != nil)
		if options[0].AutoUpdateConfig != nil {
			fmt.Printf("AutoUpdateConfig.ConfigStruct 存在: %v\n", options[0].AutoUpdateConfig.ConfigStruct != nil)
		}
	} else {
		fmt.Printf("options 为空或 nil\n")
	}
	
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
				fmt.Printf("调用 registerAutoUpdateConfig\n")
				r.registerAutoUpdateConfig(router, "POST", options[0].AutoUpdateConfig)
			} else {
				fmt.Printf("AutoUpdateConfig 为 nil，跳过注册\n")
			}
		}

		r.routerMap[key] = worker
		fmt.Printf("新路由已注册\n")
	} else {
		r.routerMap[key].Handel = handel
		fmt.Printf("路由已存在，仅更新处理器\n")
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

	// 生成配置键，只对 method 做小写
	configKey := fmt.Sprintf("function.%s.%s", safeRouter, strings.ToLower(method))

	// 获取配置管理器
	configManager := GetConfigManager()

	// 注册配置变更回调
	if autoConfig.BeforeConfigChange != nil {
		configManager.RegisterCallback(configKey, autoConfig.BeforeConfigChange)
	}

	// 注册配置结构体
	if autoConfig.ConfigStruct != nil {
		configManager.RegisterConfigStruct(configKey, autoConfig.ConfigStruct)
		// 添加调试日志
		fmt.Printf("=== 配置结构体注册 ===\n")
		fmt.Printf("配置键: %s\n", configKey)
		fmt.Printf("配置结构体类型: %T\n", autoConfig.ConfigStruct)
		fmt.Printf("配置结构体值: %+v\n", autoConfig.ConfigStruct)
	} else {
		fmt.Printf("=== 配置结构体注册失败 ===\n")
		fmt.Printf("配置键: %s\n", configKey)
		fmt.Printf("AutoUpdateConfig.ConfigStruct 为 nil\n")
	}
}
