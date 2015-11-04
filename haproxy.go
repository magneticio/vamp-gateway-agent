package main

import (
	"fmt"
	"net"
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type HAProxy struct {
	Binary     string
	PidFile    string
	ConfigFile string
	LogSocket  string
}

// returns an error if the file was already there
func (h *HAProxy) Init() {
	//Create and empty pid file on the specified location, if not already there
	if _, err := os.Stat(h.PidFile); err != nil {
		emptyPid := []byte("")
		ioutil.WriteFile(h.PidFile, emptyPid, 0644)
		Log.Info("Created new pidfile.")
	} else {
		Log.Info("Pidfile exists at %s.", HaProxy.PidFile)
	}

	Log.Info(fmt.Sprintf("Connecting to haproxy log socket: %s", HaProxy.LogSocket))

	// Set the socket HaProxy can write logs to.
	conn, err := net.ListenUnixgram("unixgram", &net.UnixAddr{HaProxy.LogSocket, "unixgram"})
	if err != nil {
		Log.Error("Error while connecting to haproxy log socket: ", err.Error())
		return
	}

	Log.Info(fmt.Sprintf("Opened Unix socket at: %s", HaProxy.LogSocket))
	logChannel := make(chan []byte, 65536)

	// Start the logging stream.
	go Reader(conn, logChannel)
	go Sender(*LogstashHost, *LogstashPort, logChannel)
}

// Reload runtime with configuration
func (h *HAProxy) Reload() {

	Log.Notice("Reloading HaProxy")

	pid, err := ioutil.ReadFile(h.PidFile)
	if err != nil {
		Log.Error("Error while reloading haproxy: %s", err.Error())
		return
	}

	Log.Info("HaProxy configuration file: %s", h.ConfigFile)

	/*
	 * Setup all the command line parameters so we get an executable similar to
	 *  /usr/local/bin/haproxy -f haproxy.cfg -p haproxy.pid -sf 1234
	 *
	 */
	arg0 := "-f"
	arg1 := h.ConfigFile
	arg2 := "-p"
	arg3 := h.PidFile
	arg4 := "-D"
	arg5 := "-sf"
	arg6 := strings.Trim(string(pid), "\n")
	var cmd *exec.Cmd

	// If this is the first run, the PID value will be empty, otherwise it will be > 0
	if len(arg6) > 0 {
		Log.Info("HaProxy command: %s %s %s %s %s %s %s %s", h.Binary, arg0, arg1, arg2, arg3, arg4, arg5, arg6)
		cmd = exec.Command(h.Binary, arg0, arg1, arg2, arg3, arg4, arg5, arg6)
	} else {
		Log.Info("HaProxy command: %s %s %s %s %s %s", h.Binary, arg0, arg1, arg2, arg3, arg4)
		cmd = exec.Command(h.Binary, arg0, arg1, arg2, arg3, arg4)
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		Log.Error("Error while reloading haproxy: %s", err.Error())
	}
}


