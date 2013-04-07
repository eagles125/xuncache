package main

import (
	json_sim "./simlejson"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

type Promise interface{}

var accept Promise
var storemap = make(map[string]map[string]interface{})
var new_config = make(map[string]string)

const (
	IP   string = "127.0.0.1"
	PORT string = "1200"
)

func main() {
	//初始内容
	fmt.Print("Server started, xuncache version 0.2\n")
	//读取配置文件
	var config = make(map[string]string)
	config_file, err := os.Open("config.conf") //打开文件
	defer config_file.Close()
	if err != nil {
		fmt.Print("Can not read configuration file. now exit\n")
		os.Exit(0)
	}
	buff := bufio.NewReader(config_file) //读入缓存
	//读取配置文件
	for {
		line, err := buff.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil {
			break
		}
		rs := []rune(line)
		if string(rs[0:1]) == "#" || len(line) < 3 {
			continue
		}
		str_type := string(rs[0:strings.Index(line, " ")])
		detail := string(rs[strings.Index(line, " ")+1 : len(rs)-1])
		config[str_type] = detail
	}
	//再次过滤 (防止没有配置文件)
	new_config := verify(config)
	//创建服务端
	tcpAddr, err := net.ResolveTCPAddr("tcp4", new_config["bind"]+":"+new_config["port"])
	fmt.Printf("The server is now ready to accept connections on %s:%s\n", new_config["bind"], new_config["port"])
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

//处理数据
func handleClient(conn net.Conn) {

	//标记结束连接
	defer conn.Close()
	defer fmt.Print("Client closed connection\n")
	ipAddr := conn.RemoteAddr()
	fmt.Printf("Accepted %s\n", ipAddr)
	for {
		var back = make(map[string]interface{})
		//获取数据
		var buf [1024]byte
		n, _ := conn.Read(buf[0:])
		b := []byte(buf[0:n])
		if len(b) < 1 {
			return
		}
		js, _ := json_sim.NewJson(b)
		pass, _ := js.Get("Pass").String()
		if pass != new_config["password"] && len(new_config["password"]) > 1 {
			fmt.Printf("Encountered a connection password is incorrect Accepted %s\n", ipAddr)
			back["error"] = true
			back["point"] = "password error!"
			rewrite(back, conn)
			return
		}
		//获取key
		key, _ := js.Get("Key").String()
		if len(key) < 1 {
			fmt.Printf("Error agreement is key %s\n", key)
			back["error"] = true
			back["point"] = "Please input Key!"
			rewrite(back, conn)
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

		switch protocol {
		case "delete":
			delete(storemap, key)
			back["status"] = true
			break
		case "set":
			storemap[key] = data
			back["status"] = true
			break
		case "find":
			back["data"] = storemap[key]
			back["status"] = true
			break
		default:
			back["status"] = false
			fmt.Print("error protocol \n")
			break
		}
		//返回内容
		rewrite(back, conn)
	}
}

//写入数据
func rewrite(back map[string]interface{}, conn net.Conn) {
	jsback, _ := json.Marshal(back)
	//返回内容
	conn.Write(jsback)
}

//验证配置文件
func verify(config map[string]string) (config_bak map[string]string) {
	if len(config["bind"]) < 3 {
		config["bind"] = IP
	}
	if len(config["port"]) < 1 {
		config["port"] = PORT
	}
	return config
}

//输出错误信息
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
