package outbound

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"space-api/constants"
	"space-api/util"

	"github.com/gin-gonic/gin"
)

func UseErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 具体的错误捕获处理
		defer func() {
			isPanic := false
			// 捕获 panic
			catchErr := recover()
			if catchErr != nil {
				l := fmt.Sprintf("%s%s%s", constants.RED, string(debug.Stack()), constants.RESET)
				log.Println(l)
			}
			if catchErr != nil {
				isPanic = true
			} else {
				if len(ctx.Errors) != 0 {
					catchErr = ctx.Errors[0].Err
				}
			}
			if catchErr == nil {
				return
			}

			code := util.TernaryExpr(isPanic, http.StatusInternalServerError, http.StatusOK)
			var restErr *util.RestResult[any]
			switch err := catchErr.(type) {
			case error: /* 确保都是实现 error 接口的结构体的引用 */
				switch err := err.(type) {
				case *util.BizErr,
					*util.LimitErr,
					*util.AuthErr:
					restErr = util.RestWithError(err.Error())
				case *util.FatalErr:
					restErr = util.RestWithError("服务内部错误, 请稍后重试或联系站长修复")
				default:
					restErr = util.RestWithError("未知的错误")
				}
			default: /* 非 error 对象 */
				restErr = util.RestWithError("未知错误, 请稍后重试")
			}
			// 根据请求头返回响应格式
			handleProduce(code, restErr, ctx)
			ctx.Abort()
		}()

		ctx.Next()
	}
}
