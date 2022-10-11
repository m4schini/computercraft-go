package cc

import (
	"fmt"
	"github.com/m4schini/cc-go/computer"
	"log"
	"testing"
	"time"
)

func TestServe(t *testing.T) {
	OnTurtleConnected(func(id string, t computer.Turtle) {
		log.Println(id)
		log.Println(t.Version())
		log.Println(t.ComputerLabel())

		scan := func() {
			x, y, z, err := t.Locate()
			if err != nil {
				return
			}
			ok, data, err := t.InspectDown()
			if ok {
				fmt.Printf("[%v,%v,%v] %v\n", x, y-1, z, data)
			}
		}

		size := 8
		nextRight := true
		for i := 0; i < size; i++ {
			for i := 0; i < size; i++ {
				scan()

				t.Forward()
			}

			if nextRight {
				t.TurnRight()
			} else {
				t.TurnLeft()
			}
			scan()
			t.Forward()
			if nextRight {
				t.TurnRight()
			} else {
				t.TurnLeft()
			}

			nextRight = !nextRight
		}

		t.TurnLeft()
		for i := 0; i < size; i++ {
			t.Forward()
		}
		t.TurnRight()
	})

	t.Fatal(Serve("0.0.0.0:8080"))
}

func TestReconnectM(t *testing.T) {
	OnTurtleConnected(func(id string, turtle computer.Turtle) {
		t.Log("new turtle connected:")
		for true {
			rid, err1 := turtle.ComputerId()
			label, err2 := turtle.ComputerLabel()
			x, y, z, err3 := turtle.LocateWithTimeout(1 * time.Second)
			if err1 != nil || err2 != nil || err3 != nil {
				t.Log(id, err1, err2, err3)
			} else {
				t.Logf("turtle is still alive: %v=%v (%v %v %v)", rid, label, x, y, z)
			}

			time.Sleep(5 * time.Second)
		}
	})
	OnTurtleDisconnected(func(uuid string, turtle computer.Turtle) {
		t.Log(uuid, "DISCONNECTED")
	})

	t.Fatal(Serve("0.0.0.0:8080"))
}
