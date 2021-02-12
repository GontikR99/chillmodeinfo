// +build web,wasm

package admin

import (
	"context"
	"errors"
	"github.com/GontikR99/chillmodeinfo/cmd/webapp/localprofile"
	"github.com/GontikR99/chillmodeinfo/internal/profile"
	"github.com/GontikR99/chillmodeinfo/internal/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/toast"
	"github.com/vugu/vugu"
	"sort"
)

type Userlist struct {
	Self profile.Entry
	Requesters []profile.Entry
	Elders []profile.Entry
	Subordinates []profile.Entry
	Denied []profile.Entry
}

type byDisplayName []profile.Entry

func (b byDisplayName) Len() int {return len(b)}
func (b byDisplayName) Less(i, j int) bool {return b[i].GetDisplayName() < b[j].GetDisplayName()}
func (b byDisplayName) Swap(i, j int) {b[i], b[j] = b[j], b[i]}

func (c *Userlist) update(env vugu.EventEnv) {
	selfProfile := localprofile.GetProfile()
	if selfProfile == nil {
		toast.Error("profile", errors.New("No current profile available"))
		return
	}
	entries, err := restidl.GetProfile().ListAdmins(context.Background())
	if err!=nil {
		toast.Error("communications", err)
		return
	}
	ul := Userlist{}
	for _, entry := range entries {
		switch entry.GetAdminState() {
		case profile.StateAdminUnrequested:
		case profile.StateAdminDenied:
			ul.Denied = append(ul.Denied, entry)
		case profile.StateAdminRequested:
			ul.Requesters = append(ul.Requesters, entry)
		case profile.StateAdminApproved:
			if entry.GetUserId()==selfProfile.GetUserId() {
				ul.Self=entry
			} else if entry.GetStartDate().Before(selfProfile.GetStartDate()) {
				ul.Elders=append(ul.Elders, entry)
			} else {
				ul.Subordinates=append(ul.Subordinates, entry)
			}
		}
	}
	sort.Sort(byDisplayName(ul.Requesters))
	sort.Sort(byDisplayName(ul.Elders))
	sort.Sort(byDisplayName(ul.Subordinates))
	sort.Sort(byDisplayName(ul.Denied))

	env.Lock()
	*c=ul
	env.UnlockRender()
}

func (c *Userlist) Init(ctx vugu.InitCtx) {
	c.Self=&profile.BasicProfile{}
	go c.update(ctx.EventEnv())
}

func (c *Userlist) ChangeState(event vugu.DOMEvent, userId string, newState profile.AdminState) {
	event.PreventDefault()
	go func() {
		err := restidl.GetProfile().UpdateAdmin(context.Background(), userId, newState)
		if err==nil {
			c.update(event.EventEnv())
		} else {
			toast.Error("admin", err)
		}
	}()
}