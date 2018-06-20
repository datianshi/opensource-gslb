package gtm

import (
	"fmt"

	"github.com/miekg/dns"
)

// LoadBalancing define how to load balance
type LoadBalancing func([]IP) string

type GetAnswer func(dns.Question) []dns.RR

type ServeDNS func(dns.ResponseWriter, *dns.Msg)

type SelectRecord func(q dns.Question, records []*Record, domain string) *Record

//LBAnswer Given ips and ttl configuration, return a Get Answer func
func LBAnswer(records []*Record, getRecord SelectRecord, domain string) func(loadBalancer LoadBalancing) GetAnswer {
	return func(loadBalancer LoadBalancing) GetAnswer {
		return func(q dns.Question) []dns.RR {
			record := getRecord(q, records, domain)
			if record == nil {
				return make([]dns.RR, 0)
			}
			ip := loadBalancer(record.HealthCheck.Receive())
			if ip == "" {
				return make([]dns.RR, 0)
			}
			rr, err := dns.NewRR(fmt.Sprintf("%s %d A %s", q.Name, record.TTL, ip))
			if err != nil {
				return make([]dns.RR, 0)
			}
			return []dns.RR{rr}
		}
	}
}

const WILD_CARD string = "*"

func DefaultSelectRecord(q dns.Question, records []*Record, domain string) *Record {
	var wildCard *Record
	for _, record := range records {
		if record.Name == WILD_CARD {
			wildCard = record
			continue
		}
		fqdn := fmt.Sprintf("%s.%s", record.Name, domain)
		if q.Name == fqdn {
			return record
		}
	}
	return wildCard
}

// func (lb LoadBalancing) withHealthCheck(frequency time.Duration, hk HealthCheck) LoadBalancing {
//
// }

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
