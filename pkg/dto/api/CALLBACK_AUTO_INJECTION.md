# 回调自动注入功能

## 功能概述

新的回调自动注入功能允许系统在解析API参数时，自动将`FunctionInfo`中定义的字段级回调（`OnInputFuzzyMap`和`OnInputValidateMap`）注入到对应字段的回调配置中，无需在结构体标签中重复定义。

## 使用方法

### 1. 传统方式（仍然支持）

```go
// 只能通过标签定义回调
type UserReq struct {
    Username string `json:"username" runner:"code:username;name:用户名" callback:"OnInputValidate(trigger:blur)"`
}

// 使用传统接口
params, err := api.NewRequestParams(UserReq{}, "form")
```

### 2. 新的自动注入方式

```go
// 结构体标签中无需定义回调
type UserReq struct {
    Username string `json:"username" runner:"code:username;name:用户名"`
    Email    string `json:"email" runner:"code:email;name:邮箱"`
}

// 在FunctionInfo中定义字段级回调
var UserConfig = &runner.FunctionInfo{
    OnInputValidateMap: map[string]runner.OnInputValidate{
        "username": func(ctx *runner.Context, req *usercall.OnInputValidateReq) (*usercall.OnInputValidateResp, error) {
            // 验证用户名是否存在
            return &usercall.OnInputValidateResp{ErrorMsg: ""}, nil
        },
        "email": func(ctx *runner.Context, req *usercall.OnInputValidateReq) (*usercall.OnInputValidateResp, error) {
            // 验证邮箱格式和唯一性
            return &usercall.OnInputValidateResp{ErrorMsg: ""}, nil
        },
    },
    OnInputFuzzyMap: map[string]runner.OnInputFuzzy{
        "company_search": func(ctx *runner.Context, req *usercall.OnInputFuzzyReq) (*usercall.OnInputFuzzyResp, error) {
            // 公司名称模糊搜索
            return &usercall.OnInputFuzzyResp{}, nil
        },
    },
}

// 使用新接口，自动注入回调
params, err := api.NewRequestParamsWithFunctionInfo(UserReq{}, "form", UserConfig)
```

## 自动注入规则

1. **字段匹配**：通过字段的`code`（从`runner`标签或`json`标签获取）与`FunctionInfo`中的`Map`键进行匹配
2. **避免重复**：如果字段标签中已经定义了相同类型的回调，不会重复注入
3. **默认参数**：自动注入的回调使用默认参数：
   - `OnInputFuzzy`: `delay:300, min:2`
   - `OnInputValidate`: `trigger:blur`

## 在buildApiInfo中的应用

修改后的`buildApiInfo`方法会自动使用新接口：

```go
// function-go/runner/default.go
if config.Request != nil {
    // 使用支持FunctionInfo的版本，自动注入回调信息
    params, err := api.NewRequestParamsWithFunctionInfo(config.Request, config.RenderType, config)
    if err != nil {
        return nil, err
    }
    apiInfo.ParamsIn = params
}
```

## 优势

1. **减少重复**：不需要在结构体标签中重复定义已经在`FunctionInfo`中实现的回调
2. **自动同步**：回调实现和字段配置自动保持一致
3. **向后兼容**：完全兼容现有的标签定义方式
4. **灵活配置**：支持标签和`FunctionInfo`两种方式的混合使用

## 实际效果

使用新功能后，前端会自动获得完整的回调配置信息，包括：
- 标签中定义的回调
- 从`FunctionInfo`中自动注入的回调

这样前端就能正确地触发相应的交互功能，如模糊搜索、输入验证等。 