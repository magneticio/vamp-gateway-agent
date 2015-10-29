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
	for {

		zks := strings.Split(servers, ",")
		conn, _, err := zk.Connect(zks, (60 * time.Second))

		if err == nil {
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
			time.Sleep(5 * time.Second)
		}

		if conn != nil {
			conn.Close()
		}
	}
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
