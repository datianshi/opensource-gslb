package gtm_test

import (
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/datianshi/simple-cf-gtm"
	"github.com/datianshi/simple-cf-gtm/fakes"
	"github.com/miekg/dns"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestGTM(t *testing.T) {
	spec.Run(t, "TestGTM", composer(testDNSRequest, testHealthCheck, testConfig), spec.Report(report.Terminal{}))
}

type specFunc func(*testing.T, spec.G, spec.S)

func composer(funcs ...specFunc) specFunc {
	return func(t *testing.T, g spec.G, s spec.S) {
		s.Before(func() {
			RegisterTestingT(t)
		})
		for _, f := range funcs {
			f(t, g, s)
		}
	}
}

func testDNSRequest(t *testing.T, when spec.G, it spec.S) {

	when("test dns request", func() {
		var dnsWriter *fakes.FakeResponseWriter
		var msgIn *dns.Msg
		var domain Domain
		var healtchCheck *fakes.FakeHealthCheck
		var ips []IP
		var msgAnswerCatcher dns.RR

		it.Before(func() {
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
			healtchCheck = &fakes.FakeHealthCheck{}
			dnsWriter.WriteMsgStub = func(msg *dns.Msg) error {
				msgAnswerCatcher = msg.Answer[0]
				return nil
			}
			ips = []IP{
				IP{Address: "192.168.0.3"},
				IP{Address: "192.168.0.2"},
			}
		})

		when("Given a domain has 2 IP Options", func() {
			it.Before(func() {
				domain = Domain{
					DomainName: ".xip.io",
					IPs:        ips,
					TTL:        5,
				}
				healtchCheck.ReceiveStub = func() []IP {
					return ips
				}
			})
			when("Given a Health function remove one IP", func() {
				var (
					catchNumberIps int
				)
				it.Before(func() {
					healtchCheck.ReceiveStub = func() []IP {
						return []IP{ips[0]}
					}
					var loadBalancer LoadBalancing = func(ips []IP) string {
						catchNumberIps = len(ips)
						return ips[0].Address
					}
					DNSRequest(LBAnswer(domain.IPs, domain.TTL, healtchCheck)(loadBalancer))(dnsWriter, msgIn)
				})
				it("Should have only one ip passed in", func() {
					Ω(catchNumberIps).Should(Equal(1))
				})
			})
			when("Given a Load Balancing function return first IP", func() {
				it.Before(func() {
					var loadBalancer LoadBalancing = func(ips []IP) string {
						return ips[0].Address
					}
					DNSRequest(LBAnswer(domain.IPs, domain.TTL, healtchCheck)(loadBalancer))(dnsWriter, msgIn)
				})
				it("Should write message: ", func() {
					Ω(msgAnswerCatcher.String()).Should(Equal("abc.xip.io.	5	IN	A	192.168.0.3"))
				})
			})
			when("Given a Load Balancing function return second IP", func() {
				it.Before(func() {
					var loadBalancer LoadBalancing = func(ips []IP) string {
						return ips[1].Address
					}
					DNSRequest(LBAnswer(domain.IPs, domain.TTL, healtchCheck)(loadBalancer))(dnsWriter, msgIn)
				})
				it("Should write message: ", func() {
					Ω(msgAnswerCatcher.String()).Should(Equal("abc.xip.io.	5	IN	A	192.168.0.2"))
				})
			})
		})
	})
}
