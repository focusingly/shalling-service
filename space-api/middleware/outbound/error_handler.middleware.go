package outbound

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"space-api/constants"
	"space-api/middleware/inbound"
	"space-api/util"
	"space-api/util/ip"
	"space-api/util/performance"
	"space-api/util/ptr"
	"space-api/util/rest"
	"space-domain/dao/extra"
	"space-domain/model"
	"time"

	"github.com/gin-gonic/gin"
)

func UseErrorHandler() gin.HandlerFunc {
	runner := performance.DefaultTaskRunner

	return func(ctx *gin.Context) {
		startTime := time.Now().UnixMilli()

		// 具体的错误捕获处理
		defer func() {
			costTime := time.Now().UnixMilli() - startTime
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

			ipv4Str := inbound.GetRealIpWithContext(ctx)
			userDetail := inbound.GetUserAgentFromContext(ctx)
			source, e := ip.GetIpSearcher().SearchByStr(ipv4Str)

			logInfo := &model.LogInfo{
				LogType: string(constants.APIRequest),
				// Message:       "",
				Level:         string(constants.Error),
				CostTime:      costTime,
				RequestMethod: &ctx.Request.Method,
				RequestURI:    &ctx.Request.RequestURI,
				StackTrace: util.TernaryExpr(
					isPanic,
					ptr.ToPtr(ptr.Bytes2String(debug.Stack())),
					nil,
				),
				IPAddr:    &ipv4Str,
				IPSource:  util.TernaryExpr(e != nil, nil, &source),
				Useragent: &userDetail.Useragent,
				CreatedAt: time.Now().UnixMilli(),
			}

			code := util.TernaryExpr(
				isPanic,
				http.StatusInternalServerError,
				http.StatusOK,
			)
			var restErr *rest.RestResult[any]
			switch err := catchErr.(type) {
			case error: /* 确保都是实现 error 接口的结构体的引用 */
				switch err := err.(type) {
				case *util.BizErr:
					restErr = rest.RestWithError(err.Error())
					logInfo.Message = fmt.Sprintf("%s : %s", err.Msg, err.Reason.Error())

				case *util.AuthErr:
					restErr = rest.RestWithError(err.Error())
					logInfo.Level = string(constants.Warn)
					logInfo.Message = fmt.Sprintf("%s : %s", err.Msg, err.Reason.Error())

				case *util.LimitErr:
					restErr = rest.RestWithError(err.Error())
					logInfo.LogType = string(constants.RequestLimit)
					logInfo.Message = fmt.Sprintf("%s : %s", err.Msg, err.Reason.Error())

				case *util.FatalErr:
					restErr = rest.RestWithError("服务内部错误, 请稍后重试或联系站长修复")
					logInfo.Level = string(constants.Fatal)
					logInfo.Message = fmt.Sprintf("%s : %s", err.Msg, err.Reason.Error())

				case *util.NotMethodOrResourceErr:
					restErr = rest.RestWithError(err.Error())
					logInfo.Level = string(constants.Warn)
					logInfo.Message = fmt.Sprintf("%s : %s", err.Msg, err.Reason.Error())
				default:
					restErr = rest.RestWithError("未知的错误")
					logInfo.Level = string(constants.Fatal)
				}
			default: /* 非 error 对象 */
				restErr = rest.RestWithError("未知错误, 请稍后重试")
				logInfo.Level = string(constants.Fatal)
				logInfo.Message = fmt.Sprintf("%#v", err)
			}

			// 根据请求头返回响应格式
			handleProduce(code, restErr, ctx)
			ctx.Abort()

			// 异步写入日志
			runner.Go(func() {
				extra.LogInfo.WithContext(context.TODO()).Create(logInfo)
			})
		}()

		ctx.Next()
	}
}
