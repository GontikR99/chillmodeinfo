// +build wasm,electron

package eqfiles

import (
	"bytes"
	"context"
	"github.com/GontikR99/chillmodeinfo/internal/settings"
	"github.com/GontikR99/chillmodeinfo/pkg/nodejs/fs"
	"github.com/GontikR99/chillmodeinfo/pkg/nodejs/path"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	filenameMatch = regexp.MustCompile("^eqlog_([A-Za-z]*)_([A-Za-z]*).txt$")
	loglineMatch  = regexp.MustCompile("^\\[([^\\]]*)] (.*)$")
)

type ListenerHandle int

var nextListenerHandler = ListenerHandle(1)
var logListeners = make(map[ListenerHandle]func([]*LogEntry))

func RegisterLogsListener(listener func([]*LogEntry)) ListenerHandle {
	curId := nextListenerHandler
	nextListenerHandler++
	logListeners[curId] = listener
	return curId
}

func (h ListenerHandle) Release() {
	delete(logListeners, h)
}

var cancelLogScans = func() {}

func RestartLogScans() {
	cancelLogScans()
	var ctx context.Context
	ctx, cancelLogScans = context.WithCancel(context.Background())
	go readAllLogsLoop(ctx)
}

func readAllLogsLoop(ctx context.Context) {
	eqDir, _, _ := settings.LookupSetting(settings.EverQuestDirectory)
	eqLogDir := path.Join(eqDir, "Logs")
	seen := make(map[string]struct{})
	for {
		entries, err := fs.ReadDir(eqLogDir)
		if err == nil {
			for _, fi := range entries {
				if _, ok := seen[fi.Name()]; !ok {
					seen[fi.Name()] = struct{}{}
					if parts := filenameMatch.FindStringSubmatch(fi.Name()); parts != nil {
						go tailLog(ctx, path.Join(eqLogDir, fi.Name()), parts[1], parts[2])
					}
				}
			}
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(1000 * time.Millisecond):
		}
	}
}

func tailLog(ctx context.Context, filename string, character string, server string) {
	fd, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	fd.Seek(-1, io.SeekEnd)
	rdbuf := make([]byte, 1024)
	buffer := new(bytes.Buffer)
	for {
		cnt, _ := fd.Read(rdbuf)
		if cnt <= 0 {
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
		buffer.Write(rdbuf[0:cnt])
		ib := bytes.IndexByte(buffer.Bytes(), '\n')
		if ib >= 0 {
			buffer.ReadBytes('\n')
			break
		}
	}

	for {
		var entries []*LogEntry
		for ib := bytes.IndexByte(buffer.Bytes(), '\n'); ib >= 0; ib = bytes.IndexByte(buffer.Bytes(), '\n') {
			line, _ := buffer.ReadString('\n')
			line = strings.ReplaceAll(line, "\r", "")
			line = strings.ReplaceAll(line, "\n", "")
			if parts := loglineMatch.FindStringSubmatch(line); parts != nil {
				parsedTime, _ := time.Parse(time.ANSIC, parts[1])
				entry := &LogEntry{
					Character: character,
					Server:    server,
					Timestamp: parsedTime,
					Message:   parts[2],
				}
				entries = append(entries, entry)
			}
		}
		if len(entries)!=0 {
			for _, callback := range logListeners {
				callback(entries)
			}
		}
		cnt, _ := fd.Read(rdbuf)
		if cnt <= 0 {
			select {
			case <-ctx.Done():
				return
			default:
				continue
			}
		}
		buffer.Write(rdbuf[0:cnt])
	}
}
