// 用户侧的回调，用户可以在ApiConfig中配置回掉的相关逻辑，不配置就不会触发
package runner

import (
	"encoding/json"
	"fmt"

	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/function-go/pkg/dto/usercall"
	consts "github.com/yunhanshu-net/pkg/constants/usercall"
	"github.com/yunhanshu-net/pkg/logger"
)

// OnPageLoad 当用户进入某个函数的页面后，函数默认调用的行为，用户可以通过这个来初始化表单数据，resetRequest可以返回初始化后的表单数据
type OnPageLoad func(ctx *Context, resp response.Response) (initData *usercall.OnPageLoadResp, err error)

// OnApiCreated 创建新的api时候的回调函数,新增一个api假如新增了一张user表， 可以在这里用gorm的db.AutoMigrate(&User)来创建表，
// 保证新版本的api可以正常使用新增的表 这个api只会在我创建这个api的时候执行一次
type OnApiCreated func(ctx *Context, req *usercall.OnApiCreatedReq) error

// OnApiUpdated 当api发生变更时候的回调函数
type OnApiUpdated func(ctx *Context, req *usercall.OnApiUpdatedReq) error

// BeforeApiDelete  api删除前触发回调，比如该api删除的话，可以备份某些数据
type BeforeApiDelete func(ctx *Context, req *usercall.BeforeApiDeleteReq) error

// AfterApiDeleted  api删除后触发回调，比如该api删除的话，可以在这里做一些操作，比如删除该api对应的表
type AfterApiDeleted func(ctx *Context, req *usercall.AfterApiDeletedReq) error

// BeforeRunnerClose 程序结束前的回调函数，可以在程序结束前做一些操作，比如上报一些数据
type BeforeRunnerClose func(ctx *Context, req *usercall.BeforeRunnerCloseReq) error

// AfterRunnerClose 程序结束后的回调函数，可以在程序结束后做一些操作，比如清理某些文件
type AfterRunnerClose func(ctx *Context, req *usercall.AfterRunnerCloseReq) error

// OnVersionChange 每次版本发生变更都会回调这个函数（新增/删除api）
type OnVersionChange func(ctx *Context, req *usercall.OnVersionChangeReq) error

// OnInputFuzzy 模糊搜索回调函数，比如搜索用户，可以在这里做一些操作，比如根据用户名模糊搜索用户，然后返回用户列表
type OnInputFuzzy func(ctx *Context, req *usercall.OnInputFuzzyReq) (*usercall.OnInputFuzzyResp, error)

// OnInputValidate 验证输入框输入的名称是否重复或者输入是否合法
type OnInputValidate func(ctx *Context, req *usercall.OnInputValidateReq) (*usercall.OnInputValidateResp, error)

// OnTableDeleteRows 当返回前端的数据是table类型时候，前端会把数据渲染成表格，这时候表格数据会有删除的行为，实现这个函数用来删除数据
type OnTableDeleteRows func(ctx *Context, req *usercall.OnTableDeleteRowsReq) (*usercall.OnTableDeleteRowsResp, error)

// OnTableUpdateRows 当返回前端的数据是table类型时候，前端会把数据渲染成表格，这时候表格数据会有更新的行为，实现这个函数用来更新数据
type OnTableUpdateRows func(ctx *Context, req *usercall.OnTableUpdateRowsReq) (*usercall.OnTableUpdateRowsResp, error)

// OnTableAddRows 当返回前端的数据是table类型时候，前端会把数据渲染成表格，这时候表格数据会有新增的行为，实现这个函数用来新增数据
type OnTableAddRows func(ctx *Context, req *usercall.OnTableAddRowsReq) (*usercall.OnTableAddRowsResp, error)

// OnTableSearch 当返回前端的数据是table类型时候，前端会把数据渲染成表格，这时候表格数据会有搜索的行为，实现这个函数用来搜索数据
type OnTableSearch func(ctx *Context, req *usercall.OnTableSearchReq) (*usercall.OnTableSearchResp, error)

// OnDryRun DryRun 回调函数，用于预览危险操作
type OnDryRun func(ctx *Context, req *usercall.OnDryRunReq) (*usercall.OnDryRunResp, error)

func (r *Runner) _callback(ctx *Context, req *usercall.Request, resp response.Response) (err error) {
	var res usercall.Response

	// 记录请求参数
	//reqJSON, _ := json.Marshal(req)
	//logger.Infof(ctx, "处理回调 [类型:%s] [路由:%s] [方法:%s] 请求参数: %s", req.Type, req.Router, req.Method, string(reqJSON))

	worker, exist := r.getRouter(req.Router, req.Method)
	if !exist {
		err = fmt.Errorf("router not found")
		logger.Infof(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
		return err
	}

	//todo
	// 检查是否是 FunctionOptions 类型
	//function, ok := worker.Option.(*FunctionOptions)
	//if !ok || function == nil {
	//	err = fmt.Errorf("router config nil or not FunctionOptions")
	//	logger.Infof(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
	//	return err
	//}
	callbacks := worker.Option.GetCallbacks()
	//baseConf:=worker.Option.GetBaseConfig()

	switch req.Type {
	case consts.CallbackTypeOnCreateTables:
		err1 := worker.CreateTables(ctx)
		if err1 != nil {
			logger.Infof(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err1
		}
	// 页面加载回调
	case consts.CallbackTypeOnPageLoad:
		callbackTypeOnPageLoad, yes := callbacks[consts.CallbackTypeOnPageLoad]
		if !yes {
			logger.Errorf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("OnPageLoad handler not configured %s不存在", req.Type)
		}
		onPageLoad, ok := callbackTypeOnPageLoad.(OnPageLoad)
		if !ok {
			logger.Errorf(ctx, "回调处理失败 [类型不匹配:%s]: %v", req.Type, err)
			return fmt.Errorf("回调%s类型不匹配", req.Type)
		}
		userResp := &response.RunFunctionResp{}
		rsp, err := onPageLoad(ctx, userResp)
		if err != nil {
			logger.Infof(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("OnPageLoad failed: %w", err)
		}
		type OnPageLoadResp struct {
			Multiple   bool                        `json:"multiple"` //是否多返回值
			Request    interface{}                 `json:"request"`  //会初始化前端的表单参数
			Response   interface{}                 `json:"response"` //会初始化前端的的响应参数
			AutoRun    bool                        `json:"auto_run"` //是否自动运行
			DisableRun bool                        `json:"disable_run"`
			Message    *usercall.OnPageLoadMessage `json:"message"`
		}

		if rsp == nil {
			return resp.Form(&OnPageLoadResp{}).Build()
		}
		rs := &OnPageLoadResp{
			Multiple:   userResp.Multiple,
			Request:    rsp.Request,
			Response:   userResp.GetData(),
			DisableRun: rsp.DisableRun,
			Message:    rsp.Message,
			AutoRun:    rsp.AutoRun,
		}
		return resp.Form(rs).Build()

	// API 生命周期回调
	case consts.UserCallTypeOnApiCreated:
		var reqData usercall.OnApiCreatedReq
		callbackOnApiCreated, yes := callbacks[consts.UserCallTypeOnApiCreated]
		if !yes {
			err = fmt.Errorf("OnApiCreated handler not configured %s不存在", req.Type)
			logger.Infof(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		onApiCreated, ok := callbackOnApiCreated.(OnApiCreated)
		if !ok {
			logger.Errorf(ctx, "回调处理失败 [类型不匹配:%s]: %v", req.Type, err)
			return fmt.Errorf("回调%s类型不匹配", req.Type)
		}

		if err := onApiCreated(ctx, &reqData); err != nil {
			logger.Infof(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("OnApiCreatedReq failed: %w", err)
		}
		logger.Infof(ctx, "回调处理成功 [类型:%s]", req.Type)

	case consts.CallbackTypeOnApiUpdated:
		var reqData usercall.OnApiUpdatedReq
		callbackOnApiUpdated, yes := callbacks[consts.CallbackTypeOnApiUpdated]
		if !yes {
			err = fmt.Errorf("OnApiUpdated handler not configured %s不存在", req.Type)
			logger.Infof(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		onApiUpdated, ok := callbackOnApiUpdated.(OnApiUpdated)
		if !ok {
			logger.Errorf(ctx, "回调处理失败 [类型不匹配:%s]: %v", req.Type, err)
			return fmt.Errorf("回调%s类型不匹配", req.Type)
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.Infof(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("OnApiUpdatedReq decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.Infof(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		if err := onApiUpdated(ctx, &reqData); err != nil {
			logger.Infof(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("OnApiUpdatedReq failed: %w", err)
		}
		logger.Infof(ctx, "回调处理成功 [类型:%s]", req.Type)

	case consts.CallbackTypeBeforeApiDelete:
		var reqData usercall.BeforeApiDeleteReq
		callbackBeforeApiDelete, yes := callbacks[consts.CallbackTypeBeforeApiDelete]
		if !yes {
			err = fmt.Errorf("BeforeApiDelete handler not configured %s不存在", req.Type)
			logger.Infof(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		beforeApiDelete, ok := callbackBeforeApiDelete.(BeforeApiDelete)
		if !ok {
			logger.Errorf(ctx, "回调处理失败 [类型不匹配:%s]: %v", req.Type, err)
			return fmt.Errorf("回调%s类型不匹配", req.Type)
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.Infof(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("BeforeApiDeleteReq decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.Infof(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		if err := beforeApiDelete(ctx, &reqData); err != nil {
			logger.Infof(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("BeforeApiDeleteReq failed: %w", err)
		}
		logger.Infof(ctx, "回调处理成功 [类型:%s]", req.Type)

	case consts.CallbackTypeAfterApiDeleted:
		var reqData usercall.AfterApiDeletedReq
		callbackAfterApiDeleted, yes := callbacks[consts.CallbackTypeAfterApiDeleted]
		if !yes {
			err = fmt.Errorf("AfterApiDeleted handler not configured %s不存在", req.Type)
			logger.Infof(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		afterApiDeleted, ok := callbackAfterApiDeleted.(AfterApiDeleted)
		if !ok {
			logger.Errorf(ctx, "回调处理失败 [类型不匹配:%s]: %v", req.Type, err)
			return fmt.Errorf("回调%s类型不匹配", req.Type)
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.Infof(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("AfterApiDeletedReq decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.Infof(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		if err := afterApiDeleted(ctx, &reqData); err != nil {
			logger.Infof(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("AfterApiDeletedReq failed: %w", err)
		}
		logger.Infof(ctx, "回调处理成功 [类型:%s]", req.Type)

	// Runner 生命周期回调
	case consts.CallbackTypeBeforeRunnerClose:
		var reqData usercall.BeforeRunnerCloseReq
		callbackBeforeRunnerClose, yes := callbacks[consts.CallbackTypeBeforeRunnerClose]
		if !yes {
			err = fmt.Errorf("BeforeRunnerClose handler not configured %s不存在", req.Type)
			logger.Infof(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		beforeRunnerClose, ok := callbackBeforeRunnerClose.(BeforeRunnerClose)
		if !ok {
			logger.Errorf(ctx, "回调处理失败 [类型不匹配:%s]: %v", req.Type, err)
			return fmt.Errorf("回调%s类型不匹配", req.Type)
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.Infof(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("BeforeRunnerCloseReq decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.Infof(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		if err := beforeRunnerClose(ctx, &reqData); err != nil {
			logger.Infof(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("BeforeRunnerCloseReq failed: %w", err)
		}
		logger.Infof(ctx, "回调处理成功 [类型:%s]", req.Type)

	case consts.CallbackTypeAfterRunnerClose:
		var reqData usercall.AfterRunnerCloseReq
		callbackAfterRunnerClose, yes := callbacks[consts.CallbackTypeAfterRunnerClose]
		if !yes {
			err = fmt.Errorf("AfterRunnerClose handler not configured %s不存在", req.Type)
			logger.Infof(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		afterRunnerClose, ok := callbackAfterRunnerClose.(AfterRunnerClose)
		if !ok {
			logger.Errorf(ctx, "回调处理失败 [类型不匹配:%s]: %v", req.Type, err)
			return fmt.Errorf("回调%s类型不匹配", req.Type)
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.Infof(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("AfterRunnerCloseReq decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.Infof(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		if err := afterRunnerClose(ctx, &reqData); err != nil {
			logger.Infof(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("AfterRunnerCloseReq failed: %w", err)
		}
		logger.Infof(ctx, "回调处理成功 [类型:%s]", req.Type)

	// 版本控制回调
	case consts.CallbackTypeOnVersionChange:
		var reqData usercall.OnVersionChangeReq
		callbackOnVersionChange, yes := callbacks[consts.CallbackTypeOnVersionChange]
		if !yes {
			err = fmt.Errorf("OnVersionChange handler not configured %s不存在", req.Type)
			logger.Infof(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		onVersionChange, ok := callbackOnVersionChange.(OnVersionChange)
		if !ok {
			logger.Errorf(ctx, "回调处理失败 [类型不匹配:%s]: %v", req.Type, err)
			return fmt.Errorf("回调%s类型不匹配", req.Type)
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.Infof(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("OnVersionChangeReq decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.Infof(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		if err := onVersionChange(ctx, &reqData); err != nil {
			logger.Infof(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("OnVersionChangeReq failed: %w", err)
		}
		logger.Infof(ctx, "回调处理成功 [类型:%s]", req.Type)

	// 输入交互回调
	case consts.CallbackTypeOnInputFuzzy:
		var reqData usercall.OnInputFuzzyReq
		if err = req.DecodeData(&reqData); err != nil {
			logger.Infof(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("OnInputFuzzyReq decode failed: %w", err)
		}
		// 对于Map类型的回调，需要特殊处理
		callbackOnInputFuzzyMap, yes := callbacks["OnInputFuzzyMap"]
		if !yes {
			err = fmt.Errorf("OnInputFuzzyMap handler not configured")
			logger.Infof(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		onInputFuzzyMap, ok := callbackOnInputFuzzyMap.(map[string]OnInputFuzzy)
		if !ok {
			logger.Errorf(ctx, "回调处理失败 [类型不匹配:%s]: %v", req.Type, err)
			return fmt.Errorf("回调%s类型不匹配", req.Type)
		}
		fuzzy := onInputFuzzyMap[reqData.Code]
		if fuzzy == nil {
			return fmt.Errorf("OnInputFuzzyReq handler not configured")
		}
		fuzzyResp, err := fuzzy(ctx, &reqData)
		if err != nil {
			logger.Errorf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("OnInputFuzzyReq failed: %w", err)
		}
		logger.Infof(ctx, "回调处理成功 [类型:%s]", req.Type)
		return resp.Form(fuzzyResp).Build()

	case consts.CallbackTypeOnInputValidate:
		var reqData usercall.OnInputValidateReq
		if err = req.DecodeData(&reqData); err != nil {
			logger.Infof(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("OnInputValidateReq decode failed: %w", err)
		}
		// 对于Map类型的回调，需要特殊处理
		callbackOnInputValidateMap, yes := callbacks["OnInputValidateMap"]
		if !yes {
			err = fmt.Errorf("OnInputValidateMap handler not configured")
			logger.Infof(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		onInputValidateMap, ok := callbackOnInputValidateMap.(map[string]OnInputValidate)
		if !ok {
			logger.Errorf(ctx, "回调处理失败 [类型不匹配:%s]: %v", req.Type, err)
			return fmt.Errorf("回调%s类型不匹配", req.Type)
		}
		validate := onInputValidateMap[reqData.Code]
		if validate == nil {
			return fmt.Errorf("OnInputValidateReq handler not configured")
		}
		validateResp, err := validate(ctx, &reqData)
		if err != nil {
			logger.Errorf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("OnInputValidateReq failed: %w", err)
		}
		logger.Infof(ctx, "回调处理成功 [类型:%s]", req.Type)
		return resp.Form(validateResp).Build()

	// 表格操作回调
	case consts.CallbackTypeOnTableDeleteRows:
		var reqData usercall.OnTableDeleteRowsReq
		if err = req.DecodeData(&reqData); err != nil {
			logger.Infof(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("OnTableDeleteRowsReq decode failed: %w", err)
		}
		tb, ok := worker.Option.(*TableFunctionOptions)
		if !ok {
			return fmt.Errorf("onTableDeleteRowsReq 类型不匹配 failed")
		}
		var respData interface{}
		if tb.OnTableDeleteRows != nil {
			rows, err := tb.OnTableDeleteRows(ctx, &reqData)
			if err != nil {
				logger.Errorf(ctx, "回调处理失败 [类型:%s]: OnTableDeleteRows 解码失败 %v", req.Type, err)
			}
			respData = rows
		} else {
			if tb.AutoCrudTable != nil {
				err = tb.defaultDeleteRows(ctx, &reqData)
				if err != nil {
					logger.Errorf(ctx, "回调处理失败 [类型:%s]: defaultDeleteRows 解码失败 %v", req.Type, err)
				}
			}
		}

		res.Response = respData

		// 记录响应参数
		respDataJSON, _ := json.Marshal(respData)
		logger.Infof(ctx, "回调处理成功 [类型:%s] 响应: %s", req.Type, respDataJSON)
		return resp.Form(&usercall.OnTableDeleteRowsResp{}).Build()
	case consts.CallbackTypeOnTableUpdateRows:
		var reqData usercall.OnTableUpdateRowsReq

		if err := req.DecodeData(&reqData); err != nil {
			logger.Infof(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("OnTableUpdateRowReq decode failed: %w", err)
		}
		tb, ok := worker.Option.(*TableFunctionOptions)
		if !ok {
			return fmt.Errorf("onTableDeleteRowsReq 类型不匹配 failed")
		}
		var respData interface{}
		if tb.OnTableUpdateRows != nil {
			rsp, err := tb.OnTableUpdateRows(ctx, &reqData)
			if err != nil {
				logger.Errorf(ctx, "回调处理失败 [类型:%s]: OnTableUpdateRows 解码失败 %v", req.Type, err)
			}
			respData = rsp
		} else {
			if tb.AutoCrudTable != nil {
				err = tb.defaultUpdateRows(ctx, &reqData)
				if err != nil {
					logger.Errorf(ctx, "回调处理失败 [类型:%s]: defaultUpdateRows 解码失败 %v", req.Type, err)
				}
			}
		}
		res.Response = respData

		// 记录响应参数
		respDataJSON, _ := json.Marshal(respData)
		logger.Infof(ctx, "回调处理成功 [类型:%s] 响应: %s", req.Type, respDataJSON)
		return resp.Form(&usercall.OnTableUpdateRowsResp{}).Build()
	case consts.CallbackTypeOnTableAddRows:
		var reqData usercall.OnTableAddRowsReq

		if err := req.DecodeData(&reqData); err != nil {
			logger.Infof(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("OnTableAddRowsReq decode failed: %w", err)
		}

		tb, ok := worker.Option.(*TableFunctionOptions)
		if !ok {
			return fmt.Errorf("CallbackTypeOnTableAddRows 类型不匹配 failed")
		}

		var respData interface{}

		if tb.OnTableAddRows != nil {
			rsp, err := tb.OnTableAddRows(ctx, &reqData)
			if err != nil {
				logger.Errorf(ctx, "回调处理失败 [类型:%s]: OnTableAddRowsReq 处理失败req:%+v %v", req.Type, req, err)
				return err
			}
			respData = rsp

		} else {
			if tb.AutoCrudTable != nil {
				err = tb.defaultAddRows(ctx, &reqData)
				if err != nil {
					logger.Errorf(ctx, "回调处理失败 [类型:%s]: defaultAddRows 执行失败 req:%+v %v", req.Type, req, err)
					return err
				}

			}
		}
		res.Response = respData
		// 记录响应参数
		respDataJSON, _ := json.Marshal(respData)
		logger.Infof(ctx, "回调处理成功 [类型:%s]请求：%+v 响应: %s", req.Type, reqData, respDataJSON)
		return resp.Form(&usercall.OnTableAddRowsResp{}).Build()
	case consts.CallbackTypeOnTableSearch:
		var reqData usercall.OnTableSearchReq
		callbackOnTableSearch, yes := callbacks[consts.CallbackTypeOnTableSearch]
		if !yes {
			err = fmt.Errorf("OnTableSearch handler not configured %s不存在", req.Type)
			logger.Infof(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		onTableSearch, ok := callbackOnTableSearch.(OnTableSearch)
		if !ok {
			logger.Errorf(ctx, "回调处理失败 [类型不匹配:%s]: %v", req.Type, err)
			return fmt.Errorf("回调%s类型不匹配", req.Type)
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.Infof(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("OnTableSearchReq decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.Infof(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		respData, err := onTableSearch(ctx, &reqData)
		if err != nil {
			logger.Infof(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("OnTableSearchReq failed: %w", err)
		}
		res.Response = respData

		// 记录响应参数
		respDataJSON, _ := json.Marshal(respData)
		logger.Infof(ctx, "回调处理成功 [类型:%s] 响应: %s", req.Type, respDataJSON)

	// DryRun 回调
	case consts.CallbackTypeOnDryRun:
		var reqData usercall.OnDryRunReq
		callbackOnDryRun, yes := callbacks[consts.CallbackTypeOnDryRun]
		if !yes {
			err = fmt.Errorf("OnDryRun handler not configured %s不存在", req.Type)
			logger.Infof(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		onDryRun, ok := callbackOnDryRun.(OnDryRun)
		if !ok {
			logger.Errorf(ctx, "回调处理失败 [类型不匹配:%s]: %v", req.Type, err)
			return fmt.Errorf("回调%s类型不匹配", req.Type)
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.Infof(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("OnDryRunReq decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.Infof(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		respData, err := onDryRun(ctx, &reqData)
		if err != nil {
			logger.Errorf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("OnDryRun failed: %w", err)
		}
		res.Response = respData

		// 记录响应参数
		respDataJSON, _ := json.Marshal(respData)
		logger.Infof(ctx, "回调处理成功 [类型:%s] 响应: %s", req.Type, respDataJSON)
		return resp.Form(respData).Build()

	// 配置管理回调
	case consts.CallbackTypeOnUpdateConfig:
		var reqData usercall.UpdateConfigReq
		if err := req.DecodeData(&reqData); err != nil {
			logger.Infof(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("UpdateConfigReq decode failed: %w", err)
		}

		// 验证参数
		if reqData.Router == "" {
			return resp.Form(&usercall.UpdateConfigResp{Success: false, Error: "router参数不能为空"}).Build()
		}

		if reqData.Method == "" {
			return resp.Form(&usercall.UpdateConfigResp{Success: false, Error: "method参数不能为空"}).Build()
		}

		// 使用请求对象生成配置键
		configKey := reqData.GenerateConfigKey()

		// 获取配置管理器
		configManager := GetConfigManager()

		// 更新配置
		err := configManager.UpdateConfig(ctx, configKey, reqData.ToConfigData())
		if err != nil {
			return resp.Form(&usercall.UpdateConfigResp{
				Success: false,
				Error:   err.Error(),
			}).Build()
		}

		return resp.Form(&usercall.UpdateConfigResp{
			Success: true,
			Message: "配置更新成功",
		}).Build()

	case consts.CallbackTypeOnGetConfig:
		var reqData usercall.GetConfigReq
		if err := req.DecodeData(&reqData); err != nil {
			logger.Infof(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("GetConfigReq decode failed: %w", err)
		}

		// 验证参数
		if reqData.Router == "" {
			return resp.Form(&usercall.GetConfigResp{
				Success: false,
				Error:   "router参数不能为空",
			}).Build()
		}

		if reqData.Method == "" {
			return resp.Form(&usercall.GetConfigResp{
				Success: false,
				Error:   "method参数不能为空",
			}).Build()
		}

		// 使用请求对象生成配置键
		configKey := reqData.GenerateConfigKey()

		// 获取配置管理器
		configManager := GetConfigManager()
		configData := configManager.GetByKey(ctx, configKey)
		if configData == nil {
			return resp.Form(&usercall.GetConfigResp{
				Success: false,
				Error:   "配置未找到",
			}).Build()
		}

		return resp.Form(&usercall.GetConfigResp{
			Success: true,
			Config:  configData,
		}).Build()

	default:
		err = fmt.Errorf("unsupported callback type: %s", req.Type)
		logger.Infof(ctx, "回调处理失败 [类型:%s]: 不支持的回调类型", req.Type)
		return err
	}

	err = resp.Form(res).Build()
	if err != nil {
		logger.Infof(ctx, "回调处理失败 [类型:%s]: 构建响应失败 %v", req.Type, err)
		return err
	}

	// 在没有响应参数被记录的情况下，记录最终响应
	if res.Response != nil && (req.Type == consts.UserCallTypeOnApiCreated ||
		req.Type == consts.CallbackTypeOnApiUpdated ||
		req.Type == consts.CallbackTypeBeforeApiDelete ||
		req.Type == consts.CallbackTypeAfterApiDeleted ||
		req.Type == consts.CallbackTypeBeforeRunnerClose ||
		req.Type == consts.CallbackTypeAfterRunnerClose ||
		req.Type == consts.CallbackTypeOnVersionChange) {
		resJSON, _ := json.Marshal(res.Response)
		logger.Infof(ctx, "回调处理完成 [类型:%s] 响应: %s", req.Type, resJSON)
	}

	return nil
}

//func (r *Runner) _syscall(ctx *Context, req *syscall.Request, resp response.Response) error {
//	s, err := __syscall(ctx, r, req)
//	if err != nil {
//		return err
//	}
//	return resp.Form(s.Data).Build()
//}
//func __syscall(ctx *Context, r *Runner, req *syscall.Request) (resp *syscall.Response, err error) {
//	resp = new(syscall.Response)
//
//	if req.CallbackType == sysconsts.TypeCreateTables {
//		tablesReq, ok := req.Data.(*syscall.OnCreateTablesReq)
//		if !ok {
//			err = fmt.Errorf("OnCreateTablesReq decode failed: %w", err)
//			return resp, err
//		}
//		function, exist := r.getRouter(tablesReq.Router, tablesReq.Method)
//		if !exist {
//			return resp, fmt.Errorf("router not found: %s", tablesReq.Router)
//		}
//		err = function.CreateTables(ctx)
//		if err != nil {
//			return resp, err
//		}
//		resp.Data = "ok"
//	}
//	return resp, err
//}
