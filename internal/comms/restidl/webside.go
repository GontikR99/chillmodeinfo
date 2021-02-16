// +build wasm

package restidl

import (
	"context"
	"encoding/json"
	"github.com/GontikR99/chillmodeinfo/internal/profile/signins"
	"github.com/GontikR99/chillmodeinfo/internal/comms/httputil"
	"net/http"
)

func wrapError(resCode int, err error) error {
	if resCode >= 400 {
		return httputil.NewError(resCode, err)
	} else {
		return err
	}
}

func call(method string, path string, request interface{}, response interface{}) error {
	packaged := &packagedRequest{
		IdToken: signins.GetToken(),
		ReqMsg:  request,
	}

	reqText, err := json.Marshal(packaged)
	if err != nil {
		return err
	}

	resBody, resCode, err := httpCall(method, path, reqText)
	if err != nil {
		return err
	}

	var responsePackage packagedResponse
	err = json.Unmarshal(resBody, &responsePackage)
	if err != nil {
		return wrapError(resCode, err)
	}

	respMsgBytes, err := json.Marshal(responsePackage.ResMsg)
	if err != nil {
		return wrapError(resCode, err)
	}
	err = json.Unmarshal(respMsgBytes, response)
	if err != nil && !responsePackage.HasError {
		return wrapError(resCode, err)
	}

	if responsePackage.HasError {
		return httputil.NewError(resCode, responsePackage.Error)
	}

	if resCode >= 400 {
		return httputil.NewError(resCode, "Unspecified error in comms")
	} else {
		return nil
	}
}

// Do nothing call so that IDL files can build client side
func serve(mux *http.ServeMux, path string, handler func(ctx context.Context, method string, request *Request) (interface{}, error)) {
}
