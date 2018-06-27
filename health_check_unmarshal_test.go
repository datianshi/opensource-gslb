package gtm_test

import (
	"bytes"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"net/http"

	"text/template"

	. "github.com/datianshi/simple-cf-gtm"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testUnmarshalRecord(t *testing.T, when spec.G, it spec.S) {
	var recordString string
	var err error
	var record *Record

	when("Given a valid record json string", func() {
		record = &Record{}
		recordString = `
    {
      "name": "*.xip.io.",
      "ips": [
          {
            "address": "192.168.0.2"
          },
          {
            "address": "192.168.0.3"
          }
      ],
      "ttl" : 5
    }
    `
		it.Before(func() {
			err = record.UnmarshalJSON([]byte(recordString))
		})
		it("Should return no error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})
		it("Should have a do nothing healthcheck returned", func() {
			var doNothing *DoNothingHealthCheck
			t := reflect.TypeOf(doNothing)
			Ω(reflect.TypeOf(record.HealthCheck).AssignableTo(t)).Should(Equal(true))
		})
	})
	when("Given a record string with layer 4 healthcheck", func() {
		record = &Record{}
		recordString = `
    {
      "name": "*.xip.io.",
      "ips": [
          {
            "address": "192.168.0.2"
          },
          {
            "address": "192.168.0.3"
          }
      ],
      "ttl" : 5,
      "health_check" : {
        "type": "layer4",
        "port": 443,
        "frequency": "5s"
      }
    }
    `
		it.Before(func() {
			err = record.UnmarshalJSON([]byte(recordString))
		})
		it("Should return no error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})
		it("Should have a layer4 healthcheck method returned", func() {
			var healthCheck *DefaultHealthCheck
			t := reflect.TypeOf(healthCheck)
			Ω(reflect.TypeOf(record.HealthCheck).AssignableTo(t)).Should(Equal(true))
		})
	})
	when.Pend("Given a record string with layer 7 healthcheck", func() {
		var (
			s            *httptest.Server
			handlerCalls int
			h            http.Handler
			url          *url.URL
		)
		it.Before(func() {
			handlerCalls = 0
			s = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				handlerCalls = handlerCalls + 1
				h.ServeHTTP(w, req)
			}))
			url, _ = url.Parse(s.URL)
		})
		record = &Record{}
		recordString = `
    {
      "name": "*.xip.io.",
      "ips": [
          {
            "address": "192.168.0.2",
            "layer7_health_check_host": "{{.Host}}"
          },
          {
            "address": "192.168.0.3",
            "layer7_health_check_host": "{{.Host}}"
          }
      ],
      "ttl" : 5,
      "health_check" : {
        "type": "layer7",
        "https": false,
        "skip_ssl": false,
        "frequency": "5s",
        "path": "{{.Path}}",
        "http_status_code": 200
      }
    }
    `
		var err error
		it.Before(func() {
			tmpl, _ := template.New("test").Parse(recordString)
			buf := new(bytes.Buffer)
			tmpl.Execute(buf, url)
			err = record.UnmarshalJSON(buf.Bytes())
			h = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				Expect(req.Host).To(Equal(url.Host))
				w.WriteHeader(http.StatusOK)
			})
		})
		it("Should return no error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})
		it("Should run the client http get once", func() {
			time.Sleep(time.Second * 1)
			Ω(handlerCalls > 0).Should(Equal(true))
			Ω(len(record.HealthCheck.Receive())).Should(Equal(2))
		})
		it("Should get health check failed when status code not matching", func() {
			h = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				Expect(req.Host).To(Equal(url.Host))
				w.WriteHeader(http.StatusNotFound)
			})
			time.Sleep(time.Second * 1)
			Ω(len(record.HealthCheck.Receive())).Should(Equal(0))
		})
	})
}
