package restidl

import (
	"context"
	"github.com/GontikR99/chillmodeinfo/internal/comms/httputil"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"net/http"
	"strings"
)

const endpointRaidV0 = "/rest/v0/raid"

type raidClientStub struct{}

var Raid = &raidClientStub{}

type RaidHandler interface {
	Add(ctx context.Context, raid record.Raid) error
	Fetch(ctx context.Context) ([]record.Raid, error)
	Delete(ctx context.Context, raidId uint64) error
	Update(ctx context.Context, raid record.Raid) (record.Raid, error)
}

type addRaidRequestV0 struct{ Raid *record.BasicRaid }
type addRaidResponseV0 struct{}

func (r raidClientStub) Add(ctx context.Context, raid record.Raid) error {
	req := &addRaidRequestV0{record.NewBasicRaid(raid)}
	res := new(addRaidResponseV0)
	return call(http.MethodPut, endpointRaidV0, req, res)
}

type getRaidsRequestV0 struct{}
type getRaidsResponseV0 struct{ Raids []*record.BasicRaid }

func (r raidClientStub) Fetch(ctx context.Context) ([]record.Raid, error) {
	req := new(getRaidsRequestV0)
	res := new(getRaidsResponseV0)
	err := call(http.MethodGet, endpointRaidV0, req, res)
	if err != nil {
		return nil, err
	}
	var raids []record.Raid
	for _, v := range res.Raids {
		raids = append(raids, record.NewBasicRaid(v))
	}
	return raids, nil
}

type deleteRaidRequestV0 struct{ RaidId uint64 }
type deleteRaidResponseV0 struct{}

func (r raidClientStub) Delete(ctx context.Context, raidId uint64) error {
	req := &deleteRaidRequestV0{raidId}
	res := new(deleteRaidResponseV0)
	return call(http.MethodDelete, endpointRaidV0, req, res)
}

type updateRaidRequestV0 struct{ Raid *record.BasicRaid }
type updateRaidResponseV0 struct{ Updated *record.BasicRaid }

func (r raidClientStub) Update(ctx context.Context, raid record.Raid) (record.Raid, error) {
	req := &updateRaidRequestV0{record.NewBasicRaid(raid)}
	res := new(updateRaidResponseV0)
	err := call(http.MethodPatch, endpointRaidV0, req, res)
	return res.Updated, err
}

func HandleRaid(handler RaidHandler) func(mux *http.ServeMux) {
	return func(mux *http.ServeMux) {
		serve(mux, endpointRaidV0, func(ctx context.Context, method string, request *Request) (interface{}, error) {
			if strings.EqualFold(http.MethodGet, method) {
				res := new(getRaidsResponseV0)
				raids, err := handler.Fetch(ctx)
				for _, v := range raids {
					res.Raids = append(res.Raids, record.NewBasicRaid(v))
				}
				return res, err
			} else if strings.EqualFold(http.MethodPut, method) {
				var req addRaidRequestV0
				request.ReadTo(&req)
				res := new(addRaidRequestV0)
				err := handler.Add(ctx, req.Raid)
				return res, err
			} else if strings.EqualFold(http.MethodDelete, method) {
				var req deleteRaidRequestV0
				request.ReadTo(&req)
				res := new(deleteRaidResponseV0)
				err := handler.Delete(ctx, req.RaidId)
				return res, err
			} else if strings.EqualFold(http.MethodPatch, method) {
				var req updateRaidRequestV0
				request.ReadTo(&req)
				res := new(updateRaidResponseV0)
				updated, err := handler.Update(ctx, req.Raid)
				res.Updated = record.NewBasicRaid(updated)
				return res, err
			} else {
				return nil, httputil.UnsupportedMethod(method)
			}
		})
	}
}
