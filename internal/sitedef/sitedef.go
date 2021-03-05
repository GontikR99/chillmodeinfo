package sitedef

import (
	"strconv"
	"time"
)

const DNSName = "chillmode.info"
const Port = 443

var SiteURL = "https://" + DNSName + ":" + strconv.Itoa(Port)

const GoogleSigninClientId = "465672423976-qn8u1junpmdanfnlbc2ne3dobfmf4nvp.apps.googleusercontent.com"

const InactiveDuration = -7 * 24 * time.Hour
