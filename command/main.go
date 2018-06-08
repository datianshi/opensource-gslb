package main

import (
	"fmt"
	"log"
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
	var SimpleLoadBalancer gtm.LoadBalancing = func(ips []gtm.IP) string {
		return ips[1].Address
	}
	// start server
	for _, domain := range config.Domains {
		dns.HandleFunc(domain.DomainName, gtm.DNSRequest(domain)(SimpleLoadBalancer))
	}
	server := &dns.Server{Addr: ":" + strconv.Itoa(config.Port), Net: "udp"}
	log.Printf("Starting at %d\n", config.Port)
	err = server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
