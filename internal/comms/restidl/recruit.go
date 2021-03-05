package restidl

import (
	"context"
	"github.com/GontikR99/chillmodeinfo/internal/comms/httputil"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"net/http"
	"strings"
)

const endpointRecruitV0="/rest/v0/recruit"

type recruitClientStub struct {}
var Recruit = &recruitClientStub{}

type RecruitHandler interface {
	Update(ctx context.Context, target record.RecruitmentTarget) error
	Fetch(ctx context.Context) ([]record.RecruitmentTarget, error)
}

type updateRecruitRequestV0 struct {Target *record.BasicRecruitmentTarget}
type updateRecruitResponseV0 struct{}

func (r recruitClientStub) Update(ctx context.Context, target record.RecruitmentTarget) error {
	req := &updateRecruitRequestV0{Target: record.NewBasicRecruitmentTarget(target)}
	res := new(updateRecruitResponseV0)
	return call(http.MethodPut, endpointRecruitV0, req, res)
}

type fetchRecruitRequestV0 struct{}
type fetchRecruitResponseV0 struct{Targets []*record.BasicRecruitmentTarget}

func (r recruitClientStub) Fetch(ctx context.Context) ([]record.RecruitmentTarget, error) {
	req := new(fetchRecruitRequestV0)
	res := new(fetchRecruitResponseV0)
	err := call(http.MethodGet, endpointRecruitV0, req, res)
	if err!=nil {
		return nil, err
	}
	var casted []record.RecruitmentTarget
	for _, v := range res.Targets {
		casted=append(casted, v)
	}
	return casted, nil
}

func HandleRecruit(handler RecruitHandler) func(mux *http.ServeMux) {
	return func(mux *http.ServeMux) {
		serve(mux, endpointRecruitV0, func(ctx context.Context, method string, request *Request) (interface{}, error) {
			if strings.EqualFold(http.MethodGet, method) {
				targets, err := handler.Fetch(ctx)
				res := new(fetchRecruitResponseV0)
				for _, v := range targets {
					res.Targets = append(res.Targets, record.NewBasicRecruitmentTarget(v))
				}
				return res, err
			} else if strings.EqualFold(http.MethodPut, method) {
				var req updateRecruitRequestV0
				request.ReadTo(&req)
				return new(updateRecruitResponseV0), handler.Update(ctx, req.Target)
			} else {
				return nil, httputil.UnsupportedMethod(method)
			}
		})
	}
}