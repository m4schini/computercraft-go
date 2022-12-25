package connection

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/m4schini/logger"
	"go.uber.org/zap"
	"io"
	"sync"
)

type DeviceType string

const (
	DeviceComputer       = "c"
	DevicePocketComputer = "p"
	DeviceTurtle         = "t"
)

var log = logger.Named("connection").Sugar()

type Connection interface {
	Execute(ctx context.Context, command string) ([]interface{}, error)
	Context() context.Context
	Device() DeviceType
	RemoteHost() string
	Id() string
	io.Closer
}

type HandshakeData struct {
	Id   string
	Host string
	Type DeviceType
}

type conn struct {
	In        <-chan []byte
	Out       chan<- []byte
	log       *zap.SugaredLogger
	closer    io.Closer
	ctx       context.Context
	handshake HandshakeData
	mu        sync.Mutex
}

func New(ctx context.Context, in <-chan []byte, out chan<- []byte, closer io.Closer) *conn {
	c := &conn{
		In:     in,
		Out:    out,
		closer: closer,
		ctx:    ctx,
	}
	c.handshake = c.doHandshake()
	c.log = log.With("host", c.handshake.Host, "id", c.handshake.Id, "type", c.handshake.Type)
	return c
}

func (c *conn) send(f string) error {
	c.log.Debugf("sending instruction: \"%v\"", f)
	bytes, err := json.Marshal(map[string]any{
		"func": fmt.Sprintf("return {%s}", f),
	})
	if err != nil {
		return err
	}

	c.Out <- bytes
	c.log.Debugf("send instruction \"%v\"", f)
	return err
}

func (c *conn) receive(ctx context.Context) ([]interface{}, error) {
	c.log.Debug("waiting for incoming message")
	res := make(chan []interface{})
	var e error
	go func() {
		buffer := <-c.In
		response := make([]interface{}, 0)
		err := json.Unmarshal(buffer, &response)
		if err != nil {
			e = err
			res <- []interface{}{}
		}

		e = nil
		res <- response
	}()

	select {
	case r := <-res:
		c.log.Debugf("received message: \"%v\"", r)
		return r, e
	case <-ctx.Done():
		err := ctx.Err()
		c.log.Debugw("receiving message failed", "err", err)
		return []interface{}{}, err
	}
}

func (c *conn) doHandshake() HandshakeData {
	var data = HandshakeData{}
	buffer := <-c.In

	var msg = make(map[string]interface{})
	err := json.Unmarshal(buffer, &msg)
	if err != nil {
		return data
	}

	_id, ok := msg["id"]
	if !ok {
		return data
	}
	id, ok := _id.(float64)
	if !ok {
		return data
	}
	data.Id = fmt.Sprintf("%v", id)

	_t, ok := msg["type"]
	if !ok {
		return data
	}
	t, ok := _t.(string)
	if !ok {
		return data
	}
	data.Type = DeviceType(t)

	return data
}

func (c *conn) Execute(ctx context.Context, command string) (response []interface{}, err error) {
	log.Infof("sending: %v", command)
	defer log.Infof("received: %v (err=%v)", response, err)
	c.mu.Lock()
	defer c.mu.Unlock()
	err = c.send(command)
	if err != nil {
		return nil, err
	}

	response, err = c.receive(ctx)
	return response, err
}

func (c *conn) Context() context.Context {
	if c.ctx == nil {
		return context.Background()
	}

	return c.ctx
}

func (c *conn) Device() DeviceType {
	return c.handshake.Type
}

func (c *conn) RemoteHost() string {
	return c.handshake.Host
}

func (c *conn) Id() string {
	return c.handshake.Id
}

func (c *conn) Close() error {
	c.log.Warn("closing connection")
	if c.closer != nil {
		return c.closer.Close()
	}

	return nil
}
