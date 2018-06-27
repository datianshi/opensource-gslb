package gtm

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type HealthCheck interface {
	Start()
	Receive() []IP
}

type Config struct {
	Domains     []*Domain `json:"domains"`
	Port        int       `json:"port"`
	RelayServer string    `json:"relay_server"`
}

type Record struct {
	Name              string            `json:"name"`
	IPs               []IP              `json:"ips"`
	TTL               int               `json:"ttl"`
	HealthCheckConfig HealthCheckConfig `json:"health_check"`
	HealthCheck       HealthCheck
}

type Domain struct {
	DomainName string    `json:"name"`
	Records    []*Record `json:"records"`
}

//HealthCheckConfig Config for health check
type HealthCheckConfig struct {
	Type           string `json:"type"`
	PORT           int    `json:"port"`
	HTTPS          bool   `json:"https"`
	SkipSSL        bool   `json:"skip_ssl"`
	PATH           string `json:"path"`
	HTTPStatusCode int    `json:"http_status_code"`
	Fequency       string `json:"frequency"`
}

type IP struct {
	Address string `json:"address"`
	Host    string `json:"layer7_health_check_host"`
}

func ParseConfig(reader io.Reader) (*Config, error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	s := []byte(bytes)
	var config Config
	err = json.Unmarshal([]byte(s), &config)
	if err != nil {
		return nil, err
	}
	return &config, err
}
