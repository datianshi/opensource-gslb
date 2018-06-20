package gtm

import (
	"log"
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

type sleep func()

func sleepSeconds(sec int64) sleep {
	return func() {
		time.Sleep(time.Duration(sec) * time.Second)
	}
}

type DefaultHealthCheck struct {
	Port            int
	EndPoints       []IP
	healthEndpoints *[]IP
	Frequency       time.Duration
	control         chan []IP
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
