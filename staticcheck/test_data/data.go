package test_data

import "github.com/fishedee/tools/query"

type User struct {
	Id uint
	id uint
}

func Test() {
	var users []User
	query.Column[User, uint](users, "Id") // exist
	query.Column[User, uint](users, "id") // exist
	query.Column[User, uint](users, "ID") // not exist

	query.Group[User, uint, []uint](users, "Id", func(t []User) uint {
		return 0
	})
	query.Group[User, uint, []uint](users, "id", func(t []User) uint {
		return 0
	})
	query.Group[User, uint, []uint](users, "ID", func(t []User) uint {
		return 0
	})
}
