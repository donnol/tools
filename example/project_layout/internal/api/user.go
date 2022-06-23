package api

import (
	"context"
	"project_layout/model/request/user"

	"github.com/gin-gonic/gin"
)

func (api *API) ModName() gin.HandlerFunc {
	return jsonHandler(func(ctx context.Context, p *Param) (interface{}, error) {
		var req user.Req
		if err := p.Parse(ctx, &req); err != nil {
			return nil, err
		}

		if err := api.UserSrv.ModName(req.Id.Id, req.Name); err != nil {
			return nil, err
		}

		return nil, nil
	})
}
