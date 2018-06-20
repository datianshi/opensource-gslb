package gtm

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"time"
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

type Domain struct {
	DomainName        string            `json:"name"`
	IPs               []IP              `json:"ips"`
	TTL               int               `json:"ttl"`
	HealthCheckConfig HealthCheckConfig `json:"health_check"`
	HealthCheck       HealthCheck
}

//HealthCheckConfig Config for health check
type HealthCheckConfig struct {
	Type           string `json:"type"`
	PORT           int    `json:"port"`
	HTTPS          bool   `json:"https"`
	PATH           string `json:"path"`
	HTTPStatusCode int    `json:"http_status_code"`
	Fequency       string `json:"frequency"`
}

type IP struct {
	Address string `json:"address"`
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
	for _, domain := range config.Domains {
		var hk HealthCheckMethod
		if domain.HealthCheckConfig.Type == "layer4" {
			frequency, err := time.ParseDuration(domain.HealthCheckConfig.Fequency)
			if err != nil {
				return nil, fmt.Errorf("format of frequency is not valid: %v", err)
			}
			hk = Layer4HealthCheck(domain.HealthCheckConfig.PORT)
			domain.HealthCheck = &DefaultHealthCheck{
				EndPoints:   domain.IPs,
				Frequency:   frequency,
				CheckHealth: hk,
				SleepFunc:   sleepDuration(frequency),
			}
		} else {
			domain.HealthCheck = &doNothingHealthCheck{
				ips: domain.IPs,
			}
		}

		domain.HealthCheck.Start()

	}
	return &config, err
}
