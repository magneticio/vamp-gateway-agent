package main

import (
	"time"
	"bytes"
	"strings"

	"github.com/samuel/go-zookeeper/zk"
)

type ZooKeeper struct {
	Servers    string
	Path       string
	Connection *zk.Conn
}

func (zooKeeper *ZooKeeper) init() {
	logger.Notice("Initializing ZooKeeper connection: %s", zooKeeper.Servers)
	servers := strings.Split(zooKeeper.Servers, ",")
	conn, _, err := zk.Connect(servers, 60 * time.Second)

	if err != nil {
		logger.Fatal("Error trying to connect to Zookeeper: %s", err.Error())
	} else {
		zooKeeper.Connection = conn
	}
}

func (zooKeeper *ZooKeeper) Watch(onChange func([]byte)) {

	zooKeeper.init()

	var err error
	var oldData, newData []byte
	for {
		if zooKeeper.Connection.State() == zk.StateHasSession {
			if *debug {
				logger.Debug("ZooKeeper connection state: %s", zk.StateHasSession)
			}
			// Using GetW(path) would crash the process due to some bug in ZooKeeper client (ZooKeeper start/stop).
			newData, _, err = zooKeeper.Connection.Get(zooKeeper.Path)

			if err != nil {
				logger.Info("Reading from ZooKeeper path %s: %s", zooKeeper.Path, err.Error())
			} else if bytes.Compare(oldData, newData) != 0 {
				logger.Notice("ZooKeeper %s data has been changed.", zooKeeper.Path)
				oldData = newData
				onChange(oldData)
			}
		} else {
			logger.Info("ZooKeeper connection state: %s", zooKeeper.Connection.State())
		}

		time.Sleep(5 * time.Second)
	}
}
