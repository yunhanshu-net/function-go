package runner

type FunctionType string
type OperateTableType string

const (
	FunctionTypeStatic  FunctionType = "static_function"  //无需参数，或者输入参数，但是结果永远恒定
	FunctionTypeDynamic FunctionType = "dynamic_function" //default，请求参数不可预测，响应参数不可预测（比如查询用户信息，用户的信息时时刻刻都可能在变化）输入zhangsan 现在可能他18岁，但是过了一会他就19岁了，所以结果不可预知
	FunctionTypePure    FunctionType = "pure_function"    //纯函数，例如1+1=2 or  sin(x) cos(x) 这种
)
const (
	OperateTableTypeCreate OperateTableType = "create" //增加数据
	OperateTableTypeUpdate OperateTableType = "update" //修改数据
	OperateTableTypeDelete OperateTableType = "delete" //删除数据
	OperateTableTypeGet    OperateTableType = "get"    //获取数据
)

type ApiInfo struct {
	Router       string       `json:"router"`
	Method       string       `json:"method"`
	ApiDesc      string       `json:"api_desc"`
	IsPublicApi  bool         `json:"is_public_api"`
	ChineseName  string       `json:"chinese_name"`
	EnglishName  string       `json:"english_name"`
	Classify     string       `json:"classify"`
	Tags         []string     `json:"tags"`
	Async        bool         `json:"async"`         //是否异步，比较耗时的api，或者需要后台慢慢处理的api
	FunctionType FunctionType `json:"function_type"` //函数类型 默认：dynamic_function
	Timeout      int          `json:"timeout"`       //超时时间，单位毫秒,0表示不超时

	RenderType    string                   `json:"widget"`         // 渲染类型	//form，table，echarts
	CreateTables  []interface{}            `json:"create_tables"`  //创建该api时候会自动帮忙创建这个数据库表gorm的model列表
	UseTables     []interface{}            `json:"use_tables"`     //这里需要记录这个函数用到的数据表，方便梳理引用关系
	OperateTables map[interface{}][]string `json:"operate_tables"` //用到了哪些表，对表进行了哪些操作方便梳理引用关系
	UseDB         []string                 `json:"use_db"`         //用到的db文件
	AutoRun       bool                     `json:"-"`              //是否自动运行，默认false，如果为true，则在用户访问这个函数时候，会自动运行一次
	Request       interface{}              `json:"-"`
	Response      interface{}              `json:"-"`

	//用map的都是字段级别的回调，其他的都是接口级别回调

	OnPageLoad OnPageLoad `json:"-"` //优先级最高，先初始化表单参数，然后再判断是否有自动运行的回调，如果有，则执行

	OnApiCreated    OnApiCreated    `json:"-"`
	OnApiUpdated    OnApiUpdated    `json:"-"`
	BeforeApiDelete BeforeApiDelete `json:"-"`
	AfterApiDeleted AfterApiDeleted `json:"-"`

	BeforeRunnerClose BeforeRunnerClose `json:"-"`
	AfterRunnerClose  AfterRunnerClose  `json:"-"`
	OnVersionChange   OnVersionChange   `json:"-"`

	OnTableDeleteRows OnTableDeleteRows `json:"-"`
	OnTableUpdateRow  OnTableUpdateRow  `json:"-"`
	OnTableSearch     OnTableSearch     `json:"-"`

	OnInputFuzzyMap    map[string]OnInputFuzzy    `json:"-"` //key是字段的code，字段级回调
	OnInputValidateMap map[string]OnInputValidate `json:"-"` //key是字段的code，字段级回调

}
