package apitest

import (
	"encoding/json"
	"testing"
)

type inner struct {
	Phone string `json:"phone"` // 手机
}

type testModel struct {
	Name string `json:"name"` // 名称
	List []struct {
		Id   uint   `json:"id"`   // id
		Name string `json:"name"` // 名字
		Age  int    `json:"age"`  // 年龄
		Addr struct {
			City string `json:"city"` // 城市
			Home string `json:"home"` // 家
		} `json:"addr"` // 地址

		inner
	} `json:"list"` // 用户列表
}

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

func Test_dataToSummary(t *testing.T) {
	tm := testModel{
		Name: "users",
		List: []struct {
			Id   uint   "json:\"id\""
			Name string "json:\"name\""
			Age  int    "json:\"age\""
			Addr struct {
				City string "json:\"city\""
				Home string "json:\"home\""
			} "json:\"addr\""
			inner
		}{
			{
				Id:   1,
				Name: "jd",
				Age:  20,
				Addr: struct {
					City string "json:\"city\""
					Home string "json:\"home\""
				}{
					City: "gd",
					Home: "gz",
				},
				inner: inner{
					Phone: "123908",
				},
			},
		},
	}
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
            "id": 1, // id
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
