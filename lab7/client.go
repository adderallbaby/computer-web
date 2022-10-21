package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"

	"github.com/gorilla/websocket"
)

func main() {
	messageOut := make(chan string)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8080", Path: "/ws"}
	log.Printf("connecting to %s", u.String())
	c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("handshake failed with status %d", resp.StatusCode)
		log.Fatal("dial:", err)
	}

	// When the program closes close the connection
	defer c.Close()
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("server: %s", message)
			if string(message) == "You connected!" || strings.Contains(string(message), "Answer") || string(message) == "Wrong input" {
				var x int
				var str string
				fmt.Scanf("%d", &x)
				if x == 0 {
					err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
					if err != nil {
						log.Println("write close:", err)
						return
					}
				} else {
					in := bufio.NewReader(os.Stdin)
					for i := 0; i < x; i++ {
						buf, _ := in.ReadString('\n')
						str += buf + " "
					}
					messageOut <- str
				}
			}
		}

	}()

	for {
		select {
		case m := <-messageOut:
			err := c.WriteMessage(websocket.TextMessage, []byte(m))
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
	}
}
