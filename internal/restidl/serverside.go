// +build server

package restidl

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GontikR99/chillmodeinfo/internal/httputil"
	"github.com/GontikR99/chillmodeinfo/internal/signins"
	"io/ioutil"
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
		userId, idErr := signins.ValidateToken(r.Context(), packaged.IdToken)

		result, err := func() (val interface{}, err error) {
			defer func() {
				if r := recover(); r != nil {
					if ev, ok := r.(error); ok {
						val = ev.Error()
						err = ev
					} else {
						val = fmt.Sprint(r)
						err = errors.New(fmt.Sprint(r))
					}
				}
			}()
			val, err = handler(r.Method, &Request{
				IdToken:       packaged.IdToken,
				UserId:        userId,
				IdentityError: idErr,
				packaged:      packaged,
			})
			return
		}()
		w.Header().Set("Content-Type", "application/json")
		var response *packagedResponse
		if err == nil {
			response = &packagedResponse{
				HasError: false,
				Error:    "",
			}
			w.WriteHeader(http.StatusOK)
		} else {
			response = &packagedResponse{
				HasError: true,
				Error:    err.Error(),
				ResMsg:   result,
			}
			if he, ok := err.(*httputil.HttpError); ok {
				w.WriteHeader(he.StatusCode)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
		response.ResMsg = result
		outbytes, _ := json.Marshal(response)
		w.Write(outbytes)
	})
}

// Server-side do nothing method so IDLs will build
func call(method string, path string, request interface{}, response interface{}) error {
	return nil
}
