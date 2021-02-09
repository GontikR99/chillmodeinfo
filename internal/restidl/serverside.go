// +build server

package restidl

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GontikR99/chillmodeinfo/internal/signins"
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
		packaged := &packagedRequest{}
		err = json.Unmarshal(bodyText, packaged)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var userId string
		var idErr error
		if packaged.IdToken!="" {
			userId, idErr = signins.ValidateToken(r.Context(), packaged.IdToken)
		} else if packaged.ClientId!="" {
			userId, idErr = signins.ValidateClientId(r.Context(), packaged.ClientId)
		} else {
			userId=""
			idErr = NewError(http.StatusUnauthorized, "No identity provided")
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
			val, err = handler(r.Method, &Request{
				UserId:        userId,
				IdentityError: idErr,
				packaged:      packaged,
			})
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
