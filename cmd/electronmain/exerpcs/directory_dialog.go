// +build wasm,electron

package exerpcs

import (
	"errors"
	"github.com/GontikR99/chillmodeinfo/internal/console"
	"github.com/GontikR99/chillmodeinfo/internal/electron/dialog"
	"github.com/GontikR99/chillmodeinfo/internal/nodejs/fs"
	"github.com/GontikR99/chillmodeinfo/internal/rpcidl"
)

func init() {
	register(rpcidl.HandleDirectoryDialog(func(initial string) (string, error) {
		filePaths, err := dialog.ShowOpenDialog(&dialog.OpenOptions{
			Title:       "Select a directory",
			DefaultPath: initial,
			Properties:  &[]string{dialog.OpenDirectory, dialog.DontAddToRecent},
		})

		if err != nil {
			return "", err
		}
		if len(filePaths) != 1 {
			return "", errors.New("Expected single path")
		}
		go func() {
			dirents, err := fs.ReadDir(filePaths[0])
			if err == nil {
				for _, dirent := range dirents {
					console.Log(dirent.Name())
				}
			}
		}()
		return filePaths[0], nil
	}))
}
