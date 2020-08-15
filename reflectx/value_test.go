package reflectx

import (
	"reflect"
	"testing"

	"github.com/donnol/tools/reflectx/testdata"
)

func TestSetAnyValue(t *testing.T) {
	ao := &testdata.Ao{}
	db := &testdata.DB{
		Name: "jd",
	}

	var refType reflect.Type
	var specType reflect.Type
	var refValue reflect.Value
	var specValue reflect.Value
	refType = reflect.TypeOf(ao)
	refValue = reflect.ValueOf(ao)
	specType = reflect.TypeOf((*testdata.DB)(nil))
	specValue = reflect.ValueOf(db)

	setAnyValue(refType, specType, refValue, specValue)

	t.Logf("%p\n", db)
	t.Logf("%+v\n", ao)
}
