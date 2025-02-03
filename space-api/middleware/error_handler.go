package middleware

import (
	"net/http"
	"space-api/util"

	"github.com/gin-gonic/gin"
)

func UseErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 具体的错误捕获处理
		defer func() {
			isPanic := false
			// 捕获 panic
			err := recover()
			if err != nil {
				isPanic = true
			} else {
				if len(ctx.Errors) != 0 {
					err = ctx.Errors[0].Err
				}
			}

			if err == nil {
				return
			}

			code := util.TernaryExpr(isPanic, http.StatusInternalServerError, http.StatusOK)
			switch err := err.(type) {
			case error: /* 确保都是实现 error 接口的结构体的引用 */
				switch err := err.(type) {
				case *util.BizErr:
					ctx.JSON(code, util.RestWithError(err.Error()))
				case *util.LimitErr:
					ctx.JSON(code, util.RestWithError("请求过于频繁: "+err.Error()))
				case *util.VerifyErr:
					ctx.JSON(code, util.RestWithError("参数不正确: "+err.Error()))
				case *util.InnerErr:
					ctx.JSON(code, util.RestWithError("服务内部错误, 请稍后重试或联系站长修复"))
				default:
					ctx.JSON(code, util.RestWithError("未知的错误..."))
				}
			default: /* 非 error 对象 */
				ctx.JSON(code, util.RestWithError("未知错误, 请稍后重试"))
			}

			ctx.Abort()
		}()

		ctx.Next()
	}
}
