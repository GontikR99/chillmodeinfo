package rpcidl

import (
	"github.com/GontikR99/chillmodeinfo/internal/eqspec"
	"net/rpc"
)

type LogEntryBufferServer interface {
	// Retrieve all buffered messages, and stop buffering
	FetchBufferedMessages() []*eqspec.LogEntry
}

type LogEntryBufferStub struct {
	leb LogEntryBufferServer
}

type FetchBufferedMessagesRequest struct{}
type FetchBufferedMessagesResponse struct {
	Messages []*eqspec.LogEntry
}

func (s *LogEntryBufferStub) FetchBufferedMessages(req *FetchBufferedMessagesRequest, res *FetchBufferedMessagesResponse) error {
	res.Messages = s.leb.FetchBufferedMessages()
	return nil
}

func FetchBufferedMessages(client *rpc.Client) ([]*eqspec.LogEntry, error) {
	req := new(FetchBufferedMessagesRequest)
	res := new(FetchBufferedMessagesResponse)
	err := client.Call("LogEntryBufferStub.FetchBufferedMessages", req, res)
	return res.Messages, err
}

func HandleLogEntryBuffer(leb LogEntryBufferServer) func(server *rpc.Server) {
	return func(server *rpc.Server) {
		server.Register(&LogEntryBufferStub{leb})
	}
}
