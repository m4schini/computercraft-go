package computercraft

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/m4schini/computercraft-go/adapter"
	"github.com/m4schini/computercraft-go/computer"
	"github.com/m4schini/computercraft-go/connection"
)

// New creates a device from a websocket. If a new device was created, the bool return is true
func New(ws *websocket.Conn) (computer.Computer, error) {
	conn := connection.New(
		adapter.ReaderFromWebsocket(ws),
		adapter.WriterFromWebsocket(ws),
		ws,
	)

	c := computer.NewComputer(conn)
	if isTurtle, err := c.IsTurtle(context.TODO()); err == nil && isTurtle {
		return computer.NewTurtle(conn), nil
	}

	if isPocket, err := c.IsPocket(context.TODO()); err == nil && isPocket {
		panic("not implemented")
	}

	return c, nil
}
