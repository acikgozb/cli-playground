package scan

import (
	"fmt"
	"net"
	"time"
)

// PortState represents the state of a single TCP port.
type PortState struct {
	Port int
	Open state
}

type Results struct {
	Host       string
	NotFound   bool
	PortStates []PortState
}

// By creating custom types, we can associate methods to it.
type state bool

// String converts the boolean value of state to a human readable string
func (s state) String() string {
	if s {
		return "open"
	}

	return "closed"
}

// scanPort performs a port scan on a single TCP port
func scanPort(host string, port int) PortState {
	p := PortState{
		Port: port,
	}

	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	scanConn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		// Assume that error means the port is not open.
		return p
	}

	scanConn.Close()
	p.Open = true
	return p
}

// Run performs a port scan on the hosts list.
func Run(hl *HostsList, ports []int) []Results {
	res := make([]Results, 0, len(hl.Hosts))

	for _, h := range hl.Hosts {
		r := Results{
			Host: h,
		}

		// Resolve the host name into a volid IP address.
		if _, err := net.LookupHost(h); err != nil {
			r.NotFound = true
			res = append(res, r)
			continue
		}

		for _, p := range ports {
			r.PortStates = append(r.PortStates, scanPort(h, p))
		}

		res = append(res, r)
	}

	return res
}
