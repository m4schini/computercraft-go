package adapter

import (
	"github.com/gorilla/websocket"
	"github.com/m4schini/logger"
	"io"
)

var log = logger.Named("adapter")

type reader struct {
	conn *websocket.Conn
}

func (r *reader) Read(p []byte) (int, error) {
	log.Debug("READING WEBSOCKET MESSAGE")
	_, reader, err := r.conn.NextReader()
	if err != nil {
		return 0, err
	}

	var size int
	for {
		log.Debug("WAITING FOR NEXT PART")
		var buffer []byte
		n, err := reader.Read(buffer)
		if err == io.EOF {
			log.Debug("RECEIVED EOF")
			return size, nil
		}
		if err != nil {
			return size, err
		}

		size = size + n
		p = append(p, buffer...)
	}
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
	log.Debug("WRITING WEBSOCKET MESSAGE")
	out, err := w.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return 0, err
	}

	log.Debug("WRITING BYTES")
	n, err = out.Write(p)
	if err != nil {
		return n, err
	}
	log.Debug("CLOSING MESSAGE")
	return n, out.Close()
}

func (w *writer) Close() error {
	return w.conn.Close()
}

func WriterFromWebsocket(conn *websocket.Conn) io.WriteCloser {
	w := &writer{}
	w.conn = conn
	return w
}
