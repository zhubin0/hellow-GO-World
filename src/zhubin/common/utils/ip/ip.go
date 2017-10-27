package ip

import (
	"fmt"
	"net"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

//get the ip address of the local machine
// sort 192 172 and 10 ,and return the first one.
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	var ips = make([]string, 0)
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.IsGlobalUnicast() && !ipnet.IP.IsLinkLocalMulticast() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	netCardNum := len(ips)
	if netCardNum <= 0 {
		return ""
	}
	sort.Sort(sort.Reverse(sort.StringSlice(ips)))
	return ips[0]
}

func GetAllLocalIp() []string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return []string{}
	}
	var ips = make([]string, 0)
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && ipnet.IP.IsGlobalUnicast() && !ipnet.IP.IsLinkLocalMulticast() {
			if ipnet.IP.To4() != nil { // only keep ipv4
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips
}

func GetPrivateIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	var ips = make([]string, 0)
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	netCardNum := len(ips)
	if netCardNum == 1 {
		return ips[0]
	}

	for i, _ := range ips {
		if strings.HasPrefix(ips[i], "192") || strings.HasPrefix(ips[i], "172") || strings.HasPrefix(ips[i], "10") {
			return ips[i]
		}
	}

	if netCardNum > 0 {
		return ips[0]
	}

	return ""
}

func ValidIP4(ipAddress string) bool {
	ipAddress = strings.Trim(ipAddress, " ")

	re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	if re.MatchString(ipAddress) {
		return true
	}
	return false
}

func GetClientIp(req *http.Request) (string, error) {
	if req.RemoteAddr == "" {
		return "", fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}

	forward := req.Header.Get("X-Forwarded-For")
	return strings.Split(forward, ",")[0], nil
}
