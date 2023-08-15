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

	var interval int

	var verboseDetail bool

	var freq int

	var bufferSize int

	flag.BoolVar(&bServer, "s", false, "默认false, tcp客户端")
	flag.StringVar(&host, "h", "", "主机名")
	flag.IntVar(&port, "p", 80, "目标端口号，默认80")
	flag.IntVar(&count, "c", 1, "建立连接数量，默认为1")
	flag.BoolVar(&verbose, "v", false, "默认false")
	flag.BoolVar(&bSendData, "d", false, "默认发送数据")
	flag.IntVar(&interval, "i", 0, "每个连接发送数据间隔")
	flag.BoolVar(&verboseDetail, "vd", false, "输出连接过程每次发送数据详情")
	flag.IntVar(&freq, "f", 0, "每个连接每秒数据发送次数")
	flag.IntVar(&bufferSize, "b", 0, "每个连接的socket buffer size")
	flag.Parse()
	if verboseDetail {
		verbose = true
	}
	if freq > 0 {
		interval = 0
	}

	if bServer {
		server.Server(host, port, verbose, verboseDetail, bufferSize)
	} else {
		client.Client(host, port, count, verbose, bSendData, interval, verboseDetail, freq, bufferSize)
	}

}
