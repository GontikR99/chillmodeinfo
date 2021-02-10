// +build wasm,electron

package settings

import (
	"encoding/json"
	"github.com/GontikR99/chillmodeinfo/pkg/console"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/application"
	"github.com/GontikR99/chillmodeinfo/pkg/nodejs/path"
	"io/ioutil"
	"os"
)

var memoizedSettings map[string]string
var dbDir = path.Join(application.GetPath("appData"), "chillmodeinfo")
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
	return sync()
}

// Remove a setting from the map
func ClearSetting(key string) error {
	delete(memoizedSettings, key)
	return sync()
}

// Write the settings to disk
func sync() error {
	data, err := json.Marshal(&memoizedSettings)
	if err!=nil {
		return err
	}
	err = ioutil.WriteFile(dbFilename+".tmp", data, 0600)
	if err!=nil {
		return err
	}
	err = os.Rename(dbFilename+".tmp", dbFilename)
	return err
}

func init() {
	memoizedSettings=make(map[string]string)
	err := os.Mkdir(dbDir, 0700)
	if err!=nil && !os.IsExist(err) {
		console.Logf("Failed to create directory: %v", err)
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