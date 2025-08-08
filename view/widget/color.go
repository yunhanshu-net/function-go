package widget

import (
	"strings"

	"github.com/yunhanshu-net/pkg/x/tagx"
)

// ColorWidget 颜色选择器组件
type ColorWidget struct {
	// 颜色格式配置
	Format string `json:"format"` // 颜色格式：hex/rgb/rgba/hsl/hsla
	// 默认值和预设
	Predefine []string `json:"predefine"` // 预定义颜色列表
	// 透明度控制
	ShowAlpha bool `json:"show_alpha"` // 是否显示透明度选择器

	// 交互配置
}

// newColorWidget 创建颜色选择器组件
func newColorWidget(info *tagx.RunnerFieldInfo) (Widget, error) {
	color := &ColorWidget{
		Format: "hex", // 默认十六进制格式
	}

	if info.Tags == nil {
		return color, nil
	}

	tag := info.Tags

	// 设置颜色格式
	if format, ok := tag["format"]; ok && format != "" {
		// 验证格式是否有效
		validFormats := map[string]bool{
			"hex": true, "rgb": true, "rgba": true, "hsl": true, "hsla": true,
		}
		if validFormats[format] {
			color.Format = format
		}
	}

	// 解析预定义颜色
	if predefine, ok := tag["predefine"]; ok && predefine != "" {
		// 按逗号分割预定义颜色
		parts := strings.Split(predefine, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				color.Predefine = append(color.Predefine, part)
			}
		}
	}

	// 解析透明度控制
	if showAlpha, ok := tag["show_alpha"]; ok && showAlpha != "" {
		color.ShowAlpha = showAlpha == "true"
	}

	return color, nil
}

func (w *ColorWidget) GetWidgetType() string {
	return WidgetColor
}
