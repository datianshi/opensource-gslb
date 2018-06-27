package gtm

import (
	"fmt"
	"log"
	"net/http"
)

type HttpCheckClient interface {
	Get(string) (*http.Response, error)
}

func Layer7HealthCheck(check HttpCheckClient, schema string, path string, statusCode int) HealthCheckMethod {
	return func(ip IP) bool {
		url := fmt.Sprintf("%s://%s%s", schema, ip.Host, path)
		resp, err := check.Get(url)
		if err != nil {
			log.Println(fmt.Sprintf("Health check to %s failed with error: %s", url, err))
			return false
		}
		if resp.StatusCode != statusCode {
			log.Println(fmt.Sprintf("Health Check failed to %s with status code %d", url, resp.StatusCode))
			return false
		}
		return true
	}
}
