package main

import (
	"fmt"
	"io"
	"net"
	"strconv"
)

func sender(host string, port int, messageChan chan []byte) {

	ServerAddress, err := net.ResolveUDPAddr("udp", host + ":" + strconv.Itoa(port))
	if err != nil {
		Logger.Error(err.Error())
		return
	}

	LocalAddress, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		Logger.Error(err.Error())
		return
	}

	Conn, err := net.DialUDP("udp", LocalAddress, ServerAddress)
	if err != nil {
		Logger.Error(err.Error())
		return
	}

	defer Conn.Close()
	for {
		select {
		case msg := <-messageChan:

			if *DebugSwitch {
				Logger.Debug(fmt.Sprintf("Writing to Logstash socket:  %s", msg[:]))
			}

			_, err := Conn.Write(msg[:])
			if err != nil {
				Logger.Error(err.Error())
			}
		}
	}
}

func simpleReader(r io.Reader, messageChan chan []byte) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			Logger.Error(err.Error())
			return
		}

		if *DebugSwitch {
			Logger.Debug(fmt.Sprintf("Receiving from HAProxy socket:  %s", buf[0:n]))
		}

		select {
		case messageChan <- buf[0:n]:
		default:
		}
	}
}

