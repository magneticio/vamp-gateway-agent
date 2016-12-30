# Vamp Gateway Agent

[![Join the chat at https://gitter.im/magneticio/vamp](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/magneticio/vamp?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Docker](https://img.shields.io/badge/docker-images-blue.svg)](https://hub.docker.com/r/magneticio/vamp-gateway-agent/tags/)
[![Download](https://api.bintray.com/packages/magnetic-io/downloads/vamp-gateway-agent/images/download.svg) ](https://bintray.com/magnetic-io/downloads/vamp-gateway-agent/_latestVersion)

[HAProxy](http://www.haproxy.org/) is a tcp/http load balancer, the purpose of this agent is to: 

- read the HAProxy configuration using [confd](https://github.com/kelseyhightower/confd) and reload HAProxy on each configuration change with as little client traffic interruption as possible.
- send HAProxy log to Logstash/Elasticsearch.
- handle and recover from ZooKeeper, etcd, Consul and Logstash outages without interrupting the haproxy process and client requests.

Vamp generated HAProxy configuration will be appended to base configuration `haproxy.basic.cnf`.
This allows using different base configuration if needed.

## Usage

Following environment variables are mandatory:

- `VAMP_KEY_VALUE_STORE_TYPE <=> confd -backend`
- `VAMP_KEY_VALUE_STORE_CONNECTION <=> confd -node`
- `VAMP_KEY_VALUE_STORE_PATH <=> key used by confd`

Example:

```
docker run -e VAMP_KEY_VALUE_STORE_TYPE=zookeeper \
           -e VAMP_KEY_VALUE_STORE_CONNECTION=localhost:2181 \
           -e VAMP_KEY_VALUE_STORE_PATH=/vamp/gateways/haproxy/1.6 \
           magneticio/vamp-gateway-agent:katana
```

Available Docker images can be found at [Docker Hub](https://hub.docker.com/r/magneticio/vamp-gateway-agent/).

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

To enable dnsmasq to resolve virtual hosts, pass the following environment variables to the Docker container:

- `VAMP_VGA_DNS_ENABLE` Set to non-empty value to enable 
- `VAMP_VGA_DNS_PORT` Listening port, default: 5353

## Building Binary

Using the `build.sh` script:
```
  ./build.sh --make
```

Alternatively:

- `go get github.com/tools/godep`
- `godep restore`
- `go install`
- `CGO_ENABLED=0 go build -v -a -installsuffix cgo`

 
## Building Docker Images

```
$ ./build.sh -h

Usage of ./build.sh:

  -h|--help   Help.
  -l|--list   List built Docker images.
  -r|--remove Remove Docker image.
  -m|--make   Make Docker image files.
  -b|--build  Build Docker image.

```

Docker images after the build (e.g. `./build.sh -b`): 

- magneticio/vamp-gateway-agent:katana
