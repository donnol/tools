package reflectx

type IStruct interface{ GetFields() []Field }

type IStructMock interface{ GetFields() []Field }
