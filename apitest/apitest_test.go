package apitest

import (
	"encoding/json"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/donnol/tools/reflectx"
	"github.com/getkin/kin-openapi/openapi2"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
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

func TestFakeStruct(t *testing.T) {
	var user reflectx.User
	if err := gofakeit.Struct(&user); err != nil {
		t.Fatal(err)
	}

	t.Logf("user with random value: %+v\n", user)
}

func TestWriteFile(t *testing.T) {
	f, err := OpenFile("testdata/user.md", "用户接口文档")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	for _, tc := range []struct {
		path    string
		method  string
		comment string
	}{
		{path: "/api/user", method: http.MethodGet, comment: "获取用户信息"},
		{path: "/api/user", method: http.MethodPost, comment: "添加用户信息"},
	} {
		t.Run(tc.method+"-"+tc.path, func(t *testing.T) {
			at := NewAT(tc.path, tc.method, tc.comment, nil, nil)
			var res = struct {
				UserId uint   `json:"userId"`
				Name   string `json:"name"`
			}{
				UserId: 1,
				Name:   "jd",
			}
			if err := at.SetParam(&struct {
				UserId uint   `json:"userId"`
				Name   string `json:"name"`
			}{UserId: 1, Name: "jd"}).
				FakeRun().
				Result(&res).
				WriteFile(f).
				Err(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestToSwagger(t *testing.T) {
	f, err := os.OpenFile("testdata/swagger_test.json", os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	at := NewAT("/api/user", http.MethodGet, "获取用户信息", nil, nil)
	var res = struct {
		UserId uint   `json:"userId"`
		Name   string `json:"name"`
	}{
		UserId: 1,
		Name:   "jd",
	}

	if err := at.SetParam(&struct {
		UserId uint   `json:"userId"`
		Name   string `json:"name"`
	}{UserId: 1, Name: "jd"}).
		FakeRun().
		Result(&res).
		toSwagger(f).
		Err(); err != nil {
		t.Fatal(err)
	}
}

func TestSwaggerJSON(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		data, err := os.ReadFile("./testdata/swagger.json")
		if err != nil {
			t.Fatal(err)
		}

		var doc openapi2.T
		if err := json.Unmarshal(data, &doc); err != nil {
			t.Fatal(err)
		}

		for _, path := range doc.Paths {
			t.Logf("doc path: %+v\n", path)
		}

		data2, err := json.Marshal(doc)
		if err != nil {
			t.Fatal(err)
		}
		JSONIndent(os.Stdout, data2)
	})

	t.Run("yaml", func(t *testing.T) {
		data, err := os.ReadFile("./testdata/swagger.json")
		if err != nil {
			t.Fatal(err)
		}

		var doc openapi2.T
		if err := json.Unmarshal(data, &doc); err != nil {
			t.Fatal(err)
		}

		data2, err := yaml.Marshal(doc)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile("./testdata/swagger.yaml", data2, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestStartSwaggerServer(t *testing.T) {
	startSwaggerServer(t)
}

func startSwaggerServer(t *testing.T) {
	r := gin.Default()

	// 配置本地生成的swagger.json
	r.StaticFile("/swaggerjson", "./testdata/swagger_test.json")

	// 从项目`github.com/swagger-api/swagger-ui`复制dist目录
	// 修改swagger-initializer.js文件里的url为上述/swaggerjson路径
	r.Static("/swaggerui/", "./testdata/swaggerui/dist/")

	port := ":8000"
	t.Logf("listen %v", port)
	if err := http.ListenAndServe(port, r); err != nil {
		t.Fatal(err)
	}
}
