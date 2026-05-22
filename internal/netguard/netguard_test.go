package netguard

import (
	"net"
	"testing"
)

func TestIsPrivate_RFC1918(t *testing.T) {
	t.Parallel()
	for _, ip := range []string{
		"10.0.0.1",
		"10.255.255.255",
		"172.16.0.1",
		"172.31.0.1",
		"192.168.1.1",
	} {
		if !IsPrivate(net.ParseIP(ip)) {
			t.Fatalf("%s should be private", ip)
		}
	}
}

func TestIsPrivate_Loopback(t *testing.T) {
	t.Parallel()
	for _, ip := range []string{"127.0.0.1", "::1"} {
		if !IsPrivate(net.ParseIP(ip)) {
			t.Fatalf("%s should be private", ip)
		}
	}
}

func TestIsPrivate_LinkLocal(t *testing.T) {
	t.Parallel()
	for _, ip := range []string{"169.254.1.1", "fe80::1"} {
		if !IsPrivate(net.ParseIP(ip)) {
			t.Fatalf("%s should be private", ip)
		}
	}
}

func TestIsPrivate_CGNAT(t *testing.T) {
	t.Parallel()
	if !IsPrivate(net.ParseIP("100.64.0.1")) {
		t.Fatalf("100.64/10 should be private")
	}
	if !IsPrivate(net.ParseIP("100.127.255.255")) {
		t.Fatalf("100.127.x.x should be private")
	}
}

func TestIsPrivate_UniqueLocalIPv6(t *testing.T) {
	t.Parallel()
	if !IsPrivate(net.ParseIP("fc00::1")) {
		t.Fatalf("fc00::/7 should be private")
	}
	if !IsPrivate(net.ParseIP("fd12:3456::1")) {
		t.Fatalf("fd00::/8 should be private")
	}
}

func TestIsPrivate_PublicAllowed(t *testing.T) {
	t.Parallel()
	for _, ip := range []string{"8.8.8.8", "1.1.1.1", "2606:4700:4700::1111"} {
		if IsPrivate(net.ParseIP(ip)) {
			t.Fatalf("%s should NOT be private", ip)
		}
	}
}

func TestIsPrivate_Multicast(t *testing.T) {
	t.Parallel()
	if !IsPrivate(net.ParseIP("224.0.0.1")) {
		t.Fatalf("multicast should be private")
	}
}

func TestIsPrivate_Unspecified(t *testing.T) {
	t.Parallel()
	if !IsPrivate(net.ParseIP("0.0.0.0")) {
		t.Fatalf("0.0.0.0 should be private")
	}
}

func TestIsPrivate_Nil(t *testing.T) {
	t.Parallel()
	if !IsPrivate(nil) {
		t.Fatalf("nil should be private (fail closed)")
	}
}
