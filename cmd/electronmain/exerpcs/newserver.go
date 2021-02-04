package exerpcs

import "net/rpc"

var installers []func(server *rpc.Server)

func register(installer func(server *rpc.Server)) {
	installers = append(installers, installer)
}

func NewServer() *rpc.Server {
	server := rpc.NewServer()
	for _, installer := range installers {
		installer(server)
	}
	return server
}
