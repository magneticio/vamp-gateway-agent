# Vamp Gateway Agent

Vamp Gateway Agent provides the following services: 

- read logs from HAProxy over sockets and push them to Logstash over UDP
- read the HAProxy configuration from ZooKeeper and reloads the HAProxy on each configuration change.

## Usage

Run `vamp-gateway-agent -h` to display the usage instructions:

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

## Building Binary Locally

- `go get github.com/tools/godep`
- `godep restore`
- `go install`
- `go build`

Alternatively using the `build.sh` script:
```
  ./build.sh --make
```
Deliverable is in `target/go` directory.
 
## Building Docker Images Locally

Directory `docker` contains `Dockerfile`s for the following:

- HAProxy 1.5.14
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

- magneticio/vamp-gateway-agent_1.5.14-ubuntu-14.04:0.8.0
- magneticio/vamp-gateway-agent_1.5.14-centos-7:0.8.0
- magneticio/vamp-gateway-agent_1.5.14-alpine-3.2:0.8.0 

## Travis CI Build

On each push to `master` branch another commit is pushed to `docker` branch.
All deliverables from `target/docker` directory are committed to `docker` directory of `docker` branch.

