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
	spec.Run(t, "TestGTM", composer(
		testDNSRequest,
		testHealthCheck,
		testConfig,
		testLayer7HealthCheck,
		testUnmarshalRecord,
		testDNSClient), spec.Report(report.Terminal{}))
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
		var record *Record
		var healtchCheck *fakes.FakeHealthCheck
		var ips []IP
		var msgAnswerCatcher dns.RR
		var dummySelect SelectRecord

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
				if len(msg.Answer) == 0 {
					return nil
				}
				msgAnswerCatcher = msg.Answer[0]
				return nil
			}
			ips = []IP{
				IP{Address: "192.168.0.3"},
				IP{Address: "192.168.0.2"},
			}
		})

		when("Given one record with 2 IP Options", func() {
			it.Before(func() {
				record = &Record{
					Name: "abc",
					IPs:  ips,
					TTL:  5,
				}

				dummySelect = func(q dns.Question, records []*Record, domain string) *Record {
					return record
				}
				healtchCheck.ReceiveStub = func() []IP {
					return ips
				}
				record.HealthCheck = healtchCheck
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
					DNSRequest(LBAnswer([]*Record{record}, dummySelect, "xip.io")(loadBalancer))(dnsWriter, msgIn)
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
					DNSRequest(LBAnswer([]*Record{record}, dummySelect, "xip.io")(loadBalancer))(dnsWriter, msgIn)
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
					DNSRequest(LBAnswer([]*Record{record}, dummySelect, "xip.io")(loadBalancer))(dnsWriter, msgIn)
				})
				it("Should write message: ", func() {
					Ω(msgAnswerCatcher.String()).Should(Equal("abc.xip.io.	5	IN	A	192.168.0.2"))
				})
			})
			when("The select record return nil", func() {
				it.Before(func() {
					var loadBalancer LoadBalancing = func(ips []IP) string {
						return ips[1].Address
					}
					dummySelect = func(q dns.Question, records []*Record, domain string) *Record {
						return nil
					}
					DNSRequest(LBAnswer([]*Record{record}, dummySelect, "xip.io")(loadBalancer))(dnsWriter, msgIn)
				})
				it("Should write empty: ", func() {
					Ω(msgAnswerCatcher).Should(BeNil())
				})
			})
		})
		when("Given serveral records, a domain", func() {
			var (
				records []*Record
				domain  string
			)
			when("records has no wild card", func() {
				it.Before(func() {
					records = []*Record{
						&Record{
							Name: "abc",
						},
						&Record{
							Name: "g",
						},
					}
					domain = "xip.io."
				})
				it("Should return nil for question does not match", func() {
					question := dns.Question{
						Name: "notfound.xip.io.",
					}
					record = DefaultSelectRecord(question, records, domain)
					Ω(record).Should(BeNil())
				})
				it("Should return match record with g", func() {
					question := dns.Question{
						Name: "g.xip.io.",
					}
					record = DefaultSelectRecord(question, records, domain)
					Ω(record.Name).Should(Equal("g"))
				})
			})
			when("records has wild card", func() {
				it.Before(func() {
					records = []*Record{
						&Record{
							Name: "abc",
						},
						&Record{
							Name: "*",
						},
						&Record{
							Name: "g",
						},
					}
					domain = "xip.io."
				})
				it("Should return match record with g", func() {
					question := dns.Question{
						Name: "g.xip.io.",
					}
					record = DefaultSelectRecord(question, records, domain)
					Ω(record.Name).Should(Equal("g"))
				})
				it("Should return match record with wild card", func() {
					question := dns.Question{
						Name: "cheer.xip.io.",
					}
					record = DefaultSelectRecord(question, records, domain)
					Ω(record.Name).Should(Equal("*"))
				})
			})
		})
	})
}
