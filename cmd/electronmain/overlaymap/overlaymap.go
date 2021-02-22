package overlaymap

type Description struct {
	Name  string
	Title string
	Page  string
}

var overlays = []*Description{
	{"bid", "Bid Overlay", "overlay_bids.html"},
	{"update", "Update Overlay", "overlay_update.html"},
}

func Lookup(name string) *Description {
	for _, desc := range overlays {
		if name == desc.Name {
			return desc
		}
	}
	return nil
}
