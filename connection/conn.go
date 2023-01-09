package connection

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/m4schini/logger"
	"go.uber.org/zap"
	"sync"
)

var log = logger.Named("connection").Sugar()

type Connection interface {
	Execute(ctx context.Context, command string) ([]interface{}, error)
}

type conn struct {
	In  <-chan []byte
	Out chan<- []byte
	log *zap.SugaredLogger
	mu  sync.Mutex
}

func New(in <-chan []byte, out chan<- []byte) *conn {
	c := &conn{
		In:  in,
		Out: out,
		log: log,
	}
	return c
}

func (c *conn) send(f string) (err error) {
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("%v", x)
		}
	}()

	c.log.Debugf("sending instruction: \"%v\"", f)
	c.Out <- []byte(fmt.Sprintf(`{"func": "return {%s}"}`, f))
	c.log.Debugf("send instruction \"%v\"", f)
	return err
}

func (c *conn) receive(ctx context.Context) ([]interface{}, error) {
	c.log.Debug("waiting for incoming message")
	res := make(chan []interface{})
	var e error
	go func() {
		buffer, ok := <-c.In
		if !ok {
			c.log.Warn("tried to receive response on closed channel")
			e = ClosedChannelErr
			res <- []interface{}{}
			return
		}

		response := make([]interface{}, 0)
		err := json.Unmarshal(buffer, &response)
		if err != nil {
			c.log.Error(err)
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
		c.log.Debugw("waiting for incoming message timed out!", "err", err)
		return []interface{}{}, err
	}
}

func (c *conn) Execute(ctx context.Context, command string) (response []interface{}, err error) {
	c.log.Infof("Execute Started: %v", command)
	c.mu.Lock()
	defer c.mu.Unlock()
	err = c.send(command)
	if err != nil {
		c.log.Errorw("Execute Failed", "err", err)
		return nil, err
	}

	response, err = c.receive(ctx)
	c.log.Infof("Execute Finished: %+v (error: %v)", response, err)
	return response, err
}
