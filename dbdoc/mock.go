package dbdoc

import "io"

type TableMock struct {
	MakeGraphFunc func() *Table

	NewFunc func() *Table

	ResolveFunc func(v interface{}) *Table

	SetCommentFunc func(comment string) *Table

	SetDescriptionFunc func(description string) *Table

	SetMapperFunc func(f Mapper) *Table

	SetTypeMapperFunc func(f Mapper) *Table

	WriteFunc func(w io.Writer) *Table
}

var _ ITable = &TableMock{}

func (mockRecv *TableMock) MakeGraph() *Table {
	return mockRecv.MakeGraphFunc()
}

func (mockRecv *TableMock) New() *Table {
	return mockRecv.NewFunc()
}

func (mockRecv *TableMock) Resolve(v interface{}) *Table {
	return mockRecv.ResolveFunc(v)
}

func (mockRecv *TableMock) SetComment(comment string) *Table {
	return mockRecv.SetCommentFunc(comment)
}

func (mockRecv *TableMock) SetDescription(description string) *Table {
	return mockRecv.SetDescriptionFunc(description)
}

func (mockRecv *TableMock) SetMapper(f Mapper) *Table {
	return mockRecv.SetMapperFunc(f)
}

func (mockRecv *TableMock) SetTypeMapper(f Mapper) *Table {
	return mockRecv.SetTypeMapperFunc(f)
}

func (mockRecv *TableMock) Write(w io.Writer) *Table {
	return mockRecv.WriteFunc(w)
}

type TableMockMock struct {
	MakeGraphFunc func() *Table

	NewFunc func() *Table

	ResolveFunc func(v interface{}) *Table

	SetCommentFunc func(comment string) *Table

	SetDescriptionFunc func(description string) *Table

	SetMapperFunc func(f Mapper) *Table

	SetTypeMapperFunc func(f Mapper) *Table

	WriteFunc func(w io.Writer) *Table
}

var _ ITableMock = &TableMockMock{}

func (mockRecv *TableMockMock) MakeGraph() *Table {
	return mockRecv.MakeGraphFunc()
}

func (mockRecv *TableMockMock) New() *Table {
	return mockRecv.NewFunc()
}

func (mockRecv *TableMockMock) Resolve(v interface{}) *Table {
	return mockRecv.ResolveFunc(v)
}

func (mockRecv *TableMockMock) SetComment(comment string) *Table {
	return mockRecv.SetCommentFunc(comment)
}

func (mockRecv *TableMockMock) SetDescription(description string) *Table {
	return mockRecv.SetDescriptionFunc(description)
}

func (mockRecv *TableMockMock) SetMapper(f Mapper) *Table {
	return mockRecv.SetMapperFunc(f)
}

func (mockRecv *TableMockMock) SetTypeMapper(f Mapper) *Table {
	return mockRecv.SetTypeMapperFunc(f)
}

func (mockRecv *TableMockMock) Write(w io.Writer) *Table {
	return mockRecv.WriteFunc(w)
}
