package common

type Id struct {
	Id uint `json:"id,string"`
}

type Page struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type Result struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Timestamp int64  `json:"timestamp,omitempty"`
	TraceId   string `json:"traceId"`

	Data interface{} `json:"data"`
}

// Fill 使用r里除success之外的非空字段填充res并返回
func (res Result) Fill(r Result) Result {
	if r.Data != nil {
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
