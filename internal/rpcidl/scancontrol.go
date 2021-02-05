package rpcidl

import "net/rpc"

type RestartScanRequest struct{}
type RestartScanResponse struct{}

type RestartScanStub struct {
	restart func()
}

func (rss *RestartScanStub) RestartScan(req *RestartScanRequest, res *RestartScanResponse) error {
	rss.restart()
	return nil
}

func RestartScan(client *rpc.Client) error {
	return client.Call("RestartScanStub.RestartScan", new(RestartScanRequest), new(RestartScanResponse))
}

func HandleRestartScan(restart func()) func(server *rpc.Server) {
	return func(server *rpc.Server) {
		server.Register(&RestartScanStub{restart})
	}
}
