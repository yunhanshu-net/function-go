package runner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/yunhanshu-net/function-go/env"
	"github.com/yunhanshu-net/function-go/pkg/dto/usercall"
	"github.com/yunhanshu-net/pkg/trace"

	"github.com/yunhanshu-net/pkg/constants"
	"github.com/yunhanshu-net/pkg/logger"
	"github.com/yunhanshu-net/pkg/typex/files"
)

// ContextLogger 封装logger，保持正确的堆栈信息
type ContextLogger struct {
	ctx context.Context
}

// newContextLogger 创建ContextLogger
func newContextLogger(ctx context.Context) *ContextLogger {
	return &ContextLogger{
		ctx: ctx,
	}
}

// Debug 调试日志
func (l *ContextLogger) Debug(msg string) {
	logger.DebugWrapped(l.ctx, msg)
}

// Debugf 格式化调试日志
func (l *ContextLogger) Debugf(format string, args ...interface{}) {
	logger.DebugfWrapped(l.ctx, format, args...)
}

// Info 信息日志
func (l *ContextLogger) Info(msg string) {
	logger.InfoWrapped(l.ctx, msg)
}

// Infof 格式化信息日志
func (l *ContextLogger) Infof(format string, args ...interface{}) {
	logger.InfofWrapped(l.ctx, format, args...)
}

// Warn 警告日志
func (l *ContextLogger) Warn(msg string) {
	logger.WarnWrapped(l.ctx, msg)
}

// Warnf 格式化警告日志
func (l *ContextLogger) Warnf(format string, args ...interface{}) {
	logger.WarnfWrapped(l.ctx, format, args...)
}

// Error 错误日志
func (l *ContextLogger) Error(msg string, err error) {
	logger.ErrorWrapped(l.ctx, msg, err)
}

// Errorf 格式化错误日志
func (l *ContextLogger) Errorf(format string, args ...interface{}) {
	logger.ErrorfWrapped(l.ctx, format, args...)
}

// Fatal 致命错误日志
func (l *ContextLogger) Fatal(msg string, err error) {
	logger.FatalWrapped(l.ctx, msg, err)
}

// Fatalf 格式化致命错误日志
func (l *ContextLogger) Fatalf(format string, args ...interface{}) {
	logger.FatalfWrapped(l.ctx, format, args...)
}

type Context struct {
	refRouter string
	refMethod string
	context.Context
	user    string
	name    string
	version string

	router string
	method string

	Locker *Lock
	// Logger 绑定的日志记录器
	Logger *ContextLogger

	runner *Runner
}

func NewContext(ctx context.Context, method string, router string, runner *Runner) *Context {
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

	// 设置多个TraceID键，确保各种场景都能正确获取
	c := context.WithValue(ctx, trace.FunctionMsgKey, functionMsg)
	c = context.WithValue(c, constants.TraceID, traceID)
	// 同时设置pkg/logger期望的键
	c = logger.WithContext(c, traceID)

	// 创建Context实例
	contextInstance := &Context{
		Context: c,
		runner:  runner,
		user:    env.User,
		name:    env.Name,
		version: env.Version,
		method:  method,
		router:  router,
		Locker:  newLock(), //分布式锁
	}

	// 初始化Logger，使用新的创建方法
	contextInstance.Logger = newContextLogger(c)

	return contextInstance
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

type UserInfo struct {
	IsLoggedIn bool   `json:"is_logged_in"` //是否已经登陆？
	Username   string `json:"username"`     //用户名
}

func (c *Context) GetUserInfo() UserInfo {

	return UserInfo{
		IsLoggedIn: true,
		Username:   c.user,
	}
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

// ===== Temp 目录工具（统一输出到 ./temp/<traceID>/...） =====

// TempBaseDir 返回当前请求的基础临时目录路径：./temp/<traceID>
// 并确保目录已创建。
func (c *Context) TempBaseDir() (string, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("获取工作目录失败: %w", err)
	}
	base := filepath.Join(workDir, "temp", c.getTraceId())
	if err := os.MkdirAll(base, 0755); err != nil {
		return "", fmt.Errorf("创建临时目录失败: %w", err)
	}
	return base, nil
}

// TempDir 在基础临时目录下拼接子路径并确保创建。
func (c *Context) TempDir(parts ...string) (string, error) {
	base, err := c.TempBaseDir()
	if err != nil {
		return "", err
	}
	full := filepath.Join(append([]string{base}, parts...)...)
	if err := os.MkdirAll(full, 0755); err != nil {
		return "", fmt.Errorf("创建临时子目录失败: %w", err)
	}
	return full, nil
}

// TempUniqueDir 在基础临时目录下创建带时间戳的唯一子目录：./temp/<traceID>/<prefix>/<ns>
func (c *Context) TempUniqueDir(prefix string, parts ...string) (string, error) {
	baseParts := []string{prefix, fmt.Sprintf("%d", time.Now().UnixNano())}
	baseParts = append(baseParts, parts...)
	return c.TempDir(baseParts...)
}

// ===== FS 聚合句柄 =====

// ContextFS 聚合：文件与临时目录相关方法
type ContextFS struct{ ctx *Context }

// FS 返回文件系统聚合句柄
func (c *Context) FS() *ContextFS { return &ContextFS{ctx: c} }

// Time 返回时间处理聚合句柄
func (c *Context) TimeUtils() *ContextTime { return &ContextTime{ctx: c} }

// 目录相关（委托 Context 内部实现）
func (fs *ContextFS) TempBaseDir() (string, error)            { return fs.ctx.TempBaseDir() }
func (fs *ContextFS) TempDir(parts ...string) (string, error) { return fs.ctx.TempDir(parts...) }
func (fs *ContextFS) TempUniqueDir(prefix string, parts ...string) (string, error) {
	return fs.ctx.TempUniqueDir(prefix, parts...)
}

// Files 构建与上传（委托 Context 内部实现）
func (fs *ContextFS) NewFiles(input interface{}) *files.Files { return fs.ctx.NewFiles(input) }
func (fs *ContextFS) NewTemporaryFiles() *files.Files         { return fs.ctx.NewTemporaryFiles() }
func (fs *ContextFS) NewExpiringFiles() *files.Files          { return fs.ctx.NewExpiringFiles() }
func (fs *ContextFS) NewPermanentFiles() *files.Files         { return fs.ctx.NewPermanentFiles() }
func (fs *ContextFS) CreateFilesFromData(filename string, data []byte) (*files.Files, error) {
	return fs.ctx.CreateFilesFromData(filename, data)
}
func (fs *ContextFS) CreateFilesFromPath(localPath string) (*files.Files, error) {
	return fs.ctx.CreateFilesFromPath(localPath)
}

// ===== Config 相关方法 =====

// ConfigManager 获取配置管理器
func (c *Context) ConfigManager() *ConfigManager {
	return GetConfigManager()
}

// GetConfig 获取当前函数的配置结构体值
func (c *Context) GetConfig() interface{} {
	configKey := c.generateConfigKey()
	c.Logger.Infof("GetConfig - 配置键: %s", configKey)

	//todo 这里如果是在回调方法里的话会导致，路由取回调路径无法获取到配置，所以我们应该获取注册的真实路径
	configData := c.ConfigManager().GetByKey(c, configKey)
	if configData == nil {
		c.Logger.Warnf("GetConfig - 配置数据为空")
		return nil
	}

	c.Logger.Infof("GetConfig - 配置数据类型: %T", configData.Data)

	// 从配置管理器获取对应的结构体类型并解析
	result := c.ConfigManager().GetConfigStruct(c, configKey)
	c.Logger.Infof("GetConfig - GetConfigStruct返回类型: %T", result)
	return result
}

// generateConfigKey 生成配置键
func (c *Context) generateConfigKey() string {
	if c.refMethod != "" && c.refRouter != "" {
		return usercall.GenerateConfigKey(c.refRouter, c.refMethod)
	}
	return usercall.GenerateConfigKey(c.router, c.method)
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
