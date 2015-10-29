package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type HaProxy struct {
	Binary     string
	PidFile    string
	ConfigFile string
}

// returns an error if the file was already there
func (h *HaProxy) SetPid() error {
	//Create and empty pid file on the specified location, if not already there
	if _, err := os.Stat(h.PidFile); err != nil {
		emptyPid := []byte("")
		ioutil.WriteFile(h.PidFile, emptyPid, 0644)
		return nil
	}
	return errors.New("file already there")
}

// Reload runtime with configuration
func (h *HaProxy) Reload() error {

	pid, err := ioutil.ReadFile(h.PidFile)
	if err != nil {
		return err
	}

	// Setup all the command line parameters so we get an executable similar to
	// /usr/local/bin/haproxy -f resources/haproxy_new.cfg -p resources/haproxy-private.pid -sf 1234
	arg0 := "-f"
	arg1 := h.ConfigFile
	arg2 := "-p"
	arg3 := h.PidFile
	arg4 := "-D"
	arg5 := "-sf"
	arg6 := strings.Trim(string(pid), "\n")
	var cmd *exec.Cmd

	// fmt.Println(r.Binary + " " + arg0 + " " + arg1 + " " + arg2 + " " + arg3 + " " + arg4 + " " + arg5 + " " + arg6)
	// If this is the first run, the PID value will be empty, otherwise it will be > 0
	if len(arg6) > 0 {
		cmd = exec.Command(h.Binary, arg0, arg1, arg2, arg3, arg4, arg5, arg6)
	} else {
		cmd = exec.Command(h.Binary, arg0, arg1, arg2, arg3, arg4)
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	cmdErr := cmd.Run()
	if cmdErr != nil {
		return cmdErr
	}

	return nil
}
