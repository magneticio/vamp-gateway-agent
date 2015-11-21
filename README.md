# Vamp Gateway Agent

HAProxy with configuration from ZooKeeper

[![Build Status](https://travis-ci.org/magneticio/vamp-gateway-agent.svg?branch=master)](https://travis-ci.org/magneticio/vamp-gateway-agent)

[HAProxy](http://www.haproxy.org/) is a tcp/http load balancer, the purpose of this agent is to: 

- read the HAProxy configuration from ZooKeeper and reloads the HAProxy on each configuration change with as close to zero client request interruption as possible.
- read the logs from HAProxy over socket and push them to Logstash over UDP.
- handle and recover from ZooKeeper and Logstash outages without interrupting the haproxy process and client requests.

## Usage

```
$ ./vamp-proxy-agent -h
                                       
Usage of ./vamp-gateway-agent:
  -debug
        Switches on extra log statements.
  -logo
        Show logo. (default true)
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
    type => haproxy
  }
}

filter {
  grok {
    match => ["message", "^.+?]: (?<metrics>{.*})$"]
    match => ["message", "^.*$"]
  }
  if [metrics] =~ /.+/ {
    json {
      source => "metrics"
    }
    if [t] =~ /.+/ {
      date {
        match => ["t", "dd/MMM/YYYY:HH:mm:ss.SSS"]
      }
    }
  }
}

output {
  elasticsearch {
    hosts => "elasticsearch:9200"
  }
  stdout {
    codec => rubydebug
  }
}
```

**Note:** Logstash configuration depends on HAProxy log configuration and that is not in the scope of the agent (HAProxy configuration is retrieved from ZooKeeper). 

## Building Binary

- `go get github.com/tools/godep`
- `godep restore`
- `go install`
- `go build`

Alternatively using the `build.sh` script:
```
  ./build.sh --make
```
Deliverable is in `target/go` directory.
 
## Building Docker Images

Directory `docker` contains `Dockerfile`s for the following:

- HAProxy 1.5.15
- Ubuntu 14.04, CentOS 7 and Alpine 3.2

```
$ ./build.sh -h

Usage of ./build.sh:

  -h|--help   Help.
  -l|--list   List all available images.
  -c|--clean  Remove all available images.
  -m|--make   Build vamp-gateway-agent binary and copy it to Docker directories.
  -b|--build  Build all available images.

```

Docker images after the build (e.g. `./build.sh -b`): 

- magneticio/vamp-gateway-agent_1.5.15-ubuntu-14.04:0.8.0
- magneticio/vamp-gateway-agent_1.5.15-centos-7:0.8.0
- magneticio/vamp-gateway-agent_1.5.15-alpine-3.2:0.8.0 

## Travis CI Build

Build is performed on each push to `master` branch and all directories from `target/docker` are pushed to `docker` branch.

## Docker Images

**Alpine**

[![](https://badge.imagelayers.io/magneticio/vamp-gateway-agent_1.5.15-alpine-3.2:0.8.0.svg)](https://imagelayers.io/?images=magneticio/vamp-gateway-agent_1.5.15-alpine-3.2:0.8.0) 1.5.15-alpine-3.2:0.8.0

e.g.

```
docker run --net=host --restart=always magneticio/vamp-gateway-agent_1.5.15-alpine-3.2:0.8.0
```

**CentOS**

[![](https://badge.imagelayers.io/magneticio/vamp-gateway-agent_1.5.15-centos-7:0.8.0.svg)](https://imagelayers.io/?images=magneticio/vamp-gateway-agent_1.5.15-centos-7:0.8.0) 1.5.15-centos-7:0.8.0

e.g.

```
docker run --net=host --restart=always magneticio/vamp-gateway-agent_1.5.15-centos-7:0.8.0
```

**Ubuntu**

[![](https://badge.imagelayers.io/magneticio/vamp-gateway-agent_1.5.15-ubuntu-14.04:0.8.0.svg)](https://imagelayers.io/?images=magneticio/vamp-gateway-agent_1.5.15-ubuntu-14.04:0.8.0) 1.5.15-ubuntu-14.04:0.8.0

e.g.

```
docker run --net=host --restart=always magneticio/vamp-gateway-agent_1.5.15-ubuntu-14.04:0.8.0
```
