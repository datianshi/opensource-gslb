package gtm

import (
	"time"
)

type HealthCheckMethod func(IP) bool

type doNothingHealthCheck struct {
	ips []IP
}

//Start empty start method
func (hk *doNothingHealthCheck) Start() {}

//Receive empty
func (hk *doNothingHealthCheck) Receive() []IP {
	return hk.ips
}

type DefaultHealthCheck struct {
	Port            int
	EndPoints       []IP
	healthEndpoints *[]IP
	Frequency       time.Duration
	control         chan []IP
	CheckHealth     HealthCheckMethod
}

//Start Start HealthCheck
func (hk *DefaultHealthCheck) Start() {
	hk.healthEndpoints = &hk.EndPoints
	hk.control = make(chan []IP)
	go func() {
		for {
			newEndpoint := make([]IP, 0)
			for _, ip := range hk.EndPoints {
				if hk.CheckHealth(ip) {
					newEndpoint = append(newEndpoint, ip)
				}
			}
			hk.healthEndpoints = &newEndpoint
			hk.control <- newEndpoint
		}
	}()
}

func (hk *DefaultHealthCheck) Receive() []IP {
	var ips []IP
	select {
	case ips = <-hk.control:
	default:
		ips = *hk.healthEndpoints
	}
	return ips
}
