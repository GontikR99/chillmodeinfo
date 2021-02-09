// +build server

package main

import (
	"log"
	"net/http"
)

type handlerWrapper struct {
	next http.Handler
}

type rwWrapper struct {
	statusCode int
	next http.ResponseWriter
}

func (r *rwWrapper) Header() http.Header {
	return r.next.Header()
}

func (r *rwWrapper) Write(bytes []byte) (int, error) {
	return r.next.Write(bytes)
}

func (r *rwWrapper) WriteHeader(statusCode int) {
	r.statusCode=statusCode
	r.next.WriteHeader(statusCode)
}

func (h *handlerWrapper) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	wrapper := &rwWrapper{next: writer}
	h.next.ServeHTTP(wrapper, request)
	log.Printf("[%v] %s %s -> %d", request.RemoteAddr, request.Method, request.RequestURI, wrapper.statusCode)
}
