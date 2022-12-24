package adapter

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/m4schini/logger"
	"io"
)

var log = logger.Named("adapter")

type reader struct {
	conn *websocket.Conn
}

func (r *reader) Read(p []byte) (int, error) {
	_, res, err := r.conn.ReadMessage()
	if err != nil {
		return len(res), err
	}
	log.Debug(fmt.Sprintf("read bytes (%v) from websocket", len(p)))

	p = append(p, res...)
	return len(res), nil
}

func (r *reader) Close() error {
	return r.conn.Close()
}

func ReaderFromWebsocket(conn *websocket.Conn) io.ReadCloser {
	r := &reader{}
	r.conn = conn
	return r
}

type writer struct {
	conn *websocket.Conn
}

func (w *writer) Write(p []byte) (n int, err error) {
	l := len(p)
	log.Debug(fmt.Sprintf("writing bytes (%v) to websocket", l))
	return l, w.conn.WriteMessage(websocket.TextMessage, p)
}

func (w *writer) Close() error {
	return w.conn.Close()
}

func WriterFromWebsocket(conn *websocket.Conn) io.WriteCloser {
	w := &writer{}
	w.conn = conn
	return w
}
