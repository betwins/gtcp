package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var connList = make([]net.Conn, 0, 0)

	// 主机名
	var host string
	// 端口号
	var port int
	// 连接数
	var count int

	var verbose bool

	flag.StringVar(&host, "h", "127.0.0.1", "主机名，默认127.0.0.1")
	flag.IntVar(&port, "p", 80, "目标端口号，默认80")
	flag.IntVar(&count, "c", 1, "建立连接数量，默认为1")
	flag.BoolVar(&verbose, "v", false, "默认false")
	flag.Parse()

	connAddr := fmt.Sprintf("%s:%d", host, port)

	go func() {
		fmt.Println("start connect to ", host, port, count, "times")
		for i := 0; i < count; i++ {
			conn, err := net.Dial("tcp", connAddr)
			if err != nil {
				if verbose {
					fmt.Println("connect fail：" + err.Error())
				}
				continue
			}
			if verbose {
				fmt.Println("connect success", i)
			}
			connList = append(connList, conn)
		}
		fmt.Println("connection establised: ", len(connList))
	}()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-signalChan
	fmt.Println("Get Signal:" + sig.String())
	fmt.Println("Quit ...")

	fmt.Println("Close connections", len(connList))
	for _, conn := range connList {
		_ = conn.Close()
	}

	time.Sleep(5 * time.Second)

	fmt.Println("Quit done")

}
