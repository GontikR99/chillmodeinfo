// +build wasm

package rpcidl

import (
	"net/rpc"
)

type SigninHandler interface {
	SignIn() error
	SignOut() error
	PollSignIn() error
}

type SigninServerStub struct {
	handler SigninHandler
}

type signinClientStub struct {
	client *rpc.Client
}

type SigninRequest struct{}
type SigninResponse struct{}

func (ss *SigninServerStub) SignIn(req *SigninRequest, res *SigninResponse) error {
	return ss.handler.SignIn()
}

func (cs *signinClientStub) SignIn() error {
	return cs.client.Call("SigninServerStub.SignIn", new(SigninRequest), new(SigninResponse))
}

type SignoutRequest struct{}
type SignoutResponse struct{}

func (ss *SigninServerStub) SignOut(req *SignoutRequest, res *SignoutResponse) error {
	return ss.handler.SignOut()
}

func (cs *signinClientStub) SignOut() error {
	return cs.client.Call("SigninServerStub.SignOut", new(SignoutRequest), new(SignoutResponse))
}

type PollSignInRequest struct{}
type PollSignInResponse struct{}

func (ss *SigninServerStub) PollSignIn(req *PollSignInRequest, res *PollSignInResponse) error {
	return ss.handler.PollSignIn()
}

func (cs *signinClientStub) PollSignIn() error {
	return cs.client.Call("SigninServerStub.PollSignIn", new(PollSignInRequest), new(PollSignInResponse))
}

func HandleSignIn(handler SigninHandler) func(server *rpc.Server) {
	return func(server *rpc.Server) {
		server.Register(&SigninServerStub{handler})
	}
}

func GetSignIn(client *rpc.Client) SigninHandler {
	return &signinClientStub{client}
}
