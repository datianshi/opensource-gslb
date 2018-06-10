# A DNS Load Balancer by golang

## Run it

```
gtm --config config.json
```

## Configuration

```
{
  "domains": [
    {
          "name" : "xip.io",
          "ips": [
              {
                "address": "192.168.0.2"
              },
              {
                "address": "192.168.0.3"
              }
          ],
          "ttl" : 5
    },
    {
          "name" : "example.io",
          "ips": [
              {
                "address": "10.0.5.4"
              },
              {
                "address": "10.0.5.5"
              }
          ],
          "ttl" : 5
    }
  ],
  "port" : 5050,
  "relay_server" : "8.8.8.8:53"
}
```

### Domains

For each domain, it will resolve every records underneath it to the provided ips

E.g.

```
dig [anything].xip.xio @localhost -p 5050
```

will results either 192.168.0.2 or 192.168.0.3

### Port

Port Number DNS server should listen on.

### Relay Server

If the records/domains not defined in the domains, they will be resolved from the relay server

## Plugin in your own Load Balancing logic

Currently, the load balancing strategy is extremely naive:

```
var simpleLoadBalancer gtm.LoadBalancing = func(ips []gtm.IP) string {
  //Simple Round Robin
  return ips[rand.Intn(len(ips))].Address
}
```

The load balancing logic is pluggable as long as developer implement another load balancer method ```func(ips[]gtm.IP) string```

##Future works

* Specify particular A records underneath a domain
* Configurable health check 
