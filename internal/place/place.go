// +build wasm

package place

import (
	"github.com/GontikR99/chillmodeinfo/internal/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/rpcidl"
	"github.com/GontikR99/chillmodeinfo/pkg/console"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/ipc/ipcrenderer"
	"github.com/vugu/vugu"
	"github.com/vugu/vugu/js"
	"log"
	"net/url"
)

func GetPlace() string {
	href := js.Global().Get("window").Get("location").Get("href").String()
	parsed, err := url.Parse(href)
	if err!=nil {
		log.Println(err)
		return ""
	}
	return parsed.Fragment
}

func NavigateTo(env vugu.EventEnv, place string) {
	if ipcrenderer.Client !=nil {
		go rpcidl.Ping(ipcrenderer.Client, "Visiting "+place)
	}
	go func() {
		val, err := restidl.PingV0("Visiting "+place)
		console.Log("pong:", val, err)
	}()

	go func() {
		env.Lock()
		defer env.UnlockRender()
		js.Global().Get("window").Get("history").Call("pushState", nil, "", "#"+place)
	}()
}
