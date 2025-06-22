package response

import (
	"encoding/json"
	"reflect"

	"github.com/pkg/errors"
)

const (
	RenderTypeForm    = "form"
	RenderTypeJSON    = "json"
	RenderTypeTable   = "table"
	RenderTypeFiles   = "files"
	RenderTypeEcharts = "echarts"
)

func build(resp *RunFunctionResp, data interface{}, renderType string) error {
	if resp == nil {
		return errors.New("resp is nil")
	}

	// 检查数据中是否包含需要上传的文件
	processedData, err := processFileUploads(data)
	if err != nil {
		return errors.Wrap(err, "处理文件上传失败")
	}

	if resp.Data != nil {
		resp.Multiple = true
		resp.DataList = append(resp.DataList, resp.Data)
		resp.DataList = append(resp.DataList, processedData)
		resp.RenderType = resp.RenderType + "," + renderType
	}
	resp.RenderType = renderType
	resp.Data = processedData
	resp.Msg = "ok"
	return nil
}

// processFileUploads 处理数据中的文件上传
func processFileUploads(data interface{}) (interface{}, error) {
	if data == nil {
		return data, nil
	}

	// 使用反射检查数据结构
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return data, nil
	}

	// 遍历结构体字段，查找实现了json.Marshaler的字段
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanInterface() {
			continue
		}

		fieldInterface := field.Interface()

		// 检查是否实现了json.Marshaler接口
		if marshaler, ok := fieldInterface.(json.Marshaler); ok {
			// 触发JSON序列化，这会调用MarshalJSON方法
			_, err := marshaler.MarshalJSON()
			if err != nil {
				return nil, errors.Wrapf(err, "序列化字段失败")
			}
		}
	}

	return data, nil
}
