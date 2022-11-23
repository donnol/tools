package basetype

import "testing"

type result[T any] struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Timestamp int64  `json:"timestamp,omitempty"`
	TraceId   string `json:"traceId"`

	Data T `json:"data"`
}

// Fill 使用r里的非空字段填充res并返回
func (res result[T]) Fill(r result[T]) result[T] {
	if any(r.Data) != nil {
		res.Data = r.Data
	}

	if r.Code != 0 {
		res.Code = r.Code
	}
	if r.Msg != "" {
		res.Msg = r.Msg
	}
	if r.Timestamp != 0 {
		res.Timestamp = r.Timestamp
	}
	if r.TraceId != "" {
		res.TraceId = r.TraceId
	}

	return res
}

func TestTypeParameterZero(t *testing.T) {
	var r result[int]
	var r1 result[int] = result[int]{Data: 1}
	r = r.Fill(r1)
	if r.Data != 1 {
		t.Errorf("got %d != %d", r.Data, 1)
	}

	// if any(r.Data) != nil {
	// 	res.Data = r.Data
	// }
	// 如果使用any(T) != nil来判断，因为传进去的是int类型(导致接口的类型部分不为nil)，也会导致值被覆盖
	var r2 result[int] = result[int]{Data: 1}
	var r3 = result[int]{Data: 0}
	r2 = r2.Fill(r3) // Data字段值变为了0
	if r2.Data != 0 {
		t.Errorf("got %d != %d", r2.Data, 0)
	}
}
