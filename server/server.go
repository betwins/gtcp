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

type ConnInfo struct {
	Conn       net.Conn
	ConnTime   time.Time
	LastOpTime time.Time
	Bytes      int
	Index      int
}

var startTime int64
var rate int64

func Server(host string, port int, verbose bool, verboseDetail bool, bufferSize int) {
	ctx, cancel := context.WithCancel(context.Background())
	log.Println("start tcp server v1.0.8")

	var connList = make([]*ConnInfo, 0, 0)
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
		i := 0
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

				//if bufferSize > 0 {
				//	fd, _ := conn.(*net.TCPConn).File()
				//	err = syscall.SetsockoptInt(syscall.Handle(fd.Fd()), syscall.SOL_SOCKET, syscall.SO_RCVBUF, bufferSize)
				//	if err != nil {
				//		log.Fatalln("set buffer error", err.Error())
				//	}
				//}

				if startTime == 0 {
					startTime = time.Now().UnixMilli()
				}

				connInfo := ConnInfo{
					Conn:       conn,
					ConnTime:   time.Now(),
					LastOpTime: time.Now(),
					Bytes:      0,
					Index:      i,
				}
				i++

				connList = append(connList, &connInfo)
				handle(ctx, &connInfo, verbose, verboseDetail)
			}
		}
	}(ctx)

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-signalChan
	log.Println("Get Signal:" + sig.String())
	cancel()
	closeAllConn(connList)
	log.Println("all connections was closed")
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	sig = <-signalChan
	log.Println("traffic rate: ", rate, "Mbps")
	time.Sleep(10 * time.Second)
	log.Println("Quit Done")
}

func closeAllConn(connList []*ConnInfo) {
	var totalBytes int64
	timeSpend := time.Now().UnixMilli() - startTime
	for _, connInfo := range connList {
		totalBytes += int64(connInfo.Bytes)
		_ = connInfo.Conn.Close()
	}
	rate = totalBytes / (timeSpend / 1000) * 8
	rate = rate / (1024 * 1024)
	log.Println("close all connection")
}

func handle(ctx context.Context, connInfo *ConnInfo, verbose bool, verboseDetail bool) {
	// create a local context which is canceled when the function returns
	// close the connection when the context is canceled
	go func() {
		defer connInfo.Conn.Close()
		var buf = make([]byte, 100000)
		for {
			bytes, err := connInfo.Conn.Read(buf)
			if bytes == -1 || err != nil {
				log.Println("connection was closed, err is: ", err.Error())
				break
			} else if bytes > 0 {
				connInfo.LastOpTime = time.Now()
				connInfo.Bytes = connInfo.Bytes + bytes
				if verboseDetail {
					log.Println("read bytes: ", bytes, "for connection: ", connInfo.Index)
					log.Println("read content: ", string(buf))
				}
			} else {
				log.Println("read zero bytes for connection: ", connInfo)
			}

			//if verbose {
			//	log.Println("this time recv bytes: ", bytes)
			//}

		}

		if verbose {
			log.Println("recv bytes: ", connInfo.Bytes)
			log.Println("close connection: ", connInfo.Conn.RemoteAddr())
			log.Println("conn connect time: ", connInfo.ConnTime)
			log.Println("conn last op time: ", connInfo.LastOpTime)
			log.Println("conn index: ", connInfo.Index)
		}

	}()

	return
}
