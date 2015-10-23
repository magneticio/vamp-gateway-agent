package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"syscall"
	"os/signal"
	"sync"
	"strconv"
)

func HandleError(err error) {
	if err  != nil {
		fmt.Println("Error: " , err)
	}
}

func sender(host string, port int, messageChan chan []byte){

	ServerAddr,err := net.ResolveUDPAddr("udp",host + ":" + strconv.Itoa(port))
	HandleError(err)

	LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	HandleError(err)

	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	HandleError(err)

	defer Conn.Close()
	for {
		select {
			case msg := <- messageChan:
				_,err := Conn.Write(msg[:])
				if err != nil {
					HandleError(err)
				}
		}
	}
}

func reader(r io.Reader, messageChan chan []byte) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return
		}
		messageChan <-buf[0:n]
	}
}

func cleanup(){
	os.Remove("/var/run/vamp.log.sock")
}

func main() {

	// waiter keeps the programming from exiting instantly
	waiter := &sync.WaitGroup{}
	waiter.Add(1)

	// catch an CTR+C exits so the cleanup routine is called
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	// set the socket Haproxy can write logs to
	conn, err := net.ListenUnixgram("unixgram", &net.UnixAddr{"/var/run/vamp.log.sock", "unixgram"})
	if err != nil {
		panic(err)
	}

	messageChan := make(chan []byte,1000000)
	defer cleanup()
	go reader(conn,messageChan)
	go sender("127.0.0.1",10002,messageChan)

	waiter.Wait()

}
