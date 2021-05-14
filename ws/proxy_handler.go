package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"os"
)

func proxyHandler(
	interrupt *chan os.Signal,
) func(http.ResponseWriter, *http.Request) {
	// TODO: kill multiple goroutines at once
	// https://stackoverflow.com/questions/61007385/golang-pattern-to-kill-multiple-goroutines-at-once

	return func(writer http.ResponseWriter, req *http.Request) {
		conn, deferFunc, err := upgradeConnection(writer, req)
		if err != nil {
			log.Print("failed to upgrade websocket:", err)
			return
		}
		defer deferFunc()

		// client channel
		fromClientCh := make(chan []byte)
		toClientCh := make(chan []byte)
		go readMessageFromClient(conn, &fromClientCh)
		go sendMessageToClient(conn, &toClientCh)


		// connect to server
		fromServerCh := make(chan []byte)
		toServerCh := make(chan []byte)

		// proxy server response
		go func() {
			for {
				msg := <-fromServerCh
				log.Printf("<- TARGET: %s", msg)
				toClientCh <- msg
			}
		}()

		// proxy client request
		go func() {
			for {
				msg := <-fromClientCh
				log.Printf("CLIENT ->:  %s", msg)
				toServerCh <- msg
			}
		}()

		go ConnectServer(url.URL{
			Scheme: "wss",
			Host:   "wss.gopax.co.kr",
			Path:   "/primus",
		}, interrupt, &fromServerCh, &toServerCh)


		// wait for interrupt
		done := make(chan struct{})
		for {
			select {
			case <-done:
				return
			case <-*interrupt:
				log.Println("interrupted.")
				return
			}
		}
	}
}

func readMessageFromClient(conn *websocket.Conn, ch *chan []byte) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			// handle server disconnection
			if _, ok := err.(*websocket.CloseError); ok {
				log.Print("disconnected from client")
				return
			}
			log.Print("failed to read from client:", err)
			return
		}

		*ch <- msg
	}
}

func sendMessageToClient(conn *websocket.Conn, ch *chan []byte) {
	for {
		msg := <- *ch

		err := conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Print("failed to write to client:", err)
			return
		}
	}
}
