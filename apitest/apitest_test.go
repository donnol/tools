package apitest

import (
	"os"
	"reflect"
	"testing"

	"github.com/donnol/tools/reflectx"
)

var (
	testStruct = reflectx.User{
		Model: &reflectx.Model{
			InnerModel: reflectx.InnerModel{
				CornerStone: "jd",
			},
		},
		Basic: reflectx.Basic{
			Name: "jd",
		},
		Age: 18,
		AddressList: []reflectx.Address{
			{
				Basic: reflectx.Basic{
					Name: "tianhe",
				},
				Position: "guangdong tianhe",
			},
		},
	}
)

func TestStructToMap(t *testing.T) {
	if m, err := structToMap(testStruct); err != nil {
		t.Fatal(err)
	} else {
		JSONIndent(os.Stdout, m)
	}
}

func TestRandomValue(t *testing.T) {
	for _, cas := range []struct {
		Kind reflect.Kind
		L    int
	}{
		{reflect.Int, 0},
		{reflect.String, 5},
		{reflect.Bool, 0},
		{reflect.Float64, 0},
	} {
		r := randomValue(cas.Kind, cas.L)
		t.Logf("%s: %+v\n", cas.Kind, r)
	}
}

func TestStructOf(t *testing.T) {
	r := structOf([]reflect.StructField{
		{
			Name: "Height",
			Type: reflect.TypeOf(float64(0)),
			Tag:  `json:"height"`,
		},
		{
			Name: "Age",
			Type: reflect.TypeOf(int(0)),
			Tag:  `json:"age"`,
		},
	})
	t.Logf("r: %+v\n", r)
	JSONIndent(os.Stdout, r)
}

func TestStructRandomValue(t *testing.T) {
	r, err := structRandomValue(testStruct)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("r: %+v\n", r)
	JSONIndent(os.Stdout, r)
}

func TestCompositeStructValue(t *testing.T) {
	v := compositeStructValue(reflect.TypeOf(reflectx.Model{}))
	t.Logf("v: %+v\n", v)
	JSONIndent(os.Stdout, v.Interface())
}

func TestRandomValueByTag(t *testing.T) {
	for _, cas := range []struct {
		Tag string
		Max int
	}{
		{"range=one(1,13)", 13},
		{"range=one(1.1,13.4)", 14},
		{"range=one(1,5,8,13)", 13},
		{"enum=one(nil,1,2,3)", 4},
		{"enum=one(a,b,c)", 4},
		{"enum=one(true,false)", 4},
		{"enum=many(1,2,3)", 4},
		{"enum=many(a,b,c)", 4},
		{"call=year(2018)", 0},
		{"call=month(2018,1)", 0},
		{"call=day(2018,1,1)", 0},
		// {"db=one(org,id)", 0},
		// {"db=many(org,id)", 0},
		// {`regexp=one("^[a-z]+$")`, 0},
	} {
		v := randomValueByTag(cas.Tag)
		if reflect.DeepEqual(v, cas.Max) {
			t.Fatalf("Bad v: %v\n", v)
		}
		t.Logf("v: %+v\n", v)
	}
}

func TestMakeFunc(t *testing.T) {
	makeFunc()
}

func TestCollectStructField(t *testing.T) {
	typ := reflect.TypeOf(testStruct)
	sf := collectStructField(typ)
	t.Logf("sf: %+v\n", sf)
}
