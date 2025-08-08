package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/yunhanshu-net/function-go/pkg/dto/api"
	"github.com/yunhanshu-net/function-go/pkg/dto/request"
	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/function-go/pkg/dto/usercall"

	constants "github.com/yunhanshu-net/pkg/constants/usercall"
	"github.com/yunhanshu-net/pkg/logger"
	"github.com/yunhanshu-net/pkg/x/jsonx"
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
	//r.post("/_syscall", r._syscall)
}
func _env(ctx *Context, req *request.NoData, resp response.Response) error {
	return resp.Form(map[string]string{"version": "1.0", "lang": "go"}).Build()
}

func (r *Runner) help(ctx *Context, req *request.NoData, resp response.Response) error {
	s := ""

	for _, info := range r.routerMap {
		s += info.Method + ":" + info.Router + "\n"
	}
	return resp.Form(map[string]string{"router": s}).Build()
}

func ping(ctx *Context, req *request.NoData, resp response.Response) error {
	return resp.Form(map[string]string{"ping": "pong"}).Build()
}

// buildApiInfo 从路由信息构建API信息
func (r *Runner) buildApiInfo(worker *routerInfo) (*api.Info, error) {
	// 检查是否是 FunctionOptions 类型
	if functionOptions, ok := worker.Option.(*FunctionOptions); ok {
		return r.buildApiInfoFromFunctionOptions(worker, functionOptions)
	}

	// 否则使用新的 Option 接口
	opt := worker.Option
	config := worker.Option.GetBaseConfig()

	if config == nil {
		return nil, fmt.Errorf("路由配置为空")
	}

	infoInterface, ok := opt.(api.FunctionInfoInterface)
	if !ok {
		return nil, fmt.Errorf("FunctionInfoInterface not impl")
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
	if config.Group != nil {
		apiInfo.Group = &api.Group{
			CnName: config.Group.CnName,
			EnName: config.Group.EnName,
		}
	}

	if config.Request != nil {
		// 获取请求参数信息 - 暂时使用简单版本
		params, err := api.NewRequestParamsWithFunctionInfo(config.Request, opt.GetRenderType(), infoInterface)
		if err != nil {
			return nil, err
		}
		apiInfo.ParamsIn = params
	}

	if config.Response != nil {
		// 获取响应参数信息
		responseParams, err := api.NewResponseParams(config.Response, opt.GetRenderType())
		if err != nil {
			return nil, err
		}
		apiInfo.ParamsOut = responseParams
	}

	logger.Infof(context.Background(), "worker %s AutoUpdateConfig ==nil:%v config: %v el:%+v ",
		worker.Router, config.AutoUpdateConfig == nil, jsonx.String(config), config)

	if opt.GetAutoCrudTable() != nil {
		apiInfo.Callbacks = append(apiInfo.Callbacks, constants.CallbackTypeOnTableAddRows)
		apiInfo.Callbacks = append(apiInfo.Callbacks, constants.CallbackTypeOnTableDeleteRows)
		apiInfo.Callbacks = append(apiInfo.Callbacks, constants.CallbackTypeOnTableUpdateRows)
	}
	// 处理配置相关
	if config.AutoUpdateConfig != nil {
		// 解析配置结构体，生成表单配置
		configParams, err := api.NewRequestParamsWithFunctionInfo(config.AutoUpdateConfig.ConfigStruct, opt.GetRenderType(), infoInterface)
		if err != nil {
			logger.Errorf(context.Background(), "autoUpdateConfig err: %v %+v", err, config.AutoUpdateConfig)
			// 记录错误但不中断API构建
			fmt.Printf("解析配置结构体失败: %v\n", err)
		} else {
			logger.Infof(context.Background(), "autoUpdateConfig config params: %+v", configParams)
			apiInfo.ParamsConfig = configParams

			// 将配置结构体转换为map类型作为初始数据
			configDataMap, err := structToMap(config.AutoUpdateConfig.ConfigStruct)
			if err != nil {
				logger.Errorf(context.Background(), "convert config struct to map failed: %v", err)
			} else {
				apiInfo.ParamsData = configDataMap
			}
		}

		// 将初始配置写入文件
		if err := r.writeInitialConfig(worker, config.AutoUpdateConfig.ConfigStruct); err != nil {
			logger.Errorf(context.Background(), "writeInitialConfig err: %v %+v", err, config.AutoUpdateConfig)
			// 记录错误但不中断API构建
			fmt.Printf("写入初始配置失败: %v\n", err)
		}
	}

	// 获取数据表信息
	for _, table := range config.CreateTables {
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

	callbacks := opt.GetCallbacks()
	for name, _ := range callbacks {
		apiInfo.Callbacks = append(apiInfo.Callbacks, name)
	}
	//// 获取回调函数信息 - 简化处理
	//apiInfo.Callbacks = []string{}

	return apiInfo, nil
}

// structToMap 将结构体转换为map[string]interface{}
func structToMap(obj interface{}) (map[string]interface{}, error) {
	// 先序列化为JSON
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("序列化结构体失败: %w", err)
	}

	// 再反序列化为map
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("反序列化为map失败: %w", err)
	}

	return result, nil
}

// writeInitialConfig 写入初始配置到文件
func (r *Runner) writeInitialConfig(worker *routerInfo, configStruct interface{}) error {
	// 生成配置键
	configKey := generateConfigKey(worker.Router, worker.Method)

	// 获取配置管理器
	configManager := GetConfigManager()

	// 创建正确初始化的上下文
	ctx := NewContext(context.Background(), worker.Method, worker.Router)

	// 检查配置是否已存在
	existingConfig := configManager.GetByKey(ctx, configKey)
	if existingConfig != nil {
		// 配置已存在，不覆盖
		return nil
	}

	// 创建配置数据，直接存储配置对象
	config := &usercall.ConfigData{
		Type: "json",
		Data: configStruct, // 直接存储配置对象，避免双重序列化
	}

	// 注册配置结构体类型
	configManager.RegisterConfigStruct(configKey, configStruct)

	// 写入配置
	return configManager.UpdateConfig(ctx, configKey, config)
}

// generateConfigKey 生成配置键
func generateConfigKey(router, method string) string {
	return usercall.GenerateConfigKey(router, method)
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

func (r *Runner) getApiInfo(req *usercall.ApiInfoRequest) (*api.Info, error) {
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

func (r *Runner) _getApiInfos(ctx *Context, req *usercall.NoData, resp response.Response) error {
	apis, err := r.getApiInfos()
	if err != nil {
		return err
	}
	return resp.Form(apis).Build()
}

func (r *Runner) _getApiInfo(ctx *Context, req *usercall.ApiInfoRequest, resp response.Response) error {
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
