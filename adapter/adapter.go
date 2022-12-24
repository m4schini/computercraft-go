package adapter

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/m4schini/logger"
)

var log = logger.Named("adapter").Sugar()

func ReaderFromWebsocket(conn *websocket.Conn) <-chan []byte {
	ch := make(chan []byte, 8)
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			if !json.Valid(msg) {
				log.Warnw("incoming websocket message is not json", "payloadSize", len(msg), "remoteAddr", conn.RemoteAddr())
			}

			ch <- msg
		}
	}()
	return ch
}

func WriterFromWebsocket(conn *websocket.Conn) chan<- []byte {
	ch := make(chan []byte, 8)
	go func() {
		for msg := range ch {
			err := conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Warnw("Websocket message failed", "err", err)
			}
		}
	}()
	return ch
}
