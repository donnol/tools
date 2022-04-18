package api

import (
	"net/http"
	"project_layout/model/request/common"
	"project_layout/model/request/user"

	"github.com/gin-gonic/gin"
)

func (api *API) ModName() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req user.Req
		if err := ctx.ShouldBindJSON(&req); err != nil {
			return
		}

		res := api.UserSrv.ModName(req.Id.Id, req.Name)
		ctx.JSON(http.StatusOK, common.Result{
			Data: res,
		})
	}
}
