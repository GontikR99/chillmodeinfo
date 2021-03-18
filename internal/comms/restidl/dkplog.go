package restidl

import (
	"context"
	"github.com/GontikR99/chillmodeinfo/internal/comms/httputil"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"net/http"
	"strings"
)

type DKPLogHandler interface {
	Append(ctx context.Context, entry record.DKPChangeEntry) error
	Retrieve(ctx context.Context, target string) ([]record.DKPChangeEntry, error)
	Update(ctx context.Context, entry record.DKPChangeEntry) (record.DKPChangeEntry, error)
	Remove(ctx context.Context, entryId uint64) error
	Sync(ctx context.Context) (string, error)
}

const endpointDKPLog = "/rest/v0/dkp"
const endpointDKPSync = "/rest/v0/dkpSync"

type dkpLogHandlerClient struct{}

type dkpLogAppendRequestV0 struct {
	Entry record.BasicDKPChangeEntry
}
type dkpLogAppendResponseV0 struct{}

func (d *dkpLogHandlerClient) Append(ctx context.Context, entry record.DKPChangeEntry) error {
	req := &dkpLogAppendRequestV0{Entry: *record.NewBasicDKPChangeEntry(entry)}
	res := new(dkpLogAppendResponseV0)
	return call(http.MethodPut, endpointDKPLog, req, res)
}

type dkpLogRetrieveRequestV0 struct {
	Target string
}
type dkpLogRetrieveResponseV0 struct {
	Entries []record.BasicDKPChangeEntry
}

func (d *dkpLogHandlerClient) Retrieve(ctx context.Context, target string) ([]record.DKPChangeEntry, error) {
	req := &dkpLogRetrieveRequestV0{Target: target}
	res := new(dkpLogRetrieveResponseV0)
	err := call(http.MethodGet, endpointDKPLog, req, res)
	if err != nil {
		return nil, err
	}
	var result []record.DKPChangeEntry
	for _, v := range res.Entries {
		result = append(result, record.NewBasicDKPChangeEntry(&v))
	}
	return result, nil
}

type dkpLogRemoveRequestV0 struct {
	EntryId uint64
}
type dkpLogRemoveResponseV0 struct{}

func (c *dkpLogHandlerClient) Remove(ctx context.Context, entryId uint64) error {
	req := &dkpLogRemoveRequestV0{entryId}
	res := new(dkpLogRemoveResponseV0)
	return call(http.MethodDelete, endpointDKPLog, req, res)
}

type dkpLogUpdateRequestV0 struct {
	NewEntry record.BasicDKPChangeEntry
}
type dkpLogUpdateResponseV0 struct {
	UpdatedEntry record.BasicDKPChangeEntry
}

func (c *dkpLogHandlerClient) Update(ctx context.Context, newEntry record.DKPChangeEntry) (record.DKPChangeEntry, error) {
	req := &dkpLogUpdateRequestV0{*record.NewBasicDKPChangeEntry(newEntry)}
	res := new(dkpLogUpdateResponseV0)
	err := call(http.MethodPatch, endpointDKPLog, req, res)
	if err != nil {
		return nil, err
	} else {
		return &res.UpdatedEntry, nil
	}
}

type dkpLogSyncRequestV0 struct {}
type dkpLogSyncResponseV0 struct {
	Message string
}

func (c *dkpLogHandlerClient) Sync(ctx context.Context) (string, error) {
	req := new(dkpLogSyncRequestV0)
	res := new(dkpLogSyncResponseV0)
	err := call(http.MethodPost, endpointDKPSync, req, res)
	return res.Message, err
}

var DKPLog = &dkpLogHandlerClient{}

func HandleDKPLog(handler DKPLogHandler) func(mux *http.ServeMux) {
	return func(mux *http.ServeMux) {
		serve(mux, endpointDKPLog, func(ctx context.Context, method string, request *Request) (interface{}, error) {
			if strings.EqualFold(http.MethodGet, method) {
				var req dkpLogRetrieveRequestV0
				res := new(dkpLogRetrieveResponseV0)
				request.ReadTo(&req)
				entries, err := handler.Retrieve(ctx, req.Target)
				for _, v := range entries {
					res.Entries = append(res.Entries, *record.NewBasicDKPChangeEntry(v))
				}
				return res, err
			} else if strings.EqualFold(http.MethodPut, method) {
				var req dkpLogAppendRequestV0
				res := new(dkpLogAppendResponseV0)
				request.ReadTo(&req)
				err := handler.Append(ctx, &req.Entry)
				return res, err
			} else if strings.EqualFold(http.MethodDelete, method) {
				var req dkpLogRemoveRequestV0
				res := new(dkpLogRemoveResponseV0)
				request.ReadTo(&req)
				err := handler.Remove(ctx, req.EntryId)
				return res, err
			} else if strings.EqualFold(http.MethodPatch, method) {
				var req dkpLogUpdateRequestV0
				res := new(dkpLogUpdateResponseV0)
				request.ReadTo(&req)
				update, err := handler.Update(ctx, &req.NewEntry)
				res.UpdatedEntry = *record.NewBasicDKPChangeEntry(update)
				return res, err
			} else {
				return nil, httputil.UnsupportedMethod(method)
			}
		})
		serve(mux, endpointDKPSync, func(ctx context.Context, method string, request *Request) (interface{}, error) {
			if strings.EqualFold(http.MethodPost, method) {
				msg, err := handler.Sync(ctx)
				return &dkpLogSyncResponseV0{Message: msg}, err
			} else {
				return nil, httputil.UnsupportedMethod(method)
			}
		})
	}
}
