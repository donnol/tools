package inject

import "testing"

type I interface {
	Name() string
	SetName(string)
}

func NewI() I {
	return &impl{}
}

type impl struct {
	name string
}

func (impl *impl) Name() string {
	return impl.name
}

func (impl *impl) SetName(name string) {
	impl.name = name
}

type I1 interface {
	Name() string
	SetName(string)
}

type Name string

func NewName() Name {
	return Name("default name")
}

func NewI1(name Name) I1 {
	return &impl1{
		name: name,
	}
}

type impl1 struct {
	name Name
}

func (impl *impl1) Name() string {
	return string(impl.name)
}

func (impl *impl1) SetName(name string) {
	impl.name = Name(name)
}

type M struct {
	I  I
	I1 I1
}

func TestInject(t *testing.T) {
	// 新建ioc
	ioc := NewIoc(false)

	// 注册
	if err := ioc.RegisterProvider(NewI); err != nil {
		t.Fatal(err)
	}
	if err := ioc.RegisterProvider(NewI1); err != nil {
		t.Fatal(err)
	}
	if err := ioc.RegisterProvider(NewName); err != nil {
		t.Fatal(err)
	}

	// 注入
	m := M{}
	if err := ioc.Inject(&m); err != nil {
		t.Fatal(err)
	}

	// 业务
	if m.I.Name() != "" {
		t.Fatalf("Bad name: %v != %v\n", m.I.Name(), "")
	}
	cas1 := "jd"
	m.I.SetName(cas1)
	if m.I.Name() != cas1 {
		t.Fatalf("Bad name: %v != %v\n", m.I.Name(), cas1)
	}
	cas2 := "jc"
	m.I.SetName(cas2)
	if m.I.Name() != cas2 {
		t.Fatalf("Bad name: %v != %v\n", m.I.Name(), cas2)
	}
	if m.I1.Name() != "default name" {
		t.Fatalf("Bad name: %v != %v\n", m.I1.Name(), "")
	}
	cas1 = "jd"
	m.I1.SetName(cas1)
	if m.I1.Name() != cas1 {
		t.Fatalf("Bad name: %v != %v\n", m.I1.Name(), cas1)
	}
	cas2 = "jc"
	m.I1.SetName(cas2)
	if m.I1.Name() != cas2 {
		t.Fatalf("Bad name: %v != %v\n", m.I1.Name(), cas2)
	}
}
