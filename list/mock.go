package list

type IStringListMock struct {
	FilterFunc func(s string) StringList
}

var _ IStringList = &IStringListMock{}

func (*IStringListMock) Filter(s string) StringList {
	panic("Need to be implement!")
}
