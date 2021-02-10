// +build server

package serverrpcs

import "net/http"

var installers []func(mux *http.ServeMux)

func register(callback func(mux *http.ServeMux)) {
	installers = append(installers, callback)
}

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	for _, installer := range installers {
		installer(mux)
	}
	return mux
}
