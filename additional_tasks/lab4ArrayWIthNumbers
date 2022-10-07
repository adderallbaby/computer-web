package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sort"
	"strings"
)

type Package struct {
	To   string
	From string
	Data string
}

type Node struct {
	Connections map[string]bool
	Nums        []string
	Address     Address
	Num         string
}

type Address struct {
	IPv4 string
	Port string
}

// ./ main :8080
func init() {
	if len(os.Args) != 3 {
		panic("len args !=3")
	}
}
func main() {
	NewNode(os.Args[1], os.Args[2]).Run(handleServer, handleClient)
}

//ipv4:port

func NewNode(address string, num string) *Node {
	splitted := strings.Split(address, ":")
	if len(splitted) != 2 {
		return nil
	}
	nums := make([]string, 0)
	nums = append(nums, num)
	return &Node{
		Connections: make(map[string]bool),
		Nums:        nums,
		Num:         num,
		Address: Address{
			IPv4: splitted[0],
			Port: ":" + splitted[1],
		},
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
	if strings.Index(pack.Data, "append") != -1 {
		num := pack.Data[7:]
		node.Nums = append(node.Nums, num)
	}
	splitted := strings.Split(pack.Data, "|")
	switch splitted[0] {
	case "/connect":
		node.ConnectTo(splitted[0:])
	case "/network":
		node.PrintNetwork()
	case "/change":
		if len(splitted) == 3 {
			for k := 0; k < len(node.Nums); k++ {
				if node.Nums[k] == splitted[1] {
					node.Nums[k] = splitted[2]
				}
			}
			sort.Strings(node.Nums)
		}
	case "/anotherchange":
		if len(splitted) == 3 {
			for k := 0; k < len(node.Nums); k++ {
				if node.Nums[k] == splitted[1] {
					node.Nums[k] = splitted[2]
				}
			}
			sort.Strings(node.Nums)
		}
	}
}
func handleClient(node *Node) {
	for {
		message := InputString()
		splitted := strings.Split(message, "|")
		switch splitted[0] {
		case "/exit":
			os.Exit(0)
		case "/connect":
			node.ConnectTo(splitted[0:])
		case "/network":
			sort.Strings(node.Nums)
			node.PrintNetwork()
		case "/change":
			if len(splitted) == 3 {
				for k := 0; k < len(node.Nums); k++ {
					if node.Nums[k] == splitted[1] {
						node.Nums[k] = splitted[2]
					}
				}
				sort.Strings(node.Nums)
				node.PrintNetwork()
				node.SendMessageToAll("/network")
				node.SendMessageToAll("/anotherchange" + "|" + splitted[1] + "|" + splitted[2])
			}
		case "/anotherchange":
			if len(splitted) == 3 {
				for k := 0; k < len(node.Nums); k++ {
					if node.Nums[k] == splitted[1] {
						node.Nums[k] = splitted[2]
					}
				}
				sort.Strings(node.Nums)
				node.PrintNetwork()
			}
		}
	}
}
func (node *Node) PrintNetwork() {
	fmt.Println(node.Nums)
}

func (node *Node) ConnectTo(addresses []string) {
	for _, addr := range addresses {
		if node.Connections[addr] == false {
			var new_pack = Package{
				From: node.Address.IPv4 + node.Address.Port,
				Data: "append " + node.Num,
			}
			new_pack.To = addr
			node.Send(&new_pack)
		}
		node.Connections[addr] = true
	}
}

func (node *Node) SendMessageToAll(message string) {
	var new_pack = Package{
		From: node.Address.IPv4 + node.Address.Port,
		Data: message,
	}
	for addr := range node.Connections {
		new_pack.To = addr
		node.Send(&new_pack)
	}
}

func (node *Node) Send(pack *Package) {
	conn, err := net.Dial("tcp", pack.To)
	if err != nil {
		delete(node.Connections, pack.To)
		return
	}
	defer conn.Close()
	json_pack, _ := json.Marshal(*pack)
	conn.Write(json_pack)
}
func InputString() string {
	var msg string
	fmt.Scanf("%s/n", &msg)
	return strings.Replace(msg, "\n", "", -1)
}
