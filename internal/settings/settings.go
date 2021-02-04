// +build wasm,electron

package settings

import (
	"encoding/json"
	"github.com/GontikR99/chillmodeinfo/internal/console"
	"github.com/GontikR99/chillmodeinfo/internal/nodejs/path"
	"github.com/GontikR99/chillmodeinfo/internal/nodejs/process"
	"io/ioutil"
	"os"
)

var memoizedSettings map[string]string
var dbDir = path.Join(process.Getenv("LOCALAPPDATA"), "chillmodeinfo")
var dbFilename = path.Join(dbDir, "settings.json")

// Set a setting value, but only if there's not one currently present
func DefaultSetting(key string, value string) {
	_, p, err := LookupSetting(key)
	if err==nil && !p {
		SetSetting(key, value)
	}
}

// Get the current value of a setting
func LookupSetting(key string) (string, bool, error) {
	value, present := memoizedSettings[key]
	return value, present, nil
}

// Upsert a value into a setting
func SetSetting(key string, value string) error {
	memoizedSettings[key]=value
	data, err := json.Marshal(&memoizedSettings)
	if err!=nil {
		return err
	}
	return ioutil.WriteFile(dbFilename, data, 0600)
}

func init() {
	memoizedSettings=make(map[string]string)
	err := os.Mkdir(dbDir, 0700)
	if err!=nil && !os.IsExist(err) {
		console.Logf("%v", err)
		return
	}

	data, err := ioutil.ReadFile(dbFilename)
	if err!=nil {
		if os.IsNotExist(err) {
			return
		} else {
			console.Logf("Failed to read settings file: %v", err)
		}
	}
	err = json.Unmarshal(data, &memoizedSettings)
	if err!=nil {
		console.Logf("Failed to unmarshal settings file: %v", err)
	}
}