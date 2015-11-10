package main

import (
	"os"
	"fmt"
	"net"
	"bytes"
	"strings"
	"os/exec"
	"io/ioutil"
)

type HAProxy struct {
	Binary     string
	PidFile    string
	ConfigFile string
	LogSocket  string
}

func (haProxy *HAProxy) Init() {
	//Create and empty pid file on the specified location, if not already there
	if _, err := os.Stat(haProxy.PidFile); err != nil {
		emptyPid := []byte("")
		ioutil.WriteFile(haProxy.PidFile, emptyPid, 0644)
		logger.Info("Created new pidfile.")
	} else {
		logger.Info("Pidfile exists at %s.", haProxy.PidFile)
	}

	logger.Info(fmt.Sprintf("Connecting to haproxy log socket: %s", haProxy.LogSocket))

	// Set the socket HAProxy can write logs to.
	conn, err := net.ListenUnixgram("unixgram", &net.UnixAddr{haProxy.LogSocket, "unixgram"})
	if err != nil {
		logger.Error("Error while connecting to haproxy log socket: ", err.Error())
		return
	}

	logger.Info(fmt.Sprintf("Opened Unix socket at: %s", haProxy.LogSocket))
	logChannel := make(chan string)

	// Start the logging stream.
	go Reader(conn, logChannel)
	go Sender(*logstashHost, *logstashPort, logChannel)
}

func (haProxy *HAProxy) Run() {

	logger.Notice("Reloading HAProxy")

	pid, err := ioutil.ReadFile(haProxy.PidFile)
	if err != nil {
		logger.Error("Error while reloading haproxy: %s", err.Error())
		return
	}

	logger.Info("HAProxy configuration file: %s", haProxy.ConfigFile)

	// Setup all command line parameters so we get an executable like: haproxy -f haproxy.cfg -p haproxy.pid -sf 1234
	arg0 := "-f"
	arg1 := haProxy.ConfigFile
	arg2 := "-p"
	arg3 := haProxy.PidFile
	arg4 := "-D"
	arg5 := "-sf"
	arg6 := strings.Trim(string(pid), "\n")
	var cmd *exec.Cmd

	// If this is the first run, the PID value will be empty, otherwise it will be > 0
	if len(arg6) > 0 {
		logger.Info("HAProxy command: %s %s %s %s %s %s %s %s", haProxy.Binary, arg0, arg1, arg2, arg3, arg4, arg5, arg6)
		cmd = exec.Command(haProxy.Binary, arg0, arg1, arg2, arg3, arg4, arg5, arg6)
	} else {
		logger.Info("HAProxy command: %s %s %s %s %s %s", haProxy.Binary, arg0, arg1, arg2, arg3, arg4)
		cmd = exec.Command(haProxy.Binary, arg0, arg1, arg2, arg3, arg4)
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		logger.Error("Error while reloading haproxy: %s", err.Error())
	}
}

func (haProxy *HAProxy) Reload(configuration []byte) {
	err := ioutil.WriteFile(haProxy.ConfigFile, configuration, 0644)
	if err != nil {
		logger.Error("Error writing to HAProxy configuration. Reloading aborted. %s", err.Error())
	} else {
		haProxy.Run()
	}
}
