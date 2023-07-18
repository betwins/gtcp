package client

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Client(host string, port, count int, verbose bool) {

	var connList = make([]net.Conn, 0, 0)

	if host == "" {
		log.Fatalln("host is not allowed blank")
		return
	}
	connAddr := fmt.Sprintf("%s:%d", host, port)

	go func() {
		log.Println("start connect to ", host, port, count, "times")
		for i := 0; i < count; i++ {
			conn, err := net.Dial("tcp", connAddr)
			if err != nil {
				if verbose {
					log.Println("connect fail：" + err.Error())
				}
				continue
			}
			if verbose {
				log.Println("connect success", i)
			}
			connList = append(connList, conn)
		}
		fmt.Println("connection establised: ", len(connList))
	}()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-signalChan
	log.Println("Get Signal:" + sig.String())
	log.Println("Quit ...")

	log.Println("Close connections", len(connList))
	for _, conn := range connList {
		_ = conn.Close()
	}

	time.Sleep(5 * time.Second)

	log.Println("Quit done")
}
