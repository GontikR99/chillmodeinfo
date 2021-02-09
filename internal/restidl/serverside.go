// +build server

package restidl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func serve(mux *http.ServeMux, path string, handler func(method string, request *Request) (interface{}, error)) {
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		bodyText, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		wrapper := &Request{}
		err = json.Unmarshal(bodyText, wrapper)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		result, err := func() (val interface{}, err error) {
			defer func() {
				if r := recover(); r != nil {
					val = nil
					if ev, ok := r.(error); ok {
						err = ev
					} else {
						err = errors.New(fmt.Sprint(r))
					}
				}
			}()
			val, err = handler(r.Method, wrapper)
			return
		}()
		if err != nil {
			log.Println(err)
			if he, ok := err.(*httpError); ok {
				http.Error(w, he.Error(), he.StatusCode)
			} else {
				http.Error(w, he.Error(), http.StatusInternalServerError)
			}
			return
		}
		outbytes, err := json.Marshal(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(outbytes)
	})
}

// Server-side do nothing method so IDLs will build
func call(method string, path string, request interface{}, response interface{}) error {
	return nil
}
