package inject

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

// Proxy 在层间依赖调用时插入钩子调用，类似AOP
type Proxy interface {
	AddHook(...Hook)                                                        // 添加钩子，支持多个，全局使用，LIFO
	Wrap(provider interface{}, mock interface{}, hooks ...Hook) interface{} // 包装provider，可指定Hook
}

func NewProxy() Proxy {
	return &proxyImpl{}
}

type ProxyContext struct {
	PkgPath       string
	InterfaceName string
	MethodName    string
}

func (pctx ProxyContext) String() string {
	return fmt.Sprintf("[PkgPath: %s, InterfaceName: %s, MethodName: %s]", pctx.PkgPath, pctx.InterfaceName, pctx.MethodName)
}

func (pctx ProxyContext) Logf(format string, args ...interface{}) {
	log.Output(1, fmt.Sprintf(pctx.String()+": "+format, args...))
}

type Hook interface {
	Before(ProxyContext)
	After(ProxyContext)
}

type Caller interface {
	Call(args []reflect.Value) []reflect.Value
}

type Around func(pctx ProxyContext, caller Caller) Caller

func (around Around) Before(pctx ProxyContext) {

}

func (around Around) After(pctx ProxyContext) {

}

type proxyImpl struct {
	hooks []Hook // 钩子列表，LIFO
}

func (impl *proxyImpl) AddHook(hooks ...Hook) {
	impl.hooks = append(impl.hooks, hooks...)
}

var (
	MockFieldNameSuffixes = [...]string{"Func", "Handler"} // mock结构体字段名称后缀
)

// Wrap 从一个provider生成一个新的provider
// 如果mock为nil或者不是结构体指针，则直接返回provider
func (impl *proxyImpl) Wrap(provider interface{}, mock interface{}, hooks ...Hook) interface{} {
	if mock == nil {
		return provider
	}

	mockValue := reflect.ValueOf(mock)
	mockType := mockValue.Type()
	if mockType.Kind() != reflect.Ptr && mockType.Elem().Kind() != reflect.Struct {
		return provider
	}

	// provider有参数，有返回值
	pv := reflect.ValueOf(provider)
	pvt := pv.Type()
	if pvt.Kind() != reflect.Func {
		panic("provider不是函数")
	}

	// 使用新类型
	return reflect.MakeFunc(pvt, func(args []reflect.Value) []reflect.Value {

		result := pv.Call(args)

		if len(result) == 0 {
			return result
		}

		firstOut := result[0]
		firstOutType := firstOut.Type()

		if !mockType.Implements(firstOutType) {
			panic(fmt.Errorf("mock not Implements interface"))
		}

		// 根据返回值的类型(mock)生成新的类型，其中新类型的方法均加上钩子
		// 注意：生成的不是接口，是实现了接口的类型
		if firstOutType.Kind() == reflect.Interface {

			newValue := reflect.New(mockType.Elem()).Elem()
			newValueType := newValue.Type()

			// field设置
			for i := 0; i < newValueType.NumField(); i++ {
				field := newValue.Field(i)
				fieldType := newValueType.Field(i)

				// 需要写死后缀，感觉不好，但是暂时没想到更好的处理办法
				// 或者可以添加一个method tag，根据这个tag指定的名称来找方法
				var name = fieldType.Name
				for _, suffix := range MockFieldNameSuffixes {
					name = strings.TrimSuffix(name, suffix)
				}

				method := firstOut.MethodByName(name)
				methodType, ok := firstOutType.MethodByName(name)
				if !ok {
					methodTag, ok := fieldType.Tag.Lookup("method")
					if !ok {
						panic(fmt.Errorf("找不到名称对应的方法"))
					}
					fmt.Printf("tag: %+v\n", methodTag)
					name = methodTag

					method = firstOut.MethodByName(name)
					methodType, ok = firstOutType.MethodByName(name)
					if !ok {
						panic(fmt.Errorf("使用tag也找不到名称对应的方法"))
					}
				}
				pctx := ProxyContext{
					PkgPath:       firstOutType.PkgPath(),
					InterfaceName: firstOutType.Name(),
					MethodName:    methodType.Name,
				}

				// 每个方法一个hook，如果同一时间有多个请求呢？
				// 每个请求一个hook，开销可能太大。
				// 并且，如果要每个请求一个hook，那么就需要通过参数传进来，怎么传呢？
				// 在请求进来，ctx初始化时，顺带初始化proxy
				globalHooks := make([]Hook, len(impl.hooks))
				copy(globalHooks, impl.hooks)
				specHooks := make([]Hook, len(hooks))
				copy(specHooks, hooks)

				newMethod := reflect.MakeFunc(methodType.Type, func(args []reflect.Value) []reflect.Value {
					// 执行前钩子
					for hi := len(hooks) - 1; hi >= 0; hi-- {
						specHooks[hi].Before(pctx)
					}
					for hi := len(impl.hooks) - 1; hi >= 0; hi-- {
						globalHooks[hi].Before(pctx)
					}

					result := method.Call(args)

					// 执行后钩子
					for hi := len(hooks) - 1; hi >= 0; hi-- {
						specHooks[hi].After(pctx)
					}
					for hi := len(impl.hooks) - 1; hi >= 0; hi-- {
						globalHooks[hi].After(pctx)
					}

					return result
				})

				field.Set(newMethod)
			}

			result[0] = newValue.Addr().Convert(firstOutType)
		}

		return result
	}).Interface()
}
