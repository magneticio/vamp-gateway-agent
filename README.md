# Vamp Gateway Agent

HAProxy with configuration from [ZooKeeper](https://zookeeper.apache.org/), [etcd](https://coreos.com/etcd/docs/latest/) or [Consul](https://consul.io/).

[![Build Status](https://travis-ci.org/magneticio/vamp-gateway-agent.svg?branch=master)](https://travis-ci.org/magneticio/vamp-gateway-agent)

[HAProxy](http://www.haproxy.org/) is a tcp/http load balancer, the purpose of this agent is to: 

- read the HAProxy configuration from [ZooKeeper](https://zookeeper.apache.org/), [etcd](https://coreos.com/etcd/docs/latest/) or [Consul](https://consul.io/) and reloads the HAProxy on each configuration change with as close to zero client request interruption as possible.
- read the logs from HAProxy over socket and push them to Logstash over UDP.
- handle and recover from ZooKeeper, etcd, Consul and Logstash outages without interrupting the haproxy process and client requests.

It is possible to specify a custom configuration (based on arguments `configurationPath/configurationBasicFile`).
In that case any configuration read from KV store is appended to the content of custom configuration, stored as `configurationPath/haproxy.cnf` and used for the next HAProxy reload.

## Usage

```
$ ./vamp-gateway-agent: -h
                                       
Usage of ./vamp-gateway-agent:
  -configurationBasicFile string
        Basic HAProxy configuration. (default "haproxy.basic.cfg")
  -configurationPath string
        HAProxy configuration path. (default "/opt/vamp/")
  -debug
        Switches on extra log statements.
  -help
        Print usage.
  -logo
        Show logo. (default true)
  -logstashHost string
        Address of the Logstash instance (default "127.0.0.1")
  -logstashPort int
        The UDP input port of the Logstash instance (default 10001)
  -retryTimeout int
        Default retry timeout in seconds. (default 5)
  -scriptPath
        HAProxy validation and reload script path. (default "/opt/vamp/")
  -storeConnection string
        Key-value store connection string.
  -storeKey string
        HAProxy configuration store key. (default "/vamp/gateways/haproxy/1.6")
  -storeType string
        zookeeper, consul or etcd.

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

**Note:** Logstash configuration depends on HAProxy log configuration and that is not in the scope of the agent (HAProxy configuration is retrieved from [ZooKeeper](https://zookeeper.apache.org/), [etcd](https://coreos.com/etcd/docs/latest/) or [Consul](https://consul.io/)). 

## Building Binary

- `go get github.com/tools/godep`
- `godep restore`
- `go install`
- `CGO_ENABLED=0 go build -v -a -installsuffix cgo`

Alternatively using the `build.sh` script:
```
  ./build.sh --make
```
Deliverable is in `target/go` directory.
 
## Building Docker Images

Directory `haproxy` contains `Dockerfile`s for HAProxy 1.5.15 and 1.6.3 **Alpine** based images.

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

- magneticio/vamp-gateway-agent_1.6.3:0.8.4
- magneticio/vamp-gateway-agent_1.5.15:0.8.4

## Travis CI Build

Build is performed on each push to `master` branch and all directories from `target/docker` are pushed to specific version branch (e.g. 0.8.4).
After that Docker Hub Automated Build is triggered.

## Docker Images

[Docker Hub Repo](https://hub.docker.com/r/magneticio/vamp-gateway-agent/)

[![](https://badge.imagelayers.io/magneticio/vamp-gateway-agent:1.6.3-0.8.4.svg)](https://imagelayers.io/?images=magneticio/vamp-gateway-agent:1.6.3-0.8.4) 1.6.3-0.8.4

[![](https://badge.imagelayers.io/magneticio/vamp-gateway-agent:1.5.15-0.8.4.svg)](https://imagelayers.io/?images=magneticio/vamp-gateway-agent:1.5.15-0.8.4) 1.5.15-0.8.4

e.g.

```
docker run --net=host --restart=always magneticio/vamp-gateway-agent:1.6.3-0.8.4
```
