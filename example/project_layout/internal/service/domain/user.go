package domain

import "project_layout/model/db/user"

// 当业务复杂时，可将逻辑提取为域(纯内容操作，方便单元测试)，一般普通的增删改查，大可不必

type User struct {
	id   uint
	name string
}

func NewUser(
	id uint,
	name string,
) *User {
	return &User{
		id:   id,
		name: name,
	}
}

func (u *User) FromTable(other user.Table) {
	u.id = other.Id
	u.name = other.Name
}

func (u *User) ToTable() user.Table {
	other := user.Table{}
	other.Id = u.id
	other.Name = u.name
	return other
}

func (u *User) ModName(name string) {
	u.name = name
}
