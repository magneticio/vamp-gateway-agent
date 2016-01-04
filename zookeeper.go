package main

import (
	"time"
	"bytes"
	"strings"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/zookeeper"
)

type ZooKeeper struct {
	Servers string
	Path    string
	Store   store.Store
}

func (zooKeeper *ZooKeeper) Init() {
	zookeeper.Register()
}

func (zooKeeper *ZooKeeper) Connect() (error) {
	logger.Notice("Openning ZooKeeper connection: %s", zooKeeper.Servers)
	zks := strings.Split(zooKeeper.Servers, ",")

	kv, err := libkv.NewStore(
		store.ZK,
		zks,
		&store.Config{
			ConnectionTimeout: 60 * time.Second,
		},
	)
	if err != nil {
		logger.Fatal("Error trying to connect to Zookeeper: ", err.Error())
	} else {
		zooKeeper.Store = kv
	}

	return err
}

func (zooKeeper *ZooKeeper) Close() {
	zooKeeper.Store.Close()
}

func (zooKeeper *ZooKeeper) Watch(onChange func([]byte)) {
	var data []byte
	for {

		err := zooKeeper.Connect()

		if err == nil {

			// Listening for the value change doesn't work if ZK server is restarted, thus using polling

			for {
				exists, err := zooKeeper.Store.Exists(zooKeeper.Path)

				if err != nil {
					logger.Info("Error trying to connect to Zookeeper: ", err.Error())
					zooKeeper.Close()
					break

				} else {
					if exists {
						pair, err := zooKeeper.Store.Get(zooKeeper.Path)

						if err != nil {
							logger.Info("Reading from ZooKeeper path %s: %s", zooKeeper.Path, err.Error())
						} else if bytes.Compare(data, pair.Value) != 0 {
							logger.Notice("ZooKeeper %s data has been changed.", zooKeeper.Path)
							data = pair.Value
							onChange(data)
						}
					} else {
						logger.Info("ZooKeeper path does not exist: %s", zooKeeper.Path)
					}
				}

				time.Sleep(1 * time.Second)
			}
		}

		time.Sleep(5 * time.Second)
	}
}
