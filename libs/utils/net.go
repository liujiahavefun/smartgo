package utils

import (
	"net"
    "os"
    "runtime"
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

func printStack() {
    var buf [4096]byte
    n := runtime.Stack(buf[:], false)
    os.Stderr.Write(buf[:n])
}