package localip

import (
	"fmt"
	"net"
)

func GetLocalIP() (r net.IP) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	r = net.IP{}

	if err == nil {
		defer conn.Close()
		localAddress := conn.LocalAddr().(*net.UDPAddr)
		r = localAddress.IP
	}

	return
}

// ValidateIPStringInSubnet validates that an ip is valid and belongs given subnet
func ValidateIPStringInSubnet(addr string, subnet *net.IPNet) error {
	ip := net.ParseIP(addr)
	if ip == nil {
		return fmt.Errorf("\"%s\" is not a valid IP address", addr)
	}
	if subnet == nil {
		return fmt.Errorf("empty subnet")
	}
	if !subnet.Contains(ip) {
		return fmt.Errorf("IP address \"%s\" is not in subnet \"%s\"", addr, subnet.String())
	}
	return nil
}

// GetFreePort asks the kernel for a free open port that is ready to use.
func GetFreePort() (port int, err error) {
	var l net.Listener
	if l, err = net.Listen("tcp", "localhost:0"); err == nil {
		defer l.Close()
		port = l.Addr().(*net.TCPAddr).Port
	}
	return
}
