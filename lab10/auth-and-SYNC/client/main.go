// forms.go
package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

type ContactDetails struct {
	login    string
	password string
}

func main() {
	fl := true
	var x, y, login, buf string

	var message []byte

	buffer, _ := ioutil.ReadFile("output.html")

	tmpl := template.Must(template.ParseFiles("forms.html"))
	var Password string
	Password = ""
	fmt.Println(Password)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if fl == true && r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		} else if fl == false && r.Method != http.MethodPost {
			tmpl = template.Must(template.ParseFiles("output.html"))
			tmpl.Execute(w, nil)
			return
		}

		if fl == true {
			details := ContactDetails{
				login:    r.FormValue("Login"),
				password: r.FormValue("Password"),
			}
			x = details.login
			y = details.password
			login = x
			Password = y
		} else {
			x = r.FormValue("Message")
		}

		tmpl.Execute(w, nil)

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
		defer c.Close()

		go func() {
			for {
				_, message, err = c.ReadMessage()
				if err != nil {
					panic(err)
				}
				log.Printf("server: %s", message)
				if string(message) == "Correct" && fl == true {
					fl = false
					tmpl = template.Must(template.ParseFiles("output.html"))
					tmpl.Execute(w, nil)
				}

				fmt.Println(fl, x, y, login)
				if fl == true && string(message) == "You connected!" || string(message) == "Wrong" {
					fmt.Println(fl)
					str := "login: " + x + " password: " + y
					messageOut <- str
				} else if x != login {
					if string(message) != "You connected!" {
						ioutil.WriteFile("output.html", []byte(string(buffer)+"<p>"+string(message)+"</p>"), 0644)
					}
					messageOut <- x
				}
			}

		}()

		for {
			select {
			case m := <-messageOut:
				if m != buf {
					err = c.WriteMessage(websocket.TextMessage, []byte(m))
					if err != nil {
						panic(err)
					}
				}
				buf = m
			}
		}
	})

	http.ListenAndServe(":8090", nil)

}
