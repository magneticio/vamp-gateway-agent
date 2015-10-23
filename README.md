# Vamp-proxy-agent

Vamp-proxy-agent is a tiny helper agent that provides the following services: 

- read logs from HAproxy over sockets and push them to Logstash over UDP
- read statistics from HAproxy and push them to Logstash over UDP

## Usage

Run `vamp-proxy-agent -h` to display usage instructions:

```
$ ./vamp-proxy-agent -h
Usage of ./vamp-proxy-agent:
  -debug
    	Switches on extra log statements
  -haproxyLogSocket string
    	The location of the socket HAproxy logs to (default "/var/run/haproxy.log.sock")
  -haproxyStatsSocket string
    	The location of the HAproxy stats socket (default "/tmp/haproxy.stats.sock")
  -haproxyStatsType string
    	Which stats to read from haproxy: all, frontend, backend or server. (default "all")
  -logHost string
    	Address of the remote Logstash instance (default "127.0.01")
  -logPort int
    	The UDP input port of the remote Logstash instance (default 10002)
  -statsHost string
    	Address of the remote Logstash instance (default "127.0.01")
  -statsPort int
    	The UDP input port of the remote Logstash instance (default 10003)
```

## Example

This example starts a socket at `/var/run/haproxy.sock`. It reads all log lines HAproxy writes to that socket
and passes them to the UDP port 10002 on server 10.4.0.100. Furthermore, it reads the "backend" statistics 
from the `/tmp/haproxy.stats.sock` socket and passes them to port 10003 on the same host.

```bash
$ ./vamp-proxy-agent 
    -logHost 10.4.0.100 \
    -logPort 10002 \
    -haproxyLogSocket /var/run/haproxy.log.sock \
    -statsHost 10.4.0.100 \
    -statsPort 10003 \
    -haproxyStatsSocket /tmp/haproxy.stats.sock \
    -haproxyStatsType backend
   
```

HAproxy should be configured with the following line in its config to actually send the logs to the socket:

```
global
  log /var/run/haproxy.log.sock local2
  stats socket /tmp/haproxy.stats.sock level admin  
```

Logstash should be configure with the following lines in its config to actually read from the UDP sockets
and splits the HAproxy stats in the proper CSV format.

```
input {
  udp {
    port => 10002
    type => haproxy_log
  }
  udp {
    port => 10003
    type => haproxy_stats
  }
}
filter {
  if [type] == "haproxy_stats" {
    csv {
      columns => [pxname,svname,qcur,qmax,scur,smax,slim,stot,bin,bout,dreq,dresp,ereq,econ,eresp,wretr,wredis,status,weight,act,bck,chkfail,chkdown,lastchg,downtime,qlimit,pid,iid,sid,throttle,lbtot,tracked,type,rate,rate_lim,rate_max,check_status,check_code,check_duration,hrsp_1xx,hrsp_2xx,hrsp_3xx,hrsp_4xx,hrsp_5xx,hrsp_other,hanafail,req_rate,req_rate_max,req_tot,cli_abrt,srv_abrt,comp_in,comp_out,comp_byp,comp_rsp,lastsess,last_chk,last_agt,qtime,ctime,rtime,ttime]
    }
  }
}
output {
  stdout { codec => rubydebug }
}
```

If you get some traffic going, you should see Logstash output something similar to this:

```
...
{
       "message" => "<150>Oct 23 15:37:57 haproxy[23433]: { \"timestamp\" : 23/Oct/2015:15:37:57.603, \"frontend\" : \"test_route_2\", \"method\" : \"GET /ping HTTP/1.0\", \"captured_request_headers\" : \"\", \"captures_response_headers\" : \"\" }\n",
      "@version" => "1",
    "@timestamp" => "2015-10-23T13:37:58.649Z",
          "type" => "haproxy_log",
          "host" => "127.0.0.1"
}
{
           "message" => [
        [0] "test_route_2::service_a,BACKEND,0,0,0,10,50000,1002,88803,123250,0,0,,0,0,0,0,UP,100,1,0,,0,96,0,,1,7,0,,1001,,1,0,,469,,,,0,1001,0,1,0,0,,,,,0,0,0,0,0,0,39,,,0,8,6,17,"
    ],
          "@version" => "1",
        "@timestamp" => "2015-10-23T13:27:46.553Z",
              "type" => "1",
              "host" => "127.0.0.1",
            "pxname" => "test_route_2::service_a",
            "svname" => "BACKEND",
              "qcur" => "0",
              "qmax" => "0",
              ...
}
```