package runner

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/yunhanshu-net/function-go/env"
	"github.com/yunhanshu-net/pkg/trace"

	"github.com/yunhanshu-net/pkg/constants"
	"github.com/yunhanshu-net/pkg/typex/files"
)

type Context struct {
	context.Context
	user    string
	name    string
	version string

	router string
	method string
}

func NewContext(ctx context.Context, method string, router string) *Context {
	// 获取trace_id
	traceID := ""
	if value := ctx.Value(constants.TraceID); value != nil {
		if tid, ok := value.(string); ok {
			traceID = tid
		}
	}
	if traceID == "" {
		// 如果没有trace_id，生成一个简单的
		traceID = fmt.Sprintf("ctx-%d", time.Now().UnixNano())
	}

	// 创建FunctionMsg
	functionMsg := &trace.FunctionMsg{
		User:         env.User,
		Runner:       env.Name,
		Version:      env.Version,
		Method:       method,
		Router:       router,
		TraceID:      traceID,
		UploadConfig: getUploadConfig(),
	}

	c := context.WithValue(ctx, trace.FunctionMsgKey, functionMsg)

	return &Context{
		Context: c,
		user:    env.User,
		name:    env.Name,
		version: env.Version,
		method:  method,
		router:  router,
	}
}

func (c *Context) getDBName() string {
	return fmt.Sprintf("%s_%s.db", c.user, c.name)
}

func (c *Context) getTraceId() string {
	value := c.Context.Value(constants.TraceID)
	if value == nil {
		return ""
	}
	v, ok := value.(string)
	if ok {
		return v
	}
	return ""
}

func (c *Context) GetUsername() string {
	return ""
}

func (c *Context) GetFile() string {
	return ""
}

// GetUploadPath 获取当前函数的上传路径
func (c *Context) GetUploadPath() string {
	value := c.Context.Value(trace.FunctionMsgKey)
	if value == nil {
		return ""
	}

	functionMsg, ok := value.(*trace.FunctionMsg)
	if !ok {
		return ""
	}

	return functionMsg.GetUploadPath()
}

// GetFunctionMsg 获取函数消息对象
func (c *Context) GetFunctionMsg() *trace.FunctionMsg {
	value := c.Context.Value(trace.FunctionMsgKey)
	if value == nil {
		return nil
	}

	functionMsg, ok := value.(*trace.FunctionMsg)
	if !ok {
		return nil
	}

	return functionMsg
}

// ===== Files 相关方法 =====

// NewFiles 创建新的文件集合，自动设置context
func (c *Context) NewFiles(input interface{}) *files.Files {
	return files.NewFiles(input).SetContext(c.Context)
}

// NewTemporaryFiles 创建临时文件集合（下载一次后删除）
func (c *Context) NewTemporaryFiles() *files.Files {
	return files.NewFiles([]string{}).
		SetContext(c.Context).
		SetTemporary()
}

// NewExpiringFiles 创建有效期文件集合（7天后过期）
func (c *Context) NewExpiringFiles() *files.Files {
	return files.NewFiles([]string{}).
		SetContext(c.Context).
		SetExpiring7Days()
}

// NewPermanentFiles 创建永久文件集合（无限制）
func (c *Context) NewPermanentFiles() *files.Files {
	return files.NewFiles([]string{}).
		SetContext(c.Context).
		SetUnlimited()
}

// CreateFilesFromData 从数据创建文件并立即上传到规范路径
func (c *Context) CreateFilesFromData(filename string, data []byte) (*files.Files, error) {
	files := c.NewFiles([]string{})
	err := files.AddFileFromData(filename, data)
	if err != nil {
		return nil, err
	}
	return files, nil
}

// CreateFilesFromPath 从本地路径创建文件并立即上传到规范路径
func (c *Context) CreateFilesFromPath(localPath string) (*files.Files, error) {
	files := c.NewFiles([]string{})
	err := files.AddFileFromPath(localPath)
	if err != nil {
		return nil, err
	}
	return files, nil
}

// ===== Config 相关方法 =====

// Config 获取配置管理器
func (c *Context) Config() *ConfigManager {
	return GetConfigManager()
}

// GetConfig 获取当前函数的配置结构体值
func (c *Context) GetConfig() interface{} {
	configKey := c.generateConfigKey()
	configData := c.Config().GetByKey(c, configKey)
	if configData == nil {
		return nil
	}

	// 从配置管理器获取对应的结构体类型并解析
	return c.Config().GetConfigStruct(c, configKey)
}

// generateConfigKey 生成配置键
func (c *Context) generateConfigKey() string {
	// 处理路由路径，将 / 替换为 . 以安全地用作配置键
	safeRouter := strings.ReplaceAll(c.router, "/", ".")
	// 移除前后的点
	safeRouter = strings.Trim(safeRouter, ".")
	return fmt.Sprintf("function.%s.%s", safeRouter, c.method)
}

// ===== 基础信息方法 =====

// User 获取用户信息
func (c *Context) User() string {
	return c.user
}

// Name 获取运行器名称
func (c *Context) Name() string {
	return c.name
}

// Version 获取版本信息
func (c *Context) Version() string {
	return c.version
}

// Router 获取路由信息
func (c *Context) Router() string {
	return c.router
}

// Method 获取HTTP方法
func (c *Context) Method() string {
	return c.method
}
