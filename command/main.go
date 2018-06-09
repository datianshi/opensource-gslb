package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"

	gtm "github.com/datianshi/simple-cf-gtm"
	"github.com/miekg/dns"
)

func main() {

	config, err := gtm.ParseConfig(os.Getenv("GTM_CONFIG"))
	if err != nil {
		fmt.Println(err)
		return
	}
	var simpleLoadBalancer gtm.LoadBalancing = func(ips []gtm.IP) string {
		//Simple Round Robin
		return ips[rand.Intn(len(ips))].Address
	}
	// start server
	for _, domain := range config.Domains {
		dns.HandleFunc(domain.DomainName, gtm.DNSRequest(gtm.LBAnswer(domain.IPs, domain.TTL)(simpleLoadBalancer)))
	}

	if config.RelayServer != "" {
		relayClient := &gtm.RelayDNSCLient{
			Client: new(dns.Client),
		}
		dns.HandleFunc(".", gtm.DNSRequest(relayClient.RelayAnswer(config.RelayServer)))
	}

	server := &dns.Server{Addr: ":" + strconv.Itoa(config.Port), Net: "udp"}
	log.Printf("Starting at %d\n", config.Port)
	err = server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
