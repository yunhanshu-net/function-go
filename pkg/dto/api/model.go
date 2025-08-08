package api

import (
	"fmt"
	"strings"
)

type Config struct {
}

type Group struct {
	CnName string `json:"cn_name"`
	EnName string `json:"en_name"`
}

type Info struct {
	Router       string   `json:"router"`
	Method       string   `json:"method"`
	User         string   `json:"user"`
	Runner       string   `json:"runner"`
	ApiDesc      string   `json:"api_desc"`
	ChineseName  string   `json:"chinese_name"`
	EnglishName  string   `json:"english_name"`
	Classify     string   `json:"classify"`
	Tags         []string `json:"tags"`
	Timeout      int      `json:"timeout"`
	AutoRun      bool     `json:"-"`             //是否自动运行，默认false，如果为true，则在用户访问这个函数时候，会自动运行一次
	Async        bool     `json:"async"`         //是否异步，比较耗时的api，或者需要后台慢慢处理的api
	FunctionType string   `json:"function_type"` //函数类型 默认：dynamic_function
	Group        *Group   `json:"group"`

	//form，table，
	RenderType    string              `json:"widget"`         // 渲染类型
	DBName        string              `json:"db_name"`        //创建该api时候会自动帮忙创建这个数据库
	CreateTables  []string            `json:"create_tables"`  //创建该api时候会自动帮忙创建这个数据库表gorm的model列表
	OperateTables map[string][]string `json:"operate_tables"` //用到了哪些表，对表进行了哪些操作方便梳理引用关系

	//输入参数
	ParamsIn interface{} `json:"params_in"`
	//输出参数
	ParamsOut interface{} `json:"params_out"`
	UseTables []string    `json:"use_tables"`
	UseDB     []string    `json:"use_db"`
	Callbacks []string    `json:"callbacks"`

	// 配置相关
	ParamsConfig interface{} `json:"params_config"` // 配置结构体
	ParamsData   interface{} `json:"params_data"`   // 配置初始值
}

func (i *Info) GetTreePath() string {
	return fmt.Sprintf("/%s/%s/%s/", i.User, i.Runner, strings.Trim(i.Router, "/"))
}

func (i *Info) GetParentTreePath() string {

	split := strings.Split(strings.Trim(i.Router, "/"), "/")
	join := strings.Join(split[:len(split)-1], "/")
	key := fmt.Sprintf("/%s/%s/%s/", i.User, i.Runner, join)
	return key
}

type ApiLogs struct {
	Version string `json:"version"`

	Apis []*Info `json:"apis"`
}
