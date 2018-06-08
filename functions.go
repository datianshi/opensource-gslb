package gtm

import (
	"fmt"
	"log"

	"github.com/miekg/dns"
)

// LoadBalancing define how to load balance
type LoadBalancing func([]IP) string

type HealthCheck func(IP) bool

type Layer7HealthCheck func(string, string) HealthCheck

type Layer4HealthCheck func(int) HealthCheck

type ServeDNS func(dns.ResponseWriter, *dns.Msg)

// DNSRequest Serve DNS Request
func DNSRequest(domain Domain) func(loadBalancer LoadBalancing) ServeDNS {
	return func(loadBalancer LoadBalancing) ServeDNS {
		return func(w dns.ResponseWriter, r *dns.Msg) {
			m := new(dns.Msg)
			m.SetReply(r)
			m.Compress = false

			switch r.Opcode {
			case dns.OpcodeQuery:
				for _, q := range m.Question {
					switch q.Qtype {
					case dns.TypeA:
						log.Printf("Query for %s\n", q.Name)
						ip := loadBalancer(domain.IPs)
						if ip != "" {
							rr, err := dns.NewRR(fmt.Sprintf("%s %d A %s", q.Name, domain.TTL, ip))
							if err == nil {
								m.Answer = append(m.Answer, rr)
							}
						}
					}
				}
			}

			w.WriteMsg(m)
		}
	}
}
