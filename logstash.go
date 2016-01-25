package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

type Logstash struct {
	Address string
	Reader  io.Reader
	channel chan string
}

func (logstash *Logstash) Pipe() {
	logstash.channel = make(chan string, 4096)

	go logstash.reader()
	go logstash.sender()
}

func (logstash *Logstash) reader() {
	buf := make([]byte, 4096)
	for {
		n, err := logstash.Reader.Read(buf[:])
		if err != nil {
			logger.Error("Can't read: %s", err.Error())
			return
		}

		message := string(buf[0:n])

		if *debug {
			logger.Debug(fmt.Sprintf("Received from socket: %s", message))
		}

		select {
		case logstash.channel <- message:
		default:
		    // dropping old messages if channel is full
			<-logstash.channel
			logstash.channel <- message
		}
	}
}

func (logstash *Logstash) sender() {
	for {
		ServerAddress, err := net.ResolveUDPAddr("udp", logstash.Address)
		if err != nil {
			logger.Error("Can't resolve remote UDP address: %s", err.Error())
		} else {
			Conn, err := net.DialUDP("udp", nil, ServerAddress)
			if err != nil {
				logger.Error("Can't dial up: %s", err.Error())
			} else {
				logger.Notice("Connected to remote UDP address: %s", ServerAddress.String())
				for {
					select {
					case message := <-logstash.channel:
						if *debug {
							logger.Debug(fmt.Sprintf("Writing to UDP socket: %s", message))
						}
						_, err := Conn.Write([]byte(message))
						if err != nil {
							logger.Error("Can't write: %s", err.Error())
							Conn.Close()
							goto RETRY
						}
					}
				}

			}
		}

		RETRY:

		time.Sleep(retryTimeout)
	}
}


