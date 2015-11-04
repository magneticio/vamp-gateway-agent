package main

import (
	"fmt"
	"io"
	"net"
	"strconv"
)

func Reader(reader io.Reader, messageChannel chan []byte) {
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf[:])
		if err != nil {
			Log.Error("Can't read: %s", err.Error())
			return
		}

		if *DebugSwitch {
			Log.Debug(fmt.Sprintf("Received from socket:  %s", buf[0:n]))
		}

		select {
		case messageChannel <- buf[0:n]:
		default:
		}
	}
}

func Sender(host string, port int, messageChannel chan []byte) {

	ServerAddress, err := net.ResolveUDPAddr("udp", host + ":" + strconv.Itoa(port))
	if err != nil {
		Log.Error("Can't resolve remote UDP address: %s", err.Error())
		return
	}

	LocalAddress, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		Log.Error("Can't resolve local UDP address: %s", err.Error())
		return
	}

	Conn, err := net.DialUDP("udp", LocalAddress, ServerAddress)
	if err != nil {
		Log.Error("Can't dial up: %s", err.Error())
		return
	}

	defer Conn.Close()

	for {
		select {
		case msg := <-messageChannel:

			if *DebugSwitch {
				Log.Debug(fmt.Sprintf("Writing to UDP socket:  %s", msg[:]))
			}

			_, err := Conn.Write(msg[:])
			if err != nil {
				Log.Error("Can't write: %s", err.Error())
			}
		}
	}
}


