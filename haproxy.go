package main

import (
	"os"
	"fmt"
	"net"
	"bytes"
	"strings"
	"strconv"
	"os/exec"
	"io/ioutil"
)

type HAProxy struct {
	Binary      string
	PidFile     string
	BasicConfig string
	ConfigFile  string
	LogSocket   string
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

	logger.Info(fmt.Sprintf("Opened Unix socket at: %s. Creating Logstash sender.", haProxy.LogSocket))

	logstash := Logstash{
		Address: *logstashHost + ":" + strconv.Itoa(*logstashPort),
		Reader: conn,
	}
	logstash.Pipe()
}

func (haProxy *HAProxy) Run() error {

	logger.Notice("Reloading HAProxy")

	pid, err := ioutil.ReadFile(haProxy.PidFile)
	if err != nil {
		logger.Error("Error while reloading haproxy: %s", err.Error())
		return err
	}

	logger.Info("HAProxy configuration file: %s", haProxy.ConfigFile)

	// Setup all command line parameters so we get an executable like: haproxy -f haproxy.cfg -p haproxy.pid -st 1234
	arg0 := "-f"
	arg1 := haProxy.ConfigFile
	arg2 := "-p"
	arg3 := haProxy.PidFile
	arg4 := "-D"
	arg5 := "-st"
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

	return err
}

func (haProxy *HAProxy) Reload(configuration []byte) error {
	basic, err := ioutil.ReadFile(haProxy.BasicConfig)
	if err != nil {
		logger.Error("Cannot read basic HAProxy configuration: %s", err.Error())
		return err
	}

	content := append(basic[:], configuration[:]...)

	if !haProxy.changed(content) {
		if *debug {
			logger.Debug("Configuration change has been triggered, but no change in data.")
		}
		return nil
	}

	err = haProxy.validate(content)

	if err == nil {
		err = ioutil.WriteFile(haProxy.ConfigFile, content, 0644)
		if err != nil {
			logger.Error("Error writing to HAProxy configuration. Reloading aborted. %s", err.Error())
			return err
		}
		return haProxy.Run()
	} else {
		logger.Error("Reloading the new HAProxy configuration has been aborted.")
	}

	return err
}

func (haProxy *HAProxy) changed(configuration []byte) bool {
	file, err := ioutil.ReadFile(haProxy.ConfigFile)

	if (err == nil) {
		if bytes.Compare(file, configuration) != 0 {
			logger.Notice("HAProxy configuration has been changed.")
			return true
		}
		return false
	}

	return true
}

func (haProxy *HAProxy) validate(configuration []byte) error {
	logger.Notice("Validating the new HAProxy configuration.")

	err := ioutil.WriteFile(haProxy.ConfigFile + ".tmp", configuration, 0644)
	if err != nil {
		logger.Error("Error writing to temp HAProxy configuration. %s", err.Error())
		return err
	}

	arg0 := "-c"
	arg1 := "-f"
	arg2 := haProxy.ConfigFile + ".tmp"

	cmd := exec.Command(haProxy.Binary, arg0, arg1, arg2)
	var out bytes.Buffer
	cmd.Stderr = &out

	err = cmd.Run()
	if err != nil {
		logger.Error("Error while validating the new HAProxy configuration: %s - %s", err.Error(), string(out.Bytes()[:]))
	}

	return err
}


