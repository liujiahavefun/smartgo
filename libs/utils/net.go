package utils

import (
	"net"
)

func GetMacAddress() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "error_when_get_mac_address"
	}

	if len(interfaces) > 0 {
		return interfaces[0].HardwareAddr.String()
	}

	return "no_mac_hardware"
}
