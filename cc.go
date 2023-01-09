package computercraft

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/m4schini/computercraft-go/computer"
	"github.com/m4schini/computercraft-go/connection"
	"github.com/m4schini/computercraft-go/connection/adapter"
)

// New creates a device from a websocket. If a new device was created, the bool return is true
func New(ws *websocket.Conn) (computer.Computer, error) {
	in, stop := adapter.ReaderFromWebsocket(ws)
	out := adapter.WriterFromWebsocket(ws)

	ws.SetCloseHandler(func(code int, text string) error {
		stop()
		close(out)
		return nil
	})

	conn := connection.New(
		in,
		out,
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
