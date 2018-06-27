package gtm

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

//UnmarshalJSON Record parsing the data
func (r *Record) UnmarshalJSON(data []byte) error {
	var v struct {
		Name              string            `json:"name"`
		IPs               []IP              `json:"ips"`
		TTL               int               `json:"ttl"`
		HealthCheckConfig HealthCheckConfig `json:"health_check"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	r.Name = v.Name
	r.IPs = v.IPs
	r.TTL = v.TTL
	r.HealthCheckConfig = v.HealthCheckConfig

	var hk HealthCheckMethod
	if r.HealthCheckConfig.Type == "layer4" {
		frequency, err := time.ParseDuration(r.HealthCheckConfig.Fequency)
		if err != nil {
			return fmt.Errorf("format of frequency is not valid: %v", err)
		}
		hk = Layer4HealthCheck(r.HealthCheckConfig.PORT)
		r.HealthCheck = &DefaultHealthCheck{
			EndPoints:   r.IPs,
			Frequency:   frequency,
			CheckHealth: hk,
			SleepFunc:   sleepDuration(frequency),
		}
	} else if r.HealthCheckConfig.Type == "layer7" {
		frequency, err := time.ParseDuration(r.HealthCheckConfig.Fequency)
		if err != nil {
			return fmt.Errorf("format of frequency is not valid: %v", err)
		}
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: r.HealthCheckConfig.SkipSSL},
		}
		client := &http.Client{Transport: tr}
		var schema string
		if r.HealthCheckConfig.HTTPS {
			schema = "https"
		} else {
			schema = "http"
		}
		hk = Layer7HealthCheck(client, schema, r.HealthCheckConfig.PATH, r.HealthCheckConfig.HTTPStatusCode)
		r.HealthCheck = &DefaultHealthCheck{
			EndPoints:   r.IPs,
			Frequency:   frequency,
			CheckHealth: hk,
			SleepFunc:   sleepDuration(frequency),
		}
	} else {
		r.HealthCheck = &DoNothingHealthCheck{
			ips: r.IPs,
		}
	}
	r.HealthCheck.Start()
	return nil
}
