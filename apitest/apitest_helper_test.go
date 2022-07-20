package apitest

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/donnol/tools/apitest/testtype"
)

var (
	tm = testtype.TestModel{
		Name: "users",
		List: []testtype.User{
			{
				Id:   1,
				Name: "jd",
				Age:  20,
				Addr: testtype.Addr{
					City: "gd",
					Home: "gz",
				},
				Inner: testtype.Inner{
					Phone: "123908",
				},
			},
		},
	}
)

var (
	kcm = map[string]string{
		"|name":           "名称",
		"|list":           "用户列表",
		"|list|id":        "id",
		"|list|name":      "名字",
		"|list|age":       "年龄",
		"|list|addr":      "地址",
		"|list|addr|city": "城市",
		"|list|addr|home": "家",
		"|list|phone":     "手机",
	}
)

func TestStructToBlock(t *testing.T) {
	line, lkcm, err := structToBlock(paramName, http.MethodGet, &testtype.TestModel{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("line:\n%s", line)

	for k, v := range kcm {
		lv, ok := lkcm[k]
		if !ok {
			t.Fatalf("cant find %s in local kcm", k)
		}
		if lv != v {
			t.Fatalf("compare value failed: %s != %s", lv, v)
		}
	}
}

func Test_dataToSummary(t *testing.T) {

	data, err := json.Marshal(tm)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		name   string
		data   []byte
		isJSON bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "return", args: args{
			name:   "return",
			data:   data,
			isJSON: true,
		}, want: `<details>
<summary>return</summary>

` + "```" + `json
{
    "name": "users", // 名称
    "list": [ // 用户列表
        {
            "id": "1", // id
            "name": "jd", // 名字
            "age": 20, // 年龄
            "addr": { // 地址
                "city": "gd", // 城市
                "home": "gz" // 家
            },
            "phone": "123908" // 手机
        }
    ]
}
` + "```" + `

</details>

`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dataToSummary(tt.args.name, tt.args.data, tt.args.isJSON, kcm); got != tt.want {
				t.Errorf("dataToSummary() = %v, want %v", got, tt.want)
			}
		})
	}
}
