// +build server

package main

import (
	"bytes"
	"crypto"
	"encoding/hex"
	"fmt"
	"hash"
	"net/http"
	"net/url"
	"strings"
)

type cacheMangler struct {
	next http.Handler
	hashcode string
}

type transformWriter struct {
	next http.ResponseWriter
	buffer bytes.Buffer
	code int
}

func (t *transformWriter) Header() http.Header {return t.next.Header()}
func (t *transformWriter) WriteHeader(statusCode int) {t.code=statusCode}
func (t *transformWriter) Write(bytes []byte) (int, error) {
	return t.buffer.Write(bytes)
}
func (t *transformWriter) Flush(hashcode string) {
	b := t.buffer.Bytes()
	b = bytes.ReplaceAll(b, []byte("bin/webapp.wasm"), []byte("bin/"+hashcode+".wasm"))
	t.Header().Set("Content-Length", fmt.Sprint(len(b)))
	t.next.WriteHeader(t.code)
	t.next.Write(b)
}

func (c *cacheMangler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if !strings.EqualFold(http.MethodGet, request.Method) {
		c.next.ServeHTTP(writer, request)
	} else if request.URL.Path=="/" || request.URL.Path=="/index.html" {
		t := &transformWriter{next: writer}
		defer t.Flush(c.hashcode)
		t.Header().Add("Cache-Control", "no-store, max-age=0")
		c.next.ServeHTTP(t, request)
	} else if request.URL.Path=="/bin/webapp.wasm" {
		http.Error(writer, "Try fetching the persistent hashed version", http.StatusNotFound)
	} else if request.URL.Path=="/bin/"+c.hashcode+".wasm" {
		writer.Header().Add("Cache-Control", "max-age=31536000, min-fresh=31536000, immutable")
		writer.Header().Add("Expires", "Sun, 17-Jan-2038 19:14:07 GMT")
		rc := request.Clone(request.Context())
		rc.RequestURI = "/bin/webapp.wasm"
		rc.URL.Path="/bin/webapp.wasm"
		c.next.ServeHTTP(writer, rc)
	} else {
		c.next.ServeHTTP(writer, request)
	}
}

type hashWriter struct {
	hasher hash.Hash
}

func (h *hashWriter) Header() http.Header {return http.Header{}}
func (h *hashWriter) Write(i []byte) (int, error) {return h.hasher.Write(i)}
func (h hashWriter) WriteHeader(statusCode int) {}

func NewCacheMangler(next http.Handler) http.Handler {
	fauxRequest := &http.Request{
		Method:     http.MethodGet,
		RequestURI: "/bin/webapp.wasm",
		URL:  &url.URL{Path: "/bin/webapp.wasm"},
	}
	hw := &hashWriter{crypto.SHA256.New()}
	next.ServeHTTP(hw, fauxRequest)
	return &cacheMangler{
		next:     next,
		hashcode: hex.EncodeToString(hw.hasher.Sum(nil)),
	}
}