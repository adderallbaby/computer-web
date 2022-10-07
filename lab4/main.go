package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

type Package struct {
	To   string
	From string
	Data string
}

type Node struct {
	Connections  map[string]bool
	Address      Address
	StartingNode bool
}

type Address struct {
	IPv4 string
	Port string
}

// ./ main :8080 or any other port as the second argument
func init() {
	if len(os.Args) != 2 {
		panic("len args !=2")
	}
}
func main() {
	NewNode(os.Args[1]).Run(handleServer, handleClient)
}

//ipv4:port

func NewNode(address string) *Node {
	splitted := strings.Split(address, ":")
	if len(splitted) != 2 {
		return nil
	}
	return &Node{
		Connections: make(map[string]bool),
		Address: Address{
			IPv4: splitted[0],
			Port: ":" + splitted[1],
		},
		StartingNode: false,
	}
}

func (node *Node) Run(handleServer func(*Node), handleClient func(*Node)) {
	go handleServer(node)
	handleClient(node)
}

func handleServer(node *Node) {
	listen, err := net.Listen("tcp", "0.0.0.0"+node.Address.Port)
	if err != nil {
		panic("listen error")
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			break
		}
		go handleConnection(node, conn)
	}
}

func handleConnection(node *Node, conn net.Conn) {
	var (
		buffer  = make([]byte, 512)
		message string
		pack    *Package
	)
	for {
		length, err := conn.Read(buffer)
		if err != nil {
			break
		}
		message += string(buffer[:length])
	}
	err := json.Unmarshal([]byte(message), &pack)
	if err != nil {
		return
	}
	node.ConnectTo([]string{pack.From})
	if pack.Data[0:13] == "/__NETWORK__/" {
		node.PrintNetwork()
		fmt.Println("|:" + pack.Data[14:])

	} else {
		fmt.Println(pack.Data)

	}

}
func handleClient(node *Node) {
	for {
		message := InputString()
		splitted := strings.Split(message, " ")
		switch splitted[0] {
		case "/exit":
			os.Exit(0)
		case "/connect":
			node.ConnectTo(splitted[1:])
		case "/network":
			node.PrintNetwork()

		default:
			node.SendMessageToAll(message)

		}
	}
}
func (node *Node) PrintNetwork() {
	if !node.StartingNode {
		node.StartingNode = true
		fmt.Println("|" + node.Address.Port)
		node.SendMessageToAll("/__NETWORK__/")
	}
}

func (node *Node) ConnectTo(addresses []string) {
	for _, addr := range addresses {
		node.Connections[addr] = true
	}

}

func (node *Node) SendMessageToAll(message string) {
	var new_pack = Package{
		From: node.Address.IPv4 + node.Address.Port,
		Data: message,
	}
	// var pack_ring = Package{
	// 	From: node.Address.IPv4 + node.Address.Port,
	// 	Data: "/network",
	// }
	for addr := range node.Connections {
		new_pack.To = addr
		// pack_ring.To = addr
		node.Send(&new_pack)
		// node.Send(&pack_ring)

	}
}

func (node *Node) Send(pack *Package) {
	conn, err := net.Dial("tcp", pack.To)
	if err != nil {
		delete(node.Connections, pack.To)
		return
	}
	// var pack_ring = Package{
	// 	From: node.Address.IPv4 + node.Address.Port,
	// 	Data: "/network",
	//}
	defer conn.Close()
	if pack.Data == "/__NETWORK__/" {
		pack.Data = pack.Data + pack.From
	}
	json_pack, _ := json.Marshal(*pack)
	// json_ring, _ := json.Marshal(*&pack_ring)

	conn.Write(json_pack)
	node.PrintNetwork()
	// conn.Write(json_ring)

}

func InputString() string {
	msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.Replace(msg, "\n", "", -1)
}
