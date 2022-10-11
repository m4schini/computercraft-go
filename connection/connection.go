package connection

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"io"
	"strings"
)

type Connection interface {
	UUID() string
	Execute(ctx context.Context, command string) ([]interface{}, error)
	Handshake() HandshakeData
	io.Closer
}

type HandshakeData struct {
	Id    string
	Label string
}

type websocketConnection struct {
	uuid        string
	ws          *websocket.Conn
	remoteAddr  string
	messageCh   chan<- *Message
	_hsData     HandshakeData
	_cancelLoop context.CancelFunc
}

func (w *websocketConnection) UUID() string {
	return w.uuid
}

func (w *websocketConnection) Execute(ctx context.Context, command string) ([]interface{}, error) {
	waitCh := make(chan []interface{})

	var response []interface{}

	w.messageCh <- &Message{
		Instruction: command,
		OnResponse:  waitCh,
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case response = <-waitCh:
		return response, nil
	}
}

func (w *websocketConnection) Handshake() HandshakeData {
	return w._hsData
}

func (w *websocketConnection) Close() error {
	if w._cancelLoop != nil {
		w._cancelLoop()
	}
	return w.ws.Close()
}

func NewWebsocketConnection(ws *websocket.Conn, remoteAddr string, onConnectionLost func()) (*websocketConnection, error) {
	c := new(websocketConnection)
	uid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	c.uuid = uid.String()
	c.ws = ws
	c.remoteAddr = remoteAddr

	hs, err := handleHandshake(c)
	if err != nil {
		return nil, err
	}
	c._hsData = hs

	ctx, cancel := context.WithCancel(context.Background())
	c._cancelLoop = cancel

	messageCh := startConnectionLoop(ctx, c.ws, onConnectionLost)
	c.messageCh = messageCh

	return c, nil
}

func handleHandshake(wsc *websocketConnection) (HandshakeData, error) {
	var handshakeMessage = make(map[string]interface{})

	err := wsc.ws.ReadJSON(&handshakeMessage)
	if err != nil {
		return HandshakeData{}, err
	}

	var idAsInt int64
	id, ok := handshakeMessage["id"]
	if ok {
		idAsFloat, ok := id.(float64)
		if !ok {
			idAsInt = int64(idAsFloat)
		}
	}

	return HandshakeData{
		Id: fmt.Sprintf("%v#%v",
			strings.Split(wsc.remoteAddr, ":")[0],
			idAsInt,
		),
		//Label: handshakeMessage["label"].(string), //TODO
	}, nil
}

type Message struct {
	Instruction string
	OnResponse  chan []interface{}
}

type ConnError struct {
	Message *Message
	Error   error
}

func startConnectionLoop(ctx context.Context, ws *websocket.Conn, onConnectionLost func()) chan<- *Message {
	messageCh := make(chan *Message, 8)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:

			}
			message := <-messageCh

			msg := make(map[string]string)
			msg["func"] = fmt.Sprintf("return {%s}", message.Instruction)

			err := ws.WriteJSON(msg)
			if err != nil {
				//errorCh <- &ConnError{
				//	Message: message,
				//	Error:   err,
				//}
				continue
			}

			//log.Printf("waiting for response to \"%v\"\n", message.Instruction)
			var responseJson = make([]interface{}, 0)
			err = ws.ReadJSON(&responseJson)
			if err != nil {
				switch err.(type) {
				case *websocket.CloseError:
					onConnectionLost()
					return
				default:
					//log.Printf("err: %T %v\n", err, err)
					//log.Printf("empty respond to \"%v\"\n", message.Instruction)
					message.OnResponse <- make([]interface{}, 0)
					continue
				}
				//errorCh <- &ConnError{
				//	Message: message,
				//	Error:   err,
				//}

			}

			//log.Printf("response to \"%v\": [%T]%v\n", message.Instruction, responseJson[0], responseJson)
			message.OnResponse <- responseJson
		}
	}()

	return messageCh
}
