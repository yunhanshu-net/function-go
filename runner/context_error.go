package runner

// 极简的结构化错误与构造器（仅在本文件实现，不改动其他文件）

// AppError 为平台统一的结构化错误
type AppError struct {
	Code        string            `json:"code"`                   // 机器可识别错误码，例如 E_VALIDATION/E_DEP_NOT_READY
	Message     string            `json:"message"`                // 面向用户的短消息（中文）
	Hint        string            `json:"hint"`                   // 修复建议
	Detail      string            `json:"detail"`                 // 诊断信息（可选，前端可折叠）
	TraceID     string            `json:"trace_id"`               // 追踪ID
	Retryable   bool              `json:"retryable"`              // 是否可重试
	FieldErrors map[string]string `json:"field_errors,omitempty"` // 表单字段级错误
}

// Error 满足 error 接口
func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	if e.Message != "" {
		return e.Message
	}
	return e.Code
}

// ErrorBuilder 结构化错误构造器（链式）
type ErrorBuilder struct {
	ctx *Context
	e   *AppError
}

// Err 创建一个构造器；建议在函数内调用：return runner.Err(ctx).Code("...").Msg("...").Build()
func Err(ctx *Context) *ErrorBuilder {
	return &ErrorBuilder{ctx: ctx, e: &AppError{}}
}

// 在 Context 上提供便捷入口：ctx.Err(...)
// - 无参：等价于 Err(ctx)
// - error：作为基础错误包装（默认 Code=E_INTERNAL，Detail=err.Error()）
// - string：作为初始 Message
// - *AppError：拷贝其字段继续编辑
func (c *Context) Error(base ...interface{}) *ErrorBuilder {
	b := Err(c)
	if len(base) == 0 || base[0] == nil {
		return b
	}
	switch v := base[0].(type) {
	case *AppError:
		if v != nil {
			cp := *v
			b.e = &cp
		}
	case error:
		if v != nil {
			if b.e.Code == "" {
				b.e.Code = "E_INTERNAL"
			}
			b.e.Detail = v.Error()
		}
	case string:
		b.e.Message = v
	}
	return b
}

func (b *ErrorBuilder) Code(code string) *ErrorBuilder     { b.e.Code = code; return b }
func (b *ErrorBuilder) Msg(message string) *ErrorBuilder   { b.e.Message = message; return b }
func (b *ErrorBuilder) Message(s string) *ErrorBuilder     { return b.Msg(s) }
func (b *ErrorBuilder) Hint(hint string) *ErrorBuilder     { b.e.Hint = hint; return b }
func (b *ErrorBuilder) Suggest(s string) *ErrorBuilder     { return b.Hint(s) }
func (b *ErrorBuilder) Detail(detail string) *ErrorBuilder { b.e.Detail = detail; return b }
func (b *ErrorBuilder) Because(s string) *ErrorBuilder     { return b.Detail(s) }
func (b *ErrorBuilder) Retryable(v bool) *ErrorBuilder     { b.e.Retryable = v; return b }
func (b *ErrorBuilder) Retry() *ErrorBuilder               { b.e.Retryable = true; return b }
func (b *ErrorBuilder) Field(name, msg string) *ErrorBuilder {
	if b.e.FieldErrors == nil {
		b.e.FieldErrors = map[string]string{}
	}
	b.e.FieldErrors[name] = msg
	return b
}

// Build 结束构造，自动注入 TraceID
func (b *ErrorBuilder) Build() error {
	if b.ctx != nil && b.e.TraceID == "" {
		b.e.TraceID = b.ctx.getTraceId()
	}
	return b.e
}

// 便捷构造函数（最小可用）

// ValidationError 返回一个表单校验错误（含字段错误）
func ValidationError(ctx *Context, fields map[string]string) error {
	eb := Err(ctx).Code("E_VALIDATION").Msg("参数校验失败")
	for k, v := range fields {
		eb.Field(k, v)
	}
	return eb.Build()
}

// DepNotReady 返回依赖未就绪错误（例如 tesseract/poppler 未安装）
func DepNotReady(ctx *Context, hint string) error {
	return Err(ctx).Code("E_DEP_NOT_READY").Msg("依赖未就绪").Hint(hint).Retryable(true).Build()
}

// ExecFailed 外部命令失败的错误（detail 可放 stderr 摘要）
func ExecFailed(ctx *Context, msg, detail string) error {
	return Err(ctx).Code("E_EXEC_FAILED").Msg(msg).Detail(detail).Retryable(false).Build()
}

// 常见错误场景的一键封装（减少样板代码）
func (b *ErrorBuilder) Validation() *ErrorBuilder {
	if b.e.Code == "" {
		b.e.Code = "E_VALIDATION"
	}
	if b.e.Message == "" {
		b.e.Message = "参数校验失败"
	}
	return b
}

func (b *ErrorBuilder) FieldRequired(name string) *ErrorBuilder {
	return b.Field(name, "必填")
}

func (b *ErrorBuilder) DepNotReady(hint string) *ErrorBuilder {
	if b.e.Code == "" {
		b.e.Code = "E_DEP_NOT_READY"
	}
	if b.e.Message == "" {
		b.e.Message = "依赖未就绪"
	}
	b.e.Hint = hint
	b.e.Retryable = true
	return b
}

func (b *ErrorBuilder) ExecFail(detail string) *ErrorBuilder {
	if b.e.Code == "" {
		b.e.Code = "E_EXEC_FAILED"
	}
	if b.e.Message == "" {
		b.e.Message = "外部命令执行失败"
	}
	b.e.Detail = detail
	return b
}
