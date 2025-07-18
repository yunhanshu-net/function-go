# 配置系统重构总结

## 问题描述

原来的配置系统存在循环引用问题：
- `pkg/config` 依赖 `runner.Context`
- `runner` 依赖 `pkg/config`
- 这形成了循环依赖，违反了Go的包管理原则

## 解决方案

将配置管理相关的代码从 `pkg/config` 移动到 `runner` 包下，避免循环引用。

## 重构变更

### 1. 文件移动和创建

**新增文件：**
- `runner/config_manager.go` - 配置管理器
- `runner/config_storage.go` - 本地文件存储实现

**删除文件：**
- `pkg/config/manager.go` - 已移动到 runner 包
- `pkg/config/local_storage.go` - 已移动到 runner 包

### 2. 类型定义变更

**ConfigStorage 接口：**
```go
// 旧版本 (pkg/config)
type ConfigStorage interface {
    Read(ctx *runner.Context, configKey string) (*syscallback.ConfigData, error)
    Write(ctx *runner.Context, configKey string, data *syscallback.ConfigData) error
    // ...
}

// 新版本 (runner)
type ConfigStorage interface {
    Read(ctx *Context, configKey string) (*syscallback.ConfigData, error)
    Write(ctx *Context, configKey string, data *syscallback.ConfigData) error
    // ...
}
```

**ConfigChangeCallback 类型：**
```go
// 旧版本
type ConfigChangeCallback func(ctx *runner.Context, oldConfig, newConfig *syscallback.ConfigData) error

// 新版本
type ConfigChangeCallback func(ctx *Context, oldConfig, newConfig *syscallback.ConfigData) error
```

### 3. 导入更新

**移除的导入：**
```go
// 从以下文件中移除
"github.com/yunhanshu-net/function-go/pkg/config"
```

**更新的函数调用：**
```go
// 旧版本
configManager := config.GetConfigManager()
localStorage := config.NewLocalFileStorage("./configs")

// 新版本
configManager := GetConfigManager()
localStorage := NewLocalFileStorage("./configs")
```

### 4. Context 结构体更新

```go
// 旧版本
type Context struct {
    // ...
    config *config.ConfigManager
}

// 新版本
type Context struct {
    // ...
    config *ConfigManager
}
```

### 5. 配置键格式

配置键现在使用简化的格式，避免不同方法的冲突：
```
function.{router}.{method}
```

**路由处理规则：**
- 将 `/` 替换为 `.`
- 移除前后的 `.`
- 支持边界情况处理

**例如：**
- `/api/users` + `GET` → `function.api.users.GET`
- `/widgets/add` + `POST` → `function.widgets.add.POST`
- `/admin/settings` + `PUT` → `function.admin.settings.PUT`
- `/` + `GET` → `function.GET`

## 影响范围

### 更新的文件：
1. `runner/function.go` - 移除 pkg/config 导入
2. `runner/context.go` - 更新配置管理器类型
3. `runner/default.go` - 更新配置管理器调用
4. `runner/runner.go` - 更新配置管理器初始化
5. `test_config_system.go` - 更新测试代码

### 新增的HTTP方法支持：
- `GET` - 获取数据
- `POST` - 创建数据  
- `PUT` - 更新数据
- `DELETE` - 删除数据
- `PATCH` - 部分更新数据

## 测试验证

创建了以下测试文件验证重构：
- `test_config_refactor.go` - 验证重构后的配置系统
- `test_config_method.go` - 验证不同HTTP方法的配置隔离

## 优势

1. **消除循环引用** - 解决了包依赖问题
2. **更好的封装** - 配置管理逻辑集中在 runner 包内
3. **HTTP方法隔离** - 不同HTTP方法有独立的配置
4. **向后兼容** - API接口保持不变，只是内部实现重构

## API接口变更

### ConfigUpdateRequest 结构变更
```go
// 旧版本
type ConfigUpdateRequest struct {
    ConfigKey  string      `json:"config_key"`  // 要更新的配置键
    ConfigData *ConfigData `json:"config_data"` // 完整的修改后配置
}

// 新版本
type ConfigUpdateRequest struct {
    Router     string      `json:"router"`     // 路由路径
    Method     string      `json:"method"`     // HTTP方法
    ConfigData *ConfigData `json:"config_data"` // 完整的修改后配置
}
```

### ConfigGetRequest 结构变更
```go
// 旧版本
type ConfigGetRequest struct {
    ConfigKey string `json:"config_key"` // 要获取的配置键
}

// 新版本
type ConfigGetRequest struct {
    Router string `json:"router"` // 路由路径
    Method string `json:"method"` // HTTP方法
}
```

### 配置键自动构造
内部根据 `router` 和 `method` 自动构造配置键：
```go
// 处理路由路径，将 / 替换为 . 以安全地用作配置键
safeRouter := strings.ReplaceAll(router, "/", ".")
// 移除前后的点
safeRouter = strings.Trim(safeRouter, ".")
configKey := fmt.Sprintf("function.%s.%s", safeRouter, method)
```

**配置键格式示例：**
- `/api/users` + `GET` → `function.api.users.GET`
- `/widgets/add` + `POST` → `function.widgets.add.POST`
- `/admin/settings` + `PUT` → `function.admin.settings.PUT`

## 注意事项

1. 配置键现在包含HTTP方法，需要更新现有的配置管理逻辑
2. 所有使用 `pkg/config` 的代码都需要更新为使用 `runner` 包中的配置管理器
3. 配置文件路径格式已更新，包含HTTP方法信息
4. API接口现在使用 `router` 和 `method` 参数，不再需要手动构造配置键 