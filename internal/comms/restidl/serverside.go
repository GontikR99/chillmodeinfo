// +build server

package restidl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GontikR99/chillmodeinfo/internal/comms/httputil"
	"github.com/GontikR99/chillmodeinfo/internal/dao"
	"github.com/GontikR99/chillmodeinfo/internal/profile/signins"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strings"
)

const TagRequest = "tagRequest"

type loggingWriter struct {
	req *http.Request
	next http.ResponseWriter
	resCode int
	authUser string
}

func (l *loggingWriter) Header() http.Header {return l.next.Header()}
func (l *loggingWriter) Write(bytes []byte) (int, error) {return l.next.Write(bytes)}
func (l *loggingWriter) WriteHeader(statusCode int) {
	l.resCode=statusCode
	l.next.WriteHeader(statusCode)
}
func (l *loggingWriter) LogResult() {
	log.Printf("%s %s [%s] -> %d", l.req.Method, l.req.RequestURI, l.authUser, l.resCode)
}

func serve(mux *http.ServeMux, path string, handler func(ctx context.Context, method string, request *Request) (interface{}, error)) {
	mux.HandleFunc(path, func(wOrig http.ResponseWriter, r *http.Request) {
		w := &loggingWriter{req: r, next: wOrig}
		defer w.LogResult()
		bodyText, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		packaged := &packagedRequest{}
		headerBodyText := strings.Join(r.Header[HeaderRequestPayload], "")
		if headerBodyText != "" {
			bodyText = []byte(headerBodyText)
		}
		err = json.Unmarshal(bodyText, packaged)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		userId, idErr := signins.ValidateToken(r.Context(), packaged.IdToken)
		if idErr==nil {
			userProfile, err := dao.LookupProfile(userId)
			if err==nil && userProfile.GetDisplayName()!="" {
				w.authUser = fmt.Sprintf("%s(%s)", userProfile.GetDisplayName(), userId)
			} else {
				w.authUser = userId
			}
		} else {
			w.authUser=idErr.Error()
		}
		result, err := func() (val interface{}, err error) {
			defer func() {
				if r := recover(); r != nil {
					stb := make([]byte, 65536)
					stbLen := runtime.Stack(stb, false)
					log.Print(string(stb[:stbLen]))

					if ev, ok := r.(error); ok {
						val = ev.Error()
						err = ev
					} else {
						val = fmt.Sprint(r)
						err = errors.New(fmt.Sprint(r))
					}
				}
			}()
			wrappedRequest := &Request{
				IdToken:       packaged.IdToken,
				UserId:        userId,
				IdentityError: idErr,
				packaged:      packaged,
			}
			subCtx := context.WithValue(r.Context(), TagRequest, wrappedRequest)
			val, err = handler(subCtx, r.Method, wrappedRequest)
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
