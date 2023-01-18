package reflectx

import (
	"database/sql"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"sync"
	"time"
	"unsafe"

	"github.com/donnol/tools/internal/utils/debug"
)

var (
	// ErrParamNotStruct 参数不是结构体
	ErrParamNotStruct = fmt.Errorf("please input struct param")
)

// InitParam 初始化-使用反射初始化param里的指定类型
func InitParam(param any, specType reflect.Type, specValue reflect.Value, copy bool) (any, error) {
	// 反射获取type和value
	refType := reflect.TypeOf(param)
	refValue := reflect.ValueOf(param)
	if refType.Kind() == reflect.Ptr {
		refType = refType.Elem()
		refValue = refValue.Elem()
	}
	if refType.Kind() != reflect.Struct {
		return param, ErrParamNotStruct
	}

	// 创建副本
	if copy {
		var sf = make([]reflect.StructField, 0)
		for i := 0; i < refType.NumField(); i++ {
			field := refType.Field(i)

			sf = append(sf, field)
		}
		newType := reflect.StructOf(sf)
		newValue := reflect.New(refType)

		// 给value赋值
		newValueElem := newValue.Elem()
		for i := 0; i < refType.NumField(); i++ {
			oldV := refValue.Field(i)
			newV := newValueElem.Field(i)
			newV.Set(oldV)
		}

		// 替换
		refType = newType
		refValue = newValue
	}

	// 注入value
	setValue(refType, specType, refValue, specValue)

	// 返回副本
	if copy {
		return refValue.Interface(), nil
	}

	return param, nil
}

func setValue(refType, specType reflect.Type, refValue, specValue reflect.Value) {
	// 忽略非结构体或者time.Time类型
	if refType.Kind() != reflect.Struct ||
		refType == reflect.TypeOf((*time.Time)(nil)).Elem() {
		return
	}

	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)

		// 忽略非导出字段-因为无法对非导出字段赋值
		if field.PkgPath != "" {
			continue
		}

		// 获取对应字段的value
		var value reflect.Value
		if refValue.Type().Kind() == reflect.Ptr {
			value = refValue.Elem().Field(i)
		} else {
			value = refValue.Field(i)
		}

		// 根据字段type判断是否可以赋值
		if field.Type == specType { // 类型相同，直接赋值
			value.Set(specValue)
		} else { // 匿名内嵌或者包含在普通字段里，继续对该字段类型遍历
			setValue(field.Type, specType, value, specValue)
		}
	}
}

// setAnyValue 设置任意结构体字段值，无论是导出还是非导出
// From https://stackoverflow.com/questions/42664837/access-unexported-fields-in-golang-reflect
func setAnyValue(refType, specType reflect.Type, refValue, specValue reflect.Value) {
	// 反射获取type和value
	if refType.Kind() == reflect.Ptr {
		refType = refType.Elem()
		refValue = refValue.Elem()
	}
	if refType.Kind() != reflect.Struct {
		return
	}

	// 忽略非结构体或者time.Time类型
	if refType.Kind() != reflect.Struct ||
		refType == reflect.TypeOf((*time.Time)(nil)).Elem() {
		return
	}

	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)

		// 获取对应字段的value
		var value reflect.Value
		if refValue.Type().Kind() == reflect.Ptr {
			value = refValue.Elem().Field(i)
		} else {
			value = refValue.Field(i)
		}

		// 根据字段type判断是否可以赋值
		if field.Type == specType { // 类型相同，直接赋值
			setStructFieldValue(value, specValue, field.PkgPath != "")
		} else { // 匿名内嵌或者包含在普通字段里，继续对该字段类型遍历
			setAnyValue(field.Type, specType, value, specValue)
		}
	}
}

// CopyStructField 复制结构体字段
func CopyStructField(refType reflect.Type, refValue reflect.Value) (reflect.Type, reflect.Value) {
	// 新建一个同类型的Value
	newType := refType
	newValue := reflect.New(newType)

	// 复制原有值到Value
	newValueElem := newValue.Elem()
	for i := 0; i < refType.NumField(); i++ {
		fieldType := refType.Field(i)
		isUnexportedField := fieldType.PkgPath != ""

		newV := newValueElem.Field(i)

		// panic: reflect: reflect.Value.Set using value obtained using unexported field
		// 从非导出字段获取的reflect.Value不能被用来给另外的value赋值
		// oldV := refValue.Field(i)
		//
		// setStructFieldValue(newV, oldV, fieldType.PkgPath != "")

		// 使用unsafe读取
		oldV := refValue.Field(i)
		if isUnexportedField {
			oldPtr := unsafe.Pointer(oldV.UnsafeAddr())
			oldV = reflect.NewAt(oldV.Type(), oldPtr).Elem()
		}

		// 给字段赋值
		setStructFieldValue(newV, oldV, isUnexportedField)
	}

	return newType, newValue
}

// SetStructRandom 设置结构体随机值
func SetStructRandom(v any) {
	var refValue reflect.Value
	var refType reflect.Type

	if vv, ok := v.(reflect.Value); ok {
		refValue = vv.Elem()
	} else if !IsStructPointer(v) {
		panic(fmt.Errorf("v is not a struct pointer"))
	} else {
		refValue = reflect.ValueOf(v).Elem()
	}

	refType = refValue.Type()

	for i := 0; i < refType.NumField(); i++ {
		fieldValue := refValue.Field(i)
		fieldType := refType.Field(i)

		value := getTypeRandomValue(fieldType.Type)

		fieldValue.Set(value)
	}
}

var (
	timeType        = reflect.TypeOf((*time.Time)(nil)).Elem()
	sqlNullTimeType = reflect.TypeOf((*sql.NullTime)(nil)).Elem()
)

func getTimeTypeValue(typ reflect.Type) reflect.Value {
	var value reflect.Value

	switch typ {
	case timeType:
		t := getRandomTime()
		value = reflect.ValueOf(t)
	case sqlNullTimeType:
		nt := sql.NullTime{}
		b := getRandomBool()
		nt.Valid = b
		if b {
			t := getRandomTime()
			nt.Time = t
		}
		value = reflect.ValueOf(nt)
	}

	return value
}

func getRandomTime() time.Time {
	y := rand.Intn(3000)
	m := rand.Intn(12) + 1
	d := rand.Intn(31) + 1
	hour := rand.Intn(24)
	min := rand.Intn(60) + 1
	sec := rand.Intn(60) + 1
	nsec := rand.Intn(1000)
	t := time.Date(y, time.Month(m), d, hour, min, sec, nsec, time.Local)
	return t
}

func getRandomBool() bool {
	var value bool
	n := rand.Intn(2)
	if n == 0 {
		value = false
	} else {
		value = true
	}
	return value
}

func toNegative(n int) int {
	b := getRandomBool()
	if b && n != 0 {
		n = -n
	}
	return n
}

type specialType struct {
	typ     reflect.Type
	handler func(typ reflect.Type) reflect.Value
}

var (
	specTypeMap = new(sync.Map)
)

func init() {
	// 注册时间类型
	RegisterSpecType(timeType, getTimeTypeValue)
	RegisterSpecType(sqlNullTimeType, getTimeTypeValue)
}

// RegisterSpecType 注册特别类型
func RegisterSpecType(typ reflect.Type, handler func(typ reflect.Type) reflect.Value) {
	specTypeMap.Store(typ, specialType{
		typ:     typ,
		handler: handler,
	})
}

func inSpecType(typ reflect.Type) (handler func(typ reflect.Type) reflect.Value, ok bool) {
	v, ok := specTypeMap.Load(typ)
	if ok {
		st, ok := v.(specialType)
		if ok {
			return st.handler, ok
		}
	}
	return nil, false
}

func getTypeRandomValue(typ reflect.Type) reflect.Value {
	var value reflect.Value

	// 特殊类型
	specHandler, isSpec := inSpecType(typ)
	if isSpec {
		return specHandler(typ)
	}

	// 常用类型
	switch typ.Kind() {
	case reflect.Invalid:

	case reflect.Bool:
		n := getRandomBool()
		value = reflect.ValueOf(n)

	case reflect.Int:
		n := rand.Int()
		n = toNegative(n)
		value = reflect.ValueOf(n)
	case reflect.Int8:
		n := rand.Intn(128)
		n = toNegative(n)
		value = reflect.ValueOf(int8(n))
	case reflect.Int16:
		n := rand.Intn(32768)
		n = toNegative(n)
		value = reflect.ValueOf(int16(n))
	case reflect.Int32:
		n := rand.Int31()
		n = int32(toNegative(int(n)))
		value = reflect.ValueOf(n)
	case reflect.Int64:
		n := rand.Int63()
		n = int64(toNegative(int(n)))
		value = reflect.ValueOf(n)

	case reflect.Uint:
		n := rand.Int()
		value = reflect.ValueOf(uint(n))
	case reflect.Uint8:
		n := rand.Intn(256)
		value = reflect.ValueOf(uint8(n))
	case reflect.Uint16:
		n := rand.Intn(65536)
		value = reflect.ValueOf(uint16(n))
	case reflect.Uint32:
		n := rand.Int31()
		value = reflect.ValueOf(uint32(n))
	case reflect.Uint64:
		n := rand.Int63()
		value = reflect.ValueOf(uint64(n))

	case reflect.Uintptr:

	case reflect.Float32:
		n := rand.Float32()
		value = reflect.ValueOf(n)
	case reflect.Float64:
		n := rand.Float64()
		value = reflect.ValueOf(n)

	case reflect.Complex64:
		real, imag := rand.Float32(), rand.Float32()
		v := complex(real, imag)
		value = reflect.ValueOf(v)
	case reflect.Complex128:
		real, imag := rand.Float64(), rand.Float64()
		v := complex(real, imag)
		value = reflect.ValueOf(v)

	case reflect.Array:
		// 先建立类型，再用类型新建值
		value = reflect.New(reflect.ArrayOf(typ.Len(), typ.Elem())).Elem()
		for i := 0; i < typ.Len(); i++ {
			arrayValue := getTypeRandomValue(typ.Elem())
			value.Index(i).Set(arrayValue)
		}

	case reflect.Chan:
		// value = reflect.New(reflect.ChanOf(typ.ChanDir(), typ.Elem())).Elem()
		n := rand.Intn(10)
		value = reflect.MakeChan(typ, n)
		// 开个线程往chan发消息？

	case reflect.Func:
		value = reflect.MakeFunc(typ, func(args []reflect.Value) []reflect.Value {
			var values []reflect.Value

			// 我怎么知道要跑什么逻辑呢？
			// 随便打印个随机数吧
			n := rand.Int()
			fmt.Println(n)

			return values
		})

	case reflect.Interface:
		// https://github.com/golang/go/issues/4146
		// 2012年的issue，2020了还没实现，好难受

		// value是个nil
		// 当调用value.Method(i)时会报错：panic: reflect: Method on nil interface value
		debug.Printf("typ: %+v\n", typ)
		value = reflect.New(typ).Elem()
		debug.Printf("value: %+v\n", value)

		// 需要借助mock struct才行，具体可参考inject的proxy
		// for i := 0; i < typ.NumMethod(); i++ {
		//	method := typ.Method(i)

		//	newmethod := getTypeRandomValue(method.Type)

		//	value.Method(i).Set(newmethod)
		// }

	case reflect.Map:
		n := rand.Intn(10)
		value = reflect.MakeMap(typ)
		keyTyp := typ.Key()
		for i := 0; i < n; i++ {
			key := getTypeRandomValue(keyTyp)
			mapValue := getTypeRandomValue(typ.Elem())

			value.SetMapIndex(key, mapValue)
		}

	case reflect.Ptr:
		// 假设typ: *string
		// reflect.New会添加多一层指针
		value = reflect.New(typ) // **string
		debug.Printf("call reflect.New: %+v\n", value.Type())

		// 取Elem解一层指针
		value = value.Elem() // *string
		debug.Printf("call value.Elem: %+v\n", value.Type())

		// 获取原始类型随机值
		e := getTypeRandomValue(typ.Elem()) // string
		debug.Printf("get random value: %+v\n", e.Type())

		// 通过中间value，给原变量赋指针值
		tv := reflect.New(typ.Elem()).Elem() // 新建原始类型指针: string
		debug.Printf("call reflect.New to new origin value: %+v\n", tv.Type())

		tv.Set(e) // 设置值

		// 将中间变量地址赋给结果
		value.Set(tv.Addr())

	case reflect.Slice:
		n := rand.Intn(64)
		value = reflect.MakeSlice(typ, 0, n)
		for i := 0; i < n; i++ {
			value = reflect.Append(value, getTypeRandomValue(typ.Elem()))
		}

	case reflect.String:
		s := ""
		n := rand.Intn(64)
		for i := 0; i < n; i++ {
			j := rand.Int()
			s += strconv.Itoa(j)
		}
		value = reflect.ValueOf(s)

	case reflect.Struct:
		s := reflect.New(typ)
		SetStructRandom(s)
		value = s.Elem()

	case reflect.UnsafePointer:

	}

	return value
}

// ToInterface 如果in数组里存在无法取Interface的Value，则会置为nil
func ToInterface(in []reflect.Value) (out []any) {
	out = make([]any, len(in))
	for i, one := range in {
		if !one.CanInterface() {
			continue
		}
		out[i] = one.Interface()
	}
	return
}

func FromAny(in []any) []reflect.Value {
	out := make([]reflect.Value, len(in))
	for i, one := range in {
		v := reflect.ValueOf(one)
		out[i] = v
	}
	return out
}

// SetStructFieldValue 设置结构体字段的值，arg必须是结构体指针，value值类型必须与结构体对应fieldName类型一致
// 没有没有找到指定字段，则返回的StructField为空
func SetStructFieldValue(arg reflect.Value, fieldName string, value any) (fieldType reflect.StructField) {
	var refValue = arg
	var refType = arg.Type()

	if !isStructPointer(refType) {
		panic(fmt.Errorf("v is not a struct pointer"))
	}

	refValue = refValue.Elem()
	refType = refValue.Type()

	targetValue := reflect.ValueOf(value)
	for i := 0; i < refType.NumField(); i++ {
		fieldValue := refValue.Field(i)
		tmpFieldType := refType.Field(i)

		// 内嵌
		if tmpFieldType.Anonymous {
			// 新建同类型value，并将原值复制到新变量rf
			_, rf := CopyStructField(fieldValue.Type(), fieldValue)

			// 递归找字段，并赋值
			ft := SetStructFieldValue(rf, fieldName, value)

			// 字段存在，给字段赋值
			if ft.Name != "" {
				setStructFieldValue(fieldValue, rf.Elem(), true)
			}

			continue
		}

		if tmpFieldType.Name != fieldName {
			continue
		}

		fieldType = tmpFieldType
		setStructFieldValue(fieldValue, targetValue, tmpFieldType.PkgPath != "")
	}

	return
}

func setStructFieldValue(fieldValue reflect.Value, targetValue reflect.Value, isUnexportedField bool) {
	if isUnexportedField {
		// 非导出字段，使用unsafe赋值
		rf := reflect.NewAt(fieldValue.Type(), unsafe.Pointer(fieldValue.UnsafeAddr())).Elem()
		rf.Set(targetValue)
	} else {
		fieldValue.Set(targetValue)
	}
}
