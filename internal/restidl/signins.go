package restidl

import (
	"github.com/GontikR99/chillmodeinfo/internal/httputil"
	"net/http"
	"strings"
)

const pathSigninsV0 = "/rest/v0/signins"


type verifySigninRequestV0 struct {}
type verifySigninResponseV0 struct {}

func VerifyLogin() error {
	return call("POST", pathSigninsV0, new(verifySigninRequestV0), new(verifySigninResponseV0))
}

func HandleLogin() func(mux *http.ServeMux) {
	return func(mux *http.ServeMux) {
		serve(mux, pathSigninsV0, func(method string, request *Request) (interface{}, error) {
			if !strings.EqualFold("POST", method) {
				return &verifySigninResponseV0{}, httputil.NewError(http.StatusBadRequest, "Unsupported method "+method)
			} else {
				return &verifySigninResponseV0{}, request.IdentityError
			}
		})
	}
}