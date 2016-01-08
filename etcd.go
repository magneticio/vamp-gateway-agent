package main

import (
	"time"
	"bytes"
	"strings"

	"golang.org/x/net/context"
	"github.com/coreos/etcd/client"
)

type Etcd struct {
	Servers string
	Path    string
	KApi    client.KeysAPI
}

func (etcd *Etcd) init() {
	logger.Notice("Initializing etcd connection: %s", etcd.Servers)
	servers := strings.Split(etcd.Servers, ",")
	cfg := client.Config{
		Endpoints:               servers,
		Transport:               client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: 5 * time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	etcd.KApi = client.NewKeysAPI(c)
}

func (etcd *Etcd) Watch(onChange func([]byte)) {

	etcd.init()

	var data []byte

	for {
		opts := &client.WatcherOptions{Recursive: false}
		key := strings.TrimPrefix(etcd.Path, "/")
		watcher := etcd.KApi.Watcher(key, opts)

		result, err := etcd.KApi.Get(context.Background(), etcd.Path, nil)
		if err == nil {
			data = etcd.change(data, []byte(result.Node.Value), onChange)
		}

		logger.Info("Watching for Etcd change of: ", etcd.Path)
		for {
			result, err := watcher.Next(context.Background())
			if err != nil {
				logger.Info("Etcd connection error: %s", err)
				break
			}
			data = etcd.change(data, []byte(result.Node.Value), onChange)
		}
		time.Sleep(5 * time.Second)
	}
}

func (etcd *Etcd) change(oldData []byte, newData []byte, onChange func([]byte)) []byte {
	if bytes.Compare(oldData, newData) != 0 {
		logger.Notice("Etcd %s data has been changed.", etcd.Path)
		onChange(oldData)
		return newData
	}
	return oldData
}

