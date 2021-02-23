package restidl

import (
	"context"
	"github.com/GontikR99/chillmodeinfo/internal/comms/httputil"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"net/http"
	"strings"
)

const endpointMembersAllV0 = "/rest/v0/members"
const endpointMemberSingleV0 = "/rest/v0/member"

type MemberHandler interface {
	GetMember(ctx context.Context, name string) (record.Member, error)
	MergeMember(ctx context.Context, member record.Member) (record.Member, error)

	GetMembers(context.Context) (map[string]record.Member, error)
	MergeMembers(context.Context, []record.Member) (map[string]record.Member, error)
}

type membersClientStub struct{}

var Members = &membersClientStub{}

type getMembersRequestV0 struct{}
type getMembersResponseV0 struct {
	Members map[string]*record.BasicMember
}

func (m membersClientStub) GetMembers(ctx context.Context) (map[string]record.Member, error) {
	req := new(getMembersRequestV0)
	res := new(getMembersResponseV0)
	err := call(http.MethodGet, endpointMembersAllV0, req, res)
	if err != nil {
		return nil, err
	}
	if res.Members == nil {
		return nil, nil
	}
	result := make(map[string]record.Member)
	for k, v := range res.Members {
		result[k] = v
	}
	return result, err
}

type mergeMembersRequestV0 struct {
	Members []*record.BasicMember
}
type mergeMembersResponseV0 struct {
	Members map[string]*record.BasicMember
}

func (m membersClientStub) MergeMembers(ctx context.Context, members []record.Member) (map[string]record.Member, error) {
	req := new(mergeMembersRequestV0)
	res := new(mergeMembersResponseV0)
	for _, v := range members {
		req.Members = append(req.Members, record.NewBasicMember(v))
	}
	err := call(http.MethodPut, endpointMembersAllV0, req, res)
	if err != nil {
		return nil, err
	}
	if res.Members == nil {
		return nil, nil
	}
	result := make(map[string]record.Member)
	for k, v := range res.Members {
		result[k] = v
	}
	return result, nil
}

type getMemberRequestV0 struct {
	Name string
}
type getMemberResponseV0 struct {
	Member *record.BasicMember
}

func (m membersClientStub) GetMember(ctx context.Context, name string) (record.Member, error) {
	req := &getMemberRequestV0{name}
	res := new(getMemberResponseV0)
	err := call(http.MethodGet, endpointMemberSingleV0, req, res)
	if err != nil {
		return nil, err
	}
	return res.Member, nil
}

type mergeMemberRequestV0 struct {
	Member *record.BasicMember
}
type mergeMemberResponseV0 struct {
	Member *record.BasicMember
}

func (m membersClientStub) MergeMember(ctx context.Context, member record.Member) (record.Member, error) {
	req := &mergeMemberRequestV0{record.NewBasicMember(member)}
	res := new(mergeMemberResponseV0)
	err := call(http.MethodPut, endpointMemberSingleV0, req, res)
	if err != nil {
		return nil, err
	}
	return res.Member, nil
}

func HandleMembers(handler MemberHandler) func(mux *http.ServeMux) {
	return func(mux *http.ServeMux) {
		serve(mux, endpointMemberSingleV0, func(ctx context.Context, method string, request *Request) (interface{}, error) {
			if strings.EqualFold(http.MethodGet, method) {
				var req getMemberRequestV0
				request.ReadTo(&req)
				member, err := handler.GetMember(ctx, req.Name)
				return &getMemberResponseV0{Member: record.NewBasicMember(member)}, err
			} else if strings.EqualFold(http.MethodPut, method) {
				var req mergeMemberRequestV0
				request.ReadTo(&req)
				member, err := handler.MergeMember(ctx, req.Member)
				return &mergeMemberResponseV0{record.NewBasicMember(member)}, err
			} else {
				return nil, httputil.UnsupportedMethod(method)
			}
		})
		serve(mux, endpointMembersAllV0, func(ctx context.Context, method string, request *Request) (interface{}, error) {
			if strings.EqualFold(http.MethodGet, method) {
				members, err := handler.GetMembers(ctx)
				res := new(getMembersResponseV0)
				if members != nil {
					res.Members = make(map[string]*record.BasicMember)
					for k, v := range members {
						res.Members[k] = record.NewBasicMember(v)
					}
				}
				return res, err
			} else if strings.EqualFold(http.MethodPut, method) {
				var req mergeMembersRequestV0
				request.ReadTo(&req)
				newMembers := []record.Member{}
				for _, v := range req.Members {
					newMembers = append(newMembers, v)
				}
				members, err := handler.MergeMembers(ctx, newMembers)
				res := new(mergeMembersResponseV0)
				if members != nil {
					res.Members = make(map[string]*record.BasicMember)
					for k, v := range members {
						res.Members[k] = record.NewBasicMember(v)
					}
				}
				return res, err
			} else {
				return nil, httputil.UnsupportedMethod(method)
			}
		})
	}
}
