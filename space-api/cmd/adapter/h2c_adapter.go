//go:build !usehttp2
// +build !usehttp2

package adapter

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"space-api/conf"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func RunAndServe(engine *gin.Engine, appConf *conf.AppConf) {
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", appConf.Port),
		Handler: h2c.NewHandler(
			NewAdaptiveCompressionHttpWriter(engine),
			&http2.Server{},
		),
		TLSConfig: &tls.Config{
			NextProtos: []string{"h2", "http/1.1"},
		},
	}
	// 使用 h2c 进行优化传输
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
