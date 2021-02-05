// +build wasm,electron

package fs

import "github.com/GontikR99/chillmodeinfo/pkg/nodejs"

var fs= nodejs.Require("fs")
var fsPromises= fs.Get("promises")