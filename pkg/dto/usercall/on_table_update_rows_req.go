package usercall

import (
	"fmt"
	"strconv"

	"github.com/yunhanshu-net/pkg/x/jsonx"
)

type OnTableUpdateRowsReq struct {
	Ids    []int                  `json:"ids"`
	Fields map[string]interface{} `json:"fields"` // 要更新的字段和值的映射
}

func (req *OnTableUpdateRowsReq) GetFieldsMap() map[string]interface{} {
	return req.Fields
}

// GetString 安全获取字段的string值
// 【框架规范】直接调用UpdateValue接口方法，简洁优雅
func (r *OnTableUpdateRowsReq) GetString(fieldName string) (string, bool, error) {
	value, exists := r.Fields[fieldName]
	if !exists {
		return "", false, nil
	}
	if value == nil {
		return "", false, fmt.Errorf("值为空")
	}

	switch val := value.(type) {
	case string:
		return val, true, nil
	case int:
		return strconv.Itoa(val), true, nil
	case int64:
		return strconv.FormatInt(val, 10), true, nil
	case int32:
		return strconv.FormatInt(int64(val), 10), true, nil
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64), true, nil
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32), true, nil
	case bool:
		return strconv.FormatBool(val), true, nil
	default:
		return "", false, fmt.Errorf("不支持的数据类型: %T, 值: %v", val, val)
	}
}

// GetBool
func (r *OnTableUpdateRowsReq) GetBool(fieldName string) (val bool, exist bool, err error) {

	value, exists := r.Fields[fieldName]
	if !exists {
		return false, false, nil
	}
	if value == nil {
		return false, false, nil
	}
	v, ok := value.(bool)
	if !ok {
		return false, false, nil
	}
	return v, true, nil
}

// GetInt 安全获取字段的int值
// 【框架规范】直接调用UpdateValue接口方法，简洁优雅
func (r *OnTableUpdateRowsReq) GetInt(fieldName string) (int, bool, error) {
	value, exists := r.Fields[fieldName]
	if !exists {
		return 0, false, nil
	}
	if value == nil {
		return 0, false, fmt.Errorf("值为空")
	}

	switch val := value.(type) {
	case int:
		return val, true, nil
	case int64:
		return int(val), true, nil
	case int32:
		return int(val), true, nil
	case float64:
		return int(val), true, nil
	case float32:
		return int(val), true, nil
	case string:
		// 尝试从字符串解析整数
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed, true, nil
		} else {
			return 0, false, fmt.Errorf("字符串解析失败: %s, 错误: %v", val, err)
		}
	default:
		return 0, false, fmt.Errorf("不支持的数据类型: %T, 值: %v", val, val)
	}
}

// GetFloat64 安全获取字段的float64值
// 【框架规范】直接调用UpdateValue接口方法，简洁优雅
func (r *OnTableUpdateRowsReq) GetFloat64(fieldName string) (float64, bool, error) {
	value, exists := r.Fields[fieldName]
	if !exists {
		return 0, false, nil
	}
	if value == nil {
		return 0, false, fmt.Errorf("值为空")
	}

	switch val := value.(type) {
	case float64:
		return val, true, nil
	case float32:
		return float64(val), true, nil
	case int:
		return float64(val), true, nil
	case int64:
		return float64(val), true, nil
	case int32:
		return float64(val), true, nil
	case string:
		// 尝试从字符串解析数值
		if parsed, err := strconv.ParseFloat(val, 64); err == nil {
			return parsed, true, nil
		} else {
			return 0, false, fmt.Errorf("字符串解析失败: %s, 错误: %v", val, err)
		}
	default:
		return 0, false, fmt.Errorf("不支持的数据类型: %T, 值: %v", val, val)
	}
}

func (r *OnTableUpdateRowsReq) DecodeBy(el interface{}) error {
	err := jsonx.Convert(r.Fields, el)
	if err != nil {
		return err
	}
	return nil
}
