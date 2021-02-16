// +build wasm,electron

package exerpcs

import (
	"errors"
	"github.com/GontikR99/chillmodeinfo/internal/profile/signins"
	"github.com/GontikR99/chillmodeinfo/internal/comms/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/comms/rpcidl"
	"github.com/GontikR99/chillmodeinfo/internal/settings"
	"github.com/GontikR99/chillmodeinfo/internal/sitedef"
	"github.com/GontikR99/chillmodeinfo/pkg/nodejs"
)

type serverLoginHandler struct{}

var electronShell = nodejs.Require("electron").Get("shell")

func (s *serverLoginHandler) SignIn() error {
	clientId, present, err := settings.LookupSetting(settings.ClientId)
	if err != nil {
		return err
	}
	if !present {
		return errors.New("No clientId")
	}

	electronShell.Call("openExternal", sitedef.SiteURL+"/associate.html?"+clientId)
	return nil
}

func (s *serverLoginHandler) SignOut() error {
	err := restidl.DisassociateClientId()
	if err != nil {
		return err
	}
	signins.ClearToken()
	return nil
}

func (s *serverLoginHandler) PollSignIn() error {
	clientId, present, err := settings.LookupSetting(settings.ClientId)
	if err != nil {
		return err
	}
	if !present {
		return errors.New("No clientId")
	}

	err = restidl.AssertAssociated(clientId)
	if err != nil {
		return err
	}
	signins.SetToken(signins.TokenClientId + clientId)
	return nil
}

func init() {
	sl := &serverLoginHandler{}
	sl.PollSignIn()
	register(rpcidl.HandleSignIn(sl))
}
