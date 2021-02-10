// +build wasm

package place

import (
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
	go func() {
		env.Lock()
		defer env.UnlockRender()
		js.Global().Get("window").Get("history").Call("pushState", nil, "", "#"+place)
	}()
}
