package eqfiles

import "time"

const ChannelEQLog = "eqLogMsg"

type LogEntry struct {
	Character string
	Server    string
	Timestamp time.Time
	Message   string
}
