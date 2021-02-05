package overlaymap

type Description struct {
	Name  string
	Title string
	Page  string
}

var overlays = []*Description{
	{"bid", "Bid Overlay", "overlay_bids.html"},
	{"guilddump", "Guild Dump Overlay", "overlay_guild_dump.html"},
	{"raiddump", "Raid Dump Overlay", "overlay_raid_dump.html"},
}

func Lookup(name string) *Description {
	for _, desc := range overlays {
		if name == desc.Name {
			return desc
		}
	}
	return nil
}
