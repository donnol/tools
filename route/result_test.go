package route

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/donnol/tools/errors"
)

type m struct {
	A int `json:"a"`
}

func TestResult(t *testing.T) {
	for i, cas := range []struct {
		data Result
		want string
	}{
		{
			data: Result{
				Error: errors.Error{
					Code: 100,
					Msg:  "failed",
				},
				Data: m{
					A: 10,
				},
			},
			want: `{"code":100,"msg":"failed","data":{"a":10},"timestamp":0,"requestID":""}`,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			r, err := json.Marshal(cas.data)
			if err != nil {
				t.Fatal(err)
			}
			if string(r) != cas.want {
				t.Fatalf("Bad result: %s != %s\n", r, cas.want)
			}
		})
	}
}

func TestResult1(t *testing.T) {
	SetResultVersion(1)

	for i, cas := range []struct {
		data Result
		want string
	}{
		{
			data: Result{
				Error: errors.Error{
					Code: 100,
					Msg:  "failed",
				},
				Data: m{
					A: 10,
				},
			},
			want: `{"header":{"code":100,"message":"failed","count":1},"data":{"a":10}}`,
		},
		{
			data: Result{
				Error: errors.Error{
					Code: 100,
					Msg:  "failed",
				},
				Data: m{},
			},
			want: `{"header":{"code":100,"message":"failed","count":0},"data":{"a":0}}`,
		},
		{
			data: Result{
				Error: errors.Error{
					Code: 100,
					Msg:  "failed",
				},
				Data: []m{{A: 1}, {A: 2}},
			},
			want: `{"header":{"code":100,"message":"failed","count":2},"data":[{"a":1},{"a":2}]}`,
		},
		{
			data: Result{
				Error: errors.Error{
					Code: 100,
					Msg:  "failed",
				},
				Data: []m{},
			},
			want: `{"header":{"code":100,"message":"failed","count":0},"data":[]}`,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			r, err := json.Marshal(cas.data)
			if err != nil {
				t.Fatal(err)
			}
			if string(r) != cas.want {
				t.Fatalf("Bad result: %s != %s\n", r, cas.want)
			}
		})
	}
}
