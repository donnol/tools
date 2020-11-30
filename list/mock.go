package list

type StringListMock struct {
	FilterFunc func(s string) StringList
}

var _ IStringList = &StringListMock{}

func (mockRecv *StringListMock) Filter(s string) StringList {
	return mockRecv.FilterFunc(s)
}

type StringListMockMock struct {
	FilterFunc func(s string) StringList
}

var _ IStringListMock = &StringListMockMock{}

func (mockRecv *StringListMockMock) Filter(s string) StringList {
	return mockRecv.FilterFunc(s)
}
