package utils

import (
	"fmt"
	"net"
	"strings"
)

func SetDefault[T comparable](v *T, defaultValue T) {
	if v == nil {
		return
	}

	var zero T
	if *v == zero {
		*v = defaultValue
	}
}

func Ternary[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}

	return falseVal
}

func CompleteAddress(address string) string {
	const localhost = "localhost"

	if strings.HasPrefix(address, ":") {
		if ip, err := GetLocalIP(); err == nil {
			return ip + address
		}

		return localhost + address
	}

	return address
}

func GetLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no valid ip address found")
}
