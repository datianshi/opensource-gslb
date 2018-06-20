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
		for _, record := range domain.Records {
			if record.HealthCheckConfig.Type == "layer4" {
				frequency, err := time.ParseDuration(record.HealthCheckConfig.Fequency)
				if err != nil {
					return nil, fmt.Errorf("format of frequency is not valid: %v", err)
				}
				hk = Layer4HealthCheck(record.HealthCheckConfig.PORT)
				record.HealthCheck = &DefaultHealthCheck{
					EndPoints:   record.IPs,
					Frequency:   frequency,
					CheckHealth: hk,
					SleepFunc:   sleepDuration(frequency),
				}
			} else {
				record.HealthCheck = &doNothingHealthCheck{
					ips: record.IPs,
				}
			}
			record.HealthCheck.Start()
		}

	}
	return &config, err
}
