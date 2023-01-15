package computercraft

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/m4schini/computercraft-go/computer"
	"github.com/m4schini/computercraft-go/connection"
	"github.com/m4schini/computercraft-go/connection/adapter"
	"time"
)

// NewFromWebsocket creates a device from a websocket. If a new device was created, the bool return is true
func NewFromWebsocket(ws *websocket.Conn, opts ...connection.Option) (_ computer.Computer, err error) {
	o := connection.ParseOptions(opts)
	log := o.Log.Desugar()

	in, stop := adapter.ReaderFromWebsocket(ws, adapter.WithLog(log))
	out := adapter.WriterFromWebsocket(ws, adapter.WithLog(log))

	ws.SetCloseHandler(func(code int, text string) error {
		stop()
		close(out)
		return nil
	})

	return New(in, out, opts...)
}

// New uses a channel for incoming messages and outgoing messages.
// A `Message` means a valid json object/array.
// Incoming Messages are an array containing the return values of a lua function.
//
//	e.g. [23, 4, 23] for Locate()
//
// Outgoing Messages are an object with a "func" key and a value that that is lua code
// e.g. {"func": "return {turtle != nil}"}
func New(incomingMessages <-chan []byte, outgoingMessages chan<- []byte, opts ...connection.Option) (_ computer.Computer, err error) {
	conn := connection.New(
		incomingMessages,
		outgoingMessages,
		opts...,
	)

	ctx, stop := context.WithTimeout(context.Background(), 4*time.Second)
	defer stop()

	c := computer.NewComputer(conn)
	if isTurtle, err := c.IsTurtle(ctx); err == nil && isTurtle {
		return computer.NewTurtle(conn), nil
	}

	if isPocket, err := c.IsPocket(ctx); err == nil && isPocket {
		panic("not implemented")
	}

	return c, nil
}
