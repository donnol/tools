package testtype

type Inner struct {
	Phone string `json:"phone"` // 手机
}

type TestModel struct {
	Name string `json:"name"` // 名称
	List []User `json:"list"` // 用户列表
}

type User struct {
	Id   uint   `json:"id,string"` // id
	Name string `json:"name"`      // 名字
	Age  int    `json:"age"`       // 年龄
	Addr Addr   `json:"addr"`      // 地址

	Inner
}

type Addr struct {
	City string `json:"city"` // 城市
	Home string `json:"home"` // 家
}
