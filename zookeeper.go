package main

import (
	"time"
	"strings"

	"github.com/samuel/go-zookeeper/zk"
)

type ZooKeeper struct {
	ConnectionString string
	Path             string
	Connection       *zk.Conn
}

func (zooKeeper *ZooKeeper) init() {
	logger.Notice("Initializing ZooKeeper connection: %s", zooKeeper.ConnectionString)
	servers := strings.Split(zooKeeper.ConnectionString, ",")
	conn, _, err := zk.Connect(servers, 60 * time.Second)

	if err != nil {
		logger.Fatal("Error trying to connect to Zookeeper: %s", err.Error())
	} else {
		zooKeeper.Connection = conn
	}
}

func (zooKeeper *ZooKeeper) Watch(onChange func([]byte) error) {

	zooKeeper.init()

	var err error
	var data []byte
	for {
		if zooKeeper.Connection.State() == zk.StateHasSession {
			if *debug {
				logger.Debug("ZooKeeper connection state: %s", zk.StateHasSession)
			}
			// Using GetW(path) would crash the process due to some bug in ZooKeeper client (ZooKeeper start/stop).
			data, _, err = zooKeeper.Connection.Get(zooKeeper.Path)

			if err != nil {
				logger.Info("Reading from ZooKeeper path %s: %s", zooKeeper.Path, err.Error())
			} else {
				onChange(data)
			}
		} else {
			logger.Info("ZooKeeper connection state: %s", zooKeeper.Connection.State())
		}

		time.Sleep(1 * time.Second)
	}
}
