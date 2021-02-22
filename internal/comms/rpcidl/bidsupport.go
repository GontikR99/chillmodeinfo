// +build wasm

package rpcidl

import "net/rpc"

type BidSupportHandler interface {
	// Retrieve the last item the player mentioned
	GetLastMentioned() (string, error)

	// Describe a new high bid
	OfferBid(bidder string, itemname string, bidValue float64) error
}

type BidSupportServerStub struct {
	handler BidSupportHandler
}

type bidSupportClientStub struct {
	client *rpc.Client
}

type BidSupportGetLastMentionedRequest struct {}
type BidSupportGetLastMentionedResponse struct {
	LastMentioned string
}

func (bss *BidSupportServerStub) GetLastMentioned(req *BidSupportGetLastMentionedRequest, res *BidSupportGetLastMentionedResponse) error {
	var err error
	res.LastMentioned, err = bss.handler.GetLastMentioned()
	return err
}

func (bsc *bidSupportClientStub) GetLastMentioned() (string, error) {
	req := new(BidSupportGetLastMentionedRequest)
	res := new(BidSupportGetLastMentionedResponse)
	err := bsc.client.Call("BidSupportServerStub.GetLastMentioned", req, res)
	return res.LastMentioned, err
}

type BidSupportOfferBidRequest struct {
	Bidder string
	ItemName string
	BidValue float64
}
type BidSupportOfferBidResponse struct{}

func (bss *BidSupportServerStub) OfferBid(req *BidSupportOfferBidRequest, res *BidSupportOfferBidResponse) error {
	err := bss.handler.OfferBid(req.Bidder, req.ItemName, req.BidValue)
	return err
}

func (bsc *bidSupportClientStub) OfferBid(bidder string, itemName string, bidValue float64) error {
	req:=&BidSupportOfferBidRequest{bidder, itemName, bidValue}
	res:=new(BidSupportOfferBidResponse)
	return bsc.client.Call("BidSupportServerStub.OfferBid", req, res)
}

func BidSupport(client *rpc.Client) BidSupportHandler {
	return &bidSupportClientStub{client}
}

func HandleBidSupport(handler BidSupportHandler) func(server *rpc.Server) {
	return func(server *rpc.Server) {
		server.Register(&BidSupportServerStub{handler})
	}
}