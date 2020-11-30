package list

type IStringList interface{ Filter(s string) StringList }

type IStringListMock interface{ Filter(s string) StringList }
