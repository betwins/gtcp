package client

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Client(host string, port, count int, verbose bool, bSendData bool, interval int) {

	ctx, cancel := context.WithCancel(context.Background())
	log.Println("start tcp client")

	var connList = make([]net.Conn, 0, 1000)

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
				log.Println("new connect success", i)
			}
			connList = append(connList, conn)

		}

		for _, conn := range connList {
			go handle(ctx, conn, verbose, bSendData, interval)
		}

		fmt.Println("connection establised: ", len(connList))
	}()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-signalChan
	log.Println("Get Signal:" + sig.String())
	log.Println("Quit ...")

	log.Println("Close connections count: ", len(connList))
	cancel()
	closeAllConn(connList)

	time.Sleep(5 * time.Second)

	log.Println("Quit done")
}

func closeAllConn(connList []net.Conn) {
	for _, conn := range connList {
		_ = conn.Close()
	}
	log.Println("close all connection")
}

func handle(ctx context.Context, conn net.Conn, verbose, bSendData bool, interval int) {
	// create a local context which is canceled when the function returns
	// close the connection when the context is canceled
	if !bSendData {
		return
	}

	var strData = "kdorkkkkkkkkkkkkkkkkkkkkxoiejrkkkkkkxoxikjekrjioddddddrhhcjkroekckowksiekdirjfkcidkeh"
	go func() {
		defer conn.Close()

		totalBytes := 0

		for {

			bytes, err := conn.Write([]byte(strData))
			if bytes == -1 || err != nil {
				log.Println("connection was closed, err: ", err.Error())
				break
			}
			//if verbose {
			//	log.Println("this time send bytes: ", bytes)
			//}
			totalBytes = totalBytes + bytes

			if interval > 0 {
				time.Sleep(time.Duration(interval) * time.Millisecond)
			}
		}

		if verbose {
			log.Println("send bytes: ", totalBytes)
			log.Println("close connection: ", conn.RemoteAddr())
		}

	}()

	return
}
