package computercraft

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/m4schini/computercraft-go/computer"
	"github.com/m4schini/computercraft-go/connection"
	"sync"
)

var knownComputer = make(map[string]connection.Client)
var knownComputerMu sync.Mutex

func New(ws *websocket.Conn, remoteAddr string) (any, error) {
	conn, err := connection.NewWebsocketConnection(ws, remoteAddr)
	if err != nil {
		return nil, err
	}

	var client connection.Client
	var ok bool
	key := fmt.Sprintf("%v#%v", conn.RemoteHost(), conn.Id())
	knownComputerMu.Lock()
	client, ok = knownComputer[key]
	if !ok {
		client = connection.NewClient()
		knownComputer[key] = client
		knownComputerMu.Unlock()
		client.SetConnection(conn)
	} else {
		knownComputerMu.Unlock()
		client.SetConnection(conn)
	}

	switch conn.Device() {
	case connection.DeviceComputer:
		return computer.NewComputer(client), nil
	case connection.DevicePocketComputer:
		panic("not implemented")
	case connection.DeviceTurtle:
		return computer.NewTurtle(client), nil
	default:
		panic("unknown device type")
	}
}
