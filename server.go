package main

import (
	json_sim "./simlejson"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Promise interface{}

var accept Promise
var storemap = make(map[string]map[string]interface{})

const (
	PORT = "1200"
	PASS = "pass"
)

func main() {
	//初始内容
	fmt.Print("Server started, xuncache version 0.1\n")
	//创建服务端
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":"+PORT)
	fmt.Printf("The server is now ready to accept connections on port %s\n", PORT)

	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {

	//标记结束连接
	defer conn.Close()
	defer fmt.Print("Client closed connection\n")
	ipAddr := conn.RemoteAddr()
	//获取数据
	var buf [1024]byte
	n, _ := conn.Read(buf[0:])
	b := []byte(buf[0:n])
	if len(b) < 1 {
		fmt.Printf("No protocol connection Accepted %s\n", ipAddr)
		return
	}
	js, _ := json_sim.NewJson(b)
	pass, _ := js.Get("Pass").String()
	if pass != PASS {
		fmt.Printf("Encountered a connection password is incorrect Accepted %s\n", ipAddr)
		return
	}
	//获取key
	key, _ := js.Get("Key").String()
	if len(key) < 1 {
		fmt.Printf("Error agreement is key %s\n", key)
		return
	}
	//获取协议
	protocol, _ := js.Get("Protocol").String()
	//数据处理
	data, _ := js.Get("Data").Map()
	if data == nil && protocol == "set" {
		fmt.Print("There is no data \n")
		return
	}
	var back = make(map[string]interface{})

	switch protocol {
	case "delete":
		delete(storemap, key)
		back["status"] = "true"
		break
	case "set":
		storemap[key] = data
		back["status"] = "true"
		break
	case "find":
		back["data"] = storemap[key]
		back["status"] = "true"
		break
	default:
		back["status"] = "false"
		fmt.Print("error protocol \n")
		break
	}
	jsback, _ := json.Marshal(back)
	//返回内容
	conn.Write(jsback)

	//输出内容
	fmt.Printf("Accepted %s\n", ipAddr)
}


func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
