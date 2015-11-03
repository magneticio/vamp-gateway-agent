package main

import (
	"github.com/samuel/go-zookeeper/zk"
	"strings"
	"time"
	"io/ioutil"
)

type ZkClient struct {
}

func (z *ZkClient) Watch(servers, path string) error {
	zks := strings.Split(servers, ",")
	conn, _, err := zk.Connect(zks, (60 * time.Second))

	if err == nil {
		Logger.Notice("ZooKeeper path: %s", path)
		for {
			payload, _, watch, err := conn.GetW(path)

			if err != nil {
				Logger.Error("Error from Zookeeper: " + err.Error())
				break
			}

			reload(payload)

			event := <-watch

			if event.Type == zk.EventNodeDataChanged {
				Logger.Notice("ZooKeeper configuration changed.")
			}
		}
	} else {
		Logger.Error("Error connecting to Zookeeper: " + err.Error())
	}

	if conn != nil {
		conn.Close()
	}

	Logger.Notice("ZooKeeper stop monitoring.")

	return nil
}

func reload(payload []byte) {
	Logger.Notice("Reloading HaProxy")
	err := ioutil.WriteFile(haProxy.ConfigFile, payload, 0644)
	if err != nil {
		Logger.Error("Writing to HaProxy configuration. Reloading abort.")
	} else {
		haProxy.Reload()
	}
}
