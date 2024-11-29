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

// GetFreePort asks the kernel for a free open port that is ready to use.
func GetFreePort() (port int, err error) {
	var l net.Listener
	if l, err = net.Listen("tcp", "localhost:0"); err == nil {
		defer l.Close()
		port = l.Addr().(*net.TCPAddr).Port
	}
	return
}
