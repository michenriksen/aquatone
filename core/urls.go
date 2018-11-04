package core

import (
	"fmt"
)

var (
	securePorts = []int{443, 832, 981, 1010, 1311, 2083, 2087, 2095, 2096, 4712,
		7000, 8172, 8243, 8333, 8443, 8834, 9443, 12443, 18091, 18092}
)

func HostAndPortToURL(host string, port int, protocol string) string {
	var url string
	if protocol != "" {
		url = fmt.Sprintf("%s://%s", protocol, host)
	} else if isSecurePort(port) {
		url = fmt.Sprintf("https://%s", host)
	} else {
		url = fmt.Sprintf("http://%s", host)
	}
	if isStandardPort(port) {
		url = fmt.Sprintf("%s/", url)
	} else {
		url = fmt.Sprintf("%s:%d/", url, port)
	}
	return url
}

func isSecurePort(port int) bool {
	for _, p := range securePorts {
		if p == port {
			return true
		}
	}
	return false
}

func isStandardPort(port int) bool {
	if port == 80 || port == 443 {
		return true
	}
	return false
}
