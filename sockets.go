package main

import (
	"fmt"
	"io"
	"net"
	"strconv"
)

func Reader(reader io.Reader, messageChannel chan []byte) {
	buf := make([]byte, 4096)
	for {
		n, err := reader.Read(buf[:])
		if err != nil {
			logger.Error("Can't read: %s", err.Error())
			return
		}

		if *debug {
			logger.Debug(fmt.Sprintf("Received from socket:  %s", buf[0:n]))
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
		logger.Error("Can't resolve remote UDP address: %s", err.Error())
		return
	}

	Conn, err := net.DialUDP("udp", nil, ServerAddress)
	if err != nil {
		logger.Error("Can't dial up: %s", err.Error())
		return
	}

	defer Conn.Close()

	for {
		select {
		case msg := <-messageChannel:

			if *debug {
				logger.Debug(fmt.Sprintf("Writing to UDP socket: %s", msg[:]))
			}

			_, err := Conn.Write(msg[:])
			if err != nil {
				logger.Error("Can't write: %s", err.Error())
			}
		}
	}
}


