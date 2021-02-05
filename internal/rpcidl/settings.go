package rpcidl

import "net/rpc"

type LookupSettingRequest struct {
	Key string
}

type LookupSettingResponse struct {
	Value   string
	Present bool
}

type SetSettingRequest struct {
	Key   string
	Value string
}

type SetSettingResponse struct{}

type ClearSettingRequest struct {
	Key string
}
type ClearSettingResponse struct{}

type SettingsServer interface {
	Lookup(key string) (value string, present bool, err error)
	Set(key string, value string) error
	Clear(key string) error
}

type SettingsClient struct {
	settings SettingsServer
}

func (ss *SettingsClient) LookupSetting(req *LookupSettingRequest, res *LookupSettingResponse) (err error) {
	res.Value, res.Present, err = ss.settings.Lookup(req.Key)
	return
}

func (ss *SettingsClient) SetSetting(req *SetSettingRequest, res *SetSettingResponse) (err error) {
	err = ss.settings.Set(req.Key, req.Value)
	return
}

func (ss *SettingsClient) ClearSetting(req *ClearSettingRequest, res *ClearSettingResponse) (err error) {
	err = ss.settings.Clear(req.Key)
	return
}

func LookupSetting(client *rpc.Client, key string) (string, bool, error) {
	req := &LookupSettingRequest{key}
	res := new(LookupSettingResponse)
	err := client.Call("SettingsClient.LookupSetting", req, res)
	return res.Value, res.Present, err
}

func SetSetting(client *rpc.Client, key string, value string) error {
	req := &SetSettingRequest{Key: key, Value: value}
	res := new(SetSettingResponse)
	return client.Call("SettingsClient.SetSetting", req, res)
}

func ClearSetting(client *rpc.Client, key string) error {
	return client.Call("SettingsClient.ClearSetting", &ClearSettingRequest{key}, new(ClearSettingResponse))
}

func HandleSetting(setting SettingsServer) func(*rpc.Server) {
	ss := &SettingsClient{setting}
	return func(server *rpc.Server) {
		server.Register(ss)
	}
}
