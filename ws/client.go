package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
)

// ConnectServer connects to a remote websocket server.
func ConnectServer(
	u url.URL,
	interrupt *chan os.Signal,
	recvCh *chan []byte,
	sendCh *chan []byte,
) {
	log.Printf("connecting to server: %s", u.String())
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("fail to connect server: ", err)
	}
	log.Printf("connected to %s", u.String())

	defer func(c *websocket.Conn) {
		_ = c.Close()
	}(conn)

	go readMessageFromServer(conn, recvCh)
	go sendMessageToServer(conn, sendCh)

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

func readMessageFromServer(conn *websocket.Conn, ch *chan []byte) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			// handle server disconnection
			if _, ok := err.(*websocket.CloseError); ok {
				log.Print("disconnected from server")
				return
			}
			log.Print("failed to read from server:", err)
			return
		}
		*ch <- msg
	}
}

func sendMessageToServer(conn *websocket.Conn, ch *chan []byte) {
	for {
		msg := <- *ch

		err := conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Print("failed to write to server:", err)
			return
		}
	}
}

