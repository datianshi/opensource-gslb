package gtm

import (
	"log"
	"time"

	"github.com/miekg/dns"
)

//DNSClient An interface of dns client
type DNSClient interface {
	Exchange(m *dns.Msg, address string) (r *dns.Msg, rtt time.Duration, err error)
}

//RelayDNSClient A client to relay dns request to another DNS server
type RelayDNSCLient struct {
	Client DNSClient
}

func (c *RelayDNSCLient) RelayAnswer(server string) GetAnswer {
	return func(q dns.Question) []dns.RR {
		m1 := new(dns.Msg)
		m1.Id = dns.Id()
		m1.RecursionDesired = true
		m1.Question = make([]dns.Question, 1)
		m1.Question[0] = q
		in, _, err := c.Client.Exchange(m1, server)
		if err != nil {
			log.Println(err)
			return make([]dns.RR, 0)
		}
		return in.Answer
	}
}
