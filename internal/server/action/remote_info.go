package action

import (
	"net"
	"sort"
)

const remotePort = "10982"

func currentRemoteInfo() RemoteInfo {
	urls := []string{"http://127.0.0.1:" + remotePort + "/remote"}

	interfaces, err := net.Interfaces()
	if err != nil {
		return RemoteInfo{Urls: urls}
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ip, ok := ipFromAddr(addr)
			if !ok || ip.IsLoopback() {
				continue
			}

			urls = append(urls, "http://"+ip.String()+":"+remotePort+"/remote")
		}
	}

	sort.Strings(urls)

	return RemoteInfo{Urls: dedupeStrings(urls)}
}

func ipFromAddr(addr net.Addr) (net.IP, bool) {
	switch v := addr.(type) {
	case *net.IPNet:
		ip := v.IP.To4()
		return ip, ip != nil
	case *net.IPAddr:
		ip := v.IP.To4()
		return ip, ip != nil
	default:
		return nil, false
	}
}

func dedupeStrings(values []string) []string {
	seen := map[string]bool{}
	result := []string{}

	for _, value := range values {
		if seen[value] {
			continue
		}

		seen[value] = true
		result = append(result, value)
	}

	return result
}
