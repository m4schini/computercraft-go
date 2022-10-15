package connection

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/m4schini/cc-go/logger"
	"io"
	"strings"
)

var log = logger.Sub("connection").Sugar()

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
	log.Infow("command execution started",
		"command", command,
		"remoteAddr", w.remoteAddr,
		"uuid", w.uuid)
	waitCh := make(chan []interface{})

	var response []interface{}

	w.messageCh <- &Message{
		Instruction: command,
		OnResponse:  waitCh,
	}

	select {
	case <-ctx.Done():
		err := ctx.Err()
		log.Errorw("command execution failed",
			"remoteAddr", w.remoteAddr,
			"uuid", w.uuid,
			"error", err)
		return nil, err
	case response = <-waitCh:
		log.Infow("command execution succeeded",
			"remoteAddr", w.remoteAddr,
			"uuid", w.uuid)
		return response, nil
	}
}

func (w *websocketConnection) Handshake() HandshakeData {
	return w._hsData
}

func (w *websocketConnection) Close() error {
	log.Debugw("trying to close websocket connection")
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

	messageCh := startConnectionLoop(ctx, hs, c.ws, onConnectionLost)
	c.messageCh = messageCh

	log.Debugw("created websocket connection",
		"remoteAddr", remoteAddr,
		"uuid", remoteAddr)
	return c, nil
}

func handleHandshake(wsc *websocketConnection) (HandshakeData, error) {
	log.Debugw("handshake started",
		"remoteAddr", wsc.remoteAddr,
		"uuid", wsc.uuid)
	var handshakeMessage = make(map[string]interface{})

	err := wsc.ws.ReadJSON(&handshakeMessage)
	if err != nil {
		log.Debugw("handshake failed",
			"remoteAddr", wsc.remoteAddr,
			"uuid", wsc.uuid,
			"error", err)
		return HandshakeData{}, err
	}
	log.Debugw("received websocket handshake",
		"remoteAddr", wsc.remoteAddr,
		"message", handshakeMessage)

	var idAsInt int64
	id, ok := handshakeMessage["id"]
	if ok {
		idAsFloat, ok := id.(float64)
		if ok {
			idAsInt = int64(idAsFloat)
		} else {
			log.Warnw("handshake message id is not a number",
				"remoteAddr", wsc.remoteAddr)
		}
	} else {
		log.Warnw("handshake message didn't contain id",
			"remoteAddr", wsc.remoteAddr)
	}

	log.Debugw("handshake succeeded",
		"remoteAddr", wsc.remoteAddr,
		"turtleId", idAsInt,
		"uuid", wsc.uuid)
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

func startConnectionLoop(ctx context.Context, hs HandshakeData, ws *websocket.Conn, onConnectionLost func()) chan<- *Message {
	messageCh := make(chan *Message, 8)
	wsRemoteAddr := ws.RemoteAddr()

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Debugw("websocket connection context was canceled",
					"remoteAddr", wsRemoteAddr,
					"turtleId", hs.Id,
					"error", ctx.Err())
				onConnectionLost()
				return
			default:

			}
			message := <-messageCh

			msg := make(map[string]string)
			msg["func"] = fmt.Sprintf("return {%s}", message.Instruction)

			log.Debugw("sending message over websocket",
				"remoteAddr", wsRemoteAddr,
				"turtleId", hs.Id,
				"message", msg["func"])
			err := ws.WriteJSON(msg)
			if err != nil {
				log.Errorw("error while trying to send message",
					"remoteAddr", wsRemoteAddr,
					"turtleId", hs.Id,
					"error", err)
				continue
			}

			log.Debugw("waiting for websocket response",
				"remoteAddr", wsRemoteAddr,
				"turtleId", hs.Id)
			var responseJson = make([]interface{}, 0)
			_, bytes, err := ws.ReadMessage()
			if err != nil {
				log.Errorw("error while trying to receive message",
					"remoteAddr", wsRemoteAddr,
					"turtleId", hs.Id,
					"error", err)
				onConnectionLost()
				return
			}

			log.Debugw("trying to parse response",
				"remoteAddr", wsRemoteAddr,
				"turtleId", hs.Id,
			)
			err = json.Unmarshal(bytes, &responseJson)
			if err != nil {
				log.Errorw("error while trying to unmarshal response",
					"remoteAddr", wsRemoteAddr,
					"turtleId", hs.Id,
					"raw", string(bytes),
					"error", err)
				message.OnResponse <- make([]interface{}, 0)
				continue
			}

			log.Debugw("received response",
				"remoteAddr", wsRemoteAddr,
				"turtleId", hs.Id,
				"response", string(bytes),
			)
			message.OnResponse <- responseJson
		}
	}()

	return messageCh
}
