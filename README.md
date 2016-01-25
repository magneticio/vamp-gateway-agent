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

Directory `docker` contains `Dockerfile`s for the following:

- HAProxy 1.6.3
- Ubuntu 14.04, CentOS 7 and Alpine 3.3

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

- magneticio/vamp-gateway-agent_1.6.3-ubuntu-14.04:0.8.2
- magneticio/vamp-gateway-agent_1.6.3-centos-7:0.8.2
- magneticio/vamp-gateway-agent_1.6.3-alpine-3.3:0.8.2 

## Travis CI Build

Build is performed on each push to `master` branch and all directories from `target/docker` are pushed to specific version branch (e.g. 0.8.2).
After that Docker Hub Automated Build is triggered.

## Docker Images

[Docker Hub Repo](https://hub.docker.com/r/magneticio/vamp-gateway-agent/)

**Alpine**

[![](https://badge.imagelayers.io/magneticio/vamp-gateway-agent:1.6.3-alpine-3.3-0.8.2.svg)](https://imagelayers.io/?images=magneticio/vamp-gateway-agent:1.6.3-alpine-3.3-0.8.2) 1.6.3-alpine-3.3-0.8.2

e.g.

```
docker run --net=host --restart=always magneticio/vamp-gateway-agent:1.6.3-alpine-3.3-0.8.2
```

**CentOS**

[![](https://badge.imagelayers.io/magneticio/vamp-gateway-agent:1.6.3-centos-7-0.8.2.svg)](https://imagelayers.io/?images=magneticio/vamp-gateway-agent:1.6.3-centos-7-0.8.2) 1.6.3-centos-7-0.8.2

e.g.

```
docker run --net=host --restart=always magneticio/vamp-gateway-agent:1.6.3-centos-7-0.8.2
```

**Ubuntu**

[![](https://badge.imagelayers.io/magneticio/vamp-gateway-agent:1.6.3-ubuntu-14.04-0.8.2.svg)](https://imagelayers.io/?images=magneticio/vamp-gateway-agent:1.6.3-ubuntu-14.04-0.8.2) 1.6.3-ubuntu-14.04-0.8.2

e.g.

```
docker run --net=host --restart=always magneticio/vamp-gateway-agent:1.6.3-ubuntu-14.04-0.8.2
```
