package gtm_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/datianshi/simple-cf-gtm"
	"github.com/datianshi/simple-cf-gtm/fakes"
	"github.com/miekg/dns"
)

var _ = Describe("DNSRequest", func() {
	var dnsWriter *fakes.FakeResponseWriter
	var msgIn *dns.Msg
	var domain Domain

	var msgAnswerCatcher dns.RR

	BeforeEach(func() {
		msgIn = &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Opcode: dns.OpcodeQuery,
			},
		}
		msgIn.Question = []dns.Question{dns.Question{
			Name:  "abc.xip.io",
			Qtype: dns.TypeA,
		},
		}
		dnsWriter = &fakes.FakeResponseWriter{}
		dnsWriter.WriteMsgStub = func(msg *dns.Msg) error {
			msgAnswerCatcher = msg.Answer[0]
			return nil
		}
	})

	Context("Given a domain has 2 IP Options", func() {
		BeforeEach(func() {
			domain = Domain{
				DomainName: ".xip.io",
				IPs: []IP{
					IP{Address: "192.168.0.3"},
					IP{Address: "192.168.0.2"},
				},
				TTL: 5,
			}

		})
		Context("Given a Load Balancing function return first IP", func() {
			BeforeEach(func() {
				var loadBalancer LoadBalancing = func(ips []IP) string {
					return ips[0].Address
				}
				DNSRequest(LBAnswer(domain.IPs, domain.TTL)(loadBalancer))(dnsWriter, msgIn)
			})
			It("Should write message: ", func() {
				Ω(msgAnswerCatcher.String()).Should(Equal("abc.xip.io.	5	IN	A	192.168.0.3"))
			})
		})
		Context("Given a Load Balancing function return second IP", func() {
			BeforeEach(func() {
				var loadBalancer LoadBalancing = func(ips []IP) string {
					return ips[1].Address
				}
				DNSRequest(LBAnswer(domain.IPs, domain.TTL)(loadBalancer))(dnsWriter, msgIn)
			})
			It("Should write message: ", func() {
				Ω(msgAnswerCatcher.String()).Should(Equal("abc.xip.io.	5	IN	A	192.168.0.2"))
			})
		})

	})

})
