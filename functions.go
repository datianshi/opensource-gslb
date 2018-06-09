package gtm

import (
	"fmt"

	"github.com/miekg/dns"
)

// LoadBalancing define how to load balance
type LoadBalancing func([]IP) string

type GetAnswer func(dns.Question) []dns.RR

type HealthCheck func(IP) bool

type Layer7HealthCheck func(string, string) HealthCheck

type Layer4HealthCheck func(int) HealthCheck

type ServeDNS func(dns.ResponseWriter, *dns.Msg)

func LBAnswer(ips []IP, ttl int) func(loadBalancer LoadBalancing) GetAnswer {
	return func(loadBalancer LoadBalancing) GetAnswer {
		return func(q dns.Question) []dns.RR {
			ip := loadBalancer(ips)
			if ip == "" {
				return make([]dns.RR, 0)
			}
			rr, err := dns.NewRR(fmt.Sprintf("%s %d A %s", q.Name, ttl, ip))
			if err != nil {
				return make([]dns.RR, 0)
			}
			return []dns.RR{rr}
		}
	}
}

// DNSRequest Serve DNS Request
func DNSRequest(answerFunc GetAnswer) ServeDNS {
	return func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		m.SetReply(r)
		m.Compress = false

		switch r.Opcode {
		case dns.OpcodeQuery:
			for _, q := range m.Question {
				switch q.Qtype {
				case dns.TypeA:
					m.Answer = append(m.Answer, answerFunc(q)...)
				}
			}
		}

		w.WriteMsg(m)
	}
}
