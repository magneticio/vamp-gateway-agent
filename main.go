package main

import (
    "os"
    "time"
    "flag"
    "syscall"
    "os/signal"
    "io/ioutil"
)

var (
    logstashHost = flag.String("logstashHost", "127.0.0.1", "Address of the Logstash instance")
    logstashPort = flag.Int("logstashPort", 10001, "The UDP input port of the Logstash instance")

    storeType = flag.String("storeType", "", "zookeeper, consul or etcd.")
    storeConnection = flag.String("storeConnection", "", "Key-value store connection string.")
    storeKey = flag.String("storeKey", "/vamp/haproxy/1.6", "HAProxy configuration store key.")

    configurationPath = flag.String("configurationPath", "/opt/vamp/", "HAProxy configuration path.")
    configurationBasicFile = flag.String("configurationBasicFile", "haproxy.basic.cfg", "Basic HAProxy configuration.")

    scriptPath = flag.String("scriptPath", "/opt/vamp/", "HAProxy validation and reload script path.")

    timeout = flag.Int("retryTimeout", 5, "Default retry timeout in seconds.")

    logo = flag.Bool("logo", true, "Show logo.")
    help = flag.Bool("help", false, "Print usage.")
    debug = flag.Bool("debug", false, "Switches on extra log statements.")

    retryTimeout = 5 * time.Second
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

type Watcher interface {
    Watch(onChange func([]byte) error)
}

func main() {

    flag.Parse()

    if *logo {
        logger.Notice(Logo("0.8.5"))
    }

    if *help {
        flag.Usage()
        return
    }

    if len(*storeType) == 0 {
        logger.Panic("Key-value store type not speciffed.")
        return
    }

    if len(*storeConnection) == 0 {
        logger.Panic("Key-value store servers not speciffed.")
        return
    }

    if _, err := os.Stat(*configurationPath + *configurationBasicFile); os.IsNotExist(err) {
        logger.Panic("No basic HAProxy configuration: ", *configurationPath, *configurationBasicFile)
        return
    }

    retryTimeout = time.Duration(*timeout) * time.Second

    logger.Notice("Starting Vamp Gateway Agent")

    haProxy := HAProxy{
        ScriptPath:    *scriptPath,
        BasicConfig:        *configurationPath + *configurationBasicFile,
        ConfigFile:         *configurationPath + "haproxy.cfg",
        LogSocket:          *configurationPath + "haproxy.log.sock",
    }

    if _, err := os.Stat(haProxy.ConfigFile); os.IsNotExist(err) {
        basic, err := ioutil.ReadFile(haProxy.BasicConfig)
        if err != nil {
            logger.Panic("Cannot read basic HAProxy configuration: ", haProxy.BasicConfig)
            return
        }
        ioutil.WriteFile(haProxy.ConfigFile, basic, 0644)
    }

    // Waiter keeps the program from exiting instantly.
    waiter := make(chan bool)

    cleanup := func() {
        os.Remove(haProxy.LogSocket)
    }

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

    keyValueWatcher := keyValueWatcher()

    if keyValueWatcher == nil {
        return
    }

    go keyValueWatcher.Watch(haProxy.Reload)

    waiter <- true
}

func keyValueWatcher() Watcher {
    if *storeType == "etcd" {
        return &Etcd{
            ConnectionString: *storeConnection,
            Path: *storeKey,
        }
    } else if *storeType == "consul" {
        return &Consul{
            ConnectionString: *storeConnection,
            Path: *storeKey,
        }
    } else if *storeType == "zookeeper" {
        return &ZooKeeper{
            ConnectionString: *storeConnection,
            Path: *storeKey,
        }
    } else {
        logger.Panic("Key-value store type not supported: ", *storeType)
        return nil
    }
}
