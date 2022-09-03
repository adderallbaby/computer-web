package main

import (
	"../github.com/mgutz/logxi/v1"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
)

import "../proto"

// Client - состояние клиента.
type Client struct {
	logger   log.Logger    // Объект для печати логов
	conn     *net.TCPConn  // Объект TCP-соединения
	enc      *json.Encoder // Объект для кодирования и отправки сообщений
	dividers string        //
	count    int64         // Количество полученных от клиента дробей

}

// NewClient - конструктор клиента, принимает в качестве параметра
// объект TCP-соединения.
func NewClient(conn *net.TCPConn) *Client {
	return &Client{
		logger:   log.New(fmt.Sprintf("client %s", conn.RemoteAddr().String())),
		conn:     conn,
		enc:      json.NewEncoder(conn),
		dividers: "",
		count:    0,
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
	case "getDivs":
		errorMsg := ""
		if req.Data == nil {
			errorMsg = "data field is absent"
		} else {
			var frac proto.Fraction
			if err := json.Unmarshal(*req.Data, &frac); err != nil {
				errorMsg = "malformed data field"
			} else {
				square := client.getSquare(frac)
				client.dividers = square
				client.logger.Info("performing addition", "value", square)
				client.count++
			}
		}
		if errorMsg == "" {
			client.respond("ok", nil)
		} else {
			client.logger.Error("addition failed", "reason", errorMsg)
			client.respond("failed", errorMsg)
		}
	case "res":
		if client.count == 0 {
			client.logger.Error("calculation failed", "reason", "division by zero")
			client.respond("failed", "division by zero")
		} else {
			ansRadius := client.dividers

			client.respond("result", &proto.Fraction{
				Ans: ansRadius,
			})
		}
	default:
		client.logger.Error("unknown command")
		client.respond("failed", "unknown command")
	}
	return false
}
func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (client *Client) getSquare(circle proto.Fraction) string {

	x, _ := strconv.Atoi(circle.X)
	dividers := make([]int, 0)
	for i := 1; i < int(math.Sqrt(float64(x))); i++ {
		if x%i == 0 {
			dividers = append(dividers, i)
			if !contains(dividers, x/i) {
				dividers = append(dividers, x/i)
			}
		}
	}

	valuesText := []string{}

	// Create a string slice using strconv.Itoa.
	// ... Append strings to it.
	for i := range dividers {
		number := dividers[i]
		text := strconv.Itoa(number)
		valuesText = append(valuesText, text)
	}

	// Join our string slice.
	result := strings.Join(valuesText, " ")
	return result
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
