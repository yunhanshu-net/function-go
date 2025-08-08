package widget

import (
	"strings"

	"github.com/yunhanshu-net/pkg/x/tagx"
)

// TagWidget 标签组件
type TagWidget struct {
	Separator string `json:"separator,omitempty"` // 分隔符，默认为逗号
	Color     string `json:"color,omitempty"`     // 标签颜色，可选：auto, primary, success, warning, danger, info
	MaxTags   int    `json:"max_tags,omitempty"`  // 最大标签数量
	Closable  bool   `json:"closable,omitempty"`  // 是否可关闭
	Editable  bool   `json:"editable,omitempty"`  // 是否可编辑
	Size      string `json:"size,omitempty"`      // 标签大小：small, default, large
}

// newTagWidget 创建标签组件
func newTagWidget(info *tagx.RunnerFieldInfo) (Widget, error) {
	tagWidget := &TagWidget{
		Separator: ",",
		Color:     "auto",
		MaxTags:   10,
		Closable:  true,
		Editable:  true,
		Size:      "default",
	}

	if info.Tags == nil {
		return tagWidget, nil
	}

	tag := info.Tags

	// 设置分隔符
	if separator, ok := tag["separator"]; ok && separator != "" {
		tagWidget.Separator = strings.TrimSpace(separator)
	}

	// 设置颜色
	if color, ok := tag["color"]; ok && color != "" {
		tagWidget.Color = strings.TrimSpace(color)
	}

	// 设置最大标签数量
	if maxTags, ok := tag["max_tags"]; ok && maxTags != "" {
		if val := parseInt(maxTags); val > 0 {
			tagWidget.MaxTags = val
		}
	}

	// 设置是否可关闭
	if closable, ok := tag["closable"]; ok {
		tagWidget.Closable = closable == "true"
	}

	// 设置是否可编辑
	if editable, ok := tag["editable"]; ok {
		tagWidget.Editable = editable == "true"
	}

	// 设置标签大小
	if size, ok := tag["size"]; ok && size != "" {
		tagWidget.Size = strings.TrimSpace(size)
	}

	return tagWidget, nil
}

func (w *TagWidget) GetWidgetType() string {
	return WidgetTag
}

// GetDefaultConfig 返回默认配置
func (t *TagWidget) GetDefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"separator": ",",
		"color":     "auto",
		"max_tags":  10,
		"closable":  true,
		"editable":  true,
		"size":      "default",
	}
}

// Validate 验证配置
func (t *TagWidget) Validate() error {
	return nil
}

// BuildConfig 构建配置
func (t *TagWidget) BuildConfig(params map[string]string) {

	if separator, ok := params["separator"]; ok {
		t.Separator = separator
	}

	if color, ok := params["color"]; ok {
		t.Color = color
	}

	if maxTags, ok := params["max_tags"]; ok {
		if val := parseInt(maxTags); val > 0 {
			t.MaxTags = val
		}
	}

	if closable, ok := params["closable"]; ok {
		t.Closable = parseBool(closable)
	}

	if editable, ok := params["editable"]; ok {
		t.Editable = parseBool(editable)
	}

	if size, ok := params["size"]; ok {
		t.Size = size
	}
}

func parseBool(str string) bool {
	return str == "true"
}
