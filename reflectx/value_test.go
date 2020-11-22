package reflectx

import (
	"database/sql"
	"math/rand"
	"reflect"
	"testing"
	"time"

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

func init() {
	// 初始化随机种子，否则生成内容不会改变
	rand.Seed(time.Now().Unix())
}

func TestSetRandom(t *testing.T) {
	type Inner struct {
		H int
	}
	for _, cas := range []struct {
		Inner
		A     int
		UA    uint
		B     bool
		C     int64
		UC    uint64
		F     float64
		N     string
		Array [4]int
		S     []int
		D     struct {
			E string
			P *string
		}
		DP *struct {
			M  int
			T  time.Time
			TP *time.Time
			Q  sql.NullTime
			QP *sql.NullTime
		}
		M map[string]struct {
			N string
		}
		Complex64  complex64
		Complex128 complex128
		I          interface {
			String() string
		}
	}{
		{},
	} {
		SetStructRandom(&cas)
		t.Logf("cas: %#v, cas.DP: %+v\n", cas, cas.DP)
	}
}
