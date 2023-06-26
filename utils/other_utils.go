package utils

import (
	"net"
)

func IsValidIP(ip string) bool {
	res := net.ParseIP(ip)
	if res == nil {
		return false
	}
	return true
}
