package main

import (
	"flag"
	"gtcp/client"
	"gtcp/server"
)

func main() {

	// 主机名
	var host string
	// 端口号
	var port int
	// 连接数
	var count int

	var verbose bool

	var bServer bool

	var bSendData bool

	flag.BoolVar(&bServer, "s", false, "默认false, tcp客户端")
	flag.StringVar(&host, "h", "", "主机名")
	flag.IntVar(&port, "p", 80, "目标端口号，默认80")
	flag.IntVar(&count, "c", 1, "建立连接数量，默认为1")
	flag.BoolVar(&verbose, "v", false, "默认false")
	flag.BoolVar(&bSendData, "d", false, "默认发送数据")
	flag.Parse()

	if bServer {
		server.Server(host, port, verbose)
	} else {
		client.Client(host, port, count, verbose, bSendData)
	}

}
