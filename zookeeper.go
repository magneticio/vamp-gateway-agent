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
		Log.Notice("ZooKeeper path: %s", path)
		for {
			payload, _, watch, err := conn.GetW(path)

			if err != nil {
				Log.Error("Error from Zookeeper: %s", err.Error())
				break
			}

			Reload(payload)

			event := <-watch

			if event.Type == zk.EventNodeDataChanged {
				Log.Notice("ZooKeeper configuration changed.")
			}
		}
	} else {
		Log.Error("Error connecting to Zookeeper: %s", err.Error())
	}

	if conn != nil {
		conn.Close()
	}

	Log.Info("ZooKeeper stop monitoring.")

	return nil
}

func Reload(payload []byte) {
	err := ioutil.WriteFile(HaProxy.ConfigFile, payload, 0644)
	if err != nil {
		Log.Error("Error writing to HaProxy configuration. Reloading aborted. %s", err.Error())
	} else {
		HaProxy.Reload()
	}
}
