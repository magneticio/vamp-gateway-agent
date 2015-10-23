package main

import(
	"time"
	"strings"
	"fmt"
)


func statsReader(socket, statsType string, statsChannel chan []byte){

	var cmd string

	switch statsType {
		case "all":
			cmd = "show stat -1\n"
		case "backend":
			cmd = "show stat -1 2 -1\n"
		case "frontend":
			cmd = "show stat -1 1 -1\n"
		case "server":
			cmd = "show stat -1 4 -1\n"
	}


	Info(fmt.Sprintf("Connecting to Unix socket: %s",socket))

	for {
		stats, err := command(socket,cmd)
		if err != nil {
			Error(err)
		}

		// split on new lines
		lines := strings.Split(stats,"\n")

		// loop over the lines, skipping the first one with the headers
		for _,line := range lines[1:] {
			if len(line) > 0 {
				statsChannel <- []byte(line)
			}
		}
		<- time.After(3 * time.Second)
	}
}