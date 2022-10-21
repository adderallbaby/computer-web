package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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
		// print out that message for clarity
		str := strings.ReplaceAll(string(p), "\n", "")
		log.Println("Client: ", str)
		arr := make([]int, 0)
		for i := 0; i < len(str); i++ {
			if (str[i] > 47 && str[i] < 58) || str[i] == 45 {
				var buf string
				buf += string(str[i])
				for i++; i < len(str) && ((str[i] > 47 && str[i] < 58) || str[i] == 45); i++ {
					buf += string(str[i])
				}
				num, _ := strconv.Atoi(buf)
				arr = append(arr, num)
			}
		}

		var res = -1
		// const EPS = 1e-12
		for i := 0; i < len(arr); i++ {
			if arr[i] > res {
				res = arr[i]
			}
		}
		answ := strconv.Itoa(res)
		if err := conn.WriteMessage(messageType, []byte("answer is "+answ)); err != nil {
			log.Println(err)
			return
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
