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
        "name" : "xip.io.",
        "records": [
          {
            "name": "abc",
            "ips": [
                {
                  "address": "10.193.190.181"
                },
                {
                  "address": "10.193.148.200"
                }
            ],
            "ttl" : 5
          },
          {
            "name": "*",
            "ips": [
                {
                  "address": "10.193.190.103"
                },
                {
                  "address": "10.193.148.251"
                }
            ],
            "ttl" : 10,
            "health_check" : {
              "type": "layer4",
              "port": 443,
              "frequency": "5s"
            }
          }
        ]
    },
    {
          "name" : "example.io.",
          "records" : [
            {
              "name": "*",
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
          ]
    }
  ],
  "port" : 5050,
  "relay_server" : "8.8.8.8:53"
}
```

### Domain and records

* It supports wild card records
E.g.
```
dig [anything].xip.io @localhost -p 5050
```
will results either 10.193.190.103 or 10.193.148.251

* It supports a specify record and will take precedence over  
E.g.

```
dig abc.xip.io @localhost -p 5050
```
will results either 10.193.190.181 or 10.193.148.200

### Port

Port Number DNS server should listen on.

### Relay Server

If the records/domains not defined in the configured domains, they will be resolved from the relay server

## Plugin in your own Load Balancing logic

Currently, the load balancing strategy is extremely naive:

```
var simpleLoadBalancer gtm.LoadBalancing = func(ips []gtm.IP) string {
  //Simple Round Robin
  return ips[rand.Intn(len(ips))].Address
}
```

The load balancing logic is pluggable as long as developer implement another load balancer method ```func(ips[]gtm.IP) string```

## Layer4 Health Check

For each record, a layer 4 health check endpoint can be configured

```
"health_check" : {
  "type": "layer4",
  "port": 5000,
  "frequency": "5s"
}
```

## Future work

* Layer 7 Health Check
