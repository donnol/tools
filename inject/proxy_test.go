package inject

import (
	"testing"
	"time"
)

// define model

type IUserModel interface {
	Add(name string) int
	Get(id int) string
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
	return impl.name
}

// UserMock mock结构体，字段数量和需要实现的接口的方法数一致，并且名称是方法名加后缀'Func'，字段类型与方法签名一致
type UserMock struct {
	AddFunc   func(name string) int
	GetHelper func(id int) string `method:"Get"` // 表示这个字段关联的方法是Get
}

func (mock *UserMock) Add(name string) int {
	return mock.AddFunc(name)
}

func (mock *UserMock) Get(id int) string {
	return mock.GetHelper(id)
}

var _ IUserModel = &UserMock{}

// define service

type IUserSrv interface {
	Add(name string) int
	Get(id int) string
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

type UserSrvMock struct {
	AddFunc func(name string) int
	GetFunc func(id int) string
}

func (mock *UserSrvMock) Add(name string) int {
	return mock.AddFunc(name)
}

func (mock *UserSrvMock) Get(id int) string {
	return mock.GetFunc(id)
}

var _ IUserSrv = &UserSrvMock{}

// hook
type timeHook struct {
	begin time.Time
	end   time.Time
}

func (hook *timeHook) Before(pctx ProxyContext) {
	hook.begin = time.Now()
	pctx.Logf("begin: %v\n", hook.begin)
}

func (hook *timeHook) After(pctx ProxyContext) {
	hook.end = time.Now()
	pctx.Logf("end: %v\n", hook.end)

	// 计算耗时
	used := hook.end.Sub(hook.begin)
	pctx.Logf("used: %v\n", used)

	// 将耗时写入到时序数据库，再利用图表展示出来
}

type userHook struct {
}

func (hook *userHook) Before(pctx ProxyContext) {
	pctx.Logf("user before\n")
}

func (hook *userHook) After(pctx ProxyContext) {
	pctx.Logf("user after\n")
}

var _ Hook = &userHook{}

func TestProxy(t *testing.T) {
	ioc := NewIoc(true)

	proxy := NewProxy()
	proxy.AddHook(&timeHook{})                                          // 设置全局hook
	userModelProvider := proxy.Wrap(NewIUser, &UserMock{}, &userHook{}) // 包装provider，并指定专用hook

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
}
