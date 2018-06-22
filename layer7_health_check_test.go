package gtm_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/datianshi/simple-cf-gtm"
	. "github.com/datianshi/simple-cf-gtm"
	"github.com/datianshi/simple-cf-gtm/fakes"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testLayer7HealthCheck(t *testing.T, when spec.G, it spec.S) {
	when("test layer 7 health check", func() {
		var (
			check       *fakes.FakeHttpCheckClient
			schema      string
			host        string
			path        string
			statusCode  int
			checkMethod HealthCheckMethod
			ip          gtm.IP
		)

		it.Before(func() {
			check = &fakes.FakeHttpCheckClient{}
			schema = "http"
			host = "abc.xip.io"
			path = "/health"
			statusCode = 200
			checkMethod = Layer7HealthCheck(check, schema, host, path, statusCode)
			ip = gtm.IP{
				Address: "x.x.x.x",
			}
		})
		when("http client return an error", func() {
			it.Before(func() {
				check.GetStub = func(url string) (*http.Response, error) {
					return nil, errors.New("")
				}
			})
			it("Should call the correct url", func() {
				checkMethod(ip)
				Ω(check.GetArgsForCall(0)).Should(Equal("http://abc.xip.io/health"))
			})
			it("Should return the false", func() {
				Ω(checkMethod(ip)).Should(Equal(false))
			})
		})
		when("http client return no error", func() {
			var (
				responseCode int
			)
			it.Before(func() {
				check.GetStub = func(url string) (*http.Response, error) {
					return &http.Response{
						StatusCode: responseCode,
					}, nil
				}
			})
			when("response code match the status code:200", func() {
				responseCode = 200
				it("Should return true", func() {
					Ω(checkMethod(ip)).Should(Equal(true))
				})
			})
			when("response code not match the status code:200", func() {
				responseCode = 400
				it("Should return false", func() {
					Ω(checkMethod(ip)).Should(Equal(false))
				})
			})
		})
	})
}
