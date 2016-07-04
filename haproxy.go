package main

import (
    "fmt"
    "net"
    "bytes"
    "os/exec"
    "io/ioutil"
)

type HAProxy struct {
    ScriptPath  string
    BasicConfig string
    ConfigFile  string
    LogSocket   string
}

func (haProxy *HAProxy) Init() {

    if len(*logstash) == 0 {
        logger.Notice("No Logstash host:port set - not sending HAProxy logs.")
        return
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
        Address: *logstash,
        Reader: conn,
    }

    logstash.Pipe()
}

func (haProxy *HAProxy) Run() error {
    script := haProxy.ScriptPath + "reload.sh"
    logger.Notice("Reloading HAProxy - configuration file: %s, reload script: %s", haProxy.ConfigFile, script)

    cmd := exec.Command("/bin/sh", script, haProxy.ConfigFile)

    var out bytes.Buffer
    cmd.Stdout = &out

    err := cmd.Run()
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

    script := haProxy.ScriptPath + "validate.sh"
    logger.Notice("Validating the new HAProxy configuration by running script: %s", script)

    tempFile := haProxy.ConfigFile + ".tmp"
    err := ioutil.WriteFile(tempFile, configuration, 0644)
    if err != nil {
        logger.Error("Error writing to temp HAProxy configuration. %s", err.Error())
        return err
    }

    cmd := exec.Command("/bin/sh", script, tempFile)
    var out bytes.Buffer
    cmd.Stderr = &out

    err = cmd.Run()
    if err != nil {
        logger.Error("Error while validating the new HAProxy configuration: %s - %s", err.Error(), string(out.Bytes()[:]))
    }

    return err
}
