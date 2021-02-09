package restidl

import (
	"net/http"
	"strings"
)

const methodPingV0 = "POST"
const pathPingV0 = "/rest/v0/ping"

type pingRequestV0 struct {
	Text string
}

type pingResponseV0 struct {
	Text string
}

func HandlePingV0(handler func(text string, req *Request) string) func(mux *http.ServeMux) {
	return func(mux *http.ServeMux) {
		serve(mux, pathPingV0, func(method string, req *Request) (interface{}, error) {
			if !strings.EqualFold(method, methodPingV0) {
				return nil, NewError(http.StatusBadRequest, "Wrong method")
			}
			reqVal := &pingRequestV0{}
			req.ReadTo(reqVal)
			respText := handler(reqVal.Text, req)
			return &pingResponseV0{respText}, nil
		})
	}
}

func PingV0(text string) (string, error) {
	req := &pingRequestV0{Text: text}
	res := new(pingResponseV0)
	err := call(methodPingV0, pathPingV0, req, res)
	return res.Text, err
}
