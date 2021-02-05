// +build wasm

package msgcomm

import (
	"bytes"
	"container/heap"
	"encoding/gob"
	"github.com/GontikR99/chillmodeinfo/pkg/console"
	"io"
	"sync"
)

type endpointWrapper struct {
	endpoint   Endpoint
	readLock   sync.Mutex
	writeLock  sync.Mutex
	channel    string
	recv       <-chan Message
	done       func()
	unread     *bytes.Buffer
	nextSeqOut int64
	nextSeqIn  int64
	msgQueue   []*SequencedMessage
	closed     bool
}

func EndpointAsStream(channel string, endpoint Endpoint) io.ReadWriteCloser {
	listener, listenerDone := endpoint.Listen(channel)
	bufferedListener := make(chan Message, 65536)
	go func() {
		for {
			inMsg := <-listener
			if inMsg == nil {
				close(bufferedListener)
				return
			} else {
				bufferedListener <- inMsg
			}
		}
	}()
	ew := &endpointWrapper{
		endpoint: endpoint,
		channel:  channel,
		recv:     bufferedListener,
		done:     listenerDone,
		unread:   new(bytes.Buffer),
	}
	return ew
}

func (ew *endpointWrapper) Write(p []byte) (int, error) {
	ew.writeLock.Lock()
	defer ew.writeLock.Unlock()
	outbuf := new(bytes.Buffer)
	encoder := gob.NewEncoder(outbuf)
	seqMsg := &SequencedMessage{
		SequenceNumber: ew.nextSeqOut,
		Content:        p,
	}
	ew.nextSeqOut++
	err := encoder.Encode(seqMsg)
	if err != nil {
		return 0, err
	}
	ew.endpoint.Send(ew.channel, outbuf.Bytes())
	return len(p), nil
}

func (ew *endpointWrapper) Read(p []byte) (int, error) {
	ew.readLock.Lock()
	defer ew.readLock.Unlock()
	for ew.unread.Len() == 0 {
		inMsg := <-ew.recv
		if inMsg == nil {
			return 0, io.EOF
		}
		decoder := gob.NewDecoder(bytes.NewReader(inMsg.Content()))
		inSeqMsg := new(SequencedMessage)
		err := decoder.Decode(inSeqMsg)
		if err != nil {
			console.Log("Decoding error", err)
			return 0, err
		}
		heap.Push(ew, inSeqMsg)
		for len(ew.msgQueue) > 0 && ew.msgQueue[0].SequenceNumber <= ew.nextSeqIn {
			nextSeqMsg := heap.Pop(ew).(*SequencedMessage)
			if nextSeqMsg.SequenceNumber == ew.nextSeqIn {
				ew.nextSeqIn++
				ew.unread.Write(nextSeqMsg.Content)
			}
		}
	}
	return ew.unread.Read(p)
}

func (ew *endpointWrapper) Close() error {
	if !ew.closed {
		ew.done()
		ew.closed = true
	}
	return nil
}

type SequencedMessage struct {
	SequenceNumber int64
	Content        []byte
}

func (ew *endpointWrapper) Len() int {
	return len(ew.msgQueue)
}

func (ew *endpointWrapper) Less(i, j int) bool {
	return ew.msgQueue[i].SequenceNumber < ew.msgQueue[j].SequenceNumber
}

func (ew *endpointWrapper) Swap(i, j int) {
	ew.msgQueue[i], ew.msgQueue[j] = ew.msgQueue[j], ew.msgQueue[i]
}

func (ew *endpointWrapper) Push(x interface{}) {
	ew.msgQueue = append(ew.msgQueue, x.(*SequencedMessage))
}

func (ew *endpointWrapper) Pop() interface{} {
	last := ew.msgQueue[len(ew.msgQueue)-1]
	ew.msgQueue = ew.msgQueue[:len(ew.msgQueue)-1]
	return last
}
