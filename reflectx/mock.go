package reflectx

type IStructMock struct {
	GetFieldsFunc func() []Field
}

var _ IStruct = &IStructMock{}

func (*IStructMock) GetFields() []Field {
	panic("Need to be implement!")
}
