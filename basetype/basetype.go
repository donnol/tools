package basetype

type Pager struct {
	Page     int `json:"page"`     // 分页页数
	PageSize int `json:"pageSize"` // 分页大小
}
