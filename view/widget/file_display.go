package widget

import (
	"strconv"

	"github.com/yunhanshu-net/pkg/x/tagx"
)

// FileDisplayWidget 文件展示组件
type FileDisplayWidget struct {
	// 展示模式：list/card
	DisplayMode string `json:"display_mode,omitempty"`
	// 是否支持预览
	Preview bool `json:"preview,omitempty"`
	// 是否支持下载
	Download bool `json:"download,omitempty"`
	// 最大预览数量
	MaxPreview int `json:"max_preview,omitempty"`
}

// newFileDisplayWidget 创建文件展示组件
func newFileDisplayWidget(info *tagx.RunnerFieldInfo) (Widget, error) {
	file := &FileDisplayWidget{}

	tag := info.Tags
	if tag == nil {
		return file, nil
	}

	// 设置展示模式
	if displayMode, ok := tag["display_mode"]; ok {
		file.DisplayMode = displayMode
	}

	// 设置是否支持预览
	if preview, ok := tag["preview"]; ok {
		file.Preview = preview == "true"
	}

	// 设置是否支持下载
	if download, ok := tag["download"]; ok {
		file.Download = download == "true"
	}

	// 设置最大预览数量
	if maxPreview, ok := tag["max_preview"]; ok {
		if num, err := strconv.Atoi(maxPreview); err == nil {
			file.MaxPreview = num
		}
	}

	return file, nil
}

func (w *FileDisplayWidget) GetValueType() string {
	return TypeFiles
}

func (w *FileDisplayWidget) GetWidgetType() string {
	return WidgetFileDisplay
}
