package inject

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"
	// "github.com/donnol/tools/reflectx"
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
		fmt.Printf("| userHook | welcome to method Add, args: %v\n", args)
	case "Get":
		fmt.Printf("| userHook | welcome to method Get, args: %v\n", args)
	case "GetContext":
		fmt.Printf("| userHook | welcome to method GetContext, args: %v\n", args)
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

// === 基于接口的Around ===
var (
	customCtxMap = make(map[string]CtxFunc)
)

type store interface {
	add(string, int) (int, error)
	add2(string, int) (int, error)
	get(id int) string
	equal(first int, args ...any)
	equalInt(first int, args ...int)
}

func NewStore(withProxy bool) store {
	base := &storeImpl{}
	if withProxy {
		return getStoreProxy(base)
	}
	return base
}

type storeImpl struct{}

func (impl *storeImpl) add(name string, id int) (int, error) {
	log.Printf("arg, name: %v, id: %v\n", name, id)
	return 1, nil
}
func (impl *storeImpl) get(id int) string {
	return "jd"
}
func (impl *storeImpl) add2(string, int) (int, error) {
	return 2, nil
}
func (impl *storeImpl) equal(first int, args ...any) {
	log.Printf("[equal] args: %v, %+v\n", first, args)
}
func (impl *storeImpl) equalInt(first int, args ...int) {
	log.Printf("[equalInt] args: %v, %+v\n", first, args)
}

// TODO: 生成mock相关代码
func getStoreProxy(base store) *storeMock {
	return &storeMock{
		addFunc: func(arg1 string, arg2 int) (int, error) {
			var r1 int
			var r2 error
			begin := time.Now()

			ctx := ProxyContext{
				PkgPath:       "github.com/donnol/tools/inject",
				InterfaceName: "store",
				MethodName:    "add",
			}
			cf, ok := customCtxMap[ctx.Uniq()]
			if ok {
				res := cf(ctx, base.add, []any{arg1, arg2})
				r1 = res[0].(int)
				if res[1] == nil {
					r2 = nil
				} else {
					r2 = res[1].(error)
				}
			} else {
				// 默认调用
				r1, r2 = base.add(arg1, arg2)
			}

			log.Printf("[ctx: %s]used time: %v\n", ctx.Uniq(), time.Since(begin))

			return r1, r2
		},
		add2Func: base.add2,
		equalFunc: func(first int, args ...any) {
			begin := time.Now()

			ctx := ProxyContext{
				PkgPath:       "github.com/donnol/tools/inject",
				InterfaceName: "store",
				MethodName:    "equal",
			}
			cf, ok := customCtxMap[ctx.Uniq()]
			if ok {
				allArg := []any{first}
				allArg = append(allArg, args...)
				cf(ctx, base.equal, allArg)
			} else {
				base.equal(first, args...)
			}

			log.Printf("[ctx: %s]used time: %v\n", ctx.Uniq(), time.Since(begin))
		},
		equalIntFunc: func(first int, args ...int) {
			begin := time.Now()

			allArg := []any{first}
			for _, arg := range args {
				allArg = append(allArg, arg)
			}

			ctx := ProxyContext{
				PkgPath:       "github.com/donnol/tools/inject",
				InterfaceName: "store",
				MethodName:    "equalInt",
			}
			cf, ok := customCtxMap[ctx.Uniq()]
			if ok {
				cf(ctx, base.equalInt, allArg)
			} else {
				base.equalInt(first, args...)
			}

			log.Printf("[ctx: %s]used time: %v\n", ctx.Uniq(), time.Since(begin))
		},
	}
}

type storeMock struct {
	addFunc      func(string, int) (int, error)
	add2Func     func(string, int) (int, error)
	getFunc      func(id int) string
	equalFunc    func(first int, args ...any)
	equalIntFunc func(first int, args ...int)
}

func (mock *storeMock) add(name string, id int) (int, error) {
	return mock.addFunc(name, id)
}
func (mock *storeMock) add2(name string, id int) (int, error) {
	return mock.add2Func(name, id)
}

func (mock *storeMock) get(id int) string {
	return mock.getFunc(id)
}
func (mock *storeMock) equal(first int, args ...any) {
	mock.equalFunc(first, args...)
}
func (mock *storeMock) equalInt(first int, args ...int) {
	mock.equalIntFunc(first, args...)
}

type src interface {
	add(string) (int, error)
	add2(string) (int, error)
}

func NewSrc(withProxy bool, store store) src {
	base := &srcImpl{
		store: store,
	}
	if withProxy {
		return getSrcProxy(base)
	}
	return base
}

type srcImpl struct {
	store store
}

func (impl *srcImpl) add(name string) (int, error) {
	impl.store.equal(1, 2, 3, 4, "hah") // 不管加多少参数，都没问题
	impl.store.equalInt(1, 2, 3, 4)     // 不管加多少参数，都没问题
	return impl.store.add(name, 1)
}
func (impl *srcImpl) add2(name string) (int, error) {
	return impl.store.add2(name, 2)
}

func getSrcProxy(base src) *srcMock {
	return &srcMock{
		addFunc: func(arg1 string) (int, error) {
			var r1 int
			var r2 error
			begin := time.Now()

			ctx := ProxyContext{
				PkgPath:       "github.com/donnol/tools/inject",
				InterfaceName: "src",
				MethodName:    "add",
			}
			cf, ok := customCtxMap[ctx.Uniq()]
			if ok {
				res := cf(ctx, base.add, []any{arg1})
				r1 = res[0].(int)
				if res[1] == nil {
					r2 = nil
				} else {
					r2 = res[1].(error)
				}
			} else {
				r1, r2 = base.add(arg1)
			}

			log.Printf("[ctx: %s]used time: %v\n", ctx.Uniq(), time.Since(begin))

			return r1, r2
		},
		add2Func: base.add2,
	}
}

type srcMock struct {
	addFunc  func(name string) (int, error)
	add2Func func(string) (int, error)
}

func (mock *srcMock) add(name string) (int, error) {
	return mock.addFunc(name)
}
func (mock *srcMock) add2(name string) (int, error) {
	return mock.add2Func(name)
}

func TestApi(t *testing.T) {
	customCtxMap["github.com/donnol/tools/inject|src|add"] = func(ctx ProxyContext, method any, args []any) (res []any) {
		log.Printf("custom call")
		f := method.(func(string) (int, error))
		a1 := args[0].(string)
		r1, r2 := f(a1)
		res = append(res, r1, r2)
		return res
	}

	store := NewStore(true) // 用mock包装storeImpl
	src := NewSrc(true, store)

	r, err := src.add("jd")
	if err != nil {
		panic(err)
	}
	log.Printf("r: %v\n", r)

	r, err = src.add2("jd")
	if err != nil {
		panic(err)
	}
	log.Printf("r: %v\n", r)
}

// === TODO: 任意函数的Around ===

// 用户编写的代码
func A(ctx any, id int, args ...string) (string, error) {
	log.Printf("arg, ctx: %v, id: %v, args: %+v\n", ctx, id, args)
	return "A", nil
}

func C() {
	// 编译之前，通过重写ast，改为调用B（B内再调用A）
	// 1 遍历源码，找到函数调用（可配置规则，以过滤出想要改变的函数）- *ast.CallExpr
	// 2 生成一个对应的附加了额外逻辑的函数B（B内调用A）- *ast.FuncDecl
	// 3 将此处对A的调用替换为对B的调用 -
	r1, err := A(1, 1, "a", "b")
	if err != nil {
		log.Printf("err is not nil: %v\n", err)
		return
	}
	log.Printf("r1: %v\n", r1)
}

// 生成出来的代码
// 有没有办法生成一个B，使得C调B，跟C调A一样，但是又可以在B里添加额外的逻辑呢？
func B(ctx any, id int, args ...string) (string, error) {
	// 为了要支持添加额外的逻辑，显然不能直接返回
	// return A(ctx, id, args...)

	// 添加逻辑
	begin := time.Now()

	// 根据签名，不难生成出以下返回值定义
	var r1 string
	var r2 error

	r1, r2 = A(ctx, id, args...)

	// 添加逻辑
	log.Printf("used time: %v\n", time.Since(begin))

	return r1, r2
}
