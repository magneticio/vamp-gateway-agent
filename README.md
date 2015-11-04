# Vamp Gateway Agent

Vamp Gateway Agent provides the following services: 

- read logs from HAProxy over sockets and push them to Logstash over UDP
- read the HAProxy configuration from ZooKeeper and reloads the HAProxy on configuration change.

## Usage

Run `vamp-gateway-agent -h` to display usage instructions:

```
$ ./vamp-proxy-agent -h

██╗   ██╗ █████╗ ███╗   ███╗██████╗
██║   ██║██╔══██╗████╗ ████║██╔══██╗
██║   ██║███████║██╔████╔██║██████╔╝
╚██╗ ██╔╝██╔══██║██║╚██╔╝██║██╔═══╝
 ╚████╔╝ ██║  ██║██║ ╚═╝ ██║██║
  ╚═══╝  ╚═╝  ╚═╝╚═╝     ╚═╝╚═╝
                       gateway agent
                       version 0.8.0
                       by magnetic.io
                                       
Usage of /opt/vamp/vamp-gateway-agent:
  -debug
    	Switches on extra log statements
  -logstashHost string
    	Address of the remote Logstash instance (default "127.0.0.1")
  -logstashPort int
    	The UDP input port of the remote Logstash instance (default 10001)
  -zooKeeperPath string
    	ZooKeeper HAProxy configuration path. (default "/vamp/gateways/haproxy")
  -zooKeeperServers string
    	ZooKeeper servers. (default "127.0.0.1:2181")
```

Logstash example configuration:

```
input {
  udp {
    port => 10001
    type => haproxy_log
  }
}
output {
  stdout { codec => rubydebug }
}
```

## Building Docker images

Directory `docker` contains `Dockerfile`s for following:

- HAProxy 1.5.14
- Ubuntu 14.04, CentOS 7 and Alpine 3.2

```
$ ./docker.sh -h

Usage of ./docker.sh:
  -h|--help   Help.
  -l|--list   List all available images.
  -c|--clean  Remove all available images.
  -b|--build  Build all available images.

```

Available Docker images after build: 

- magneticio/vamp-gateway-agent_1.5.14-ubuntu-14.04:0.8.0
- magneticio/vamp-gateway-agent_1.5.14-centos-7:0.8.0
- magneticio/vamp-gateway-agent_1.5.14-alpine-3.2:0.8.0 
