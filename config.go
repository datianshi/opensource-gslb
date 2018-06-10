package gtm

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type Config struct {
	Domains     []Domain `json:"domains"`
	Port        int      `json:"port"`
	RelayServer string   `json:"relay_server"`
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

func ParseConfig(reader io.Reader) (*Config, error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	s := []byte(bytes)
	var domain Config
	err = json.Unmarshal([]byte(s), &domain)
	if err != nil {
		return nil, err
	}
	return &domain, err
}
