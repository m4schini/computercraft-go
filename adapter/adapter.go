package adapter

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/m4schini/logger"
	"io"
)

var log = logger.Named("adapter")

type reader struct {
	conn        *websocket.Conn
	buf         bytes.Buffer
	readMessage bool
}

func (r *reader) Read(p []byte) (int, error) {
	n, err := r.buf.Read(p)
	if err == io.EOF {
		if r.readMessage {
			r.readMessage = false
			return n, err
		}
		_, res, err := r.conn.ReadMessage()
		r.readMessage = true
		if err == io.EOF {
			return n, err
		}
		if err != nil {
			return n, err
		}
		r.buf.Write(res)
		return n, nil
	}
	return n, err
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
