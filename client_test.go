package gtm_test

import (
	"errors"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/datianshi/simple-cf-gtm"
	"github.com/datianshi/simple-cf-gtm/fakes"
	"github.com/miekg/dns"
	"github.com/sclevine/spec"
)

func testDNSClient(t *testing.T, when spec.G, it spec.S) {
	var client *fakes.FakeDNSClient
	var question dns.Question
	var server string
	var relay *RelayDNSCLient
	var answer []dns.RR

	it.Before(func() {
		question = dns.Question{}
		server = "example.com"
		client = &fakes.FakeDNSClient{}
		relay = &RelayDNSCLient{
			Client: client,
		}
	})
	when("Given the exchange return error", func() {
		BeforeEach(func() {
			client.ExchangeStub = func(m *dns.Msg, address string) (r *dns.Msg, rtt time.Duration, err error) {
				return nil, 5 * time.Second, errors.New("")
			}
			answer = relay.RelayAnswer(server)(question)
		})
		it("Should have empty answers", func() {
			Ω(len(answer)).Should(Equal(0))
		})
	})
	when("Given the exchange return one answer", func() {
		var msg *dns.Msg
		it.Before(func() {
			msg = new(dns.Msg)
			record, _ := dns.NewRR("fake rr")
			msg.Answer = []dns.RR{record}
			client.ExchangeStub = func(m *dns.Msg, address string) (r *dns.Msg, rtt time.Duration, err error) {
				return msg, 5 * time.Second, nil
			}
			answer = relay.RelayAnswer(server)(question)
		})
		it("Should have one answer returned", func() {
			Ω(len(answer)).Should(Equal(1))
		})
	})
}
