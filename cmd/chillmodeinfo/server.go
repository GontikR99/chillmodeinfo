package main

import (
	"github.com/GontikR99/chillmodeinfo/web/static"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(static.StaticFiles))

	server := &http.Server{
		Addr:    ":8123",
		Handler: mux,
	}

	log.Println("Server start")
	log.Fatal(server.ListenAndServe())
}
