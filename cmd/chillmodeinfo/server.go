// +build server

package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/GontikR99/chillmodeinfo/cmd/chillmodeinfo/serverrpcs"
	"github.com/GontikR99/chillmodeinfo/internal/sitedef"
	"github.com/GontikR99/chillmodeinfo/web/bin"
	"github.com/GontikR99/chillmodeinfo/web/static"
	"github.com/NYTimes/gziphandler"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"time"
)

func makeHTTPServer() *http.Server {
	baseMux := &http.ServeMux{}
		serverrpcs.HandleRest(baseMux)
		handleAssociatePage(baseMux)
		baseMux.Handle("/", http.FileServer(static.StaticFiles))
		baseMux.Handle("/bin/", http.StripPrefix("/bin", http.FileServer(bin.BinFiles)))

		muxWithGzip := gziphandler.GzipHandler(baseMux)

	return makeServerFromMux(muxWithGzip)
}

func makeHTTPToHTTPSRedirectServer() *http.Server {
	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
		newURI := "https://" + r.Host + r.URL.String()
		http.Redirect(w, r, newURI, http.StatusFound)
	}
	mux := &http.ServeMux{}
	mux.HandleFunc("/", handleRedirect)
	return makeServerFromMux(mux)
}

func makeServerFromMux(mux http.Handler) *http.Server {
	// set timeouts so that a slow or malicious client doesn't
	// hold resources forever
	return &http.Server{
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 480 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
}

func main() {
	var httpsSrv *http.Server
	var m *autocert.Manager

	// Note: use a sensible value for data directory
	// this is where cached certificates are stored
	dataDir := "."
	hostPolicy := func(ctx context.Context, host string) error {
		allowedHost := sitedef.DNSName
		if host == allowedHost {
			return nil
		}
		return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
	}

	httpsSrv = makeHTTPServer()
	m = &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: hostPolicy,
		Cache:      autocert.DirCache(dataDir),
	}
	httpsSrv.Addr = fmt.Sprintf(":%d", sitedef.Port)
	httpsSrv.TLSConfig = &tls.Config{
		GetCertificate: m.GetCertificate,
		MinVersion:               tls.VersionTLS12,
		//CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		//PreferServerCipherSuites: true,
		//CipherSuites: []uint16{
		//	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		//	tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		//	tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		//	tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		//},
	}

	go func() {
		err := httpsSrv.ListenAndServeTLS("", "")
		if err != nil {
			log.Fatalf("httpsSrv.ListendAndServeTLS() failed with %s", err)
		}
	}()

	httpSrv := makeHTTPToHTTPSRedirectServer()
	if m != nil {
		// allow autocert handle Let's Encrypt auth callbacks over HTTP.
		// it'll pass all other urls to our hanlder
		httpSrv.Handler = m.HTTPHandler(httpSrv.Handler)
	}
	httpSrv.Addr = ":80"
	err := httpSrv.ListenAndServe()
	if err != nil {
		log.Fatalf("httpSrv.ListenAndServe() failed with %s", err)
	}
}