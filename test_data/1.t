这是改造前的
Phone     string `json:"phone" gorm:"column:phone;comment:手机号" runner:"code:phone;name:手机号" widget:"type:input;placeholder:请输入11位手机号" search:"like" validate:"required,len=11,numeric"`
这是改造后的
Phone     string  `json:"phone" gorm:"column:phone" runner:"name:手机号" widget:"type:input" search:"like" validate:"required,len=11,numeric" msg:"手机号必须为11位数字"`

