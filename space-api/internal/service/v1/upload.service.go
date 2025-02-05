package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type (
	UploadService interface {
		Upload(ctx *gin.Context) (err error)
	}

	defaultImpl struct{}
)

// Upload implements UploadService.
func (d *defaultImpl) Upload(ctx *gin.Context) (err error) {
	form, err := ctx.MultipartForm()
	if err != nil {
		return
	}

	files := form.File["files"]
	for _, file := range files {
		fmt.Println(file.Filename, "---------> ", file.Size)
	}

	return
}

func newUploadService() UploadService {
	return &defaultImpl{}
}

var DefaultUploadService UploadService

func init() {
	DefaultUploadService = newUploadService()
}
