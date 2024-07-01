package ginx

import (
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

var L logger.LoggerV1 = logger.NewNopLogger()

func WrapBody[Req any](
	bizFn func(ctx *gin.Context, req Req) (Result, error),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.Bind(&req); err != nil {
			L.Error("输入错误", logger.Error(err))
			return
		}
		L.Debug("输入参数", logger.Field{Key: "req", Val: req})
		res, err := bizFn(ctx, req)
		if err != nil {
			L.Error("业务逻辑执行失败", logger.Error(err))
		}

		ctx.JSON(http.StatusOK, res)
	}
}
