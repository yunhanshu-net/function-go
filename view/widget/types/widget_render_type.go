package types

// 基础输入组件类型
const (
	// WidgetInput 文本输入框
	WidgetInput = "input"
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
	// WidgetFile 文件上传组件
	WidgetFile = "file"
	// WidgetListInput 列表输入组件 - 新增
	WidgetListInput = "list_input"
	// WidgetFormInput 表单输入组件 - 新增
WidgetFormInput = "form_input"
)

// 显示组件类型（输出时使用）
const (
	// WidgetListDisplay 列表显示组件
	WidgetListDisplay = "list_display"
	// WidgetFormDisplay 表单显示组件
	WidgetFormDisplay = "form_display"
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
