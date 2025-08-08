package widget

// 转义功能使用示例
// 这个文件展示了如何在options中使用转义字符

/*
转义规则说明：

1. 基本转义：
   - \, 表示字面量的逗号
   - \\ 表示字面量的反斜杠

2. 使用示例：

// 包含逗号的选项
widget:"type:select;options:逗号(\,),分号(;),制表符(\\t),竖线(|)"

// 包含反斜杠的选项
widget:"type:select;options:路径\\folder,文件\\name.txt,配置\\config.json"

// 混合转义
widget:"type:select;options:普通选项,包含\,的选项,包含\\的选项"

3. 实际应用场景：

// CSV分隔符选择
CSVDelimiter string `widget:"type:select;options:逗号(\,),分号(;),制表符(\\t),竖线(|)"`

// 文件路径选择
FilePath string `widget:"type:select;options:C:\\Users\\name,D:\\data\\file.txt,/home/user/file"`

// 配置选项
Config string `widget:"type:select;options:默认配置,高级配置\,包含逗号,特殊配置\\包含反斜杠"`

4. 注意事项：
   - 在Go字符串字面量中，需要双反斜杠来表示单反斜杠
   - 例如：`"逗号(\\,)"` 在运行时会被解析为 `"逗号(,)"`
   - 例如：`"路径\\\\folder"` 在运行时会被解析为 `"路径\\folder"`
*/ 