package gtm

import (
	"time"
)

type HealthCheck interface {
	Start() error
	Receive() []IP
}

type Layer4HealthCheckMethod func(string, int) bool

type Layer4HealthCheck struct {
	Port            int
	EndPoints       []IP
	healthEndpoints *[]IP
	Frequency       time.Duration
	control         chan []IP
	CheckHealth     Layer4HealthCheckMethod
}

//Start Start HealthCheck
func (hk *Layer4HealthCheck) Start() {
	hk.healthEndpoints = &hk.EndPoints
	hk.control = make(chan []IP)
	go func() {
		for {
			newEndpoint := make([]IP, 0)
			for _, ip := range hk.EndPoints {
				if hk.CheckHealth(ip.Address, hk.Port) {
					newEndpoint = append(newEndpoint, ip)
				}
			}
			hk.healthEndpoints = &newEndpoint
			hk.control <- newEndpoint
		}
	}()
}

func (hk *Layer4HealthCheck) Receive() []IP {
	var ips []IP
	select {
	case ips = <-hk.control:
	default:
		ips = *hk.healthEndpoints
	}
	return ips
}
