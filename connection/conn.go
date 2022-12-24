package connection

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

type conn struct {
	In        io.Reader
	Out       io.Writer
	closer    io.Closer
	ctx       context.Context
	handshake HandshakeData
}

func New(ctx context.Context, in io.Reader, out io.Writer, closer io.Closer) *conn {
	c := &conn{
		In:     in,
		Out:    out,
		closer: closer,
		ctx:    ctx,
	}
	c.handshake = c.doHandshake()
	return c
}

func (c *conn) send(f string) error {
	bytes, err := json.Marshal(map[string]any{
		"func": fmt.Sprintf("return {%s}", f),
	})
	if err != nil {
		return err
	}

	_, err = c.Out.Write(bytes)
	return err
}

func (c *conn) receive(ctx context.Context) ([]interface{}, error) {
	res := make(chan []interface{})
	var e error
	go func() {
		buffer, err := io.ReadAll(c.In)
		if err != nil {
			e = err
			res <- []interface{}{}
		}

		response := make([]interface{}, 0)
		err = json.Unmarshal(buffer, &response)
		if err != nil {
			e = err
			res <- []interface{}{}
		}
	}()

	select {
	case r := <-res:
		return r, e
	case <-ctx.Done():
		return []interface{}{}, ctx.Err()
	}
}

func (c *conn) doHandshake() HandshakeData {
	var data = HandshakeData{}
	var buffer []byte
	_, err := c.In.Read(buffer)
	if err != nil {
		return data
	}

	var msg = make(map[string]interface{})
	err = json.Unmarshal(buffer, &msg)
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

func (c *conn) Execute(ctx context.Context, command string) ([]interface{}, error) {
	err := c.send(command)
	if err != nil {
		return nil, err
	}

	return c.receive(ctx)
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
	if c.closer != nil {
		return c.closer.Close()
	}

	return nil
}
