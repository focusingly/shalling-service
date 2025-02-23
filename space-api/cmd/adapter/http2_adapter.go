//go:build usehttp2
// +build usehttp2

package adapter

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"space-api/conf"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/http2"
)

func RunAndServe(engine *gin.Engine, appConf *conf.AppConf) {
	tlsCfg := &tls.Config{
		CurvePreferences:         []tls.CurveID{tls.CurveP256, tls.X25519},
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
	}

	h2Server := &http2.Server{}
	server := &http.Server{
		TLSConfig: tlsCfg,
		Addr:      fmt.Sprintf(":%d", appConf.Port),
		Handler:   NewAdaptiveHttpWriter(engine),
	}

	http2.ConfigureServer(server, h2Server)
	server.ListenAndServeTLS(appConf.Certs.Pem, appConf.Certs.Key)
}
