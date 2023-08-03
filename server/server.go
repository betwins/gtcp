package server

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

func Server(host string, port int, verbose bool) {
	ctx, cancel := context.WithCancel(context.Background())
	log.Println("start tcp server...")

	var connList = make([]net.Conn, 0, 0)
	go func(ctx context.Context) {
		var listenAddr string
		if host != "" {
			listenAddr = fmt.Sprintf("%s:%d", host, port)
		} else {
			listenAddr = fmt.Sprintf("0.0.0.0:%d", port)
		}
		log.Println("start listen at ", listenAddr)
		listener, err := net.Listen("tcp", listenAddr)
		if err != nil {
			log.Fatalln("can't listen: ", err.Error())
		}
		for {
			select {
			case <-ctx.Done():
				_ = listener.Close()
				log.Println("listener closed")
				return
			default:
				conn, err := listener.Accept()
				if err != nil {
					log.Println("accept error: ", err.Error())
				}
				if verbose {
					log.Println("accept new connection:", conn.RemoteAddr())
				}
				connList = append(connList, conn)
				go handle(ctx, conn, verbose)
			}
		}
	}(ctx)

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-signalChan
	log.Println("Get Signal:" + sig.String())
	closeAllConn(connList)
	cancel()
	time.Sleep(10 * time.Second)
	log.Println("Quit Done")
}

func closeAllConn(connList []net.Conn) {
	for _, conn := range connList {
		_ = conn.Close()
	}
}

func handle(ctx context.Context, conn net.Conn, verbose bool) {
	// create a local context which is canceled when the function returns
	// close the connection when the context is canceled
	go func() {
		defer conn.Close()
		totalBytes := 0
		var buf = make([]byte, 1000)
		for {
			bytes, err := conn.Read(buf)
			if bytes == -1 || err != nil {
				log.Println("connection was closed, err:", err.Error())
				break
			}
			//if verbose {
			//	log.Println("this time recv bytes: ", bytes)
			//}
			totalBytes = totalBytes + bytes
		}

		if verbose {
			log.Println("recv bytes: ", totalBytes)
			log.Println("close connection: ", conn.RemoteAddr())
		}

	}()

	return
}
