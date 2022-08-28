package cc

import (
	"cc-go/computer"
	"cc-go/connection"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{} // use default options

var onNewTurtles = make([]func(id string, t computer.Turtle), 0)

func onNewTurtle(t computer.Turtle) {
	id, _ := t.ComputerId()

	for _, f := range onNewTurtles {
		go f(id, t)
	}
}

func Serve(addr string) error {
	http.HandleFunc("/connect/turtle", connectTurtleHandler)
	return http.ListenAndServe(addr, nil)
}

func OnTurtleConnected(f func(id string, t computer.Turtle)) {
	onNewTurtles = append(onNewTurtles, f)
}

func connectTurtleHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("incoming connection")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	conn, err := connection.NewWebsocketConnection(c, r.RemoteAddr)
	if err != nil {
		log.Println(err)
		return
	}

	t := computer.MakeTurtle(conn)
	onNewTurtle(t)
}
