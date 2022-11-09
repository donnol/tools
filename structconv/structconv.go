package structconv

import (
	"fmt"
	"reflect"
)

func ConvByFunc[F, T any](from F, conv func(F) T) T {
	return ConvSliceByFunc([]F{from}, conv)[0]
}

func ConvSliceByFunc[F, T any](from []F, conv func(F) T) []T {
	r := make([]T, len(from))
	for i, f := range from {
		r[i] = conv(f)
	}
	return r
}

// ConvSliceByFieldName use reflect to fill 'to' slice with from by name
func ConvSliceByFieldName[F, T any](from []F, to []T) {
	for i, f := range from {
		ConvByFieldName(f, to[i])
	}
}

// ConvByFieldName use reflect to fill 'to' with from by name
// like to.Name = from.Name
// to must be a struct pointer
func ConvByFieldName[F, T any](from F, to T) {
	toByFieldNameReflect(from, to)
}

var (
	emptyValue = reflect.Value{}
)

func toByFieldNameReflect[F, T any](from F, to T) {
	fromValue := reflect.ValueOf(from)
	if fromValue.Type().Kind() == reflect.Pointer {
		fromValue = fromValue.Elem()
	}

	toValue := reflect.ValueOf(to)
	toType := toValue.Type()
	if toType.Kind() != reflect.Pointer {
		panic(fmt.Errorf("to is not a pointer"))
	}
	toElemValue := toValue.Elem()
	toElemType := toType.Elem()
	if toElemType.Kind() != reflect.Struct {
		panic(fmt.Errorf("to is not a struct"))
	}

	for i := 0; i < toElemType.NumField(); i++ {
		field := toElemType.Field(i)
		fieldValue := toElemValue.Field(i)

		fromFieldValue := fromValue.FieldByName(field.Name)
		if fromFieldValue == emptyValue {
			continue
		}
		if fromFieldValue.Type().Kind() != fieldValue.Type().Kind() {
			continue
		}

		fieldValue.Set(fromFieldValue)
	}
}

// MakeSliceWithNew 以from的长度新建to，并为元素分配内存
// 因为make([]*T, len(from))只会把元素都设为nil
func MakeSliceWithNew[F, T any](from []F) (to []*T) {
	for range from {
		t := new(T)
		to = append(to, t)
	}
	return
}
