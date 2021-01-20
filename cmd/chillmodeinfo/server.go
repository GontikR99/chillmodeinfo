package main

import (
	"github.com/GontikR99/chillmodeinfo/web/static"
	"github.com/NYTimes/gziphandler"
	"log"
	"net/http"
)

func main() {
	baseMux := http.NewServeMux()
	baseMux.Handle("/", http.FileServer(static.StaticFiles))

	muxWithGzip := gziphandler.GzipHandler(baseMux)

	server := &http.Server{
		Addr:    ":8123",
		Handler: muxWithGzip,
	}

	log.Println("Server start")
	log.Fatal(server.ListenAndServe())
}
