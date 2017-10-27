package ip

import (
	"testing"
)

// NOTE for test on other machines
//
//package main
//
//import (
//"fmt"
//"net"
//)
//
//
//
//func GetAllLocalIp() []string {
//	addrs, err := net.InterfaceAddrs()
//	if err != nil {
//		return []string{}
//	}
//	var ips = make([]string, 0)
//	for _, address := range addrs {
//		// check the address type and if it is not a loopback the display it
//		if ipnet, ok := address.(*net.IPNet); ok && ipnet.IP.IsGlobalUnicast() && !ipnet.IP.IsLinkLocalMulticast() {
//			if ipnet.IP.To4() != nil { // only keep ipv4
//				ips = append(ips, ipnet.IP.String())
//			}
//		}
//	}
//	return ips
//}
//
//func main() {
//	fmt.Println(GetAllLocalIp())
//}
func TestGetAllLocalIp(t *testing.T) {
	t.Log(GetAllLocalIp())
}
