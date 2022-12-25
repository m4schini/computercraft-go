package computercraft

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/m4schini/computercraft-go/adapter"
	"github.com/m4schini/computercraft-go/computer"
	"github.com/m4schini/computercraft-go/connection"
	"sync"
)

var knownComputer = make(map[string]connection.Client)
var knownComputerMu sync.Mutex

// New creates a device from a websocket. If a new device was created, the bool return is true
func New(ctx context.Context, ws *websocket.Conn) (any, bool, error) {
	conn := connection.New(
		ctx,
		adapter.ReaderFromWebsocket(ws),
		adapter.WriterFromWebsocket(ws),
		ws,
	)

	var client connection.Client
	var ok bool
	var unknown bool
	key := fmt.Sprintf("%v#%v", conn.RemoteHost(), conn.Id())
	knownComputerMu.Lock()
	client, ok = knownComputer[key]
	if !ok {
		client = connection.NewClient()
		knownComputer[key] = client
		knownComputerMu.Unlock()
		unknown = true
		client.SetConnection(conn)
	} else {
		knownComputerMu.Unlock()
		unknown = false
		client.SetConnection(conn)
	}

	switch conn.Device() {
	case connection.DeviceComputer:
		return computer.NewComputer(client), unknown, nil
	case connection.DevicePocketComputer:
		panic("not implemented")
	case connection.DeviceTurtle:
		return computer.NewTurtle(client), unknown, nil
	default:
		panic("unknown device type")
	}
}
