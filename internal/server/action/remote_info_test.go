package action

import (
	"net"
	"testing"
)

func TestRemoteInfoPrefersPrivateLANOverLoopback(t *testing.T) {
	info := remoteInfoFromCandidates([]remoteCandidate{
		{
			url:      "http://127.0.0.1:10982/remote",
			ip:       net.ParseIP("127.0.0.1").To4(),
			loopback: true,
		},
		{
			url:       "http://192.168.178.30:10982/remote",
			ip:        net.ParseIP("192.168.178.30").To4(),
			ifaceName: "Wi-Fi",
		},
	})

	if got := info.Urls[0]; got != "http://192.168.178.30:10982/remote" {
		t.Fatalf("expected private LAN first, got %q", got)
	}
}

func TestRemoteInfoPenalizesVirtualInterfaces(t *testing.T) {
	info := remoteInfoFromCandidates([]remoteCandidate{
		{
			url:       "http://172.20.0.2:10982/remote",
			ip:        net.ParseIP("172.20.0.2").To4(),
			ifaceName: "Docker",
		},
		{
			url:       "http://10.0.0.12:10982/remote",
			ip:        net.ParseIP("10.0.0.12").To4(),
			ifaceName: "Ethernet",
		},
	})

	if got := info.Urls[0]; got != "http://10.0.0.12:10982/remote" {
		t.Fatalf("expected physical LAN first, got %q", got)
	}
}

func TestRemoteInfoKeepsLoopbackFallback(t *testing.T) {
	info := remoteInfoFromCandidates([]remoteCandidate{
		{
			url:      "http://127.0.0.1:10982/remote",
			ip:       net.ParseIP("127.0.0.1").To4(),
			loopback: true,
		},
	})

	if got := info.Urls[0]; got != "http://127.0.0.1:10982/remote" {
		t.Fatalf("expected loopback fallback, got %q", got)
	}
}
