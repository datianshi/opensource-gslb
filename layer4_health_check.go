package gtm

import (
	"fmt"
	"log"
	"net"
)

func Layer4HealthCheck(port int) HealthCheckMethod {
	return func(ip IP) bool {
		con, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip.Address, port))
		if err != nil {
			log.Printf("Can not connect to %s:%d, error:%v", ip.Address, port, err)
			return false
		}
		defer con.Close()
		return true
	}
}
