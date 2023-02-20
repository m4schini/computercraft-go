package test

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/m4schini/logger"
	"time"
)

var log = logger.Named("mock").Sugar()

type ComputerMock struct {
	Url string
}

func (c *ComputerMock) Boot(ctx context.Context) error {
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, c.Url, nil)
	if err != nil {
		return err
	}

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Error("Error in receive:", err)
				return
			}
			log.Info("Received: %s\n", msg)
		}
	}()

	time.Sleep(1 * time.Minute)
	return nil
}
