package signin

import (
	"syscall/js"
)

type UserInfo struct {
	ID      string
	Name    string
	Email   string
	jsValue js.Value
}

func (i *UserInfo) JSValue() js.Value {
	return i.jsValue
}

func (i *UserInfo) IdToken() string {
	return i.JSValue().Call("getAuthResponse").Get("id_token").String()
}

var currentSignIn *UserInfo

func CurrentSignIn() *UserInfo {
	return currentSignIn
}

func SignedIn() bool {
	return currentSignIn != nil
}
