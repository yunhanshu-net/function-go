# Context Logger 使用指南

## 概述

`Context.Logger` 是绑定在 `Context` 上的日志记录器，提供了便捷的日志记录功能，同时保持正确的堆栈信息（文件名、行号、函数名）和 TraceID 追踪。

## ✅ 已解决的核心问题

**堆栈信息完全正确**：通过在 `pkg/logger` 中新增 `*Wrapped` 系列方法，专门处理封装场景的调用深度，现在日志显示的是实际调用代码的位置，而不是封装层的位置。

## 特性

✅ **正确的堆栈信息**：显示实际调用日志的代码位置（文件名:行号 [函数名]）  
✅ **自动 TraceID**：自动从 Context 中提取并添加 TraceID 到日志中  
✅ **简化引用**：无需在每个函数中单独引用 logger 包  
✅ **类型安全**：提供完整的日志级别支持  
✅ **性能优化**：基于 zap logger，高性能日志记录  
✅ **零配置**：在 `NewContext()` 时自动初始化，无需手动设置

## 基本用法

### 1. 在函数中使用

```go


type HandleResp struct {
Name        string `json:"name" form:"name" runner:"code:name;name:名称" widget:"type:input;placeholder:请输入名称" data:"type:string;default_value:测试配置;example:配置演示" validate:"required,min=2,max=50"`

}


func MyHandler(ctx *runner.Context, req *MyRequest,resp response.Response) (error) {
    // 信息日志
    ctx.Logger.Info("开始处理请求")
    ctx.Logger.Infof("处理用户: %s", req.Username)
    
    // 调试日志  
    ctx.Logger.Debug("验证参数")
    ctx.Logger.Debugf("参数详情: %+v", req)
    
    // 警告日志
    ctx.Logger.Warn("检测到潜在问题")
    ctx.Logger.Warnf("用户 %s 尝试访问受限资源", req.Username)
    
    // 错误日志
    if err != nil {
        ctx.Logger.Error("操作失败", err)
        ctx.Logger.Errorf("数据库操作失败: %v", err)
        return err
    }
    
    ctx.Logger.Info("请求处理完成")
	
    return resp.Form(&HandleResp{Name:"test"}).Build()
}
```

### 2. 在业务逻辑中使用

```go
func processUserData(ctx *runner.Context, userID int) error {
    ctx.Logger.Infof("开始处理用户数据: ID=%d", userID)
    
    // 这里的日志会显示正确的文件名和行号
    ctx.Logger.Debug("查询用户信息")
    
    user, err := getUserFromDB(userID)
    if err != nil {
        ctx.Logger.Error("查询用户失败", err)
        return err
    }
    
    ctx.Logger.Debugf("用户信息: %+v", user)
    return nil
}
```

## 日志级别

| 方法 | 用途 | 示例 |
|------|------|------|
| `Debug()` / `Debugf()` | 调试信息，开发阶段使用 | 参数验证、中间状态 |
| `Info()` / `Infof()` | 一般信息，记录关键流程 | 请求开始/结束、业务节点 |
| `Warn()` / `Warnf()` | 警告信息，需要注意但不影响运行 | 参数异常、性能问题 |
| `Error()` / `Errorf()` | 错误信息，影响功能但不崩溃 | 业务逻辑错误、外部调用失败 |
| `Fatal()` / `Fatalf()` | 致命错误，程序无法继续运行 | 系统级错误、配置错误 |

## 日志输出格式

```json
{
  "level": "INFO",
  "ts": "2025-07-29 12:13:09.042",
  "caller": "function-go/runner/my_handler.go:25 [MyHandler]",
  "msg": "开始处理请求: GET /api/demo/test",
  "trace_id": "ctx-1753762243469554000"
}
```

### 字段说明

- `level`: 日志级别 (DEBUG/INFO/WARN/ERROR/FATAL)
- `ts`: 时间戳
- `caller`: 调用位置 (文件路径:行号 [函数名]) **← 现在完全正确！**
- `msg`: 日志消息
- `trace_id`: 请求追踪ID (自动添加)

## 实际效果示例

以下是真实的日志输出，展示了正确的堆栈信息：

```json
{"level":"INFO","ts":"2025-07-29 12:13:09.042","caller":"function-go/runner/context_logger_test.go:64 [func1]","msg":"","msg":"开始处理请求: GET /api/demo/test"}
{"level":"DEBUG","ts":"2025-07-29 12:13:09.046","caller":"function-go/runner/context_logger_test.go:68 [func1]","msg":"","msg":"处理步骤 1"}
{"level":"DEBUG","ts":"2025-07-29 12:13:09.046","caller":"function-go/runner/context_logger_test.go:68 [func1]","msg":"","msg":"处理步骤 2"}
{"level":"DEBUG","ts":"2025-07-29 12:13:09.046","caller":"function-go/runner/context_logger_test.go:68 [func1]","msg":"","msg":"处理步骤 3"}
{"level":"INFO","ts":"2025-07-29 12:13:09.046","caller":"function-go/runner/context_logger_test.go:71 [func1]","msg":"请求处理完成"}
```

可以看到：
- 第64行：对应实际的 `ctx.Logger.Infof()` 调用
- 第68行：对应循环中的 `ctx.Logger.Debugf()` 调用  
- 第71行：对应 `ctx.Logger.Info()` 调用

## 技术实现

### 核心解决方案

通过在 `pkg/logger/logger.go` 中新增专门的 `*Wrapped` 系列方法：

```go
// 专门用于封装场景的日志方法，会额外跳过一层调用栈
func InfofWrapped(ctx context.Context, format string, args ...interface{}) {
    fields := []zap.Field{zap.String("msg", fmt.Sprintf(format, args...))}
    logger.WithOptions(zap.AddCallerSkip(1)).Info("", withTraceID(ctx, fields)...)
}
```

### 调用链分析

1. **用户代码**：`ctx.Logger.Infof("message")` ← 这是我们希望显示的位置
2. **ContextLogger**：`logger.InfofWrapped(l.ctx, format, args...)`
3. **pkg/logger**：`logger.WithOptions(zap.AddCallerSkip(1)).Info(...)` ← 跳过一层
4. **zap logger**：实际输出

通过 `AddCallerSkip(1)` 跳过 `ContextLogger` 这一层，直接定位到用户代码。

## 最佳实践

### 1. 日志级别选择

```go
// ✅ 好的做法
ctx.Logger.Info("用户登录成功")           // 关键业务事件
ctx.Logger.Debug("验证用户密码")          // 调试信息
ctx.Logger.Warn("用户密码即将过期")        // 需要注意的情况
ctx.Logger.Error("数据库连接失败", err)    // 错误情况

// ❌ 不好的做法
ctx.Logger.Info("循环第 %d 次", i)        // 过多的细节信息
ctx.Logger.Error("用户输入为空", nil)      // 不是真正的错误
```

### 2. 结构化日志

```go
// ✅ 好的做法
ctx.Logger.Infof("用户操作: 用户=%s, 操作=%s, 耗时=%dms", 
    userID, operation, duration)

// ✅ 更好的做法（如果需要复杂结构）
ctx.Logger.Infof("用户操作完成: %+v", map[string]interface{}{
    "user_id": userID,
    "operation": operation,
    "duration_ms": duration,
    "success": true,
})
```

### 3. 错误处理

```go
// ✅ 好的做法
if err != nil {
    ctx.Logger.Error("操作失败", err)
    ctx.Logger.Errorf("详细错误: 用户=%s, 操作=%s, 错误=%v", 
        userID, operation, err)
    return nil, err
}
```

### 4. 性能考虑

```go
// ✅ 好的做法 - 避免在循环中打印大量日志
ctx.Logger.Infof("开始处理 %d 个用户", len(users))
for i, user := range users {
    // 只在关键点或出错时记录
    if err := processUser(user); err != nil {
        ctx.Logger.Errorf("处理用户失败: %s, 错误: %v", user.ID, err)
    }
}
ctx.Logger.Info("用户处理完成")

// ❌ 避免的做法
for i, user := range users {
    ctx.Logger.Debugf("处理用户 %d: %s", i, user.ID) // 过多日志
}
```

## 与旧方式的对比

### 旧方式
```go
import "github.com/yunhanshu-net/pkg/logger"

func MyHandler(ctx *runner.Context, req *MyRequest) error {
    logger.Infof(ctx, "开始处理请求: %s", req.Username)
    // 需要每次传入 ctx，容易遗漏
    // 需要引入额外的包
    // 堆栈信息可能不准确
}
```

### 新方式
```go
func MyHandler(ctx *runner.Context, req *MyRequest) error {
    ctx.Logger.Infof("开始处理请求: %s", req.Username)
    // 直接使用，无需额外引入包
    // 自动包含 TraceID 和正确的堆栈信息
    // 堆栈信息完全准确
}
```

## 注意事项

1. **初始化**：`Context.Logger` 在 `NewContext()` 时自动初始化，无需手动设置
2. **TraceID**：会自动从 Context 中提取 TraceID 并添加到日志中
3. **性能**：基于高性能的 zap logger，适合生产环境使用
4. **堆栈信息**：通过专门的 `*Wrapped` 方法确保显示正确的调用位置
5. **线程安全**：zap logger 本身是线程安全的
6. **零配置**：开箱即用，无需任何额外配置

## 故障排查

如果遇到日志问题，可以检查：

1. **日志不显示**：确认 logger 已正确初始化
2. **堆栈信息错误**：现在已经完全解决，如有问题请检查是否有其他封装层
3. **TraceID 缺失**：确认 Context 中包含正确的 TraceID
4. **性能问题**：避免在高频循环中使用 Debug 级别日志

## 示例代码

完整的使用示例请参考：
- `context_logger_test.go` - 基础功能测试
- `context_logger_example.go` - 实际使用示例

## 总结

通过在 `pkg/logger` 层面解决调用深度问题，`Context.Logger` 现在提供了：

✅ **完美的堆栈信息** - 精确显示实际调用位置  
✅ **简洁的API** - 直接使用 `ctx.Logger.Infof()`  
✅ **自动TraceID** - 无需手动传递context  
✅ **高性能** - 基于zap的高效日志记录  
✅ **零配置** - 开箱即用

这是一个完美的日志封装解决方案！🎉 