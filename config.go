package gtm

import "encoding/json"

type Config struct {
	Domains []Domain `json:"domains"`
	Port    int      `json:"port"`
}

type Domain struct {
	DomainName string `json:"name"`
	IPs        []IP   `json:"ips"`
	TTL        int    `json:"ttl"`
}
type IP struct {
	Address      string `json:"address"`
	HealthCheckM HealthCheck
}

func ParseConfig(s string) (*Config, error) {
	var domain Config
	err := json.Unmarshal([]byte(s), &domain)
	if err != nil {
		return nil, err
	}
	return &domain, err
}
