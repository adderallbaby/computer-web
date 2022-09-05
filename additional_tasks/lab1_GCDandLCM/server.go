package main

import (
	"../github.com/mgutz/logxi/v1"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
)

import "../proto"

// Client - состояние клиента.
type Client struct {
	logger    log.Logger    // Объект для печати логов
	conn      *net.TCPConn  // Объект TCP-соединения
	enc       *json.Encoder // Объект для кодирования и отправки сообщений
	solutions string        //
}

// NewClient - конструктор клиента, принимает в качестве параметра
// объект TCP-соединения.
func NewClient(conn *net.TCPConn) *Client {
	return &Client{
		logger:    log.New(fmt.Sprintf("client %s", conn.RemoteAddr().String())),
		conn:      conn,
		enc:       json.NewEncoder(conn),
		solutions: "",
	}
}

// serve - метод, в котором реализован цикл взаимодействия с клиентом.
// Подразумевается, что метод serve будет вызаваться в отдельной go-программе.
func (client *Client) serve() {
	defer client.conn.Close()
	decoder := json.NewDecoder(client.conn)
	for {
		var req proto.Request
		if err := decoder.Decode(&req); err != nil {
			client.logger.Error("cannot decode message", "reason", err)
			break
		} else {
			client.logger.Info("received command", "command", req.Command)
			fmt.Printf("GOT IT\n")
			if client.handleRequest(&req) {
				client.logger.Info("shutting down connection")
				break
			}
		}
	}
}

// handleRequest - метод обработки запроса от клиента. Он возвращает true,
// если клиент передал команду "quit" и хочет завершить общение.
func (client *Client) handleRequest(req *proto.Request) bool {
	switch req.Command {
	case "quit":
		client.respond("ok", nil)
		return true
	case "solve":
		errorMsg := ""
		if req.Data == nil {
			errorMsg = "data field is absent"
		} else {
			var frac proto.Fraction
			if err := json.Unmarshal(*req.Data, &frac); err != nil {
				errorMsg = "malformed data field"
			} else {

				square := client.getSquare(frac)
				client.solutions = square
				client.logger.Info("performing addition", "value", square)
			}
		}
		if errorMsg == "" {
			client.respond("ok", nil)
		} else {
			client.logger.Error("addition failed", "reason", errorMsg)
			client.respond("failed", errorMsg)
		}
	case "res":
		ansRadius := client.solutions
		client.respond("result", &proto.Fraction{
			Ans: ansRadius,
		})
	default:
		client.logger.Error("unknown command")
		client.respond("failed", "unknown command")
	}
	return false
}

func gcd(temp1 int, temp2 int) int {
	var gcdnum int
	/* Use of And operator in For Loop */
	for i := 1; i <= temp1 && i <= temp2; i++ {
		if temp1%i == 0 && temp2%i == 0 {
			gcdnum = i
		}
	}
	return gcdnum
}

func (client *Client) getSquare(circle proto.Fraction) string {

	a, _ := strconv.Atoi(circle.A)
	b, _ := strconv.Atoi(circle.B)
	result := a * b / gcd(a, b)
	solutions := make([]int, 0)
	solutions = append(solutions, gcd(a, b))
	solutions = append(solutions, result)

	valuesText := []string{}

	// Create a string slice using strconv.Itoa.
	// ... Append strings to it.
	valuesText = append(valuesText, "GCD:")
	fl := false
	for i := range solutions {
		number := solutions[i]
		text := strconv.Itoa(number)
		valuesText = append(valuesText, text)
		if fl == false {
			valuesText = append(valuesText, "LCM:")
			fl = true
		}
	}
	resultt := strings.Join(valuesText, " ")

	return resultt
}

// respond - вспомогательный метод для передачи ответа с указанным статусом
// и данными. Данные могут быть пустыми (data == nil).
func (client *Client) respond(status string, data interface{}) {
	var raw json.RawMessage
	raw, _ = json.Marshal(data)
	client.enc.Encode(&proto.Response{status, &raw})
}

func main() {
	// Работа с командной строкой, в которой может указываться необязательный ключ -addr.
	var addrStr string
	flag.StringVar(&addrStr, "addr", "127.0.0.1:6000", "specify ip address and port")
	flag.Parse()

	// Разбор адреса, строковое представление которого находится в переменной addrStr.
	if addr, err := net.ResolveTCPAddr("tcp", addrStr); err != nil {
		log.Error("address resolution failed", "address", addrStr)
	} else {
		log.Info("resolved TCP address", "address", addr.String())

		// Инициация слушания сети на заданном адресе.
		if listener, err := net.ListenTCP("tcp", addr); err != nil {
			log.Error("listening failed", "reason", err)
		} else {
			// Цикл приёма входящих соединений.
			for {
				if conn, err := listener.AcceptTCP(); err != nil {
					log.Error("cannot accept connection", "reason", err)
				} else {
					log.Info("accepted connection", "address", conn.RemoteAddr().String())

					// Запуск go-программы для обслуживания клиентов.
					go NewClient(conn).serve()
				}
			}
		}
	}
}
