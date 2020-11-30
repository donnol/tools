package inject

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/donnol/tools/reflectx"
)

// define model

type IUserModel interface {
	Add(name string) int
	Get(id int) string
	GetContext(ctx context.Context, id int) string
}

func NewIUser() IUserModel {
	return &userImpl{}
}

type userImpl struct {
	id   int
	name string
}

func (impl *userImpl) Add(name string) int {
	impl.name = name
	return 1
}

func (impl *userImpl) Get(id int) string {
	if impl.id == id {
		return impl.name
	}
	return impl.name
}

func (impl *userImpl) GetContext(ctx context.Context, id int) string {
	return impl.Get(id)
}

// UserMock mock结构体，字段数量和需要实现的接口的方法数一致，并且名称是方法名加后缀'Func'，字段类型与方法签名一致
type UserMock struct {
	AddFunc        func(name string) int
	GetHelper      func(id int) string `method:"Get"` // 表示这个字段关联的方法是Get
	GetContextFunc func(ctx context.Context, id int) string
}

func (mock *UserMock) Add(name string) int {
	return mock.AddFunc(name)
}

func (mock *UserMock) Get(id int) string {
	return mock.GetHelper(id)
}

func (mock *UserMock) GetContext(ctx context.Context, id int) string {
	return mock.GetContextFunc(ctx, id)
}

var _ IUserModel = &UserMock{}

// define service

type IUserSrv interface {
	Add(name string) int
	Get(id int) string
	GetContext(ctx context.Context, id int) string
}

func NewIUserSrv(
	userModel IUserModel,
) IUserSrv {
	return &userSrvImpl{
		userModel: userModel,
	}
}

// 可以看到，这里依赖了model接口，现在我们想在这里调用model方法的前后调用钩子，也就是方便地插入代码
type userSrvImpl struct {
	userModel IUserModel
}

func (impl *userSrvImpl) Add(name string) int {
	return impl.userModel.Add(name)
}

func (impl *userSrvImpl) Get(id int) string {
	return impl.userModel.Get(id)
}

func (impl *userSrvImpl) GetContext(ctx context.Context, id int) string {
	return impl.userModel.GetContext(ctx, id)
}

type UserSrvMock struct {
	AddFunc        func(name string) int
	GetFunc        func(id int) string
	GetContextFunc func(ctx context.Context, id int) string
}

func (mock *UserSrvMock) Add(name string) int {
	return mock.AddFunc(name)
}

func (mock *UserSrvMock) Get(id int) string {
	return mock.GetFunc(id)
}

func (mock *UserSrvMock) GetContext(ctx context.Context, id int) string {
	return mock.GetContextFunc(ctx, id)
}

var _ IUserSrv = &UserSrvMock{}

type userHook struct {
}

type testKeyType string

const (
	testKey testKeyType = "testKey"
)

// 在这个方法里的变量就是随请求而变的
func (hook *userHook) Around(pctx ProxyContext, method reflect.Value, args []reflect.Value) []reflect.Value {
	pctx.LogShortf("Around\n")

	// 如果正在执行的方法不同，可以做不同的处理
	// 特别是在需要结合参数来处理时
	switch pctx.MethodName {
	case "Add":
		fmt.Printf("| userHook | welcome to method Add\n")
		fmt.Printf("| userHook | args: %+v\n", reflectx.ToInterface(args))
	case "Get":
		fmt.Printf("| userHook | welcome to method Get\n")
		fmt.Printf("| userHook | args: %+v\n", reflectx.ToInterface(args))
	case "GetContext":
		fmt.Printf("| userHook | welcome to method GetContext\n")
		iargs := reflectx.ToInterface(args)
		fmt.Printf("| userHook | args: %+v\n", iargs)
		ctx := iargs[0].(context.Context)
		fmt.Printf("| userHook | args context: %+v, value: %v\n", ctx, ctx.Value(testKey))
	}

	begin := time.Now()
	fmt.Printf("| userHook | begin: %p\n", &begin)

	// 不调原有的方法，重写为其它方法都行
	result := method.Call(args)

	pctx.Logf("| userHook | used time: %v\n", time.Since(begin))

	return result
}

func TestProxy(t *testing.T) {
	ctx := context.Background()

	ioc := NewIoc(true)

	proxy := NewProxy()
	userModelProvider := proxy.Around(NewIUser, &UserMock{}, &userHook{})

	// 注册
	if err := ioc.RegisterProvider(userModelProvider); err != nil {
		t.Fatal(err)
	}
	if err := ioc.RegisterProvider(NewIUserSrv); err != nil {
		t.Fatal(err)
	}

	// 依赖注入
	type insImpl struct {
		userSrv IUserSrv
	}
	var ins insImpl
	if err := ioc.Inject(&ins); err != nil {
		t.Fatal(err)
	}

	// 实际执行
	name := "jd"
	id := ins.userSrv.Add(name)
	gname := ins.userSrv.Get(id)
	if gname != name {
		t.Fatalf("Bad result, %v != %v\n", gname, name)
	}
	ins.userSrv.GetContext(ctx, id)

	ctx = context.WithValue(ctx, testKey, "testValue")
	// 再执行一次同样的，Around方法里的变量是不一样的
	id = ins.userSrv.Add(name)
	gname = ins.userSrv.Get(id)
	if gname != name {
		t.Fatalf("Bad result, %v != %v\n", gname, name)
	}
	ins.userSrv.GetContext(ctx, id)
}
