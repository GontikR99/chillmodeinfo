// +build wasm,electron

package updateoverlay

import (
	"bytes"
	"errors"
	"github.com/GontikR99/chillmodeinfo/internal/eqspec"
	"github.com/GontikR99/chillmodeinfo/internal/overlay/update"
	"github.com/GontikR99/chillmodeinfo/internal/settings"
	"github.com/GontikR99/chillmodeinfo/pkg/nodejs/path"
	"github.com/GontikR99/chillmodeinfo/pkg/toast"
	"io/ioutil"
	"regexp"
)

var filenameRe=regexp.MustCompile("^Outputfile Complete: (.*)$")

func init() {
	eqspec.RegisterLogsListener(func(entries []*eqspec.LogEntry) {
		for _, entry := range entries {
			if m := filenameRe.FindStringSubmatch(entry.Message); m!=nil {
				filename := m[1]
				go func() {
					eqDir, present, err := settings.LookupSetting(settings.EverQuestDirectory)
					if err!=nil || !present {
						toast.Error("dump scanning", errors.New("Hrm, I knew how to read your log files, but now can't remember where your EQ directory is."))
						return
					}
					dumpFullPath := path.Join(eqDir, filename)
					contents, err := ioutil.ReadFile(dumpFullPath)
					if err!=nil {
						Enqueue(update.NewError("Couldn't read dump: "+err.Error()))
						return
					}

					attendees, err := eqspec.ParseRaidDump(bytes.NewReader(contents))
					if err==nil {
						if len(attendees)==0 {
							Enqueue(update.NewError("Empty raid dump: "+filename))
							return
						} else {
							Enqueue(update.NewRaidDump(attendees))
							return
						}
					}

					members, err := eqspec.ParseGuildDump(bytes.NewReader(contents))
					if err==nil {
						Enqueue(update.NewGuildDump(members))
						return
					}

					Enqueue(update.NewError("Unrecognized Outputfile: "+filename))
				}()
			}
		}
	})
}