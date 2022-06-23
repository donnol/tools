package api

import (
	"project_layout/internal/service"
)

type API struct {
	PingSrv service.PingSrv
	UserSrv service.UserSrv

	// TODO: add more field...
}
