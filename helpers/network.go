package helpers

import (
	"net"
)

func IsValidIP4(s string) bool {
	ip := net.ParseIP(s)
	if ip == nil {
		return false
	}
	return (ip.To4() != nil)
}

func IsValidIP6(s string) bool {
	ip := net.ParseIP(s)
	if ip == nil {
		return false
	}
	return (ip.To16() != nil)
}
