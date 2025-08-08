package runner

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yunhanshu-net/function-go/pkg/dto/usercall"
	"github.com/yunhanshu-net/pkg/logger"
)

// LocalFileStorage 本地文件存储实现
type LocalFileStorage struct {
	basePath string // 配置文件基础路径
}

// NewLocalFileStorage 创建本地文件存储
func NewLocalFileStorage(basePath string) *LocalFileStorage {
	if basePath == "" {
		basePath = "./configs" // 默认路径
	}

	// 确保目录存在
	if err := os.MkdirAll(basePath, 0755); err != nil {
		panic(fmt.Errorf("创建配置目录失败: %w", err))
	}

	return &LocalFileStorage{
		basePath: basePath,
	}
}

// getConfigFilePath 获取配置文件路径
func (lfs *LocalFileStorage) getConfigFilePath(configKey string) string {
	// 将配置键转换为文件路径
	// 例如: function.cmp.config_demo.POST -> function.cmp.config_demo.POST.json
	// 直接使用配置键作为文件名，避免路径分隔符问题
	safeKey := configKey + ".json"
	return filepath.Join(lfs.basePath, safeKey)
}

// Read 读取配置
func (lfs *LocalFileStorage) Read(ctx *Context, configKey string) (*usercall.ConfigData, error) {
	filePath := lfs.getConfigFilePath(configKey)

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Debugf(ctx, "配置文件不存在: %s", filePath)
			return nil, nil
		}
		return nil, fmt.Errorf("读取配置文件失败 %s: %w", filePath, err)
	}

	// 解析配置文件
	var configData usercall.ConfigData
	if err := json.Unmarshal(data, &configData); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &configData, nil
}

// Write 写入配置
func (lfs *LocalFileStorage) Write(ctx *Context, configKey string, configData *usercall.ConfigData) error {
	filePath := lfs.getConfigFilePath(configKey)

	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}

	// 序列化配置数据
	data, err := json.Marshal(configData)
	if err != nil {
		return fmt.Errorf("序列化配置数据失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败 %s: %w", filePath, err)
	}

	logger.Debugf(ctx, "配置文件已保存: %s data：%s", filePath, string(data))
	return nil
}

// Exists 检查配置是否存在
func (lfs *LocalFileStorage) Exists(ctx *Context, configKey string) (bool, error) {
	filePath := lfs.getConfigFilePath(configKey)

	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("检查配置文件失败 %s: %w", filePath, err)
	}

	return true, nil
}

// Delete 删除配置
func (lfs *LocalFileStorage) Delete(ctx *Context, configKey string) error {
	filePath := lfs.getConfigFilePath(configKey)

	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			logger.Debugf(ctx, "配置文件不存在，无需删除: %s", filePath)
			return nil
		}
		return fmt.Errorf("删除配置文件失败 %s: %w", filePath, err)
	}

	logger.Debugf(ctx, "配置文件已删除: %s", filePath)
	return nil
}
