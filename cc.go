package computercraft

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/m4schini/computercraft-go/computer"
	"github.com/m4schini/computercraft-go/connection"
	"github.com/m4schini/computercraft-go/connection/adapter"
)

// NewConnectionFromWebsocket creates a device from a websocket. If a new device was created, the bool return is true
func NewConnectionFromWebsocket(ws *websocket.Conn, opts ...connection.Option) (_ connection.Connection, err error) {
	if ws == nil {
		return nil, fmt.Errorf("ws cannot be nil")
	}
	o := connection.ParseOptions(opts)
	log := o.Log.Desugar()

	in, stop := adapter.ReaderFromWebsocket(ws, adapter.WithLog(log))
	out := adapter.WriterFromWebsocket(ws, adapter.WithLog(log))

	ws.SetCloseHandler(func(code int, text string) error {
		stop()
		close(out)
		return nil
	})

	return NewConnection(in, out, opts...)
}

// NewConnection uses a channel for incoming messages and outgoing messages.
// A `Message` means a valid json object/array.
// Incoming Messages are an array containing the return values of a lua function.
//
//	e.g. [23, 4, 23] for Locate()
//
// Outgoing Messages are an object with a "func" key and a value that that is lua code
// e.g. {"func": "return {turtle != nil}"}
func NewConnection(incomingMessages <-chan []byte, outgoingMessages chan<- []byte, opts ...connection.Option) (_ connection.Connection, err error) {
	if incomingMessages == nil || outgoingMessages == nil {
		return nil, fmt.Errorf("required parameter was nil")
	}
	conn := connection.New(
		incomingMessages,
		outgoingMessages,
		opts...,
	)

	return conn, nil
}

func NewComputer(conn connection.Connection) (computer.Computer, error) {
	if conn == nil {
		return nil, fmt.Errorf("connection is nil")
	}

	panic("not implemented")
}

func NewTurtle(conn connection.Connection) (computer.Turtle, error) {
	if conn == nil {
		return nil, fmt.Errorf("connection is nil")
	}

	panic("not implemented")
}

func NewPocket(conn connection.Connection) (computer.Computer, error) {
	if conn == nil {
		return nil, fmt.Errorf("connection is nil")
	}

	panic("not implemented")
}
