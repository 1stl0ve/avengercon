package agent

import (
	"fmt"
	"net"
)

// create a list of ip addresses for the system
func getIPAddresses() ([]string, error) {
	ips := []string{}

	ifaces, err := net.Interfaces()
	if err != nil {
		return ips, err
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return ips, err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// ignore loop back and IPv6 addresses
			if !ip.IsLoopback() && ip.To4() != nil {
				ips = append(ips, fmt.Sprintf("%s: %s", i.Name, ip))
			}
		}
	}
	return ips, nil
}
