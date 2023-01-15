package adapter

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

func ReaderFromWebsocket(conn *websocket.Conn, opts ...Option) (in <-chan []byte, stop func()) {
	o := parseOptions(opts)
	log := o.log

	var closed = false
	ch := make(chan []byte, 8)

	closeF := func() {
		close(ch)
		closed = true
	}

	go func() {
		defer func() {
			x := recover()
			if x != nil {
				log.Warnf("%v: %v", conn.RemoteAddr(), x)
			}
		}()

		for !closed {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				//TODO check if this is enough error handling
				log.Warnw("Websocket message failed", "err", err)
				break
			}

			if !json.Valid(msg) {
				log.Warnw("incoming websocket message is not json", "payloadSize", len(msg), "remoteAddr", conn.RemoteAddr())
			}

			ch <- msg
		}

		closeF()
	}()
	return ch, closeF
}

func WriterFromWebsocket(conn *websocket.Conn, opts ...Option) (out chan<- []byte) {
	o := parseOptions(opts)
	log := o.log

	ch := make(chan []byte, 8)
	go func() {
		for msg := range ch {
			err := conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				//TODO check if this is enough error handling
				log.Warnw("Websocket message failed", "err", err)
				return
			}
		}
	}()
	return ch
}
