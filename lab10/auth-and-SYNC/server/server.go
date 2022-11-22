package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		str := strings.ReplaceAll(string(p), "\n", "")
		log.Println("Client: ", str)

		if str == "login: styopa password: 0000" {
			if err := conn.WriteMessage(messageType, []byte("Correct")); err != nil {
				log.Println(err)
				return
			}
		} else if strings.Contains(str, "login") {
			if err := conn.WriteMessage(messageType, []byte("Wrong")); err != nil {
				log.Println(err)
				return
			}
		} else {
			//line := "ping -c 1 151.248.113.144"
			splitted := strings.Split(str, " ")
			command := exec.Command(splitted[0], splitted[1:]...)
			var out bytes.Buffer
			command.Stdout = &out
			if err := command.Run(); err != nil {
				log.Println(err)
				return
			}
			if err := conn.WriteMessage(messageType, out.Bytes()); err != nil {
				log.Println(err)
				return
			}
		}
	}

}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected!")
	err = ws.WriteMessage(1, []byte("You connected!"))
	if err != nil {
		log.Println(err)
	}

	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
