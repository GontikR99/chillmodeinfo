// +build dev

package static

import "net/http"

//go:generate go run -mod=vendor github.com/shurcooL/vfsgen/cmd/vfsgendev -source="github.com/GontikR99/chillmodeinfo/web/static".StaticFiles

var StaticFiles = http.Dir("data")

