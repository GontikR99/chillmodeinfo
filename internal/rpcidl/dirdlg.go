package rpcidl

import (
	"net/rpc"
)

type DirectoryDialogRequest struct {
	InitialDirectory string
}

type DirectoryDialogResponse struct {
	ChosenDirectory string
}

type DirectoryDialogStub struct {
	callback func(initial string) (chosen string, err error)
}

func (d *DirectoryDialogStub) DirectoryDialog(req *DirectoryDialogRequest, res *DirectoryDialogResponse) error {
	chosen, err := d.callback(req.InitialDirectory)
	*res=DirectoryDialogResponse{chosen}
	return err
}

func DirectoryDialog(client *rpc.Client, initial string) (string,error) {
	req := &DirectoryDialogRequest{initial}
	res := &DirectoryDialogResponse{}
	err := client.Call("DirectoryDialogStub.DirectoryDialog", req, res)
	return res.ChosenDirectory, err
}

func HandleDirectoryDialog(callback func(string)(string, error)) func(*rpc.Server) {
	dds := &DirectoryDialogStub{callback}
	return func(server *rpc.Server) {
		server.Register(dds)
	}
}