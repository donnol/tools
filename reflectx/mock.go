package reflectx

type StructMockMock struct {
	GetFieldsFunc func() []Field
}

var _ IStructMock = &StructMockMock{}

func (mockRecv *StructMockMock) GetFields() []Field {
	return mockRecv.GetFieldsFunc()
}

type StructMock struct {
	GetFieldsFunc func() []Field
}

var _ IStruct = &StructMock{}

func (mockRecv *StructMock) GetFields() []Field {
	return mockRecv.GetFieldsFunc()
}
