package reflectx

import "reflect"

// IsStructPointer 是否结构体指针
func IsStructPointer(v any) bool {
	refType := reflect.TypeOf(v)

	return isStructPointer(refType)
}

func isStructPointer(refType reflect.Type) bool {
	// 是否指针
	if refType.Kind() != reflect.Ptr {
		return false
	}
	// 是否结构体
	refType = refType.Elem()

	return refType.Kind() == reflect.Struct
}
