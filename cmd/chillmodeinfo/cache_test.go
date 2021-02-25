// +build server

package main

import (
	"bytes"
	"github.com/GontikR99/chillmodeinfo/web/bin"
	"github.com/GontikR99/chillmodeinfo/web/static"
	"net/http"
	"net/url"
	"regexp"
	"testing"
)

type bufferWriter struct {
	header http.Header
	body bytes.Buffer
	code int
}

func (b *bufferWriter) Header() http.Header {return b.header}
func (b *bufferWriter) Write(i []byte) (int, error) {return b.body.Write(i)}
func (b *bufferWriter) WriteHeader(statusCode int) {b.code=statusCode}

func newBufferWriter() *bufferWriter {
	return &bufferWriter{header: http.Header{}}
}

var rex=regexp.MustCompile("\"(bin/.*[.]wasm)\"")

func TestCacheMangler(t *testing.T) {
	baseMux := http.NewServeMux()
	baseMux.Handle("/", http.FileServer(static.StaticFiles))
	baseMux.Handle("/bin/", http.StripPrefix("/bin", http.FileServer(bin.BinFiles)))

	cm := NewCacheMangler(baseMux)
	fauxRequest := &http.Request{
		Method:     http.MethodGet,
		RequestURI: "/",
		URL:  &url.URL{Path: "/"},
	}
	bw := newBufferWriter()
	cm.ServeHTTP(bw, fauxRequest)
	if bw.code!=200 {
		t.Fatal(bw.code)
	}
	var parts []string
	if parts = rex.FindStringSubmatch(bw.body.String()); parts==nil {
		t.Fatal("no hash")
	}

	fauxRequest = &http.Request{
		Method:     http.MethodGet,
		RequestURI: "/"+parts[1],
		URL:  &url.URL{Path: "/"+parts[1]},
	}
	bw = newBufferWriter()
	cm.ServeHTTP(bw, fauxRequest)
	if bw.code!=200 {
		t.Fatal(bw.code)
	}
}
