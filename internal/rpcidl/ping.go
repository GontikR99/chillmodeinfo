package rpcidl

import "net/rpc"

type PingRequest struct {
	Message string
}

type PingResponse struct {
}

type PingStub struct {
	callback func(string)
}

func (ps *PingStub) Ping(req *PingRequest, resp *PingResponse) error {
	ps.callback(req.Message)
	return nil
}

func PingHandler(callback func(string)) func(server *rpc.Server) {
	ps := &PingStub{callback: callback}
	return func(server *rpc.Server) {
		server.Register(ps)
	}
}

func Ping(client *rpc.Client, message string) error {
	return client.Call("PingStub.Ping", &PingRequest{message}, &PingResponse{})
}
