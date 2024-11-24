package util

import (
	"fmt"
	"net"
)

func GetInterfaceInfo(name string) (string, string, error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return "", "", fmt.Errorf("failed to get interface: %v", err)
	}

	// Get MAC address
	macAddr := iface.HardwareAddr.String()

	// Get IP address
	addrs, err := iface.Addrs()
	if err != nil {
		return "", "", fmt.Errorf("failed to get addresses: %v", err)
	}

	var ipAddr string
	for _, addr := range addrs {
		// Check if address is IP network
		if ipnet, ok := addr.(*net.IPNet); ok {
			// Use only IPv4 address
			if ip4 := ipnet.IP.To4(); ip4 != nil {
				ipAddr = ip4.String()
				break
			}
		}
	}

	if ipAddr == "" {
		return "", "", fmt.Errorf("no IPv4 address found")
	}

	return ipAddr, macAddr, nil
}
