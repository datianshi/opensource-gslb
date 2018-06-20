package gtm_test

import (
	"testing"
	"time"

	. "github.com/datianshi/simple-cf-gtm"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testHealthCheck(t *testing.T, when spec.G, it spec.S) {

	when("test generic health check", func() {
		var healthCheck *DefaultHealthCheck
		var endpoints []IP
		var ip1, ip2, ip3 IP
		var frequency time.Duration
		var port int
		var wakeUp chan bool

		it.Before(func() {
			port = 8080
			ip1, ip2, ip3 = IP{Address: "192.168.0.1"}, IP{Address: "192.168.0.2"}, IP{Address: "192.168.0.3"}
			endpoints = []IP{ip1, ip2, ip3}
			wakeUp = make(chan bool)
		})
		when("Given A health check delete unhealth endpoint ip2", func() {
			var removeIP2 bool
			it.Before(func() {
				removeIP2 = true
				frequency = 2 * time.Millisecond
				healthCheck = &DefaultHealthCheck{
					Port:      port,
					EndPoints: endpoints,
					Frequency: frequency,
					CheckHealth: func(ip IP) bool {
						if ip == ip2 && removeIP2 {
							return false
						}
						return true
					},
					SleepFunc: func() {
						<-wakeUp
					},
				}
				healthCheck.Start()
			})
			it("should have 3 ips returns", func() {
				result := healthCheck.Receive()
				Ω(len(result)).Should(Equal(3))
			})
			it("should have 3 ips returns on second run after", func() {
				healthCheck.Receive()
				result := healthCheck.Receive()
				Ω(len(result)).Should(Equal(3))
			})
			it("should only have 2 ips returns on second run after health check is done", func() {
				healthCheck.Receive()
				wakeUp <- true
				result := healthCheck.Receive()
				Ω(len(result)).Should(Equal(2))
			})
			it("should only have 2 ips returns if run again ", func() {
				healthCheck.Receive()
				wakeUp <- true
				healthCheck.Receive()
				result := healthCheck.Receive()
				Ω(len(result)).Should(Equal(2))
			})
			it("should come back healthy with all three nodes once health check successful", func() {
				healthCheck.Receive()
				wakeUp <- true
				result := healthCheck.Receive()
				Ω(len(result)).Should(Equal(2))
				removeIP2 = false
				wakeUp <- true
				time.Sleep(1 * time.Second)
				result = healthCheck.Receive()
				Ω(len(result)).Should(Equal(3))
			})
		})
	})
}
