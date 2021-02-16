package restidl

import (
	"encoding/json"
	"github.com/GontikR99/chillmodeinfo/internal/comms/httputil"
	"net/http"
)

const HeaderRequestPayload = "Request-Payload"

type packagedRequest struct {
	IdToken string
	ReqMsg  interface{}
}

type Request struct {
	IdToken       string
	UserId        string
	IdentityError error
	packaged      *packagedRequest
}

func (rr *Request) ReadTo(reqStruct interface{}) {
	jsonText, err := json.Marshal(rr.packaged.ReqMsg)
	if err != nil {
		panic(httputil.NewError(http.StatusInternalServerError, err))
	}
	err = json.Unmarshal(jsonText, reqStruct)
	if err != nil {
		panic(httputil.NewError(http.StatusBadRequest, err))
	}
}

type packagedResponse struct {
	HasError bool
	Error    string
	ResMsg   interface{}
}
