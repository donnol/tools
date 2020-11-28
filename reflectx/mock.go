package reflectx

type StructMock struct {
	GetFieldsFunc func() []Field
}

var _ IStruct = &StructMock{}

func (mockRecv *StructMock) GetFields() []Field {
	return mockRecv.GetFieldsFunc()
}
