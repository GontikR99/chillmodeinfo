package ipc

const Prefix ="golang-ipc-"

type Message interface {
	Content() []byte
	Reply(channelName string, data []byte)
}

type Endpoint interface {
	Listen(channelName string) <- chan Message
	Send(channelName string, data []byte)
}
