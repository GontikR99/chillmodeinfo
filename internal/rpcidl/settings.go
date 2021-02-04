package rpcidl

import "net/rpc"

type LookupSettingRequest struct {
	Key string
}

type LookupSettingResponse struct {
	Value string
	Present bool
}

type SetSettingRequest struct {
	Key string
	Value string
}

type SetSettingResponse struct {}

type SettingStub struct {
	lookup func(key string) (value string, present bool, err error)
	set    func(key string, value string) error
}

func (ss *SettingStub) LookupSetting(req *LookupSettingRequest, res *LookupSettingResponse) (err error) {
	res.Value, res.Present, err = ss.lookup(req.Key)
	return
}

func (ss *SettingStub) SetSetting(req *SetSettingRequest, res *SetSettingResponse) (err error) {
	err = ss.set(req.Key, req.Value)
	return
}

func LookupSetting(client *rpc.Client, key string) (string, bool, error) {
	req:=&LookupSettingRequest{key}
	res:=new(LookupSettingResponse)
	err := client.Call("SettingStub.LookupSetting", req, res)
	return res.Value, res.Present, err
}

func SetSetting(client *rpc.Client, key string, value string) error {
	req := &SetSettingRequest{Key: key, Value: value}
	res := new(SetSettingResponse)
	return client.Call("SettingStub.SetSetting", req, res)
}

func HandleSetting(lookup func(string)(string,bool,error), set func(string, string)error) func (*rpc.Server) {
	ss := &SettingStub{lookup, set}
	return func(server *rpc.Server) {
		server.Register(ss)
	}
}