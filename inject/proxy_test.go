package inject

import (
	"context"
	"fmt"
	"log"
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

// === Around ===
var (
	customCtxMap = make(map[string]CtxFunc)
)

// 我只要用ProxyHelper来调用A，就可以在调用A前后添加逻辑
// 可以看到，这里的关键是f的类型是不确定的，如果把f的类型改为any，就必须断言得到具体的类型，这个需要用代码生成在编译前完成
// 记录所有方法的函数类型签名
// 记录参数个数以及参数类型，并做断言
// 记录返回值个数，并append到结果里，res是[]any类型，那怎么得到具体的返回值呢
//
// 这个函数在另外的函数看来是透明的
func ProxyHelper(ctx ProxyContext, method any, args []any) (res []any) {
	// 执行前可以做很多东西

	begin := time.Now()

	// TODO: 生成这部分代码
	// 如果仅仅用签名来区分，遇到签名一样的不同方法时怎么区分呢？
	// 所以，必须加入ctx的Uniq部分来区分
	switch ctx.Uniq() {
	case "github.com/donnol/tools/inject|store|add":
		cf, ok := customCtxMap[ctx.Uniq()]
		if ok {
			// 自定义around，需要自己决定怎么调用，一般需要包含下面的默认调用
			res = cf(ctx, method, args)
		} else {
			// 默认调用
			f := method.(func(string, int) (int, error))
			a1 := args[0].(string)
			a2 := args[1].(int)
			r1, r2 := f(a1, a2)
			res = append(res, r1, r2)
		}

	case "github.com/donnol/tools/inject|src|add":
		cf, ok := customCtxMap[ctx.Uniq()]
		if ok {
			// 自定义around，需要自己决定怎么调用，一般需要包含下面的默认调用
			res = cf(ctx, method, args)
		} else {
			f := method.(func(string) (int, error))
			a1 := args[0].(string)
			r1, r2 := f(a1)
			res = append(res, r1, r2)
		}

	case "github.com/donnol/tools/inject|store|equal":
		cf, ok := customCtxMap[ctx.Uniq()]
		if ok {
			// 自定义around，需要自己决定怎么调用，一般需要包含下面的默认调用
			res = cf(ctx, method, args)
		} else {
			f := method.(func(int, ...any))
			// 参数数量是不确定的，怎么知道要做几次断言呢？
			// 因为不定参数一定是最后一个，所以可以用减法，先把前面确定的断言出来，剩下的用args[n:]...传进去就好了
			a1 := args[0].(int)
			f(a1, args[1:]...)
		}

	case "github.com/donnol/tools/inject|store|equalInt":
		cf, ok := customCtxMap[ctx.Uniq()]
		if ok {
			// 自定义around，需要自己决定怎么调用，一般需要包含下面的默认调用
			res = cf(ctx, method, args)
		} else {
			f := method.(func(int, ...int))
			// 参数数量是不确定的，怎么知道要做几次断言呢？
			// 因为不定参数一定是最后一个，所以可以用减法，先把前面确定的断言出来，剩下的用args[n:]...传进去就好了
			a1 := args[0].(int)
			// 如果不定参数不是any类型，需要断言
			left := make([]int, 0, len(args)-1)
			for i := 1; i < len(args); i++ {
				left = append(left, args[i].(int))
			}
			f(a1, left...)
		}

		// more case ...
	}

	// 执行后可以做很多东西

	used := time.Since(begin)
	log.Printf("[ctx: %s]used time: %v\n", ctx.Uniq(), used)

	return
}

type store interface {
	add(string, int) (int, error)
	add2(string, int) (int, error)
	get(id int) string
	equal(first int, args ...any)
	equalInt(first int, args ...int)
}

func NewStore(useMock bool) store {
	base := &storeImpl{}
	if useMock {
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
			res := ProxyHelper(ProxyContext{
				PkgPath:       "github.com/donnol/tools/inject",
				InterfaceName: "store",
				MethodName:    "add",
			}, base.add, []any{arg1, arg2})
			if res[1] == nil {
				return res[0].(int), nil
			}
			return res[0].(int), res[1].(error)
		},
		add2Func: base.add2,
		equalFunc: func(first int, args ...any) {
			allArg := []any{first}
			allArg = append(allArg, args...)
			ProxyHelper(ProxyContext{
				PkgPath:       "github.com/donnol/tools/inject",
				InterfaceName: "store",
				MethodName:    "equal",
			}, base.equal, allArg)
		},
		equalIntFunc: func(first int, args ...int) {
			allArg := []any{first}
			for _, arg := range args {
				allArg = append(allArg, arg)
			}
			ProxyHelper(ProxyContext{
				PkgPath:       "github.com/donnol/tools/inject",
				InterfaceName: "store",
				MethodName:    "equalInt",
			}, base.equalInt, allArg)
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

func NewSrc(useMock bool, store store) src {
	base := &srcImpl{
		store: store,
	}
	if useMock {
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
			res := ProxyHelper(ProxyContext{
				PkgPath:       "github.com/donnol/tools/inject",
				InterfaceName: "src",
				MethodName:    "add",
			}, base.add, []any{arg1})
			if res[1] == nil {
				return res[0].(int), nil
			}
			return res[0].(int), res[1].(error)
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
