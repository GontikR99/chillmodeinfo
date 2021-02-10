package signins

const (
	TokenGoogle   = "google_token:"
	TokenClientId = "client_id:"

	IdTypeGoogle = "google_id:"
)

const ChannelSignins = "signin-updates"

type signinAnnouncement struct {
	IsSignedIn bool
	Token      string
}
