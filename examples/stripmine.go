package main

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/m4schini/computercraft-go/computer"
	"github.com/m4schini/computercraft-go/connection"
	"github.com/m4schini/computercraft-go/connection/adapter"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{} // use default options

func main() {
	http.HandleFunc("/api/ws", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		reader, _ := adapter.ReaderFromWebsocket(ws)
		conn := connection.New(reader, adapter.WriterFromWebsocket(ws))
		if deviceType, err := computer.GetDeviceType(ctx, conn); err != nil || deviceType != computer.DeviceTurtle {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		go StripMine(computer.NewTurtle(conn), 30)
		w.WriteHeader(http.StatusOK)
	})
	err := http.ListenAndServe("[::]:8080", nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func sliceForward(ctx context.Context, t computer.Turtle) error {
	dig := func() error {
		_, err := t.Dig(ctx)
		return err
	}

	digUp := func() error {
		_, err := t.DigUp(ctx)
		return err
	}

	forward := func() error {
		ok, err := t.Forward(ctx)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("couldn't move")
		}
		return nil
	}

	err := dig()
	if err != nil {
		return err
	}

	err = forward()
	if err != nil {
		return err
	}

	err = digUp()
	return err
}

func strip(ctx context.Context, t computer.Turtle, length int) (int, error) {
	for i := 0; i < length; i++ {
		err := sliceForward(ctx, t)
		if err != nil {
			return length - (i + 1), err
		}
	}

	return length, nil
}

func StripMine(t computer.Turtle, iterations int) {
	ctx := context.TODO()

	for i := 0; i < iterations; i++ {
		_, err := strip(ctx, t, 3)
		if err != nil {
			panic(err)
		}
		t.TurnRight(ctx)
		strip(ctx, t, 5)
		t.TurnRight(ctx)
		t.TurnRight(ctx)
		strip(ctx, t, 11)
		t.TurnRight(ctx)
		t.TurnRight(ctx)
		strip(ctx, t, 5)
		t.TurnLeft(ctx)
	}
}
