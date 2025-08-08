package widget

// 基础输入组件类型
const (
	// WidgetInput 文本输入框
	WidgetInput = "input"
	// WidgetNumber 数字输入框
	WidgetNumber = "number"
	// WidgetCheckbox 多选框
	WidgetCheckbox = "checkbox"
	// WidgetRadio 单选框
	WidgetRadio = "radio"
	// WidgetSelect 下拉框
	WidgetSelect = "select"
	// WidgetSwitch 开关
	WidgetSwitch = "switch"
	// WidgetSlider 滑块
	WidgetSlider = "slider"
	// WidgetColor 颜色选择器
	WidgetColor = "color"
	// WidgetMultiSelect 多选下拉框
	WidgetMultiSelect = "multiselect"
	// WidgetTag 标签组件
	WidgetTag = "tag"
	// WidgetFileUpload 文件上传组件
	WidgetFileUpload = "file_upload"
	// WidgetFileDisplay 文件展示组件
	WidgetFileDisplay = "file_display"
	// WidgetListInput 列表输入组件
	WidgetListInput = "list_input"

	//WidgetList 等价list_input，输入输出通用
	WidgetList = "list"

	// WidgetFormInput 表单输入组件
	WidgetFormInput = "form_input"

	//WidgetForm 和form_input等价
	WidgetForm = "form"
)

// 表格组件类型
const (
	// WidgetTable 表格
	WidgetTable = "table"
)

// 图表组件类型
const (
	// WidgetLineChart 折线图
	WidgetLineChart = "line"
	// WidgetBarChart 柱状图
	WidgetBarChart = "bar"
	// WidgetPieChart 饼图
	WidgetPieChart = "pie"
	// WidgetScatterChart 散点图
	WidgetScatterChart = "scatter"
	// WidgetRadarChart 雷达图
	WidgetRadarChart = "radar"
)

// 日期时间组件类型
const (
	// WidgetDate 日期选择器
	WidgetDate = "date"
	// WidgetTime 时间选择器
	WidgetTime = "time"
	// WidgetDateTime 日期时间选择器
	WidgetDateTime = "datetime"
	// WidgetDateRange 日期范围选择器
	WidgetDateRange = "daterange"
	// WidgetTimeRange 时间范围选择器
	WidgetTimeRange = "timerange"
)

// 数据类型
const (
	// TypeString 字符串类型
	TypeString = "string"
	// TypeNumber 数字类型
	TypeNumber = "number"
	// TypeBoolean 布尔类型
	TypeBoolean = "boolean"
	// TypeArray 数组类型
	TypeArray = "array"
	// TypeObject 对象类型
	TypeObject = "object"
	// TypeTime 时间类型
	TypeTime = "time"
	// TypeFloat 浮点数类型
	TypeFloat = "float"
	// TypeFiles 文件类型（复数形式，表示可能是多个文件）
	TypeFiles = "files"
	// TypeStruct 结构体类型
	TypeStruct = "struct"
	// TypeListStruct 结构体数组类型
	TypeListStruct = "[]struct"
)
