package action

import (
	"encoding/base64"
	"net"
	"sort"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

const remotePort = "10982"

type remoteCandidate struct {
	url       string
	ip        net.IP
	ifaceName string
	loopback  bool
}

func currentRemoteInfo() RemoteInfo {
	candidates := []remoteCandidate{{
		url:      "http://127.0.0.1:" + remotePort + "/remote",
		ip:       net.ParseIP("127.0.0.1").To4(),
		loopback: true,
	}}

	interfaces, err := net.Interfaces()
	if err != nil {
		return remoteInfoFromCandidates(candidates)
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

			candidates = append(candidates, remoteCandidate{
				url:       "http://" + ip.String() + ":" + remotePort + "/remote",
				ip:        ip,
				ifaceName: iface.Name,
			})
		}
	}

	return remoteInfoFromCandidates(candidates)
}

func remoteInfoFromCandidates(candidates []remoteCandidate) RemoteInfo {
	sort.SliceStable(candidates, func(i, j int) bool {
		left := remoteCandidateScore(candidates[i])
		right := remoteCandidateScore(candidates[j])
		if left != right {
			return left < right
		}

		return candidates[i].url < candidates[j].url
	})

	urls := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		urls = append(urls, candidate.url)
	}

	return remoteInfoFromUrls(dedupeStrings(urls))
}

func remoteCandidateScore(candidate remoteCandidate) int {
	if candidate.loopback {
		return 900
	}

	score := 500
	if isPrivateLAN(candidate.ip) {
		score -= 300
	}
	if isCommonHomeLAN(candidate.ip) {
		score -= 80
	}
	if isLinkLocal(candidate.ip) {
		score += 250
	}

	name := strings.ToLower(candidate.ifaceName)
	if strings.Contains(name, "wi-fi") || strings.Contains(name, "wifi") || strings.Contains(name, "wlan") || strings.Contains(name, "ethernet") || strings.Contains(name, "lan") {
		score -= 40
	}
	if strings.Contains(name, "virtual") || strings.Contains(name, "vmware") || strings.Contains(name, "virtualbox") || strings.Contains(name, "hyper-v") || strings.Contains(name, "docker") || strings.Contains(name, "vpn") || strings.Contains(name, "tailscale") || strings.Contains(name, "zerotier") || strings.Contains(name, "bluetooth") {
		score += 160
	}

	return score
}

func isPrivateLAN(ip net.IP) bool {
	ip = ip.To4()
	if ip == nil {
		return false
	}

	return ip[0] == 10 ||
		(ip[0] == 172 && ip[1] >= 16 && ip[1] <= 31) ||
		(ip[0] == 192 && ip[1] == 168)
}

func isCommonHomeLAN(ip net.IP) bool {
	ip = ip.To4()
	if ip == nil {
		return false
	}

	return ip[0] == 192 && ip[1] == 168
}

func isLinkLocal(ip net.IP) bool {
	ip = ip.To4()
	if ip == nil {
		return false
	}

	return ip[0] == 169 && ip[1] == 254
}

func remoteInfoFromUrls(urls []string) RemoteInfo {
	targets := make([]RemoteTarget, 0, len(urls))
	for _, url := range urls {
		targets = append(targets, RemoteTarget{
			Url:    url,
			QrCode: qrCodeDataURI(url),
		})
	}

	return RemoteInfo{
		Urls:    urls,
		Targets: targets,
	}
}

func qrCodeDataURI(value string) string {
	png, err := qrcode.Encode(value, qrcode.Medium, 320)
	if err != nil {
		return ""
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
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
