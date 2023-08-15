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

type ConnInfo struct {
	Conn        net.Conn
	ConnTime    time.Time
	LastOpTime  time.Time
	Bytes       int
	Index       int
	FirstOpTime time.Time
	OpCount     int
}

var startTime int64
var rate int64

func Client(host string, port, count int, verbose bool, bSendData bool, interval int, verboseDetail bool, freq int, bufferSize int) {

	ctx, cancel := context.WithCancel(context.Background())
	log.Println("start tcp client v1.0.8")

	var connList = make([]*ConnInfo, 0, 1000)

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

			//if bufferSize > 0 {
			//	fd, _ := conn.(*net.TCPConn).File()
			//	fHandle := fd.Fd()
			//	err = syscall.SetsockoptInt((syscall.Handle)(fHandle), syscall.SOL_SOCKET, syscall.SO_SNDBUF, bufferSize)
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
			connList = append(connList, &connInfo)
		}

		for _, connInfo := range connList {
			handle(ctx, connInfo, verbose, bSendData, interval, verboseDetail, freq)
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
	log.Println("all connections was closed")
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	sig = <-signalChan
	log.Println("traffic rate: ", rate, "Mbps")
	time.Sleep(5 * time.Second)
	log.Println("Quit done")
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

func handle(ctx context.Context, connInfo *ConnInfo, verbose, bSendData bool, interval int, verboseDetail bool, freq int) {
	// create a local context which is canceled when the function returns
	// close the connection when the context is canceled
	if !bSendData {
		return
	}

	go func() {
		defer connInfo.Conn.Close()

		var strData = "yourkkkkkkkkkkkkkkkkkkkkxoiejrkkkkkkxoxikjekrjioddddddrhhcjkroekckowksiekdirjfkcimine"
		IntervalPerWrite := 0
		if freq > 0 {
			IntervalPerWrite = 1000 / freq
		}

		for {
			if connInfo.Bytes == 0 {
				connInfo.FirstOpTime = time.Now()
			}
			if freq > 0 {
				elapse := time.Now().UnixMilli() - connInfo.FirstOpTime.UnixMilli()
				writePoint := int64((connInfo.OpCount + 1) * IntervalPerWrite)
				if writePoint > elapse {
					waitTime := int64(20)
					if writePoint > elapse+20 {
						waitTime = writePoint - elapse
					}
					time.Sleep(time.Duration(waitTime) * time.Millisecond)
				}
			}
			bytes, err := connInfo.Conn.Write([]byte(strData))
			if bytes == -1 || err != nil {
				log.Println("connection was closed, err: ", err.Error())
				break
			} else if bytes > 0 {
				connInfo.OpCount++
				connInfo.LastOpTime = time.Now()
				connInfo.Bytes = connInfo.Bytes + bytes
				if verboseDetail {
					log.Println("write bytes: ", bytes, "for connection: ", connInfo.Index)
				}
			} else {
				log.Println("write zero bytes for connection: ", connInfo)
			}
			//if verbose {
			//	log.Println("this time send bytes: ", bytes)
			//}

			if interval > 0 {
				time.Sleep(time.Duration(interval) * time.Millisecond)
			}
		}

		if verbose {
			log.Println("send bytes: ", connInfo.Bytes)
			log.Println("close connection: ", connInfo.Conn.RemoteAddr())
			log.Println("conn connect time: ", connInfo.ConnTime)
			log.Println("conn last op time: ", connInfo.LastOpTime)
			log.Println("conn index: ", connInfo.Index)
		}

	}()

	return
}
