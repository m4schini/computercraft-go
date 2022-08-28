package cc

import (
	"cc-go/computer"
	"fmt"
	"log"
	"testing"
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
