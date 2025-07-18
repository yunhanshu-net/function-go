package runner

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"

	"github.com/yunhanshu-net/function-go/pkg/dto/syscallback"
	"github.com/yunhanshu-net/pkg/logger"
)

// ConfigStorage 配置存储接口
type ConfigStorage interface {
	// 读取配置
	Read(ctx *Context, configKey string) (*syscallback.ConfigData, error)

	// 写入配置
	Write(ctx *Context, configKey string, data *syscallback.ConfigData) error

	// 检查配置是否存在
	Exists(ctx *Context, configKey string) (bool, error)

	// 删除配置
	Delete(ctx *Context, configKey string) error
}

type AutoUpdateConfig struct {
	ConfigStruct       interface{}          `json:"config_struct"` // 配置结构体指针
	BeforeConfigChange ConfigChangeCallback `json:"-"`             // 配置变更前回调
}

// ConfigChangeCallback 配置变更回调函数类型
type ConfigChangeCallback func(ctx *Context, oldConfig, newConfig *syscallback.ConfigData) error

// ConfigManager 配置管理器
type ConfigManager struct {
	cache         map[string]*syscallback.ConfigData
	storage       ConfigStorage
	mutex         sync.RWMutex
	callbacks     map[string]ConfigChangeCallback // 配置键到回调函数的映射
	configStructs map[string]interface{}          // 配置键到解析后的结构体指针的映射
}

var (
	globalConfigManager *ConfigManager
	configManagerOnce   sync.Once
)

// GetConfigManager 获取全局配置管理器单例
func GetConfigManager() *ConfigManager {
	configManagerOnce.Do(func() {
		globalConfigManager = &ConfigManager{
			cache:         make(map[string]*syscallback.ConfigData),
			callbacks:     make(map[string]ConfigChangeCallback),
			configStructs: make(map[string]interface{}),
		}
	})
	return globalConfigManager
}

// SetStorage 设置存储方式
func (cm *ConfigManager) SetStorage(storage ConfigStorage) {
	cm.storage = storage
}

// RegisterCallback 注册配置变更回调
func (cm *ConfigManager) RegisterCallback(configKey string, callback ConfigChangeCallback) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.callbacks[configKey] = callback
}

// RegisterConfigStruct 注册配置结构体
func (cm *ConfigManager) RegisterConfigStruct(configKey string, configStruct interface{}) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	// 存储结构体的类型，用于后续创建实例
	cm.configStructs[configKey] = reflect.TypeOf(configStruct)
}

// GetByKey 根据配置键获取配置
func (cm *ConfigManager) GetByKey(ctx *Context, configKey string) *syscallback.ConfigData {
	cm.mutex.RLock()
	if config, exists := cm.cache[configKey]; exists {
		cm.mutex.RUnlock()
		return config
	}
	cm.mutex.RUnlock()

	// 缓存未命中，从存储加载
	return cm.loadConfig(ctx, configKey)
}

// loadConfig 从存储加载配置
func (cm *ConfigManager) loadConfig(ctx *Context, configKey string) *syscallback.ConfigData {
	if cm.storage == nil {
		logger.Warnf(ctx, "配置存储未设置，无法加载配置: %s", configKey)
		return nil
	}

	data, err := cm.storage.Read(ctx, configKey)
	if err != nil {
		logger.Warnf(ctx, "加载配置失败 %s: %v", configKey, err)
		return nil
	}

	// 深拷贝配置数据以确保安全
	var configCopy *syscallback.ConfigData
	if data != nil {
		configCopy = &syscallback.ConfigData{
			Type: data.Type,
			Data: data.Data,
		}
	}

	// 缓存配置
	cm.mutex.Lock()
	cm.cache[configKey] = configCopy
	cm.mutex.Unlock()

	return configCopy
}

// UpdateConfig 更新配置
func (cm *ConfigManager) UpdateConfig(ctx *Context, configKey string, newConfig *syscallback.ConfigData) error {
	cm.mutex.RLock()
	oldConfig := cm.cache[configKey]
	cm.mutex.RUnlock()

	// 触发 BeforeConfigChange 回调
	if callback := cm.getBeforeConfigChangeCallback(configKey); callback != nil {
		if err := callback(ctx, oldConfig, newConfig); err != nil {
			return fmt.Errorf("配置变更验证失败: %w", err)
		}
	}

	// 深拷贝配置数据以确保安全
	var configCopy *syscallback.ConfigData
	if newConfig != nil {
		configCopy = &syscallback.ConfigData{
			Type: newConfig.Type,
			Data: newConfig.Data,
		}
	}

	// 验证通过，更新配置
	cm.mutex.Lock()
	cm.cache[configKey] = configCopy
	cm.mutex.Unlock()

	// 同时更新缓存的结构体指针
	cm.updateCachedStruct(ctx, configKey, configCopy)

	// 同时更新到存储
	if cm.storage != nil {
		if err := cm.storage.Write(ctx, configKey, configCopy); err != nil {
			logger.Errorf(ctx, "保存配置到存储失败 %s: %v", configKey, err)
			// 不返回错误，因为内存缓存已经更新
		}
	}

	logger.Infof(ctx, "配置 %s 更新成功", configKey)
	return nil
}

// updateCachedStruct 更新缓存的结构体指针
func (cm *ConfigManager) updateCachedStruct(ctx *Context, configKey string, configData *syscallback.ConfigData) {
	cm.mutex.RLock()
	cachedStruct, exists := cm.configStructs[configKey]
	cm.mutex.RUnlock()

	if !exists {
		return
	}

	// 检查是否是已解析的结构体指针
	if reflect.TypeOf(cachedStruct).Kind() != reflect.Ptr {
		return // 还没有解析过，等待下次获取时解析
	}

	// 更新结构体指针的值
	switch configData.Type {
	case "json":
		if err := json.Unmarshal([]byte(configData.Data), cachedStruct); err != nil {
			logger.Warnf(ctx, "更新缓存结构体失败: %v", err)
			// 解析失败时清除缓存，下次重新解析
			cm.mutex.Lock()
			delete(cm.configStructs, configKey)
			cm.mutex.Unlock()
		}
	}
}

// getBeforeConfigChangeCallback 获取配置变更前回调
func (cm *ConfigManager) getBeforeConfigChangeCallback(configKey string) ConfigChangeCallback {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.callbacks[configKey]
}

// ClearCache 清空缓存
func (cm *ConfigManager) ClearCache() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.cache = make(map[string]*syscallback.ConfigData)
}

// GetCacheSize 获取缓存大小
func (cm *ConfigManager) GetCacheSize() int {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return len(cm.cache)
}

// GetConfigStruct 获取配置结构体指针
func (cm *ConfigManager) GetConfigStruct(ctx *Context, configKey string) interface{} {
	// 先检查是否已有缓存的解析后结构体
	cm.mutex.RLock()
	if cachedStruct, exists := cm.configStructs[configKey]; exists {
		// 检查是否是已解析的结构体指针（不是reflect.Type）
		if reflect.TypeOf(cachedStruct).Kind() != reflect.Ptr {
			// 这是类型，需要解析
			cm.mutex.RUnlock()
		} else {
			// 这是已解析的指针，直接返回
			cm.mutex.RUnlock()
			return cachedStruct
		}
	} else {
		cm.mutex.RUnlock()
		return nil
	}

	// 获取配置数据
	configData := cm.GetByKey(ctx, configKey)
	if configData == nil {
		return nil
	}

	// 从注册的结构体中获取对应的类型
	cm.mutex.RLock()
	configStructType, exists := cm.configStructs[configKey]
	cm.mutex.RUnlock()

	if !exists {
		// 如果没有注册的结构体，返回原始数据
		return configData
	}

	// 解析配置数据为结构体
	switch configData.Type {
	case "json":
		// 创建结构体实例
		instance := reflect.New(configStructType.(reflect.Type)).Interface()
		if err := json.Unmarshal([]byte(configData.Data), instance); err != nil {
			logger.Warnf(ctx, "解析配置失败: %v", err)
			return nil
		}
		
		// 缓存解析后的结构体指针
		cm.mutex.Lock()
		cm.configStructs[configKey] = instance
		cm.mutex.Unlock()
		
		return instance
	default:
		return configData
	}
}
