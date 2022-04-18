package api

import (
	"net/http"
	"project_layout/model/request/common"

	"github.com/gin-gonic/gin"
)

func (api *API) Ping() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res := api.PingSrv.Ping()
		ctx.JSON(http.StatusOK, common.Result{
			Data: res,
		})
	}
}
