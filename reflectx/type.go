package reflectx

import "reflect"

// IsStructPointer 是否结构体指针
func IsStructPointer(v interface{}) bool {
	refType := reflect.TypeOf(v)
	// 是否指针
	if refType.Kind() != reflect.Ptr {
		return false
	}
	// 是否结构体
	refType = refType.Elem()
	if refType.Kind() != reflect.Struct {
		return false
	}

	return true
}
