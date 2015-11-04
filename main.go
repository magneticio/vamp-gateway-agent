package main

import (
	"os"
	"syscall"
	"os/signal"
	"flag"
)

const (
	HaProxyPath = "/opt/vamp/"
)

var (
	LogstashHost = flag.String("logHost", "127.0.0.1", "Address of the remote Logstash instance")
	LogstashPort = flag.Int("logPort", 10001, "The UDP input port of the remote Logstash instance")

	ZooKeeperServers = flag.String("zkServers", "127.0.0.1:2181", "ZooKeeper servers.")
	ZooKeeperPath = flag.String("zkPath", "/vamp/gateways/haproxy", "ZooKeeper HAProxy configuration path.")

	DebugSwitch = flag.Bool("debug", false, "Switches on extra log statements")

	Log = CreateLogger()

	HaProxy = HAProxy{
		Binary:     "haproxy",
		ConfigFile: HaProxyPath + "haproxy.cfg",
		PidFile:    HaProxyPath + "haproxy.pid",
		LogSocket:  HaProxyPath + "haproxy.log.sock",
	}
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
	Log.Notice(Logo("0.8.0"))
	flag.Parse()
	Log.Notice("Starting Vamp Gateway Agent")

	// Waiter keeps the program from exiting instantly.
	waiter := make(chan bool)

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

	HaProxy.Init()
	HaProxy.Reload()

	watchZooKeeper()

	waiter <- true
}

func cleanup() {
	os.Remove(HaProxy.LogSocket)
}

func watchZooKeeper() {
	Log.Notice("Initializing Zookeeper connection to " + *ZooKeeperServers)
	zkClient := ZkClient{}
	go zkClient.Watch(*ZooKeeperServers, *ZooKeeperPath)
}
