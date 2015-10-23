package main

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"os/signal"
	"flag"
)

var(
	logHost = flag.String("logHost", "127.0.01", "Address of the remote Logstash instance")
	statsHost = flag.String("statsHost", "127.0.01", "Address of the remote Logstash instance")
	logPort = flag.Int("logPort", 10002, "The UDP input port of the remote Logstash instance")
	statsPort = flag.Int("statsPort", 10003, "The UDP input port of the remote Logstash instance")
	HaproxyLogSocket = flag.String("haproxyLogSocket","/var/run/haproxy.log.sock","The location of the socket HAproxy logs to")
	HaproxyStatsSocket = flag.String("haproxyStatsSocket","/tmp/haproxy.stats.sock","The location of the HAproxy stats socket")
	HaproxyStatsType = flag.String("haproxyStatsType","all","Which stats to read from haproxy: all, frontend, backend or server.")
	DebugSwitch = flag.Bool("debug",false,"Switches on extra log statements")
)

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

	logChannel := make(chan []byte,1000000)
	statsChannel := make(chan []byte,1000000)
	defer cleanup()

	// start the logging stream
	go simpleReader(conn,logChannel)
	go sender(*logHost,*logPort,logChannel)

	// start the stats stream
	go statsReader(*HaproxyStatsSocket,*HaproxyStatsType,statsChannel)
	go sender(*statsHost,*statsPort,statsChannel)

	waiter <- true

}
