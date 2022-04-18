package user

import "project_layout/model/db/common"

type Table struct {
	common.ColumnMajor

	Name string `json:"name"`
}

const (
	TableName = "user"
)

func (Table) TableName() string {
	return TableName
}
