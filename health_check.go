package gtm

import (
	"log"
	"time"
)

type HealthCheckMethod func(IP) bool

type DoNothingHealthCheck struct {
	ips []IP
}

//Start empty start method
func (hk *DoNothingHealthCheck) Start() {}

//Receive empty
func (hk *DoNothingHealthCheck) Receive() []IP {
	return hk.ips
}

type sleep func()

func sleepDuration(d time.Duration) sleep {
	return func() {
		time.Sleep(d)
	}
}

type DefaultHealthCheck struct {
	EndPoints       []IP
	healthEndpoints *[]IP
	Frequency       time.Duration
	CheckHealth     HealthCheckMethod
	SleepFunc       sleep
	started         bool
}

//Start Start HealthCheck
func (hk *DefaultHealthCheck) Start() {
	if hk.started {
		log.Println("Health Check already started")
		return
	}
	hk.healthEndpoints = &hk.EndPoints
	go func() {
		for {
			newEndpoint := make([]IP, 0)
			for _, ip := range hk.EndPoints {
				if hk.CheckHealth(ip) {
					newEndpoint = append(newEndpoint, ip)
				}
			}
			hk.healthEndpoints = &newEndpoint
			hk.SleepFunc()
		}
	}()
	hk.started = true
}

func (hk *DefaultHealthCheck) Receive() []IP {
	var ips []IP
	ips = *hk.healthEndpoints
	return ips
}
