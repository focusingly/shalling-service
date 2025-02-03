package common

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path"
	"space-api/pack"
	"time"

	"github.com/gin-gonic/gin"
)

// ReadSeekerWrapper wraps an embed.FS file to provide ReadSeeker functionality
type ReadSeekerWrapper struct {
	file     io.Reader
	size     int64
	modTime  time.Time
	position int64
}

var _ io.ReadSeeker = (*ReadSeekerWrapper)(nil)

func (r *ReadSeekerWrapper) Read(p []byte) (n int, err error) {
	return r.file.Read(p)
}
func (r *ReadSeekerWrapper) Seek(offset int64, whence int) (int64, error) {
	var newPos int64
	switch whence {
	case io.SeekStart:
		newPos = offset
	case io.SeekCurrent:
		newPos = r.position + offset
	case io.SeekEnd:
		newPos = r.size + offset
	default:
		return 0, fmt.Errorf("invalid whence")
	}
	if newPos < 0 || newPos > r.size {
		return 0, fmt.Errorf("seek out of bounds")
	}
	r.position = newPos
	return r.position, nil
}

func UseSpaPageController(routeGroup *gin.RouterGroup) {
	const indexName = "static/dist/index.html"

	routeGroup.GET("/*filename", func(ctx *gin.Context) {
		fmt.Println(ctx.Request.URL, ctx.Request.RequestURI)

		filename := path.Join("static", "dist", ctx.Param("filename"))
		if ctx.Request.URL.Path == "/" {
			filename = indexName
		}
		var info fs.FileInfo
		file, err := pack.SpaResource.Open(filename)
		// 资源处理
		if err == nil {
			info, _ = file.Stat()
			// 重定向回根文件
			if info.IsDir() {
				file, _ = pack.SpaResource.Open(indexName)
				info, _ = file.Stat()
			}
		} else {
			file, _ = pack.SpaResource.Open(indexName)
			info, _ = file.Stat()
		}
		defer file.Close()
		wrap := &ReadSeekerWrapper{
			file:    file,
			size:    info.Size(),
			modTime: info.ModTime(),
		}
		http.ServeContent(ctx.Writer, ctx.Request, filename, wrap.modTime, wrap)
	})
}
