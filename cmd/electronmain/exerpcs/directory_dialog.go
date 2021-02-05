// +build wasm,electron

package exerpcs

import (
	"errors"
	"github.com/GontikR99/chillmodeinfo/internal/rpcidl"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/dialog"
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
		return filePaths[0], nil
	}))
}
