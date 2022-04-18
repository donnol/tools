package common

import (
	"database/sql"
	"time"
)

type ColumnMajor struct {
	ColumnId
	ColumnTime
	ColumnSource
	ColumnOperator
}

type ColumnId struct {
	Id uint `json:"id,string"`
}

type ColumnTime struct {
	Created time.Time    `json:"created"`
	Updated sql.NullTime `json:"updated"`
}

type ColumnSource struct {
	Source string `json:"source"` // 来源
}

type ColumnOperator struct {
	OperatorId uint `json:"operatorId,string"` // 操作人
}
