package inject

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
)

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

type I2 interface {
	Say(string)
}

func NewI2() (I2, error) {
	return &impl2{}, nil
}

type impl2 struct {
}

func (impl *impl2) Say(a string) {
	fmt.Println(a)
}

type I3 interface {
	Say(string)
}

func NewI3() (I3, error) {
	return &impl3{}, errors.Errorf("case: new return non nil error")
}

type impl3 struct {
}

func (impl *impl3) Say(a string) {
	fmt.Println(a)
}

type I4 interface {
	Say(string)
}

func NewI4() (I4, error, string) {
	return &impl4{}, nil, "case: new return non error type"
}

type impl4 struct {
}

func (impl *impl4) Say(a string) {
	fmt.Println(a)
}

type M struct {
	I  I
	I1 I1
	I2 I2
	// I3 I3 // 测试非nil error
	// I4 I4 // 测试非error类型
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
	if err := ioc.RegisterProvider(NewI2); err != nil {
		t.Fatal(err)
	}
	if err := ioc.RegisterProvider(NewI3); err != nil {
		t.Fatal(err)
	}
	if err := ioc.RegisterProvider(NewI4); err != nil {
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
	m.I2.Say(cas2)
}
