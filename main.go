package main

import (
	"os"
	"flag"
	"time"
	"syscall"
	"strings"
	"io/ioutil"
	"os/signal"
	"github.com/samuel/go-zookeeper/zk"
)

const (
	HaProxyPath = "/opt/vamp/"
)

var (
	LogstashHost = flag.String("logstashHost", "127.0.0.1", "Address of the remote Logstash instance")
	LogstashPort = flag.Int("logstashPort", 10001, "The UDP input port of the remote Logstash instance")

	ZooKeeperServers = flag.String("zooKeeperServers", "127.0.0.1:2181", "ZooKeeper servers.")
	ZooKeeperPath = flag.String("zooKeeperPath", "/vamp/gateways/haproxy", "ZooKeeper HAProxy configuration path.")

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

	go ConfigurationChangeWatch(*ZooKeeperServers, *ZooKeeperPath)

	waiter <- true
}

func cleanup() {
	os.Remove(HaProxy.LogSocket)
}

func ConfigurationChangeWatch(servers, path string) {
	Log.Notice("Initializing Zookeeper connection to " + *ZooKeeperServers)
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

			err = ioutil.WriteFile(HaProxy.ConfigFile, payload, 0644)
			if err != nil {
				Log.Error("Error writing to HaProxy configuration. Reloading aborted. %s", err.Error())
			} else {
				HaProxy.Reload()
			}

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
}

