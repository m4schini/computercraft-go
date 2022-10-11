package cc

import (
	"github.com/gorilla/websocket"
	"github.com/m4schini/cc-go/computer"
	"github.com/m4schini/cc-go/connection"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{} // use default options

var onNewTurtles = make([]func(uuid string, t computer.Turtle), 0)
var onLostTurtles = make([]func(uuid string, t computer.Turtle), 0)

func onNewTurtle(uuid string, t computer.Turtle) {
	for _, f := range onNewTurtles {
		go f(uuid, t)
	}
}

func onLostTurtle(uuid string, t computer.Turtle) {
	for _, f := range onLostTurtles {
		go f(uuid, t)
	}
}

func Serve(addr string) error {
	http.HandleFunc("/connect/turtle", connectTurtleHandler)
	return http.ListenAndServe(addr, nil)
}

func OnTurtleConnected(f func(uuid string, t computer.Turtle)) {
	onNewTurtles = append(onNewTurtles, f)
}

func OnTurtleDisconnected(f func(uuid string, t computer.Turtle)) {
	onLostTurtles = append(onLostTurtles, f)
}

func connectTurtleHandler(w http.ResponseWriter, r *http.Request) {
	//log.Println("incoming connection")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		//log.Print("upgrade:", err)
		return
	}

	var turtle computer.Turtle

	conn, err := connection.NewWebsocketConnection(c, r.RemoteAddr, func() {
		onLostTurtle(turtle.UUID(), turtle)
	})
	if err != nil {
		log.Println(err)
		return
	}

	turtle = computer.MakeTurtle(conn)
	if err != nil {
		log.Println(err)
		return
	}
	onNewTurtle(turtle.UUID(), turtle)
}
