package main

import (
	"os"
	"flag"
	"syscall"
	"os/signal"
)

const (
	HAProxyPath = "/opt/vamp/"
)

var (
	logstashHost = flag.String("logstashHost", "127.0.0.1", "Address of the remote Logstash instance")
	logstashPort = flag.Int("logstashPort", 10001, "The UDP input port of the remote Logstash instance")

	storeType = flag.String("storeType", "zookeeper", "zookeeper, consul or etcd.")
	storeServers = flag.String("storeServers", "127.0.0.1:2181", "Key-value store servers.")
	configurationPath = flag.String("configurationPath", "vamp/gateways/haproxy/1.6", "HAProxy configuration path.")

	logo = flag.Bool("logo", true, "Show logo.")

	debug = flag.Bool("debug", false, "Switches on extra log statements.")

	logger = CreateLogger()
)

func Logo(version string) string {
	return `
██╗   ██╗ █████╗ ███╗   ███╗██████╗
██║   ██║██╔══██╗████╗ ████║██╔══██╗
██║   ██║███████║██╔████╔██║██████╔╝
╚██╗ ██╔╝██╔══██║██║╚██╔╝██║██╔═══╝
 ╚████╔╝ ██║  ██║██║ ╚═╝ ██║██║
  ╚═══╝  ╚═╝  ╚═╝╚═╝     ╚═╝╚═╝
                       gateway agent
                       version ` + version + `
                       by magnetic.io
                                      `
}

func main() {

	flag.Parse()

	if *logo {
		logger.Notice(Logo("0.8.2"))
	}

	logger.Notice("Starting Vamp Gateway Agent")

	haProxy := HAProxy{
		Binary:     "haproxy",
		ConfigFile: HAProxyPath + "haproxy.cfg",
		PidFile:    HAProxyPath + "haproxy.pid",
		LogSocket:  HAProxyPath + "haproxy.log.sock",
	}

	// Waiter keeps the program from exiting instantly.
	waiter := make(chan bool)

	cleanup := func() { os.Remove(haProxy.LogSocket) }

	// Catch a CTR+C exits so the cleanup routine is called.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	defer cleanup()

	haProxy.Init()
	haProxy.Run()

	keyValue(haProxy)

	waiter <- true
}

func keyValue(haProxy HAProxy) {
	if *storeType == "zookeeper" {
		zooKeeper := ZooKeeper{
			Servers: *storeServers,
			Path: *configurationPath,
		}
		zooKeeper.Init()
		go zooKeeper.Watch(haProxy.Reload)
	} else if *storeType == "etcd" {
		etcd := Etcd{
			Servers: *storeServers,
			Path: *configurationPath,
		}
		etcd.Init()
		go etcd.Watch(haProxy.Reload)
	} else {
		logger.Panic("Key-value store type not supported: ", *storeType)
	}
}

