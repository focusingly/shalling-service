package netadapter

import (
	"crypto/tls"
	"fmt"
	"log"
	"space-api/conf"

	"github.com/gin-gonic/gin"
	"github.com/quic-go/quic-go/http3"
)

func RunAndServe(engine *gin.Engine, appConf *conf.AppConf) {
	tlsCfg := &tls.Config{
		Certificates: []tls.Certificate{
			mustLoadCert(
				appConf.Certs.Pem,
				appConf.Certs.Key,
			)},
		MinVersion: tls.VersionTLS13,
	}
	server := &http3.Server{
		Handler:   engine,
		Addr:      fmt.Sprintf(":%d", appConf.Port),
		TLSConfig: tlsCfg,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func mustLoadCert(certFile, keyFile string) tls.Certificate {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}
	return cert
}
