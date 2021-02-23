// +build wasm

package rpcidl

import (
	"github.com/GontikR99/chillmodeinfo/internal/overlay/update"
	"net/rpc"
)

type UpdateQueueHandler interface {
	Poll() (map[int]*update.UpdateEntry, error)
	Enqueue(map[int]*update.UpdateEntry) error
}

type UpdateQueueServerStub struct {
	handler UpdateQueueHandler
}

type updateQueueClientStub struct {
	client *rpc.Client
}

type UpdateQueuePollRequest struct{}
type UpdateQueuePollResponse struct {
	Entries map[int]*update.UpdateEntry
}

func (uss *UpdateQueueServerStub) Poll(req *UpdateQueuePollRequest, res *UpdateQueuePollResponse) error {
	var err error
	res.Entries, err = uss.handler.Poll()
	return err
}

func (ucs *updateQueueClientStub) Poll() (map[int]*update.UpdateEntry, error) {
	req := new(UpdateQueuePollRequest)
	res := new(UpdateQueuePollResponse)
	err := ucs.client.Call("UpdateQueueServerStub.Poll", req, res)
	return res.Entries, err
}

type UpdateQueueEnqueueRequest struct {
	Entries map[int]*update.UpdateEntry
}
type UpdateQueueEnqueueResponse struct{}

func (uss *UpdateQueueServerStub) Enqueue(req *UpdateQueueEnqueueRequest, res *UpdateQueueEnqueueResponse) error {
	return uss.handler.Enqueue(req.Entries)
}

func (ucs *updateQueueClientStub) Enqueue(entries map[int]*update.UpdateEntry) error {
	req := &UpdateQueueEnqueueRequest{entries}
	res := new(UpdateQueueEnqueueResponse)
	return ucs.client.Call("UpdateQueueServerStub.Enqueue", req, res)
}

func UpdateQueue(client *rpc.Client) UpdateQueueHandler {
	return &updateQueueClientStub{client}
}

func HandleUpdateQueue(handler UpdateQueueHandler) func(server *rpc.Server) {
	return func(server *rpc.Server) {
		server.Register(&UpdateQueueServerStub{handler})
	}
}
