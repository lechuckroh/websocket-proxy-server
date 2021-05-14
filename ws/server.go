package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
)

var upgrader = websocket.Upgrader{}

func upgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, func(), error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, nil, err
	}
	deferFunc := func() {
		_ = conn.Close
	}
	return conn, deferFunc, nil

}

// StartServer starts a websocket server.
func StartServer(
	port uint16,
	path string,
	interrupt *chan os.Signal,
) {
	http.HandleFunc(path, proxyHandler(interrupt))

	addr := fmt.Sprintf("localhost:%d", port)
	log.Printf("starting server on %s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
