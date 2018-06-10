package gtm_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/datianshi/simple-cf-gtm"
)

var _ = Describe("Healthcheck", func() {
	var healthCheck Layer4HealthCheck
	var endpoints []IP
	var ip1, ip2, ip3 IP
	var frequency time.Duration
	var port int

	BeforeEach(func() {
		port = 8080
		ip1, ip2, ip3 = IP{Address: "192.168.0.1"}, IP{Address: "192.168.0.2"}, IP{Address: "192.168.0.3"}
		endpoints = []IP{ip1, ip2, ip3}
	})

	Context("Given A health check with 2s frequency And Health check method delete unhealth endpoint ip2", func() {
		var removeIP2 bool
		BeforeEach(func() {
			removeIP2 = true
			frequency = 2 * time.Second
			healthCheck = Layer4HealthCheck{
				Port:      port,
				EndPoints: endpoints,
				Frequency: frequency,
				CheckHealth: func(ip string, port int) bool {
					if ip == ip2.Address && removeIP2 {
						return false
					}
					return true
				},
			}
			healthCheck.Start()
		})
		It("should have 3 ips returns on first run", func() {
			result := healthCheck.Receive()
			Ω(len(result)).Should(Equal(3))
		})
		It("should have 3 ips returns on second run after 1s", func() {
			result := healthCheck.Receive()
			Ω(len(result)).Should(Equal(3))
		})
		It("should only have 2 ips returns on second run after 3s", func() {
			time.Sleep(2 * time.Second)
			result := healthCheck.Receive()
			Ω(len(result)).Should(Equal(2))
		})
		It("should only have 2 ips returns if run again ", func() {
			time.Sleep(1 * time.Second)
			result := healthCheck.Receive()
			Ω(len(result)).Should(Equal(2))
		})
		It("should come back healthy with all three nodes once health check successful", func() {
			removeIP2 = false
			time.Sleep(3 * time.Second)
			result := healthCheck.Receive()
			Ω(len(result)).Should(Equal(3))
		})
	})
})
