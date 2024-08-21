package controller

import (
	"net/http"
	"service/protos/Common"
	"service/utils/ginResult"
	"github.com/gin-gonic/gin"
)

func Echo(ctx *gin.Context) {
	reply := ginResult.CommonResult{}

	reply.Code = Common.ErrorCodes_SUCCESS
	ctx.JSON(http.StatusOK, ginResult.GetGinResult(reply))
}
