// +build wasm,web

package signin

type ListenerHandle int

var nextHandler = ListenerHandle(1)
var signinHandlers = make(map[ListenerHandle]func(*UserInfo))

func OnStateChange(callback func(*UserInfo)) ListenerHandle {
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
		v(currentSignIn)
	}
}
