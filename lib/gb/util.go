package gb

import (
	"errors"
	"net"
)

// GetLocalIP get the IP address of the current system
//	physical or virtual machine
func GetLocalIP() (string, error) {
	ip, err := externalIP()
	return ip.String(), err
}

func externalIP() (net.IP, error) {
	iFaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, item := range iFaces {
		if item.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if item.Flags&net.FlagLoopback != 0 {
			continue // loop back interface
		}
		addrs, err := item.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}

	return nil, errors.New("get ip address error")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}
