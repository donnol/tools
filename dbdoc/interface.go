package dbdoc

import "io"

type ITable interface {
	MakeGraph() *Table
	New() *Table
	Resolve(v interface{}) *Table
	SetComment(comment string) *Table
	SetDescription(description string) *Table
	SetMapper(f Mapper) *Table
	SetTypeMapper(f Mapper) *Table
	Write(w io.Writer) *Table
}
