package gtm_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/datianshi/simple-cf-gtm"
)

var _ = Describe("Config", func() {
	var configString string
	var err error
	var config *Config

	Context("Given a valid config string", func() {
		BeforeEach(func() {
			configString = `
			{
			  "domains": [
			    {
			          "name" : ".xip.io",
			          "ips": [
			              {
			                "address": "192.168.0.2"
			              },
			              {
			                "address": "192.168.0.3"
			              }
			          ],
			          "ttl" : 5
			    }
			  ],
			  "port" : 5050,
				"relay_server" : "8.8.8.8:53"
			}
			`
			config, err = ParseConfig(strings.NewReader(configString))
		})
		It("Should not have err happen", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})
		It("Should have a valid domain config", func() {
			Ω(len(config.Domains)).Should(Equal(1))
		})
		It("The port should be 5050", func() {
			Ω(config.Port).Should(Equal(5050))
		})
		It("The domain ttl should be 5", func() {
			Ω(config.Domains[0].TTL).Should(Equal(5))
		})
		It("The domain's name should be .xip.io", func() {
			Ω(config.Domains[0].DomainName).Should(Equal(".xip.io"))
		})
		It("The domain have two ips", func() {
			Ω(len(config.Domains[0].IPs)).Should(Equal(2))
		})
		It("The first ip is 192.168.0.2", func() {
			Ω(config.Domains[0].IPs[0].Address).Should(Equal("192.168.0.2"))
		})
		It("Relay Server should be 8.8.8.8:53", func() {
			Ω(config.RelayServer).Should(Equal("8.8.8.8:53"))
		})
	})
})
