package exerpcs

import (
	"errors"
	"github.com/GontikR99/chillmodeinfo/cmd/electronmain/overlaymap"
	"github.com/GontikR99/chillmodeinfo/internal/comms/rpcidl"
	"github.com/GontikR99/chillmodeinfo/internal/overlay"
)

type overlayServer struct{}

func (o overlayServer) CloseOverlay(name string) error {
	desc := overlaymap.Lookup(name)
	if desc == nil {
		return errors.New("No such overlay: " + name)
	} else {
		overlay.Close(desc.Page)
		return nil
	}
}

func (o overlayServer) PositionOverlay(name string) error {
	desc := overlaymap.Lookup(name)
	if desc == nil {
		return errors.New("No such overlay: " + name)
	} else {
		overlay.UpdateSizing(desc.Title, desc.Page)
		return nil
	}
}

func (o overlayServer) ResetOverlay(name string) error {
	desc := overlaymap.Lookup(name)
	if desc == nil {
		return errors.New("No such overlay: " + name)
	} else {
		overlay.Close(desc.Page)
		overlay.ClearPreferredSizing(desc.Page)
		return nil
	}
}

func init() {
	register(rpcidl.HandleOverlay(overlayServer{}))
}
