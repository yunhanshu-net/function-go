package widget

import (
	"strings"
)

// parseOptionsWithEscape 解析支持转义的选项字符串
// 支持 \, 转义逗号，\\ 转义反斜杠
func parseOptionsWithEscape(optionsStr string) []string {
	if optionsStr == "" {
		return nil
	}

	var result []string
	var current strings.Builder
	escaped := false

	for i := 0; i < len(optionsStr); i++ {
		char := optionsStr[i]

		if escaped {
			// 处理转义字符
			switch char {
			case '\\':
				current.WriteByte('\\')
			case ',':
				current.WriteByte(',')
			default:
				// 无效的转义序列，保持原样
				current.WriteByte('\\')
				current.WriteByte(char)
			}
			escaped = false
		} else if char == '\\' {
			// 开始转义
			escaped = true
		} else if char == ',' {
			// 分隔符
			option := strings.TrimSpace(current.String())
			if option != "" {
				result = append(result, option)
			}
			current.Reset()
		} else {
			// 普通字符
			current.WriteByte(char)
		}
	}

	// 处理最后一个选项
	if escaped {
		// 如果以反斜杠结尾，添加反斜杠
		current.WriteByte('\\')
	}
	option := strings.TrimSpace(current.String())
	if option != "" {
		result = append(result, option)
	}

	return result
}

// parseInt 解析整数，忽略错误
func parseInt(s string) int {
	if s == "" {
		return 0
	}
	
	result := 0
	for _, r := range s {
		if r >= '0' && r <= '9' {
			result = result*10 + int(r-'0')
		} else {
			return 0 // 非数字字符，返回0
		}
	}
	return result
} 