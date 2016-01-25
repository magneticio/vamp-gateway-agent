package main

import (
	"time"
	"strings"

	"golang.org/x/net/context"
	"github.com/coreos/etcd/client"
)

type Etcd struct {
	ConnectionString string
	Path             string
	KApi             client.KeysAPI
}

func (etcd *Etcd) init() {
	logger.Notice("Initializing etcd connection: %s", etcd.ConnectionString)
	servers := strings.Split(etcd.ConnectionString, ",")
	cfg := client.Config{
		Endpoints:               servers,
		Transport:               client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: retryTimeout,
	}
	c, err := client.New(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	etcd.KApi = client.NewKeysAPI(c)
}

func (etcd *Etcd) Watch(onChange func([]byte) error) {
	etcd.init()

	for {
		opts := &client.WatcherOptions{Recursive: false}
		key := strings.TrimPrefix(etcd.Path, "/")
		watcher := etcd.KApi.Watcher(key, opts)

		result, err := etcd.KApi.Get(context.Background(), etcd.Path, nil)
		if err == nil {
			onChange([]byte(result.Node.Value))
		}

		logger.Infof("Watching for Etcd change of: %s", etcd.Path)
		for {
			result, err := watcher.Next(context.Background())
			if err != nil {
				logger.Info("Etcd connection error: %s", err)
				break
			}
			onChange([]byte(result.Node.Value))
		}
		time.Sleep(retryTimeout)
	}
}
