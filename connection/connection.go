package connection

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/m4schini/logger"
	"go.uber.org/zap"
	"io"
	"strings"
)

var log = logger.Named("connection").Sugar()

type Connection interface {
	Execute(ctx context.Context, command string) ([]interface{}, error)
	Context() context.Context
	io.Closer
}

type HandshakeData struct {
	Id    string
	Host string
}

type websocketConnection struct {
	log *zap.SugaredLogger
	ws          *websocket.Conn
	messageCh   chan<- *Message
	ctx context.Context
	cancelFunc context.CancelFunc
	_hsData    HandshakeData
}

func NewWebsocketConnection(ws *websocket.Conn, remoteAddr string) (*websocketConnection, error) {
	c := new(websocketConnection)
	c.log = log.With("remoteAddr", remoteAddr)
	c.ws = ws

	hs, err := handleHandshake(c.log, c, remoteAddr)
	if err != nil {
		return nil, err
	}
	c._hsData = hs

	ctx, cancel := context.WithCancel(context.Background())
	loopCtx, loopCancel := context.WithCancel(ctx)
	messageCh := startConnectionLoop(loopCtx, hs, c.ws, func() {
		cancel()
	})
	c.ctx = ctx
	c.cancelFunc = func() {
		loopCancel()
		cancel()
	}
	c.messageCh = messageCh

	log.Debug("created websocket connection")
	return c, nil
}

func (w *websocketConnection) Context() context.Context {
	return w.ctx
}

func (w *websocketConnection) Execute(ctx context.Context, command string) ([]interface{}, error) {
	log.Debugw("command execution started","command", command)
	waitCh := make(chan []interface{})

	var response []interface{}

	w.messageCh <- &Message{
		Instruction: command,
		OnResponse:  waitCh,
	}

	select {
	case <-ctx.Done():
		err := ctx.Err()
		log.Errorw("command execution failed","error", err)
		return nil, err
	case response = <-waitCh:
		log.Debugw("command execution succeeded")
		return response, nil
	}
}

func (w *websocketConnection) Handshake() HandshakeData {
	return w._hsData
}

func (w *websocketConnection) Close() error {
	log.Debugw("trying to close websocket connection")
	if w.cancelFunc != nil {
		w.cancelFunc()
	}
	return w.ws.Close()
}

func handleHandshake(log *zap.SugaredLogger, wsc *websocketConnection, remoteAddr string) (HandshakeData, error) {
	var handshakeMessage = make(map[string]interface{})

	err := wsc.ws.ReadJSON(&handshakeMessage)
	if err != nil {
		return HandshakeData{}, err
	}

	var idAsInt int64
	id, ok := handshakeMessage["id"]
	if ok {
		idAsFloat, ok := id.(float64)
		if ok {
			idAsInt = int64(idAsFloat)
		} else {
			log.Warnw("handshake message id is not a number")
		}
	} else {
		log.Warnw("handshake message didn't contain id")
	}

	log.Debugw("handshake succeeded", "turtleId", idAsInt)
	return HandshakeData{
		Id: fmt.Sprintf("%v", idAsInt),
		Host: strings.Split(remoteAddr, ":")[0],
	}, nil
}

type Message struct {
	Instruction string
	OnResponse  chan []interface{}
}

type ConnError struct {
	Message *Message
	err   error
}

func (c *ConnError) Error() string {
	return c.err.Error()
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

