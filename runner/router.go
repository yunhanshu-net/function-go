package runner

import (
	"github.com/yunhanshu-net/function-go/pkg/dto/request"
	"github.com/yunhanshu-net/function-go/pkg/dto/response"
	"github.com/yunhanshu-net/pkg/logger"
	"gorm.io/gorm/schema"
	"strings"
)

type routerInfo struct {
	Handel interface{} `json:"-"`
	key    string
	Router string
	Method string
	//FunctionInfo *FunctionOptions
	Option Option
}

func (r *routerInfo) CreateTables(ctx *Context) error {
	if r.Option == nil {
		return nil
	}
	if r.Option.GetCreateTables() == nil {
		return nil
	}
	for _, table := range r.Option.GetCreateTables() {
		err := ctx.MustGetOrInitDB().AutoMigrate(table)
		if err != nil {
			logger.Errorf(ctx, "create table %+v  error: %v", table, err)
		}
		tabler, ok := table.(schema.Tabler)
		if !ok {
			logger.Errorf(ctx, "create table %+v  error: %v", table, err)
			continue
		}
		logger.Infof(ctx, "CreateTables create table %s success", tabler.TableName())
	}
	return nil

}

func fmtKey(router string, method string) string {
	if !strings.HasPrefix(router, "/") {
		router = "/" + router
	}
	router = strings.TrimSuffix(router, "/")
	return router + "." + strings.ToUpper(method)
}

// getRouter 获取路由
func (r *Runner) getRouter(router string, method string) (worker *routerInfo, exist bool) {
	worker, ok := r.routerMap[fmtKey(router, method)]
	return worker, ok
}

func (r *routerInfo) call(ctx *Context, reqBody interface{}) (req *request.RunFunctionReq, resp *response.RunFunctionResp, err error) {
	// 使用读锁访问缓存
	handlerCacheMux.RLock()
	meta, ok := handlerCacheMap[r.key]
	handlerCacheMux.RUnlock()

	if !ok {
		// 需要创建新的缓存项，使用写锁
		handlerCacheMux.Lock()
		// 双重检查锁定模式
		meta, ok = handlerCacheMap[r.key]
		if !ok {
			h := &handlerMeta{}
			handlerCacheMap[r.key] = h
			meta = h
		}
		handlerCacheMux.Unlock()
	}

	//metaVal, _ := handlerCache.LoadOrStore(r.key, &handlerMeta{})
	//meta := metaVal.(*handlerMeta)
	// 确保只初始化一次
	meta.once.Do(func() {
		meta.meta, meta.initError = buildRuntimeMeta(r.Handel)
	})
	if meta.initError != nil {
		return nil, nil, meta.initError
	}
	req = new(request.RunFunctionReq)
	resp = new(response.RunFunctionResp)
	//ctx1 := &Context{Context: ctx}
	err = doCall(r.Method, meta.meta, ctx, resp, reqBody)
	if err != nil {
		return nil, nil, err
	}
	return req, resp, nil
}
