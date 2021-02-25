// +build web,wasm

package admin

import (
	"context"
	"github.com/GontikR99/chillmodeinfo/internal/comms/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/profile"
	"github.com/GontikR99/chillmodeinfo/internal/profile/localprofile"
	"github.com/GontikR99/chillmodeinfo/pkg/modal"
	"github.com/GontikR99/chillmodeinfo/pkg/toast"
	"github.com/vugu/vugu"
	"sort"
	"time"
)

type Userlist struct {
	Self         profile.Entry
	Requesters   []profile.Entry
	Elders       []profile.Entry
	Subordinates []profile.Entry
	Denied       []profile.Entry
	Done         chan struct{}
}

type byDisplayName []profile.Entry

func (b byDisplayName) Len() int           { return len(b) }
func (b byDisplayName) Less(i, j int) bool { return b[i].GetDisplayName() < b[j].GetDisplayName() }
func (b byDisplayName) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

func (c *Userlist) update(env vugu.EventEnv) {
	selfProfile := localprofile.GetProfile()
	if selfProfile == nil || selfProfile.GetAdminState() != profile.StateAdminApproved {
		return
	}
	entries, err := restidl.GetProfile().ListAdmins(context.Background())
	if err != nil {
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
			if entry.GetUserId() == selfProfile.GetUserId() {
				ul.Self = entry
			} else if entry.GetStartDate().Before(selfProfile.GetStartDate()) {
				ul.Elders = append(ul.Elders, entry)
			} else {
				ul.Subordinates = append(ul.Subordinates, entry)
			}
		}
	}
	sort.Sort(byDisplayName(ul.Requesters))
	sort.Sort(byDisplayName(ul.Elders))
	sort.Sort(byDisplayName(ul.Subordinates))
	sort.Sort(byDisplayName(ul.Denied))

	env.Lock()
	c.Requesters = ul.Requesters
	c.Elders = ul.Elders
	c.Self = ul.Self
	c.Subordinates = ul.Subordinates
	c.Denied = ul.Denied
	env.UnlockRender()
}

func (c *Userlist) Init(ctx vugu.InitCtx) {
	c.Self = &profile.BasicProfile{}
	c.Done = make(chan struct{})
	go func() {
		for {
			c.update(ctx.EventEnv())
			select {
			case <-c.Done:
				return
			case <-time.After(10 * time.Second):
			}
		}
	}()
}

func (c *Userlist) Destroy(ctx vugu.DestroyCtx) {
	if c.Done != nil {
		close(c.Done)
		c.Done = nil
	}
}

func (c *Userlist) ChangeState(event vugu.DOMEvent, userId string, newState profile.AdminState) {
	event.PreventDefault()
	go func() {
		if !modal.Verify("members", "Change Admin State", "Are you sure you want to update the admin state of this user?", "Change") {
			return
		}
		err := restidl.GetProfile().UpdateAdmin(context.Background(), userId, newState)
		if err == nil {
			selfProfile := localprofile.GetProfile()
			if selfProfile != nil && userId != selfProfile.GetUserId() {
				c.update(event.EventEnv())
			}
		} else {
			toast.Error("admin", err)
		}
	}()
}
