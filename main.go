package main

import (
	"github.com/lechuckroh/websocket-proxy-server/ws"
	"os"
	"os/signal"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	ws.StartServer(8080, "/", &interrupt)
}
