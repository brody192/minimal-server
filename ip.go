package main

import (
	"fmt"
	"net"
	"strings"
)

func init() {
	ips, err := GetInterfaceIPs("railnet0", "ipv6")
	if err != nil {
		fmt.Println("Error getting interface IPs:", err)
	}

	HTTP_RESP_IP = ips[0].String()
}

func GetInterfaceIPs(networkInterface string, networkType string) ([]net.IP, error) {
	if !strings.EqualFold(networkType, "ipv4") && !strings.EqualFold(networkType, "ipv6") {
		return nil, fmt.Errorf("invalid network type: %s", networkType)
	}

	is, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	interfaceFound := false

	for _, ifi := range is {
		if strings.EqualFold(ifi.Name, networkInterface) {
			interfaceFound = true
			break
		}
	}

	if !interfaceFound {
		return nil, fmt.Errorf("interface %s not found", networkInterface)
	}

	ret := []net.IP{}

	for _, ifi := range is {
		// skip down interfaces
		if ifi.Flags&net.FlagUp == 0 {
			continue
		}

		// skip loopback interfaces
		if ifi.Flags&net.FlagLoopback != 0 {
			continue
		}

		// skip unnamed interfaces
		if ifi.Name == "" {
			continue
		}

		// skip interfaces that are not the one we are looking for
		if !strings.EqualFold(ifi.Name, networkInterface) {
			continue
		}

		addrs, err := ifi.Addrs()
		if err != nil {
			return nil, err
		}

		// skip interfaces without addresses
		if len(addrs) == 0 {
			continue
		}

		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				return nil, err
			}

			// skip link-local addresses
			if ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
				continue
			}

			if strings.EqualFold(networkType, "ipv4") && ip.To4() == nil {
				continue
			}

			if strings.EqualFold(networkType, "ipv6") && ip.To4() != nil {
				continue
			}

			ret = append(ret, ip)
		}
	}

	if len(ret) == 0 {
		return nil, fmt.Errorf("no %s addresses found for interface %s", networkType, networkInterface)
	}

	return ret, nil
}
