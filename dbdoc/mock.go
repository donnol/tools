package dbdoc

import "io"

type ITableMock struct {
	MakeGraphFunc func() *Table

	NewFunc func() *Table

	ResolveFunc func(v interface{}) *Table

	SetCommentFunc func(comment string) *Table

	SetDescriptionFunc func(description string) *Table

	SetMapperFunc func(f Mapper) *Table

	SetTypeMapperFunc func(f Mapper) *Table

	WriteFunc func(w io.Writer) *Table
}

var _ ITable = &ITableMock{}

func (*ITableMock) MakeGraph() *Table {
	panic("Need to be implement!")
}

func (*ITableMock) New() *Table {
	panic("Need to be implement!")
}

func (*ITableMock) Resolve(v interface{}) *Table {
	panic("Need to be implement!")
}

func (*ITableMock) SetComment(comment string) *Table {
	panic("Need to be implement!")
}

func (*ITableMock) SetDescription(description string) *Table {
	panic("Need to be implement!")
}

func (*ITableMock) SetMapper(f Mapper) *Table {
	panic("Need to be implement!")
}

func (*ITableMock) SetTypeMapper(f Mapper) *Table {
	panic("Need to be implement!")
}

func (*ITableMock) Write(w io.Writer) *Table {
	panic("Need to be implement!")
}
