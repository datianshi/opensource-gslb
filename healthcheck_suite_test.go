package gtm_test

import (
	"testing"
	"time"

	. "github.com/datianshi/simple-cf-gtm"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func TestHealthcheck(t *testing.T) {
	// spec.Run(t, "TestHealthCheck", testHealthCheck, spec.Report(report.Terminal{}))
}

func testHealthCheck(t *testing.T, when spec.G, it spec.S) {
	it.Before(func() {
		RegisterTestingT(t)
	})

	when("test generic health check", func() {
		var healthCheck DefaultHealthCheck
		var endpoints []IP
		var ip1, ip2, ip3 IP
		var frequency time.Duration
		var port int

		it.Before(func() {
			port = 8080
			ip1, ip2, ip3 = IP{Address: "192.168.0.1"}, IP{Address: "192.168.0.2"}, IP{Address: "192.168.0.3"}
			endpoints = []IP{ip1, ip2, ip3}
		})
		when("Given A health check with 2s frequency And Health check method delete unhealth endpoint ip2", func() {
			var removeIP2 bool
			it.Before(func() {
				removeIP2 = true
				frequency = 2 * time.Second
				healthCheck = DefaultHealthCheck{
					Port:      port,
					EndPoints: endpoints,
					Frequency: frequency,
					CheckHealth: func(ip IP) bool {
						if ip == ip2 && removeIP2 {
							return false
						}
						return true
					},
				}
				healthCheck.Start()
			})
			it("should have 3 ips returns on first run", func() {
				result := healthCheck.Receive()
				Ω(len(result)).Should(Equal(3))
			})
			it("should have 3 ips returns on second run after 1s", func() {
				result := healthCheck.Receive()
				Ω(len(result)).Should(Equal(3))
			})
			it("should only have 2 ips returns on second run after 3s", func() {
				time.Sleep(2 * time.Second)
				result := healthCheck.Receive()
				Ω(len(result)).Should(Equal(2))
			})
			it("should only have 2 ips returns if run again ", func() {
				time.Sleep(1 * time.Second)
				result := healthCheck.Receive()
				Ω(len(result)).Should(Equal(2))
			})
			it("should come back healthy with all three nodes once health check successful", func() {
				removeIP2 = false
				time.Sleep(3 * time.Second)
				result := healthCheck.Receive()
				Ω(len(result)).Should(Equal(3))
			})
		})
	})
}
