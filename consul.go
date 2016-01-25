package main

import (
	"time"
	"github.com/hashicorp/consul/api"
)

type Consul struct {
	ConnectionString string
	Path             string
}

func (consul *Consul) Watch(onChange func([]byte) error) {
	for {
		conf := api.DefaultConfig()
		conf.Address = consul.ConnectionString
		client, err := api.NewClient(conf)
		opts := &api.QueryOptions{WaitTime: retryTimeout}

		if err != nil {
			logger.Error("Error connecting to Consul: %s", err.Error())
		} else {
			for {
				pair, meta, err := client.KV().Get(consul.Path, opts)

				if err != nil {
					logger.Error("Error getting a value: %s", err.Error())
					break
				} else if opts.WaitIndex == meta.LastIndex || pair == nil {
					continue
				}

				opts.WaitIndex = meta.LastIndex
				onChange(pair.Value)
			}
		}

		time.Sleep(retryTimeout)
	}
}
