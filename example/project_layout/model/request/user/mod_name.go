package user

import "project_layout/model/request/common"

type Req struct {
	common.Id
	Name string `json:"name"`
}

type Resp struct {
}
