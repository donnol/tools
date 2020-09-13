package reflectx

import "time"

// User 用户
// 这是一个表示用户的结构体
type User struct {
	*Model
	Basic      // 基础信息
	Age   int  // 年纪
	Sex   *int // 性别
	// Name         *string   // 名字 FIXME 与匿名结构体字段重名
	password     string    // 密码，非导出字段
	Phone        *string   // 电话号码
	Salary       float64   // 薪水
	SalaryPtr    *float64  // 薪水
	IsManager    bool      // 是否管理层
	IsManagerPtr *bool     // 是否管理层
	AddressList  []Address // 地址列表
	CreatedAt    time.Time // 创建时间

	Test       Model   `json:"test"`      // 测试
	TestSlice  []Model `json:"testSlice"` // 测试数组
	TestStruct struct {
		Model
		InnerTest      Model   `json:"innerTest"`      // 内部测试
		InnerTestSlice []Model `json:"innerTestSlice"` // 内部测试数组
	}

	// 随机赋值时，若有at tag，则使用tag里指定的值
	Range    int    `json:"range" at:"range=one(1,5,8,13)"`     // 范围，值是[1, 5)或[8, 13)，支持多个范围
	Regexp   string `json:"regexp" at:"regexp=one('^[a-z]+$')"` // 正则表达式，值必须是满足它的字符串
	EnumOne  *int   `json:"enumOne" at:"enum=one(1,2,3)"`       // 枚举，值必须是1,2,3里的一个，不传用nil表示
	EnumMany []int  `json:"enumMany" at:"enum=many(1,2,3)"`     // 枚举，值可以是1,2,3里的多种组合，不传用nil表示
	DBOne    string `json:"dbOne" at:"db=one(org,id)"`          // 值是来自数据库的数据，org表示表名，id表示表字段
	DBMany   string `json:"dbMany" at:"db=many(org,id)"`        // 值是来自数据库的数据组合，org表示表名，id表示表字段
}

// Address 地址
type Address struct {
	Basic           // 基础信息
	Position string // 位置
}

// Basic 基础信息
type Basic struct {
	Name string `json:"name"` // 名字
}

// 结构体必须这样定义，否则拿不到它的doc
type (
	// Model 模型
	//
	// 模型注释
	//
	// 模型描述.
	// 模型描述2.
	Model struct {
		InnerModel
		// 级别类型:1 P1;2 P2;3 P3
		Level       int      `json:"level"` // 级别
		LevelPtr    *int     // 级别
		Position    float64  // 位置
		PositionPtr *float64 // 位置
		Head        string   // 头
		HeadPtr     *string  // 头
		IsTurbo     bool     // 是否涡轮
		IsTurboPtr  *bool    // 是否涡轮
		// Help func() // 帮助函数
	}
)

// InnerModel 内部模型
type InnerModel struct {
	CornerStone string `json:"stone,omitempty"` // 基石
}
