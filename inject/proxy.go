package inject

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/donnol/tools/internal/utils/debug"
)

// Proxy 在层间依赖调用时插入钩子调用，类似AOP
type Proxy interface {
	Around(provider any, mock any, arounder Arounder) any
}

func NewProxy() Proxy {
	return &proxyImpl{}
}

// 每个包、每个接口、每个方法唯一对应一个方法
type ProxyContext struct {
	PkgPath       string
	InterfaceName string
	MethodName    string
}

func (pctx ProxyContext) String() string {
	return fmt.Sprintf(pctx.bracket("PkgPath: %s InterfaceName: %s MethodName: %s"), pctx.PkgPath, pctx.InterfaceName, pctx.MethodName)
}

func (pctx ProxyContext) Logf(format string, args ...any) {
	pctx.logf(pctx.String()+": "+format, args...)
}

func (pctx ProxyContext) LogShortf(format string, args ...any) {
	pctx.logf(pctx.bracket(pctx.MethodName)+": "+format, args...)
}

func (pctx ProxyContext) logf(format string, args ...any) {
	err := log.Output(3, fmt.Sprintf(format, args...))
	if err != nil {
		fmt.Printf("Output failed: %+v\n", err)
	}
}

func (pctx ProxyContext) bracket(s string) string {
	return "[" + s + "]"
}

type Arounder interface {
	Around(pctx ProxyContext, method reflect.Value, args []reflect.Value) []reflect.Value
}

type AroundFunc func(pctx ProxyContext, method reflect.Value, args []reflect.Value) []reflect.Value

func (around AroundFunc) Around(pctx ProxyContext, method reflect.Value, args []reflect.Value) []reflect.Value {
	return around(pctx, method, args)
}

type ArounderMap map[ProxyContext]AroundFunc

func (m ArounderMap) Merge(nm ArounderMap) (mr ArounderMap) {
	mr = m
	for k, v := range nm {
		mr[k] = v
	}

	return
}

type proxyImpl struct {
}

var (
	MockFieldNameSuffixes = [...]string{"Func", "Handler"} // mock结构体字段名称后缀
)

func (impl *proxyImpl) Around(provider any, mock any, arounder Arounder) any {
	return impl.around(provider, mock, arounder)
}

func (impl *proxyImpl) around(provider any, mock any, arounder Arounder) any {
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

	// 使用新的类型一样的函数
	// 在注入的时候会被调用
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
					debug.Debug("tag: %+v\n", methodTag)
					name = methodTag

					method = firstOut.MethodByName(name)
					methodType, ok = firstOutType.MethodByName(name)
					if !ok {
						panic(fmt.Errorf("使用tag也找不到名称对应的方法"))
					}
				}
				debug.Debug("method: %+v\n", method)

				pctx := ProxyContext{
					PkgPath:       firstOutType.PkgPath(),
					InterfaceName: firstOutType.Name(),
					MethodName:    methodType.Name,
				}
				debug.Debug("pctx: %+v\n", pctx)

				// newMethod会在实际请求时被调用
				// 当被调用时，newMethod内部就会调用绑定好的Arounder，然后将原函数method和参数args传入
				// 在Around方法执行完后即可获得结果
				newMethod := reflect.MakeFunc(methodType.Type, func(args []reflect.Value) []reflect.Value {
					var result []reflect.Value

					debug.Debug("args: %+v\n", args)

					// Around是对整个结构的统一包装，如果需要对不同方法做不同处理，可以根据pctx里的方法名在Around接口的实现里做处理
					result = arounder.Around(pctx, method, args)

					debug.Debug("result: %+v\n", result)

					return result
				})

				field.Set(newMethod)
			}

			result[0] = newValue.Addr().Convert(firstOutType)
		}

		return result
	}).Interface()
}
