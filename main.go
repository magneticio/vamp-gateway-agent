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
	HaProxyStatsSocket = HaProxyPath + "haproxy.stats.sock"

	LogPath = HaProxyPath + "gateway-agent.log"
)

var (
	LogHost = flag.String("logHost", "127.0.0.1", "Address of the remote Logstash instance")
	LogPort = flag.Int("logPort", 10002, "The UDP input port of the remote Logstash instance")
	StatsHost = flag.String("statsHost", "127.0.0.1", "Address of the remote Logstash instance")
	StatsPort = flag.Int("statsPort", 10003, "The UDP input port of the remote Logstash instance")
	StatsPollInterval = flag.Int("pollInterval", 5, "How often (in seconds) to poll for statistics.")

	Logger = ConfigureLog(LogPath, true)
	DebugSwitch = flag.Bool("debug", false, "Switches on extra log statements")
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
	collectLogs()
	collectStats()

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
	haProxy := HaProxy{
		Binary:     HaProxyBinary,
		ConfigFile: HaProxyConfigFile,
		PidFile:    HaProxyPidFile,
	}

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

func collectLogs() {
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

func collectStats() {
	statsChannel := make(chan []byte, 1000000)
	// start the stats stream
	go statsReader(HaProxyStatsSocket, *StatsPollInterval, statsChannel)
	go sender(*StatsHost, *StatsPort, statsChannel)
}
