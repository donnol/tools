package inject

import (
	"fmt"
	"reflect"
	"unsafe"
)

// Ioc 控制反转，Inversion of Control
type Ioc struct {
	enableUnexportedFieldSetValue bool // 开启对非导出字段的赋值

	providerMap map[reflect.Type]typeInfo
	cache       map[reflect.Type]reflect.Value
}

type typeInfo struct {
	depType  []reflect.Type
	provider reflect.Value
}

func NewIoc(
	enableUnexportedFieldSetValue bool,
) *Ioc {
	return &Ioc{
		enableUnexportedFieldSetValue: enableUnexportedFieldSetValue,
		providerMap:                   make(map[reflect.Type]typeInfo),
		cache:                         make(map[reflect.Type]reflect.Value),
	}
}

// RegisterProvider 注册provider函数，形如`func New(fielda TypeA, fieldb TypeB) (T)`
func (ioc *Ioc) RegisterProvider(v any) (err error) {
	refValue := reflect.ValueOf(v)
	refType := refValue.Type()
	if refType.Kind() != reflect.Func {
		return fmt.Errorf("Please input func")
	}

	// 分析函数的参数和返回值
	ti := typeInfo{
		depType:  make([]reflect.Type, 0, refType.NumIn()),
		provider: refValue,
	}
	for i := 0; i < refType.NumIn(); i++ {
		in := refType.In(i)
		ti.depType = append(ti.depType, in)
	}
	// 返回：instance
	min := 1
	if refType.NumOut() == 0 {
		return fmt.Errorf("can't find result in func")
	}
	if refType.NumOut() < min {
		return fmt.Errorf("too little result in func, min is %d", min)
	}
	for i := 0; i < refType.NumOut(); i++ {
		out := refType.Out(i)
		ioc.providerMap[out] = ti
	}

	return
}

// Inject 注入依赖，传入结构体，根据结构体的字段类型找到对应的provider，执行后将获得的值赋予字段
//
// 遍历v的字段，找到字段类型，再根据字段类型找到provider，调用provider获得实例，再把实例值赋给该字段
// provider需要在接口定义处注册，注册到一个统一管理的地方
// 如果provider需要参数，则根据参数类型继续找寻相应的provider，直至初始化完成
func (ioc *Ioc) Inject(v any) (err error) {
	refValue := reflect.ValueOf(v)
	refType := refValue.Type()
	if refType.Kind() != reflect.Ptr {
		return fmt.Errorf("v is not a pointer")
	}
	eleValue := refValue.Elem()
	eleType := eleValue.Type()
	if eleType.Kind() != reflect.Struct {
		return fmt.Errorf("v is not a struct")
	}

	// 遍历field
	for i := 0; i < eleValue.NumField(); i++ {
		field := eleValue.Field(i)
		fieldType := field.Type()

		// 根据类型查找值
		var value reflect.Value
		value, err = ioc.find(fieldType)
		if err != nil {
			return
		}

		// 给字段赋值
		if ioc.enableUnexportedFieldSetValue {
			rf := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
			rf.Set(value)
		} else {
			field.Set(value)
		}
	}

	return
}

var (
	emptyStruct         = reflect.TypeOf((*struct{})(nil))
	emptyStructValue    = reflect.New(emptyStruct.Elem()).Elem()
	emptyStructPtrValue = reflect.New(emptyStruct).Elem()

	errorType = reflect.TypeOf((*error)(nil)).Elem()
)

func (ioc *Ioc) find(typ reflect.Type) (r reflect.Value, err error) {
	value, ok := ioc.cache[typ]
	if ok {
		return value, nil
	}

	// 在provider里寻找初始化函数
	provider, ok := ioc.providerMap[typ]
	if !ok {
		// 检查类型是否是struct{}
		if typ.ConvertibleTo(emptyStruct.Elem()) {
			return emptyStructValue, nil
		}
		if typ.ConvertibleTo(emptyStruct) {
			return emptyStructPtrValue, nil
		}
		return r, fmt.Errorf("can't find provider of %+v", typ)
	}

	// 调用
	in := make([]reflect.Value, 0, len(provider.depType))
	for _, dep := range provider.depType {
		// 在已有provider里找指定类型
		var value reflect.Value
		if value, ok = ioc.cache[dep]; !ok {
			value, err = ioc.find(dep)
			if err != nil {
				return r, err
			}
			ioc.cache[dep] = value
		}

		in = append(in, value)
	}
	newValues := provider.provider.Call(in)
	if len(newValues) == 0 {
		return r, fmt.Errorf("can't get new value by provider")
	}

	// 返回值里，第一个必须是实例，最后一个必须是error，中间的忽略
	newValue := newValues[0]
	ioc.cache[typ] = newValue

	if len(newValues) > 1 {
		lastValue := newValues[len(newValues)-1]
		last := lastValue.Interface()
		if lastValue.Type().Implements(errorType) {
			if v, ok := last.(error); ok {
				if v != nil {
					return r, fmt.Errorf("call failed, err is %+v", v)
				}
			}
		} else {
			return r, fmt.Errorf("last return value is not error, is %+v", last)
		}
	}

	return newValue, nil
}
