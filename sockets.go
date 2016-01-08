package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"strconv"
)

func Reader(reader io.Reader, messageChannel chan string) {
	buf := make([]byte, 4096)
	for {
		n, err := reader.Read(buf[:])
		if err != nil {
			logger.Error("Can't read: %s", err.Error())
			return
		}

		message := string(buf[0:n])

		if *debug {
			logger.Debug(fmt.Sprintf("Received from socket:  %s", message))
		}

		select {
		case messageChannel <- message:
		default:
		}
	}
}

func Sender(host string, port int, messageChannel chan string) {
	for {
		ServerAddress, err := net.ResolveUDPAddr("udp", host + ":" + strconv.Itoa(port))
		if err != nil {
			logger.Error("Can't resolve remote UDP address: %s", err.Error())
		} else {
			Conn, err := net.DialUDP("udp", nil, ServerAddress)
			if err != nil {
				logger.Error("Can't dial up: %s", err.Error())
			} else {
				defer Conn.Close()

				for {
					select {
					case message := <-messageChannel:

						if *debug {
							logger.Debug(fmt.Sprintf("Writing to UDP socket: %s", message))
						}

						_, err := Conn.Write([]byte(message))
						if err != nil {
							logger.Error("Can't write: %s", err.Error())
						}
					}
				}
			}
		}

		time.Sleep(5 * time.Second)
	}
}


