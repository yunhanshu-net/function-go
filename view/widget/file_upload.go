package widget

import (
	"strconv"

	"github.com/yunhanshu-net/pkg/x/tagx"
)

// FileUploadWidget 文件上传组件
type FileUploadWidget struct {
	// 接受的文件类型
	Accept string `json:"accept,omitempty"`
	// 是否支持多文件上传
	Multiple bool `json:"multiple,omitempty"`
	// 单个文件大小限制（单位：KB）
	MaxSize int64 `json:"max_size,omitempty"`
	// 文件数量限制
	Limit int `json:"limit,omitempty"`
	// 上传按钮文字
	Placeholder string `json:"placeholder,omitempty"`
	// 是否自动上传
	AutoUpload bool `json:"auto_upload,omitempty"`
	// 上传接口地址
	Action string `json:"action,omitempty"`
	// 列表展示方式：text/picture/picture-card
	ListType string `json:"list_type,omitempty"`
	// 是否支持拖拽上传
	Drag bool `json:"drag,omitempty"`
	// 上传按钮文字
	ButtonText string `json:"button_text,omitempty"`
	// 提示文字
	Tip string `json:"tip,omitempty"`
	// 是否禁用
	Disabled bool `json:"disabled,omitempty"`
}

// newFileUploadWidget 创建文件上传组件
func newFileUploadWidget(info *tagx.RunnerFieldInfo) (Widget, error) {
	file := &FileUploadWidget{}

	tag := info.Tags
	if tag == nil {
		return file, nil
	}

	// 设置接受的文件类型
	if accept, ok := tag["accept"]; ok {
		file.Accept = accept
	}

	// 设置是否支持多文件上传
	if multiple, ok := tag["multiple"]; ok {
		file.Multiple = multiple == "true"
	}

	// 设置文件大小限制
	if maxSize, ok := tag["max_size"]; ok {
		if size, err := strconv.ParseInt(maxSize, 10, 64); err == nil {
			file.MaxSize = size
		}
	}

	// 设置文件数量限制
	if limit, ok := tag["limit"]; ok {
		if num, err := strconv.Atoi(limit); err == nil {
			file.Limit = num
		}
	}

	// 设置占位符
	if placeholder, ok := tag["placeholder"]; ok {
		file.Placeholder = placeholder
	}

	// 设置是否自动上传
	if autoUpload, ok := tag["auto_upload"]; ok {
		file.AutoUpload = autoUpload == "true"
	}

	// 设置上传接口地址
	if action, ok := tag["action"]; ok {
		file.Action = action
	}

	// 设置列表类型
	if listType, ok := tag["list_type"]; ok {
		file.ListType = listType
	}

	// 设置是否支持拖拽上传
	if drag, ok := tag["drag"]; ok {
		file.Drag = drag == "true"
	}

	// 设置上传按钮文字
	if buttonText, ok := tag["button_text"]; ok {
		file.ButtonText = buttonText
	}

	// 设置提示文字
	if tip, ok := tag["tip"]; ok {
		file.Tip = tip
	}

	// 设置是否禁用
	if disabled, ok := tag["disabled"]; ok {
		file.Disabled = disabled == "true"
	}

	return file, nil
}

func (w *FileUploadWidget) GetValueType() string {
	return TypeFiles
}

func (w *FileUploadWidget) GetWidgetType() string {
	return WidgetFileUpload
}
