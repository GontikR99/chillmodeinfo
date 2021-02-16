// +build wasm

package signins

var currentToken *string

func SignedIn() bool {
	return currentToken != nil
}

func GetToken() string {
	if SignedIn() {
		return *currentToken
	} else {
		return ""
	}
}

func SetToken(newToken string) {
	if GetToken() != newToken {
		currentToken = &newToken
		go invokeHandlers()
	}
}

func ClearToken() {
	if SignedIn() {
		currentToken = nil
		go invokeHandlers()
	}
}

type ListenerHandle int

var nextHandler = ListenerHandle(1)
var signinHandlers = make(map[ListenerHandle]func())

func OnStateChange(callback func()) ListenerHandle {
	handle := nextHandler
	nextHandler++
	signinHandlers[handle] = callback
	return handle
}

func (lh ListenerHandle) Release() {
	delete(signinHandlers, lh)
}

func invokeHandlers() {
	for _, v := range signinHandlers {
		v()
	}
}
