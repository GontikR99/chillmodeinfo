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
}

const endpointDKPLog="/rest/v0/dkp"

type dkpLogHandlerClient struct {}

type dkpLogAppendRequestV0 struct {
	Entry record.BasicDKPChangeEntry
}
type dkpLogAppendResponseV0 struct{}

func (d *dkpLogHandlerClient) Append(ctx context.Context, entry record.DKPChangeEntry) error {
	req:=&dkpLogAppendRequestV0{Entry: *record.NewBasicDKPChangeEntry(entry)}
	res:=new(dkpLogAppendResponseV0)
	return call(http.MethodPut, endpointDKPLog, req, res)
}

type dkpLogRetrieveRequestV0 struct {
	Target string
}
type dkpLogRetrieveResponseV0 struct {
	Entries []record.BasicDKPChangeEntry
}

func (d *dkpLogHandlerClient) Retrieve(ctx context.Context, target string) ([]record.DKPChangeEntry, error) {
	req:=&dkpLogRetrieveRequestV0{Target: target}
	res:=new(dkpLogRetrieveResponseV0)
	err := call(http.MethodGet, endpointDKPLog, req, res)
	if err!=nil {
		return nil, err
	}
	var result []record.DKPChangeEntry
	for _, v := range res.Entries {
		result=append(result,record.NewBasicDKPChangeEntry(&v))
	}
	return result, nil
}

var DKPLog=&dkpLogHandlerClient{}

func HandleDKPLog(handler DKPLogHandler) func(mux *http.ServeMux) {
	return func(mux *http.ServeMux) {
		serve(mux, endpointDKPLog, func(ctx context.Context, method string, request *Request) (interface{}, error) {
			if strings.EqualFold(http.MethodGet, method) {
				var req dkpLogRetrieveRequestV0
				res:=new(dkpLogRetrieveResponseV0)
				request.ReadTo(&req)
				entries, err := handler.Retrieve(ctx, req.Target)
				for _, v:=range entries {
					res.Entries=append(res.Entries, *record.NewBasicDKPChangeEntry(v))
				}
				return res, err
			} else if strings.EqualFold(http.MethodPut, method) {
				var req dkpLogAppendRequestV0
				res:=new(dkpLogAppendResponseV0)
				request.ReadTo(&req)
				err := handler.Append(ctx, &req.Entry)
				return res, err
			} else {
				return nil, httputil.UnsupportedMethod(method)
			}
		})
	}
}