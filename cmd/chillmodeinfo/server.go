// +build server

package main

import (
	"crypto/tls"
	"flag"
	"github.com/GontikR99/chillmodeinfo/cmd/chillmodeinfo/serverrpcs"
	"github.com/GontikR99/chillmodeinfo/internal/sitedef"
	"github.com/GontikR99/chillmodeinfo/web/bin"
	"github.com/GontikR99/chillmodeinfo/web/static"
	"github.com/NYTimes/gziphandler"
	"log"
	"net/http"
	"strconv"
)

func main() {
	certPath := flag.String("cert", "", "Path to PEM format certificates")
	keyPath := flag.String("key", "", "Path to PEM format private key")
	flag.Parse()

	baseMux := serverrpcs.NewMux()
	baseMux.Handle("/", http.FileServer(static.StaticFiles))
	baseMux.Handle("/bin/", http.StripPrefix("/bin", http.FileServer(bin.BinFiles)))

	muxWithGzip := gziphandler.GzipHandler(baseMux)

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	srv := &http.Server{
		Addr:         ":"+strconv.Itoa(sitedef.Port),
		Handler:      &handlerWrapper{muxWithGzip},
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	log.Println("Server start")
	log.Fatal(srv.ListenAndServeTLS(*certPath, *keyPath))
}
