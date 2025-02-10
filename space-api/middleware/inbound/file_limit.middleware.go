package inbound

import (
	"io"
	"net/http"
	"space-api/constants"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var _rawBodyMark = uuid.NewString()

// ResetUploadFileLimitSize 重置文件大小限制
func ResetUploadFileLimitSize(ctx *gin.Context) {
	if val, ok := ctx.Get(_rawBodyMark); ok {
		if rc, ok2 := val.(io.ReadCloser); ok2 {
			ctx.Request.Body = rc
		}
	}
}

func UseUploadFileLimitMiddleware(maxFileSize constants.MemoryByteSize) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rawBody := ctx.Request.Body
		ctx.Set(_rawBodyMark, rawBody)

		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, int64(maxFileSize))
		ctx.Next()
	}
}
