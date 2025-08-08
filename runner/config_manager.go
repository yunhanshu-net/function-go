package runner

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/yunhanshu-net/function-go/pkg/dto/usercall"
	"github.com/yunhanshu-net/pkg/logger"
)

// ConfigStorage 配置存储接口
type ConfigStorage interface {
	// Read 读取配置
	Read(ctx *Context, configKey string) (*usercall.ConfigData, error)

	// Write 写入配置
	Write(ctx *Context, configKey string, data *usercall.ConfigData) error

	// Exists 检查配置是否存在
	Exists(ctx *Context, configKey string) (bool, error)

	// Delete 删除配置
	Delete(ctx *Context, configKey string) error
}

type AutoUpdateConfig struct {
	ConfigStruct       interface{}                `json:"config_struct"` // 配置结构体值（用于类型注册）
	BeforeConfigChange BeforeConfigChangeCallback `json:"-"`             // 配置变更前回调
}

// ConfigChangeCallback 配置变更回调函数类型
type ConfigChangeCallback func(ctx *Context, oldConfig, newConfig *usercall.ConfigData) error

// BeforeConfigChangeCallback oldConfig和newConfig都是AutoUpdateConfig.ConfigStruct注册的结构体的值类型（值类型）
type BeforeConfigChangeCallback func(ctx *Context, oldConfig, newConfig interface{}) error

// ConfigManager 配置管理器
type ConfigManager struct {
	cache         map[string]*usercall.ConfigData
	storage       ConfigStorage
	mutex         sync.RWMutex
	callbacks     map[string]BeforeConfigChangeCallback // 配置键到回调函数的映射
	configStructs map[string]reflect.Type               // 配置键到结构体类型的映射
}

var (
	globalConfigManager *ConfigManager
	configManagerOnce   sync.Once
)

// GetConfigManager 获取全局配置管理器单例
func GetConfigManager() *ConfigManager {
	configManagerOnce.Do(func() {
		globalConfigManager = &ConfigManager{
			cache:         make(map[string]*usercall.ConfigData),
			callbacks:     make(map[string]BeforeConfigChangeCallback),
			configStructs: make(map[string]reflect.Type),
		}
	})
	return globalConfigManager
}

// SetStorage 设置存储方式
func (cm *ConfigManager) SetStorage(storage ConfigStorage) {
	cm.storage = storage
}

// RegisterCallback 注册配置变更回调
func (cm *ConfigManager) RegisterCallback(configKey string, callback BeforeConfigChangeCallback) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.callbacks[configKey] = callback
}

// RegisterConfigStruct 注册配置结构体
func (cm *ConfigManager) RegisterConfigStruct(configKey string, configStruct interface{}) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// 获取结构体类型
	var structType reflect.Type
	if reflect.TypeOf(configStruct).Kind() == reflect.Ptr {
		// 如果是指针，获取指针指向的类型
		structType = reflect.TypeOf(configStruct).Elem()
	} else {
		// 如果不是指针，直接使用类型
		structType = reflect.TypeOf(configStruct)
	}

	// 存储结构体的类型，用于后续创建实例
	cm.configStructs[configKey] = structType
}

// GetByKey 根据配置键获取配置
func (cm *ConfigManager) GetByKey(ctx *Context, configKey string) *usercall.ConfigData {
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
func (cm *ConfigManager) loadConfig(ctx *Context, configKey string) *usercall.ConfigData {
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
	var configCopy *usercall.ConfigData
	if data != nil {
		configCopy = &usercall.ConfigData{
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
func (cm *ConfigManager) UpdateConfig(ctx *Context, configKey string, newConfig *usercall.ConfigData) error {
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
	var configCopy *usercall.ConfigData
	if newConfig != nil {
		configCopy = &usercall.ConfigData{
			Type: newConfig.Type,
			Data: newConfig.Data,
		}
	}

	// 验证通过，更新配置
	cm.mutex.Lock()
	cm.cache[configKey] = configCopy
	cm.mutex.Unlock()

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

// getBeforeConfigChangeCallback 获取配置变更前回调
func (cm *ConfigManager) getBeforeConfigChangeCallback(configKey string) BeforeConfigChangeCallback {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.callbacks[configKey]
}

// ClearCache 清空缓存
func (cm *ConfigManager) ClearCache() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.cache = make(map[string]*usercall.ConfigData)
}

// GetCacheSize 获取缓存大小
func (cm *ConfigManager) GetCacheSize() int {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return len(cm.cache)
}

// GetConfigStruct 获取配置结构体值
func (cm *ConfigManager) GetConfigStruct(ctx *Context, configKey string) interface{} {
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
		return configData.Data
	}

	// 如果Data已经是结构体类型，直接返回
	if reflect.TypeOf(configData.Data) == configStructType {
		return configData.Data
	}

	// 如果Data是字符串，需要解析
	if dataStr, ok := configData.Data.(string); ok {
		// 创建结构体实例
		instance := reflect.New(configStructType).Interface()
		if err := json.Unmarshal([]byte(dataStr), instance); err != nil {
			return nil
		}
		// 返回结构体的值（不是指针）
		result := reflect.ValueOf(instance).Elem().Interface()
		return result
	}

	// 如果Data是map或其他类型，尝试直接转换
	instance := reflect.New(configStructType).Interface()
	// 使用mapstructure进行转换，配置使用JSON标签
	config := &mapstructure.DecoderConfig{
		TagName: "json",
		Result:  instance,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil
	}
	if err := decoder.Decode(configData.Data); err != nil {
		// 如果mapstructure失败，尝试JSON序列化再反序列化
		if dataBytes, err := json.Marshal(configData.Data); err == nil {
			if err := json.Unmarshal(dataBytes, instance); err != nil {
				return nil
			}
		} else {
			return nil
		}
	}
	// 返回结构体的值（不是指针）
	result := reflect.ValueOf(instance).Elem().Interface()
	return result
}
