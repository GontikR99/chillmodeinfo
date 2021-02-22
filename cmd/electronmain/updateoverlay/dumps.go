// +build wasm,electron

package updateoverlay

import (
	"bytes"
	"errors"
	"github.com/GontikR99/chillmodeinfo/internal/eqspec"
	"github.com/GontikR99/chillmodeinfo/internal/overlay"
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
						Enqueue(overlay.NewError("Couldn't read dump: "+err.Error()))
						return
					}

					raid, err := eqspec.ParseRaidDump(bytes.NewReader(contents))
					if err==nil {
						Enqueue(overlay.NewRaidDump(raid))
						return
					}

					members, err := eqspec.ParseGuildDump(bytes.NewReader(contents))
					if err==nil {
						Enqueue(overlay.NewGuildDump(members))
						return
					}

					Enqueue(overlay.NewError("Unrecognized Outputfile: "+filename))
				}()
			}
		}
	})
}