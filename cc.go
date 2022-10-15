package cc

import (
	"github.com/gorilla/websocket"
	"github.com/m4schini/cc-go/computer"
	"github.com/m4schini/cc-go/connection"
	"github.com/m4schini/cc-go/logger"
	"go.uber.org/zap"
	"net/http"
)

var log = logger.Sub("main").Sugar()

var upgrader = websocket.Upgrader{} // use default options

var onNewTurtles = make([]func(remoteAddr, uuid string, turtle computer.Turtle), 0)
var onLostTurtles = make([]func(remoteAddr, uuid string, turtle computer.Turtle), 0)

func onNewTurtle(remoteAddr, uuid string, t computer.Turtle) {
	log.Infow("turtle connected",
		"remoteAddr", remoteAddr,
		"uuid", uuid)
	for _, f := range onNewTurtles {
		go f(remoteAddr, uuid, t)
	}
}

func onLostTurtle(remoteAddr, uuid string, t computer.Turtle) {
	log.Infow("turtle disconnected",
		"remoteAddr", remoteAddr,
		"uuid", uuid)
	for _, f := range onLostTurtles {
		go f(remoteAddr, uuid, t)
	}
}

func Serve(addr string) error {
	http.HandleFunc("/connect/turtle", connectTurtleHandler)
	return http.ListenAndServe(addr, nil)
}

func OnTurtleConnected(f func(remoteAddr, uuid string, t computer.Turtle)) {
	onNewTurtles = append(onNewTurtles, f)
}

func OnTurtleDisconnected(f func(remoteAddr, uuid string, t computer.Turtle)) {
	onLostTurtles = append(onLostTurtles, f)
}

func UseLogger(newLogger *zap.Logger) {
	logger.UseLogger(newLogger)
}

func connectTurtleHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugw("new incoming websocket connection request",
		"remoteAddr", r.RemoteAddr)
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		//log.Print("upgrade:", err)
		return
	}

	var turtle computer.Turtle

	conn, err := connection.NewWebsocketConnection(c, r.RemoteAddr, func() {
		onLostTurtle(r.RemoteAddr, turtle.UUID(), turtle)
	})
	if err != nil {
		log.Error(err)
		return
	}

	turtle = computer.MakeTurtle(conn)
	if err != nil {
		log.Error(err)
		return
	}
	onNewTurtle(r.RemoteAddr, turtle.UUID(), turtle)
}
