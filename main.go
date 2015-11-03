package main

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"os/signal"
	"flag"
)

const (
	HaProxyPath = "/opt/vamp/"
	HaProxyBinary = "haproxy"
	HaProxyConfigFile = HaProxyPath + "haproxy.cfg"
	HaProxyPidFile = HaProxyPath + "haproxy.pid"
	HaProxyLogSocket = HaProxyPath + "haproxy.log.sock"

	ZooKeeperPath = "/vamp/gateways/haproxy"

	LogPath = HaProxyPath + "gateway-agent.log"
)

var (
	LogHost = flag.String("logHost", "127.0.0.1", "Address of the remote Logstash instance")
	LogPort = flag.Int("logPort", 10002, "The UDP input port of the remote Logstash instance")
	ZooKeeperServers = flag.String("zkServers", "127.0.0.1:2181", "ZooKeeper servers.")

	Logger = ConfigureLog(LogPath, true)
	DebugSwitch = flag.Bool("debug", true, "Switches on extra log statements")

	haProxy = HaProxy{
		Binary:     HaProxyBinary,
		ConfigFile: HaProxyConfigFile,
		PidFile:    HaProxyPidFile,
	}
)

func main() {

	Logger.Info(PrintLogo("0.8.0"))

	flag.Parse()

	Logger.Info("Starting Vamp Gateway Agent")

	// waiter keeps the program from exiting instantly
	waiter := make(chan bool)

	// catch an CTR+C exits so the cleanup routine is called
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	defer cleanup()

	loadHaProxy()
	collectHaProxyLogs()
	watchZooKeeper()

	waiter <- true
}

func PrintLogo(version string) string {
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

func cleanup() {
	os.Remove(HaProxyLogSocket)
}

func loadHaProxy() {
	err := haProxy.SetPid()
	if err != nil {
		Logger.Notice("Pidfile exists at %s, proceeding...", haProxy.PidFile)
	} else {
		Logger.Notice("Created new pidfile...")
	}

	err = haProxy.Reload()
	if err != nil {
		Logger.Fatal("Error while reloading haproxy: " + err.Error())
		panic(err)
	}
}

func collectHaProxyLogs() {
	// set the socket HaProxy can write logs to
	conn, err := net.ListenUnixgram("unixgram", &net.UnixAddr{HaProxyLogSocket, "unixgram"})
	if err != nil {
		Logger.Fatal("Error while connecting to haproxy log socket: " + err.Error())
	}

	Logger.Info(fmt.Sprintf("Opened Unix socket at: %s", HaProxyLogSocket))

	logChannel := make(chan []byte, 1000000)

	// start the logging stream
	go simpleReader(conn, logChannel)
	go sender(*LogHost, *LogPort, logChannel)
}

func watchZooKeeper() {
	Logger.Info("Initializing Zookeeper connection to " + *ZooKeeperServers)
	zkClient := ZkClient{}
	go zkClient.Watch(*ZooKeeperServers, ZooKeeperPath)
}