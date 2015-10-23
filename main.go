package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"syscall"
	"os/signal"
	"strconv"
	"flag"
	"strings"
)

var(
	LogstashHost  = flag.String("logstashHost", "127.0.01", "Address of the remote Logstash instance")
	LogstashPort  = flag.Int("logstashPort", 10002, "The UDP input port of the remote Logstash instance")
	HaproxyLogSocket = flag.String("haproxyLogSocket","/var/run/haproxy.log.sock","The file location of the socket HAproxy logs to")
	DebugSwitch = flag.Bool("debug",false,"Switches on extra log statements")
)

func Info(msg string){
	fmt.Println("info  ==> ", msg)
}

func Error(err error) {
	if err  != nil {
		fmt.Println("error ==> " , err)
	}
}

func Debug(msg string) {
	fmt.Println("debug ==> ",strings.TrimSpace(msg))
}

func sender(host string, port int, messageChan chan []byte){

	ServerAddr,err := net.ResolveUDPAddr("udp",host + ":" + strconv.Itoa(port))
	Error(err)

	LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	Error(err)

	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	Error(err)

	defer Conn.Close()
	for {
		select {
			case msg := <- messageChan:

				if *DebugSwitch {
					Debug(fmt.Sprintf("Writing to Logstash socket:  %s",msg[:]))
				}

				_,err := Conn.Write(msg[:])
				if err != nil {
					Error(err)
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

		if *DebugSwitch {
			Debug(fmt.Sprintf("Receiving from HAproxy socket:  %s",buf[0:n]))
		}

		messageChan <-buf[0:n]
	}
}

func cleanup(){
	os.Remove(*HaproxyLogSocket)
}

func main() {

	flag.Parse()

	Info("Starting Vamp proxy agent")

	// waiter keeps the programming from exiting instantly
	waiter := make(chan bool)

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
	conn, err := net.ListenUnixgram("unixgram", &net.UnixAddr{*HaproxyLogSocket, "unixgram"})
	if err != nil {
		panic(err)
	}

	Info(fmt.Sprintf("Opened Unix socket at: %s",*HaproxyLogSocket))

	messageChan := make(chan []byte,1000000)
	defer cleanup()
	go reader(conn,messageChan)
	go sender(*LogstashHost,*LogstashPort,messageChan)

	waiter <- true

}
