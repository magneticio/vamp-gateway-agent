package main

import (
	"time"
	"strings"
	"fmt"
)

func statsReader(socket string, interval int, statsChannel chan []byte) {

	var cmd string = "show stat -1\n"

	Logger.Info(fmt.Sprintf("Connecting to Unix socket: %s", socket))

	for {
		stats, err := command(socket, cmd)
		if err != nil {
			Logger.Panic(err)
		}

		// split on new lines
		lines := strings.Split(stats, "\n")

		// loop over the lines, skipping the first one with the headers
		for _, line := range lines[1:] {
			if len(line) > 0 {
				statsChannel <- []byte(line)
			}
		}
		<-time.After(time.Duration(interval) * time.Second)
	}
}