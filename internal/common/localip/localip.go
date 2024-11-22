package localip

import (
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
