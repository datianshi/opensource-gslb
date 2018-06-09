package gtm_test

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/datianshi/simple-cf-gtm"
	"github.com/datianshi/simple-cf-gtm/fakes"
	"github.com/miekg/dns"
)

var _ = Describe("RelayClient", func() {
	var client *fakes.FakeDNSClient
	var question dns.Question
	var server string
	var relay *RelayDNSCLient
	var answer []dns.RR

	BeforeEach(func() {
		question = dns.Question{}
		server = "example.com"
		client = &fakes.FakeDNSClient{}
		relay = &RelayDNSCLient{
			Client: client,
		}
	})
	Context("Given the exchange return error", func() {
		BeforeEach(func() {
			client.ExchangeStub = func(m *dns.Msg, address string) (r *dns.Msg, rtt time.Duration, err error) {
				return nil, 5 * time.Second, errors.New("")
			}
			answer = relay.RelayAnswer(server)(question)
		})
		It("Should have empty answers", func() {
			Ω(len(answer)).Should(Equal(0))
		})
	})
	Context("Given the exchange return one answer", func() {
		var msg *dns.Msg
		BeforeEach(func() {
			msg = new(dns.Msg)
			record, _ := dns.NewRR("fake rr")
			msg.Answer = []dns.RR{record}
			client.ExchangeStub = func(m *dns.Msg, address string) (r *dns.Msg, rtt time.Duration, err error) {
				return msg, 5 * time.Second, nil
			}
			answer = relay.RelayAnswer(server)(question)
		})
		It("Should have one answer returned", func() {
			Ω(len(answer)).Should(Equal(1))
		})
	})
})
