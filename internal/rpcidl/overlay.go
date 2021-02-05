package rpcidl

import "net/rpc"

type OverlayServer interface {
	CloseOverlay(name string) error
	PositionOverlay(name string) error
	ResetOverlay(name string) error
}

type OverlayStub struct {
	overlayServer OverlayServer
}

type CloseOverlayRequest struct {
	Name string
}
type CloseOverlayResponse struct{}

func (os OverlayStub) CloseOverlay(req *CloseOverlayRequest, res *CloseOverlayResponse) error {
	return os.overlayServer.CloseOverlay(req.Name)
}

func CloseOverlay(client *rpc.Client, name string) error {
	return client.Call("OverlayStub.CloseOverlay", &CloseOverlayRequest{name}, new(CloseOverlayResponse))
}

type PositionOverlayRequest struct {
	Name string
}
type PositionOverlayResponse struct{}

func (os OverlayStub) PositionOverlay(req *PositionOverlayRequest, res *PositionOverlayResponse) error {
	return os.overlayServer.PositionOverlay(req.Name)
}
func PositionOverlay(client *rpc.Client, name string) error {
	return client.Call("OverlayStub.PositionOverlay", &PositionOverlayRequest{name}, new(PositionOverlayResponse))
}

type ResetOverlayRequest struct {
	Name string
}
type ResetOverlayResponse struct{}

func (os OverlayStub) ResetOverlay(req *ResetOverlayRequest, res *ResetOverlayResponse) error {
	return os.overlayServer.ResetOverlay(req.Name)
}
func ResetOverlay(client *rpc.Client, name string) error {
	return client.Call("OverlayStub.ResetOverlay", &ResetOverlayRequest{name}, new(ResetOverlayResponse))
}

func HandleOverlay(overlayServer OverlayServer) func(server *rpc.Server) {
	return func(server *rpc.Server) {
		server.Register(OverlayStub{overlayServer})
	}
}
