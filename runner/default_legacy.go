package runner

import (
	"context"
	"fmt"
	"github.com/yunhanshu-net/function-go/pkg/dto/api"
	"github.com/yunhanshu-net/pkg/logger"
	"gorm.io/gorm/schema"
)

// buildApiInfoFromFunctionOptions 从旧版FunctionOptions构建API信息
func (r *Runner) buildApiInfoFromFunctionOptions(worker *routerInfo, config *FunctionOptions) (*api.Info, error) {
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
		responseParams, err := api.NewResponseParams(config.Response, config.RenderType, config)
		if err != nil {
			return nil, err
		}
		apiInfo.ParamsOut = responseParams
	}

	//logger.Infof(context.Background(), "worker %s AutoUpdateConfig ==nil:%v config: %v el:%+v ",
	//	worker.Router, config.AutoUpdateConfig == nil, jsonx.String(config), config)

	// 处理配置相关
	if config.AutoUpdateConfig != nil {
		// 解析配置结构体，生成表单配置
		configParams, err := api.NewRequestParamsWithFunctionInfo(config.AutoUpdateConfig.ConfigStruct, config.RenderType, config)
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
