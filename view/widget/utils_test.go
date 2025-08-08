package widget

import (
	"reflect"
	"testing"
)

func TestParseOptionsWithEscape(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "基本选项",
			input:    "选项1,选项2,选项3",
			expected: []string{"选项1", "选项2", "选项3"},
		},
		{
			name:     "包含逗号的选项",
			input:    "逗号(\\,),分号(;),制表符(\\\\t)",
			expected: []string{"逗号(,)", "分号(;)", "制表符(\\t)"},
		},
		{
			name:     "包含反斜杠的选项",
			input:    "路径\\\\folder,文件\\\\name.txt",
			expected: []string{"路径\\folder", "文件\\name.txt"},
		},
		{
			name:     "混合转义",
			input:    "普通选项,包含\\,的选项,包含\\\\的选项",
			expected: []string{"普通选项", "包含,的选项", "包含\\的选项"},
		},
		{
			name:     "空字符串",
			input:    "",
			expected: nil,
		},
		{
			name:     "只有空格",
			input:    "   ,  ,  ",
			expected: nil,
		},
		{
			name:     "以反斜杠结尾",
			input:    "选项1,选项2\\\\",
			expected: []string{"选项1", "选项2\\"},
		},
		{
			name:     "无效转义序列",
			input:    "选项1,选项\\\\x,选项3",
			expected: []string{"选项1", "选项\\x", "选项3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseOptionsWithEscape(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parseOptionsWithEscape(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
} 