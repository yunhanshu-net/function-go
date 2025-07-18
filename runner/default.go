package runner

import (
	"fmt"
	"github.com/yunhanshu-net/function-go/pkg/dto/api"
	"github.com/yunhanshu-net/function-go/pkg/dto/request"
	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/function-go/pkg/dto/syscallback"
	constants "github.com/yunhanshu-net/pkg/constants/usercall"
	"gorm.io/gorm/schema"
	"strings"
)

func (r *Runner) registerBuiltInRouters() {
	r.get("/_env", _env)
	r.get("/_help", r.help)
	r.get("/_ping", ping)
	r.get("/_getApiInfos", r._getApiInfos)
	r.get("/_getApiInfo", r._getApiInfo)
	r.post("/_callback", r._callback)
	r.post("/_updateConfig", r._updateConfig)
	r.get("/_getConfig", r._getConfig)
	//r.post("/_syscall", r._syscall)
}
func _env(ctx *Context, req *request.NoData, resp response.Response) error {
	return resp.Form(map[string]string{"version": "1.0", "lang": "go"}).Build()
}

func (r *Runner) help(ctx *Context, req *request.NoData, resp response.Response) error {
	s := ""

	for _, info := range r.routerMap {
		s += info.Method + ":" + info.Router
	}
	return resp.Form(map[string]string{"router": s}).Build()
}

func ping(ctx *Context, req *request.NoData, resp response.Response) error {
	return resp.Form(map[string]string{"ping": "pong"}).Build()
}

// buildApiInfo 从路由信息构建API信息
func (r *Runner) buildApiInfo(worker *routerInfo) (*api.Info, error) {
	config := worker.FunctionInfo
	if config == nil {
		return nil, fmt.Errorf("路由配置为空")
	}

	// 构建API信息
	apiInfo := &api.Info{
		Method:       worker.Method,
		Router:       worker.Router,
		User:         r.detail.User,
		Runner:       r.detail.Name,
		ApiDesc:      config.ApiDesc,
		Async:        config.Async,
		Timeout:      config.Timeout,
		FunctionType: string(config.FunctionType),
		AutoRun:      config.AutoRun,
		Tags:         config.Tags,
		Classify:     config.Classify,
		ChineseName:  config.ChineseName,
		EnglishName:  config.EnglishName,
	}

	if config.Request != nil {
		// 获取请求参数信息 - 使用支持FunctionInfo的版本
		params, err := api.NewRequestParamsWithFunctionInfo(config.Request, config.RenderType, config)
		if err != nil {
			return nil, err
		}
		apiInfo.ParamsIn = params
	}

	if config.Response != nil {
		// 获取响应参数信息
		responseParams, err := api.NewResponseParams(config.Response, config.RenderType)
		if err != nil {
			return nil, err
		}
		apiInfo.ParamsOut = responseParams
	}

	// 获取数据表信息
	for _, table := range config.UseTables {
		if tb, ok := table.(schema.Tabler); ok {
			apiInfo.UseTables = append(apiInfo.UseTables, tb.TableName())
		}
	}
	// 获取数据表信息
	for _, table := range config.CreateTables { //记录函数创建的表
		if tb, ok := table.(schema.Tabler); ok {
			apiInfo.CreateTables = append(apiInfo.UseTables, tb.TableName())
		}
	}
	for table, crud := range config.OperateTables { //记录函数对表的crud操作
		apiInfo.OperateTables = make(map[string][]string)
		if tb, ok := table.(schema.Tabler); ok {
			for _, c := range crud {
				apiInfo.OperateTables[tb.TableName()] = append(apiInfo.OperateTables[tb.TableName()], string(c))
			}
		}
	}

	// 获取回调函数信息
	apiInfo.Callbacks = getCallbacks(config)

	return apiInfo, nil
}
func (r *routerInfo) IsDefaultRouter() bool {
	return strings.HasPrefix(strings.TrimPrefix(r.Router, "/"), "_")
}
func (r *Runner) getApiInfos() ([]*api.Info, error) {
	functions := r.routerMap
	var apis []*api.Info
	for _, worker := range functions {
		if worker.IsDefaultRouter() {
			continue
		}
		apiInfo, err := r.buildApiInfo(worker)
		if err != nil {
			fmt.Println("apiInfo err:", err)
			continue // 跳过有错误的API
		}

		apis = append(apis, apiInfo)
	}
	return apis, nil
}

func (r *Runner) getApiInfo(req *syscallback.ApiInfoRequest) (*api.Info, error) {
	// 参数验证
	if req.Router == "" {
		return nil, fmt.Errorf("router参数不能为空")
	}

	// 如果没有指定Method，默认为GET
	if req.Method == "" {
		req.Method = "GET"
	}

	// 获取指定的路由信息
	worker, exist := r.getRouter(req.Router, req.Method)
	if !exist {
		return nil, fmt.Errorf("未找到路由: %s [%s]", req.Router, req.Method)
	}

	apiInfo, err := r.buildApiInfo(worker)
	if err != nil {
		return nil, err
	}
	return apiInfo, nil

}

func (r *Runner) _getApiInfos(ctx *Context, req *syscallback.NoData, resp response.Response) error {
	apis, err := r.getApiInfos()
	if err != nil {
		return err
	}
	return resp.Form(apis).Build()
}

func (r *Runner) _getApiInfo(ctx *Context, req *syscallback.ApiInfoRequest, resp response.Response) error {
	apiInfo, err := r.getApiInfo(req)
	if err != nil {
		return err
	}
	// 返回API信息
	return resp.Form(apiInfo).Build()
}

func getCallbacks(config *FunctionOptions) []string {
	var callbacks []string
	if config == nil {
		return nil
	}
	if config.OnPageLoad != nil {
		callbacks = append(callbacks, constants.CallbackTypeOnPageLoad)
	}

	// API 生命周期回调
	if config.OnApiCreated != nil {
		callbacks = append(callbacks, constants.UserCallTypeOnApiCreated)
	}
	if config.OnApiUpdated != nil {
		callbacks = append(callbacks, constants.CallbackTypeOnApiUpdated)
	}
	if config.BeforeApiDelete != nil {
		callbacks = append(callbacks, constants.CallbackTypeBeforeApiDelete)
	}
	if config.AfterApiDeleted != nil {
		callbacks = append(callbacks, constants.CallbackTypeAfterApiDeleted)
	}

	// 运行器(Runner)生命周期回调
	if config.BeforeRunnerClose != nil {
		callbacks = append(callbacks, constants.CallbackTypeBeforeRunnerClose)
	}
	if config.AfterRunnerClose != nil {
		callbacks = append(callbacks, constants.CallbackTypeAfterRunnerClose)
	}

	// 版本控制回调
	if config.OnVersionChange != nil {
		callbacks = append(callbacks, constants.CallbackTypeOnVersionChange)
	}

	if config.AutoCrudTable != nil {
		callbacks = append(callbacks, constants.CallbackTypeOnTableAddRows)
		callbacks = append(callbacks, constants.CallbackTypeOnTableDeleteRows)
		callbacks = append(callbacks, constants.CallbackTypeOnTableUpdateRows)
	} else {
		// 表格操作回调
		if config.OnTableDeleteRows != nil {
			callbacks = append(callbacks, constants.CallbackTypeOnTableDeleteRows)
		}
		if config.OnTableUpdateRows != nil {
			callbacks = append(callbacks, constants.CallbackTypeOnTableUpdateRows)
		}
		if config.OnTableAddRows != nil {
			callbacks = append(callbacks, constants.CallbackTypeOnTableAddRows)
		}
	}

	if config.OnTableSearch != nil {
		callbacks = append(callbacks, constants.CallbackTypeOnTableSearch)
	}

	// DryRun 回调
	if config.OnDryRun != nil {
		callbacks = append(callbacks, constants.CallbackTypeOnDryRun)
	}

	return callbacks
}

// _updateConfig 配置更新路由
func (r *Runner) _updateConfig(ctx *Context, req *syscallback.ConfigUpdateRequest, resp response.Response) error {
	// 验证参数
	if req.Router == "" {
		return resp.Form(&syscallback.ConfigUpdateResponse{
			Success: false,
			Error:   "router参数不能为空",
		}).Build()
	}
	
	if req.Method == "" {
		return resp.Form(&syscallback.ConfigUpdateResponse{
			Success: false,
			Error:   "method参数不能为空",
		}).Build()
	}

	// 使用请求对象生成配置键
	configKey := req.GenerateConfigKey()

	// 获取配置管理器
	configManager := GetConfigManager()

	// 更新配置
	err := configManager.UpdateConfig(ctx, configKey, req.ConfigData)
	if err != nil {
		return resp.Form(&syscallback.ConfigUpdateResponse{
			Success: false,
			Error:   err.Error(),
		}).Build()
	}

	return resp.Form(&syscallback.ConfigUpdateResponse{
		Success: true,
		Message: "配置更新成功",
	}).Build()
}

// _getConfig 配置获取路由
func (r *Runner) _getConfig(ctx *Context, req *syscallback.ConfigGetRequest, resp response.Response) error {
	// 验证参数
	if req.Router == "" {
		return resp.Form(&syscallback.ConfigGetResponse{
			Success: false,
			Error:   "router参数不能为空",
		}).Build()
	}
	
	if req.Method == "" {
		return resp.Form(&syscallback.ConfigGetResponse{
			Success: false,
			Error:   "method参数不能为空",
		}).Build()
	}

	// 使用请求对象生成配置键
	configKey := req.GenerateConfigKey()

	// 获取配置管理器
	configManager := GetConfigManager()
	configData := configManager.GetByKey(ctx, configKey)
	if configData == nil {
		return resp.Form(&syscallback.ConfigGetResponse{
			Success: false,
			Error:   "配置未找到",
		}).Build()
	}

	return resp.Form(&syscallback.ConfigGetResponse{
		Success: true,
		Config:  configData,
	}).Build()
}
